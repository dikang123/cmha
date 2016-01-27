package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	consulapi "github.com/hashicorp/consul/api"
)

func Set(ip, port, username, password string, issynchronous int) {
	servicename = beego.AppConfig.String("servicename")
	service_ip = beego.AppConfig.Strings("service_ip")
	other_hostname = beego.AppConfig.String("otherhostname")
	tag = beego.AppConfig.String("tag")
	config := &consulapi.Config{
		Datacenter: beego.AppConfig.String("datacenter"),
		Token:      beego.AppConfig.String("token"),
	}
	var healthpair []*consulapi.ServiceEntry
	var kvPair *consulapi.KVPair
	var client *consulapi.Client
	var health *consulapi.Health
	var err error
	c := len(service_ip)
	for i, _ := range service_ip {
		config.Address = service_ip[i] + ":" + beego.AppConfig.String("service_port")
		client, err = consulapi.NewClient(config)
		if err != nil {
			c--
			logger.Println("[E] Create consul-api client failed! CS ip = "+service_ip[i], err)
			timestamp := time.Now().Unix()
			logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + consulapi_failed + "{{" + service_ip[i] + "{{" + fmt.Sprintf("%s", err)
			continue
		}
		logger.Println("[I] Create consul-api client successfully! CS ip = " + service_ip[i])
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + consulapi_success + "{{" + service_ip[i]
		break
	}
	if c == 0 {
		logger.Println("[I] Give up switching to async replication!")
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + give_async_repl
		//	UploadLog(logkey,logvalue)
		return
	}
	kv = client.KV()
	health = client.Health()
	SetRepl_err_counter(other_hostname)
	kvPair, _, err = kv.Get("service/"+servicename+"/leader", nil)
	if err != nil {
		logger.Println("[E] Get and check current service leader from CS failed!", err)
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + get_leader_failed + "{{" + fmt.Sprintf("%s", err)
		count := len(service_ip)
		for i, _ := range service_ip {
			config.Address = service_ip[i] + ":" + beego.AppConfig.String("service_port")
			client, err = consulapi.NewClient(config)
			if err != nil {
				count--
				logger.Println("[E] Create consul-api client failed! CS ip= "+service_ip[i], err)
				timestamp := time.Now().Unix()
				logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + consulapi_failed + "{{" + service_ip[i] + "{{" + fmt.Sprintf("%s", err)
				continue
			}
			logger.Println("[I] Create consul-api client successfully! CS ip= " + service_ip[i])
			timestamp := time.Now().Unix()
			logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + consulapi_success + "{{" + service_ip[i]
			kvPair, _, err = kv.Get("service/"+servicename+"/leader", nil)
			if err != nil {
				count--
				logger.Println("[E] Get and check current service leader from CS failed! CS ip = "+service_ip[i], err)
				timestamp := time.Now().Unix()
				logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + get_leader_failed + "{{" + service_ip[i] + "{{" + fmt.Sprintf("%s", err)
				continue
			}
			break
		}
		if count == 0 {
			logger.Println("[I] Give up switching to async replication!")
			logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + give_async_repl
			//	UploadLog(logkey,logvalue)
			return
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
		//		UploadLog(logkey,logvalue)
		return
	}
	logger.Println("[I] " + ip + " is service leader!")
	timestamp = time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + is_leader + "{{" + ip
	healthpair, _, err = health.Service(servicename, "", false, nil)
	if err != nil {
		logger.Println("[E] Get "+servicename+" service health status from CS failed!", err)
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + get_service_status_failed + "{{" + servicename + "{{" + fmt.Sprintf("%s", err)
		count := len(service_ip)
		for i, _ := range service_ip {
			config.Address = service_ip[i] + ":" + beego.AppConfig.String("service_port")
			client, err = consulapi.NewClient(config)
			if err != nil {
				count--
				logger.Println("[E] Create consul-api client failed! CS ip= "+service_ip[i], err) //////
				timestamp := time.Now().Unix()
				logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + consulapi_failed + "{{" + service_ip[i] + "{{" + fmt.Sprintf("%s", err)
				continue
			}
			logger.Println("[I] Create consul-api client successfully! CS ip= " + service_ip[i])
			timestamp := time.Now().Unix()
			logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + consulapi_success + "{{" + service_ip[i]
			healthpair, _, err = health.Service(servicename, "", false, nil)
			if err != nil {
				count--
				logger.Println("[E] Get "+servicename+" service health status from CS failed! CS ip = "+service_ip[i], err)
				timestamp := time.Now().Unix()
				logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + get_service_status_failed + "{{" + servicename + "{{" + service_ip[i] + "{{" + fmt.Sprintf("%s", err)
				continue
			}
			break
		}
		if count == 0 {
			logger.Println("[I] Give up switching to async replication!")
			logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + give_async_repl
			//		UploadLog(logkey,logvalue)
			return
		}
	}
	logger.Println("[I] Get " + servicename + " service health status from CS successfully!")
	timestamp = time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + get_service_status_success + "{{" + servicename
	var ishealthy = true
	var status string
	for index := range healthpair {
		if healthpair[index].Node.Address == ip {
			for checkindex := range healthpair[index].Checks {
				if healthpair[index].Checks[checkindex].Status == "critical" {
					ishealthy = false
					status = healthpair[index].Checks[checkindex].Status
					break
				} else if healthpair[index].Checks[checkindex].Status == "warning" {
					status = healthpair[index].Checks[checkindex].Status
					break
				} else {
					status = healthpair[index].Checks[checkindex].Status
					break
				}
			}
		}
	}
	if !ishealthy {
		logger.Println("[I] " + ip + " service health status is " + status)
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + service_status + "{{" + ip + "{{" + status
		logger.Println("[I] Give up switching to async replication!")
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + give_async_repl
		//		UploadLog(logkey,logvalue)
		return
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
		//		UploadLog(logkey,logvalue)
		return
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
		//		UploadLog(logkey,logvalue)
		return
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
		//		UploadLog(logkey,logvalue)
		return
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
		//		UploadLog(logkey,logvalue)
		return
	}
	logger.Println("[I] Set rpl_semi_sync_master_trysyncrepl=0 successfully!")
	timestamp = time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + set_trysyncrepl_success
	logger.Println("[I] Switching local database to async replication!")
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + switch_local_async_repl
}

func SetRepl_err_counter(hostname string) {
	config := &consulapi.Config{
		Datacenter: beego.AppConfig.String("datacenter"),
		Token:      beego.AppConfig.String("token"),
	}
	cnt := 0
	var put string
	put = "1"
	kvvalue := []byte(put)
	kvotherhostname := consulapi.KVPair{
		Key:   "monitor/" + hostname,
		Value: kvvalue,
	}
try:
	_, err := kv.Put(&kvotherhostname, nil)
	if err != nil {
		logger.Println("[E] Set peer database repl_err_counter to 1 in CS failed!", err)
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + set_counter_failed + "{{" + fmt.Sprintf("%s", err)
		if cnt == 2 {
			//return
			//		UploadLog(logkey,logvalue)
			count := len(service_ip)
			for i, _ := range service_ip {
				config.Address = service_ip[i] + ":" + beego.AppConfig.String("service_port")
				client, err = consulapi.NewClient(config)
				if err != nil {
					count--
					logger.Println("[E] Create consul-api client failed! CS ip= "+service_ip[i], err) //////
					timestamp := time.Now().Unix()
					logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + consulapi_failed + "{{" + service_ip[i] + "{{" + fmt.Sprintf("%s", err)
					continue
				}
				logger.Println("[I] Create consul-api client successfully! CS ip= " + service_ip[i])
				timestamp := time.Now().Unix()
				logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + consulapi_success + "{{" + service_ip[i]
				ct := 0
			tr:
				_, err = kv.Put(&kvotherhostname, nil)
				if err != nil {
					count--
					logger.Println("[E] Set peer database repl_err_counter to 1 in CS failed! CS ip = "+service_ip[i], err)
					timestamp := time.Now().Unix()
					logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + set_counter_failed + "{{" + service_ip[i] + "{{" + fmt.Sprintf("%s", err)
					if ct == 2 {
						continue
					} else {
						ct++
						goto tr
					}
				}
				break
			}
			if count == 0 {
				UploadLog(logkey, logvalue)
				os.Exit(1)
			}
		} else {
			cnt++
			goto try
		}
	}
	logger.Println("[I] Set peer database repl_err_counter to 1 in CS successfully!")
	timestamp := time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + set_counter_success
}

func checkio_thread(ip, port, username, password, addr string) {
	logger.Println("[I] Monitor Handler Sleep 60s!")
	timestamp := time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + handlers_sleep
	time.Sleep(60000 * time.Millisecond)
	logger.Println("[I] Connecting to peer database......")
	timestamp = time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + connecting_database
	dsName := username + ":" + password + "@tcp(" + addr + ":" + port + ")/"
	db, err := sql.Open("mysql", dsName)
	if err != nil {
		logger.Println("[E] Create connection object to peer database failed!", err)
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + create_peer_database_object_failed + "{{" + fmt.Sprintf("%s", err)
		logger.Println("[I] Give up switching to async replication!")
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + give_async_repl
		//	UploadLog(logkey,logvalue)
		return
	}
	logger.Println("[I] Create connection object to peer database successfully!")
	timestamp = time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + create_peer_database_object_success
	defer db.Close()
	err = db.Ping()
	if err != nil {
		logger.Println("[E] Connected to the peer database failed!", err)
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + connected_peer_database_failed + "{{" + fmt.Sprintf("%s", err)
		Set(ip, port, username, password, 0)
		//		UploadLog(logkey,logvalue)
		return
	}
	logger.Println("[I] Connected to the peer database successfully!")
	timestamp = time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + connected_peer_database_success
	row, err := db.Query("show slave status")
	if err != nil {
		logger.Println("[I] Checking peer database I/O thread status. Failed!", err)
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + check_perr_dabases_io_failed + "{{" + fmt.Sprintf("%s", err)
		Set(ip, port, username, password, 0)
		//		UploadLog(logkey,logvalue)
		return
	}
	logger.Println("[I] Checking peer database I/O thread status. Successfully!")
	timestamp = time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + check_perr_dabases_io_success
	cols, _ := row.Columns()
	buffer := make([]interface{}, len(cols))
	data := make([]interface{}, len(cols))
	for i, _ := range buffer {
		buffer[i] = &data[i]
	}
	for row.Next() {
		err = row.Scan(buffer...)
		if err != nil {
			logger.Println("[E] Resolve slave status failed!", err)
			timestamp := time.Now().Unix()
			logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + resolve_slave_status_failed + "{{" + fmt.Sprintf("%s", err)
			logger.Println("[I] Give up switching to async replication!")
			logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + give_async_repl
			//		UploadLog(logkey,logvalue)
			return
		}
	}
	mapField2Data := make(map[string]interface{}, len(cols))
	for k, col := range data {
		mapField2Data[cols[k]] = col
	}
	Slave_IO_Running := mapField2Data["Slave_IO_Running"]
	if string(Slave_IO_Running.([]uint8)) != "Yes" {
		logger.Println("[I] The I/O thread status is " + string(Slave_IO_Running.([]uint8)) + "!")
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + io_status + "{{" + string(Slave_IO_Running.([]uint8))
		Set(ip, port, username, password, 0)
		//		UploadLog(logkey,logvalue)
		return
	}
	logger.Println("[I] The I/O thread status is " + string(Slave_IO_Running.([]uint8)) + "!")
	timestamp = time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + io_status + "{{" + string(Slave_IO_Running.([]uint8))
	logger.Println("[I] Give up switching to async replication!")
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + give_async_repl
}
