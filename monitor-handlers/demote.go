package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	consulapi "github.com/hashicorp/consul/api"
)

var address string
var consul_agent_ip string
var servicename string
var hostname string
var other_hostname string
var tag string
var client *consulapi.Client
var kv *consulapi.KV
var logvalue string
var logkey string

func CheckService() {
	var err error
	logger.Println("[I] Monitor Handler Triggered")
	timestamp := time.Now().Unix()
	logvalue = strconv.FormatInt(timestamp, 10) + triggered
	time.Sleep(10000 * time.Millisecond)
	servicename = ReadServiceName()
	address := ReadCaAddress()
	hostname, other_hostname = ReadHostName()
	tag = ReadTag()
	logkey = GetLogKey(servicename, hostname, timestamp)
	config := GetConfig(address)
	client, logvalue, err = GetClient(config, logvalue, consulapi_failed)
	if err != nil {
		UploadLog(logkey, logvalue)
	}
	logger.Println("[I] Create consul-api client successfully!")
	logvalue, timestamp = AddLogValue(consulapi_success, logvalue)
	var healthpair []*consulapi.ServiceEntry
	var health *consulapi.Health

	c := 0
	for i := 0; i < 3; i++ {
		health = client.Health()
		healthpair, _, err = health.Service(servicename, tag, false, nil)
		if err != nil {
			c += 1
			logger.Println("[E] Get peer service "+servicename+" health status from CS failed!", err)
			timestamp := time.Now().Unix()
			logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + peer_service_health_failed + "{{" + servicename + "{{" + fmt.Sprintf("%s", err)
			continue
		}
		break
	}
	if c >= 3 {
		logger.Println("[I] Give up switching to async replication!")
		logvalue, timestamp = AddLogValue(give_async_repl, logvalue)
		UploadLog(logkey, logvalue)
		return
	}
	logger.Println("[I] Get peer service " + servicename + " health status from CS successfully!")
	timestamp = time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + peer_service_health_success + "{{" + servicename
	kv = client.KV()
	if len(healthpair) <= 0 {
		logger.Println("[E] " + servicename + " peer service not exist in CS!")
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + service_noexist + "{{" + servicename
		logger.Println("[I] Give up switching to async replication!")
		logvalue, timestamp = AddLogValue(give_async_repl, logvalue)
		UploadLog(logkey, logvalue)
		return
	}
	logger.Println("[I] " + servicename + " peer service exist in CS!")
	timestamp = time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + service_exist + "{{" + servicename
	addr, status := GetAddrAndStatus(healthpair)
	ip, port, username, password := ReadDatabaseCfg()
	if status == "passing" {
		logger.Println("[I] Service health status is passing in CS!")
		logvalue, timestamp = AddLogValue(service_passing, logvalue)
		logger.Println("[I] Give up switching to async replication!")
		logvalue, timestamp = AddLogValue(give_async_repl, logvalue)
		UploadLog(logkey, logvalue)
		return
	} else if status == "warning" {
		logger.Println("[W] Warning! Peer database " + other_hostname + " replicaton error. Service health status is warning in CS!")
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + repl_err_service_warning + "{{" + other_hostname
		Switch := beego.AppConfig.String("switch_async")
		if strings.EqualFold(Switch, "off") {
			logger.Println("[I] Current switch_async value is " + Switch)
			timestamp := time.Now().Unix()
			logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + switch_async + "{{" + Switch
			logger.Println("[I] Give up switching to async replication!")
			logvalue, timestamp = AddLogValue(give_async_repl, logvalue)
			UploadLog(logkey, logvalue)
			return
		} else if strings.EqualFold(Switch, "on") {
			logger.Println("[I] Current switch_async value is " + Switch)
			timestamp := time.Now().Unix()
			logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + switch_async + "{{" + Switch
			logvalue = checkio_thread(ip, port, username, password, addr,logvalue)
			UploadLog(logkey, logvalue)
			return
		} else {
			logger.Println("[I] Config file switch_async format error,off or on!")
			logvalue, timestamp = AddLogValue(switch_async_format_err, logvalue)
			logger.Println("[I] Give up switching to async replication!")
			logvalue, timestamp = AddLogValue(give_async_repl, logvalue)
			UploadLog(logkey, logvalue)
			return
		}
	} else if status == "critical" {
		logger.Println("[I] Service health status is critical in CS!")
		logvalue, timestamp = AddLogValue(service_critical, logvalue)
		cnt := 0
		kvotherhostname := PutRepl(servicename, other_hostname)
	try:
		_, err := kv.Put(&kvotherhostname, nil)
		if err != nil {
			logger.Println("[E] Set peer database repl_err_counter to 1 in CS failed!", err)
			timestamp := time.Now().Unix()
			logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + set_counter_failed + "{{" + fmt.Sprintf("%s", err)
			if cnt == 2 {
				count := 0
				for i := 0; i < 3; i++ {
					ct := 0
				tr:
					_, err = kv.Put(&kvotherhostname, nil)
					if err != nil {
						count += 1
						logger.Println("[E] Set peer database repl_err_counter to 1 in CS failed!", err)
						timestamp := time.Now().Unix()
						logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + set_counter_failed + "{{" + fmt.Sprintf("%s", err)
						if ct == 2 {
							continue
						} else {
							ct++
							goto tr
						}
					}
					break
				}
				if count >= 3 {
					UploadLog(logkey, logvalue)
					os.Exit(1)
				}
			} else {
				cnt++
				goto try
			}
		}
		logger.Println("[I] Set peer database repl_err_counter to 1 in CS successfully!")
		logvalue, timestamp = AddLogValue(set_counter_success, logvalue)
		logvalue = checkio_thread(ip, port, username, password, addr,logvalue)
		UploadLog(logkey, logvalue)
		return
	} else {
		logger.Println("[E] Not passing,not waring,not critical ,is invalid state!")
		logvalue, timestamp = AddLogValue(service_invalid, logvalue)
		logger.Println("[I] Give up switching to async replication!")
		logvalue, timestamp = AddLogValue(give_async_repl, logvalue)
		UploadLog(logkey, logvalue)
		return
	}
}
