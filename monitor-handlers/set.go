package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	consulapi "github.com/hashicorp/consul/api"
)

var leader string

func Set(ip, port, username, password string, issynchronous int,logvalue string) string {
	leader = "cmha/service/" + servicename + "/db/leader"
	config := &consulapi.Config{
		Datacenter: beego.AppConfig.String("datacenter"),
		Token:      beego.AppConfig.String("token"),
		Address:    address,
	}
	var healthpair []*consulapi.ServiceEntry
	var kvPair *consulapi.KVPair
	var client *consulapi.Client
	var health *consulapi.Health
	var err error
	c := 0
	for i := 0; i < 3; i++ {
		client, err = consulapi.NewClient(config)
		if err != nil {
			c += 1
			logger.Println("[E] Create consul-api client failed!", err)
			timestamp := time.Now().Unix()
			logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + consulapi_failed + "{{" + fmt.Sprintf("%s", err)
			continue
		}
		logger.Println("[I] Create consul-api client successfully!")
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + consulapi_success
		break
	}
	if c >= 3 {
		logger.Println("[I] Give up switching to async replication!")
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + give_async_repl
		return logvalue
	}
	kv = client.KV()
	health = client.Health()
	logvalue = SetRepl_err_counter(other_hostname, client,logvalue)
	kvPair, _, err = kv.Get(leader, nil)
	if err != nil {
		logger.Println("[E] Get and check current service leader from CS failed!", err)
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + get_leader_failed + "{{" + fmt.Sprintf("%s", err)
		count := 0
		for i := 0; i < 3; i++ {
			client, err = consulapi.NewClient(config)
			if err != nil {
				count += 1
				logger.Println("[E] Create consul-api client failed!", err)
				timestamp := time.Now().Unix()
				logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + consulapi_failed + "{{" + fmt.Sprintf("%s", err)
				continue
			}
			logger.Println("[I] Create consul-api client successfully!")
			timestamp := time.Now().Unix()
			logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + consulapi_success
			kvPair, _, err = kv.Get(leader, nil)
			if err != nil {
				count += 1
				logger.Println("[E] Get and check current service leader from CS failed!", err)
				timestamp := time.Now().Unix()
				logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + get_leader_failed + "{{" + fmt.Sprintf("%s", err)
				continue
			}
			break
		}
		if count >= 3 {
			logger.Println("[I] Give up switching to async replication!")
			logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + give_async_repl
			return logvalue
		}
	}
	logger.Println("[I] Get and check current service leader from CS successfully!")
	timestamp := time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + get_leader_success
	content := string(kvPair.Value)
	istrue := strings.Contains(content, ip)
	if !istrue {
		logger.Println("[I] " + ip + " is not service leader!")
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + no_leader + "{{" + ip
		logger.Println("[I] Give up switching to async replication!")
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + give_async_repl
		return logvalue
	}
	logger.Println("[I] " + ip + " is service leader!")
	timestamp = time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + is_leader + "{{" + ip
	healthpair, _, err = health.Service(servicename, "", false, nil)
	if err != nil {
		logger.Println("[E] Get "+servicename+" service health status from CS failed!", err)
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + get_service_status_failed + "{{" + servicename + "{{" + fmt.Sprintf("%s", err)
		count := 0
		for i := 0; i < 3; i++ {
			client, err = consulapi.NewClient(config)
			if err != nil {
				count += 1
				logger.Println("[E] Create consul-api client failed!", err)
				timestamp := time.Now().Unix()
				logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + consulapi_failed + "{{" + fmt.Sprintf("%s", err)
				continue
			}
			logger.Println("[I] Create consul-api client successfully!")
			timestamp := time.Now().Unix()
			logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + consulapi_success
			healthpair, _, err = health.Service(servicename, "", false, nil)
			if err != nil {
				count += 1
				logger.Println("[E] Get "+servicename+" service health status from CS failed!", err)
				timestamp := time.Now().Unix()
				logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + get_service_status_failed + "{{" + servicename + "{{" + fmt.Sprintf("%s", err)
				continue
			}
			break
		}
		if count >= 3 {
			logger.Println("[I] Give up switching to async replication!")
			logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + give_async_repl
			return logvalue
		}
	}
	logger.Println("[I] Get " + servicename + " service health status from CS successfully!")
	timestamp = time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + get_service_status_success + "{{" + servicename
	ishealthy, status := GetIsHealthAndStatus(ip, healthpair)
	if !ishealthy {
		logger.Println("[I] " + ip + " service health status is " + status)
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + service_status + "{{" + ip + "{{" + status
		logger.Println("[I] Give up switching to async replication!")
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + give_async_repl
		return logvalue
	}
	logger.Println("[I] " + ip + " service health status is " + status)
	timestamp = time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + service_status + "{{" + ip + "{{" + status
	dsName := username + ":" + password + "@tcp(" + "localhost" + ":" + port + ")/"
	db, err := sql.Open("mysql", dsName)
	if err != nil {
		logger.Println("[E] Create connection object to local database failed!", err)
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + create_database_object_failed + "{{" + fmt.Sprintf("%s", err)
		logger.Println("[I] Give up switching to async replication!")
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + give_async_repl
		return logvalue
	}
	logger.Println("[I] Create connection object to local database successfully!")
	timestamp = time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + create_database_object_success
	defer db.Close()
	err = db.Ping()
	if err != nil {
		logger.Println("[E] Connected to local database failed!", err)
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + connected_database_failed + "{{" + fmt.Sprintf("%s", err)
		logger.Println("[I] Give up switching to async replication!")
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + give_async_repl
		return logvalue
	}
	logger.Println("[I] Connected to local database successfully!")
	timestamp = time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + connected_database_success
	keepsyncrepl := "set global rpl_semi_sync_master_keepsyncrepl=" + strconv.Itoa(issynchronous)
	_, err = db.Query(keepsyncrepl)
	if err != nil {
		logger.Println("[E] Set rpl_semi_sync_master_keepsyncrepl=0 failed!", err)
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + set_keepsyncrepl_failed + "{{" + fmt.Sprintf("%s", err)
		return logvalue
	}
	logger.Println("[I] Set rpl_semi_sync_master_keepsyncrepl=0 successfully!")
	timestamp = time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + set_keepsyncrepl_success
	trysyncrepl := "set global rpl_semi_sync_master_trysyncrepl=" + strconv.Itoa(issynchronous)
	_, err = db.Query(trysyncrepl)
	if err != nil {
		logger.Println("[E] Set rpl_semi_sync_master_trysyncrepl=0 failed!", err)
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + set_trysyncrepl_failed + "{{" + fmt.Sprintf("%s", err)
		return logvalue
	}
	logger.Println("[I] Set rpl_semi_sync_master_trysyncrepl=0 successfully!")
	timestamp = time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + set_trysyncrepl_success
	logger.Println("[I] Switching local database to async replication!")
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + switch_local_async_repl
	return logvalue
}
