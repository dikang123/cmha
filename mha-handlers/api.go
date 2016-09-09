package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	//	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	consulapi "github.com/hashicorp/consul/api"
)

var address string
var servicename string
var hostname string
var port string
var username string
var password string
var kv *consulapi.KV
var client *consulapi.Client
var logvalue string
var logkey string
var leader string
var err error


func SessionAndChecks() {
	logger.Println("[I] MHA Handler Triggered")
	logvalue, timestamp := GetLogValue(triggered)
	address := ReadCaAddress()
	ip, port, username, password := ReadDatabaseCfg()
	servicename = ReadServiceName()
	hostname = ReadHostName()
	logkey = GetLogKey(servicename, hostname, timestamp)
	config := GetConfig(address)
	client, logvalue, err = GetClient(config, logvalue, consulapi_failed)
	if err != nil {
		UploadLog(logkey, logvalue)
	}
	logger.Println("[I] Create consul-api client successfully!")

	logvalue, timestamp = AddLogValue(consulapi_success, logvalue)
	var kvPair *consulapi.KVPair
	var kvMonitor *consulapi.KVPair
	c := 0
	leader = GetLeaderKey(servicename)
	for i := 0; i < 3; i++ {
		//KV is used to return a handle to the K/V apis
		kv = client.KV()
		//Get is used to lookup a single key
		kvPair, _, err = kv.Get(leader, nil)
		if err != nil {
			c += 1
			logger.Println("[E] Get and check current service leader from CS failed!", err)
			timestamp := time.Now().Unix()
			logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + current_check_failed + "{{" + fmt.Sprintf("%s", err)
			continue
		}
		break
	}
	if c >= 3 {
		logvalue, timestamp = AddLogValue(give_election, logvalue)
		logger.Println("[I] Give up leader election")
		UploadLog(logkey, logvalue)
		return
	}
	logger.Println("[I] Get and check current service leader from CS successfully!")
	logvalue, timestamp = AddLogValue(current_check_success, logvalue)
	kvMonitor, _, err = kv.Get("cmha/service/"+servicename+"/db/repl_err_counter/"+hostname, nil)
	var kvValue string
	kvValue = string(kvMonitor.Value)
	if err != nil {
		logger.Println("[E] Get "+ip+" repl_err_counter = "+kvValue+" failed!", err)
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + get_counter_failed + "{{" + ip + "{{" + kvValue + "{{" + fmt.Sprintf("%s", err)
		count := 0
		for i := 0; i < 3; i++ {
			client, err = consulapi.NewClient(config)
			if err != nil {
				count += 1
				logger.Println("[E] Create consul-api client failed!", err) //////
				timestamp := time.Now().Unix()
				logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + consulapi_failed + "{{" + fmt.Sprintf("%s", err)
				continue
			}
			logger.Println("[I] Create consul-api client successfully!")
			logvalue, timestamp = AddLogValue(consulapi_success, logvalue)
			//Get is used to lookup a single key
			kvMonitor, _, err = kv.Get("cmha/service/"+servicename+"/db/repl_err_counter/"+hostname, nil)
			kvValue = string(kvMonitor.Value)
			if err != nil {
				count += 1
				logger.Println("[E] Get "+ip+" repl_err_counter="+kvValue+" failed!", err)
				timestamp := time.Now().Unix()
				logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + get_counter_failed + "{{" + ip + "{{" + kvValue + "{{" + fmt.Sprintf("%s", err)
				continue
			}
			break
		}
		if count >= 3 {
			logger.Println("[I] Give up leader election")
			logvalue, timestamp = AddLogValue(give_election, logvalue)
			UploadLog(logkey, logvalue)
			return
		}
	}
	logger.Println("[I] Get " + ip + " repl_err_counter=" + kvValue + " successfully!")
	timestamp = time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + get_counter_success + "{{" + ip + "{{" + kvValue
	if kvValue != "0" {
		logger.Println("[E] " + ip + " give up leader election")
		timestamp = time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + ip_election + "{{" + ip
		UploadLog(logkey, logvalue)
		return
	}
	//NewClient returns a new client
	if kvPair == nil {
		logger.Println("[E] Not service leader,Please create!")
		logvalue, timestamp = AddLogValue(create_counter, logvalue)
		logger.Println("[I] Give up leader election")
		logvalue, timestamp = AddLogValue(give_election, logvalue)
		UploadLog(logkey, logvalue)
		return
	}
	//Are there external connection string provided
	if kvPair.Session != "" {
		logger.Println("[I] Leader exist!")
		logvalue, timestamp = AddLogValue(leader_exist, logvalue)
		logger.Println("[I] Give up leader election")
		logvalue, timestamp = AddLogValue(give_election, logvalue)
		UploadLog(logkey, logvalue)
		return
	}
	logger.Println("[I] Leader does not exist!")
	logvalue, timestamp = AddLogValue(leader_noexist, logvalue)
	SetRead_only(username, password, port, 1)
	//Health returns a handle to the health endpoints
	health := client.Health()
	//Checks is used to return the checks associated with a service
	healthvalue, _, err := health.Checks(servicename, nil)
	if err != nil {
		logger.Println("[E] Get and check "+ip+" service health status failed!", err)
		timestamp = time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + get_health_failed + "{{" + ip + "{{" + fmt.Sprintf("%s", err)
		UploadLog(logkey, logvalue)
		return
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
			logvalue, timestamp = AddLogValue(consulapi_success, logvalue)
			healthvalue, _, err = health.Checks(servicename, nil)
			if err != nil {
				count += 1
				logger.Println("[E] Get and check "+ip+" service health status failed!", err)
				timestamp = time.Now().Unix()
				logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + get_health_failed + "{{" + ip + "{{" + fmt.Sprintf("%s", err)
				continue
			}
			break
		}
		if count >= 3 {
			logger.Println("[I] Give up leader election")
			logvalue, timestamp = AddLogValue(give_election, logvalue)
			UploadLog(logkey, logvalue)
			return
		}
	}
	logger.Println("[I] Get and check " + ip + " service health status successfully!")
	timestamp = time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + get_health_success + "{{" + ip
	if len(healthvalue) <= 0 {
		logger.Println("[I] " + servicename + " service does not exist!")
		timestamp = time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + service_noexist + "{{" + servicename
		logger.Println("[I] Give up leader election")
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + give_election
		UploadLog(logkey, logvalue)
		return
	}
	var islocal bool
	for index := range healthvalue {
		if healthvalue[index].Node == hostname {
			islocal = true
			logger.Println("[I] " + servicename + " service exist!")
			timestamp = time.Now().Unix()
			logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + service_exist + "{{" + servicename
			break
		}

	}
	if !islocal {
		logger.Println("[E] " + ip + " not is " + servicename + "!")
		timestamp = time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + ip_noservice + "{{" + ip + "{{" + servicename
		logger.Println("[I] Give up leader election")
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + give_election
		UploadLog(logkey, logvalue)
		return
	} else {
		updatevalue := consulapi.KVPair{
			Key:   leader,
			Value: []byte(""),
		}
		_, err = kv.Put(&updatevalue, nil)
		if err != nil {
			logger.Println("[E] Clean service leader value in CS failed!", err)
			timestamp = time.Now().Unix()
			logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + clean_kv_failed + "{{" + fmt.Sprintf("%s", err)
			UploadLog(logkey, logvalue)
			return
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
				logvalue, timestamp = AddLogValue(consulapi_success, logvalue)
				_, err = kv.Put(&updatevalue, nil)
				if err != nil {
					count += 1
					logger.Println("[E] Clean service leader value in CS failed!", err)
					timestamp = time.Now().Unix()
					logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + clean_kv_failed + "{{" + fmt.Sprintf("%s", err)
					continue
				}
				break
			}
			if count >= 3 {
				logger.Println("[I] Give up leader election")
				logvalue, timestamp = AddLogValue(give_election, logvalue)
				UploadLog(logkey, logvalue)
				return
			}
		}
		logger.Println("[I] Clean service leader value in CS successfully!")
		logvalue, timestamp = AddLogValue(clean_kv_success, logvalue)
		healthpair, _, err := health.Service(servicename, "", false, nil)
		if err != nil {
			logger.Println("[E] Get and check "+ip+" service health status failed!", err)
			timestamp = time.Now().Unix()
			logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + get_health_failed + "{{" + ip + "{{" + fmt.Sprintf("%s", err)
			UploadLog(logkey, logvalue)
			return
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
				logvalue, timestamp = AddLogValue(consulapi_success, logvalue)
				healthpair, _, err = health.Service(servicename, "", false, nil)
				if err != nil {
					count += 1
					logger.Println("[E] Get and check "+ip+" service health status failed!", err)
					timestamp = time.Now().Unix()
					logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + get_health_failed + "{{" + ip + "{{" + fmt.Sprintf("%s", err)
					continue
				}
				break
			}
			if count >= 3 {
				logger.Println("[I] Give up leader election")
				logvalue, timestamp = AddLogValue(give_election, logvalue)
				UploadLog(logkey, logvalue)
				return
			}
		}
		logger.Println("[I] Get and check " + ip + " service health status successfully!")
		timestamp = time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + get_health_success + "{{" + ip
		_, _ = IsStatus(healthpair, ip)
	}
}


func SetRead_only(username, password, port string, value int) {
	var timestamp int64
	dsName := username + ":" + password + "@tcp(" + "localhost" + ":" + port + ")/"
	db, err := sql.Open("mysql", dsName)
	if err != nil {
		logger.Println("[E] Create connection object to local database failed!", err)
		timestamp = time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + create_database_object_failed + "{{" + fmt.Sprintf("%s", err)
		logger.Println("[I] Give up leader election")
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + give_election
		UploadLog(logkey, logvalue)
		os.Exit(1)
	}
	logger.Println("[I] Create connection object to local database successfully!")
	timestamp = time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + create_database_object_success
	defer db.Close()
	err = db.Ping()
	if err != nil {
		logger.Println("[E] Connected to local database failed!", err)
		timestamp = time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + connected_database_failed + "{{" + fmt.Sprintf("%s", err)
		logger.Println("[I] Give up leader election")
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + give_election
		UploadLog(logkey, logvalue)
		os.Exit(1)
	}
	logger.Println("[I] Connected to local database successfully!")
	timestamp = time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + connected_database_success
	read_only := "set global read_only=" + strconv.Itoa(value)
	_, err = db.Query(read_only)
	if err != nil {
		logger.Println("[E] Set local database Read_only mode failed!", err)
		timestamp = time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + set_read_only_failed + "{{" + fmt.Sprintf("%s", err)
		logger.Println("[I] Local database downgrade failed!")
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + database_downgrade_failed
		UploadLog(logkey, logvalue)
		os.Exit(1)
	}
	logger.Println("[I] Set local database Read_only mode successfully!")
	timestamp = time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + set_read_only_success
	logger.Println("[I] Local database downgrade successfully!")
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + database_downgrade_success
}
