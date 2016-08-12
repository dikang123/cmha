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

var service_ip []string
var servicename string
var hostname string
var other_hostname string
var tag string
var client *consulapi.Client
var kv *consulapi.KV
var logvalue string
var logkey string

const (
	triggered        = "001"
	consulapi_failed = "002"
	give_async_repl  = "003"
	//	completed = "004"
	consulapi_success                   = "005"
	peer_service_health_failed          = "006"
	peer_service_health_success         = "007"
	service_noexist                     = "008"
	service_exist                       = "009"
	service_passing                     = "010"
	repl_err_service_warning            = "011"
	switch_async                        = "012"
	switch_async_format_err             = "013"
	service_critical                    = "014"
	set_counter_failed                  = "015"
	set_counter_success                 = "016"
	service_invalid                     = "017"
	get_leader_failed                   = "018"
	get_leader_success                  = "019"
	no_leader                           = "020"
	is_leader                           = "021"
	get_service_status_failed           = "022"
	get_service_status_success          = "023"
	service_status                      = "024"
	create_database_object_failed       = "025"
	create_database_object_success      = "026"
	connected_database_failed           = "027"
	connected_database_success          = "028"
	set_keepsyncrepl_failed             = "029"
	set_keepsyncrepl_success            = "030"
	set_trysyncrepl_failed              = "031"
	set_trysyncrepl_success             = "032"
	switch_local_async_repl             = "033"
	handlers_sleep                      = "034"
	connecting_database                 = "035"
	create_peer_database_object_failed  = "036"
	create_peer_database_object_success = "037"
	connected_peer_database_failed      = "038"
	connected_peer_database_success     = "039"
	check_perr_dabases_io_failed        = "040"
	check_perr_dabases_io_success       = "041"
	resolve_slave_status_failed         = "042"
	io_status                           = "043"
)

func CheckService() {
	logger.Println("[I] Monitor Handler Triggered")
	timestamp := time.Now().Unix()
	logvalue = strconv.FormatInt(timestamp, 10) + triggered
	time.Sleep(10000 * time.Millisecond)
	servicename = beego.AppConfig.String("servicename")
	service_ip = beego.AppConfig.Strings("service_ip")
	hostname = beego.AppConfig.String("hostname")
	other_hostname = beego.AppConfig.String("otherhostname")
	tag = beego.AppConfig.String("tag")
	logkey = "cmha/service/" + servicename + "/log/" + hostname + "/monitor-handlers/" + strconv.FormatInt(timestamp, 10)
	config := &consulapi.Config{
		Datacenter: beego.AppConfig.String("datacenter"),
		Token:      beego.AppConfig.String("token"),
	}
	var healthpair []*consulapi.ServiceEntry
	var health *consulapi.Health
	var err error
	c := len(service_ip)
	for i, _ := range service_ip {
		config.Address = service_ip[i] + ":" + beego.AppConfig.String("service_port")
		client, err = consulapi.NewClient(config)
		if err != nil {
			c--
			logger.Println("[E] Create consul-api client failed!", err)
			timestamp := time.Now().Unix()
			logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + consulapi_failed + "{{" + fmt.Sprintf("%s", err)
			continue
		}
		logger.Println("[I] Create consul-api client successfully!")
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + consulapi_success
		health = client.Health()
		healthpair, _, err = health.Service(servicename, tag, false, nil)
		if err != nil {
			c--
			logger.Println("[E] Get peer service "+servicename+" health status from CS failed!", err)
			timestamp := time.Now().Unix()
			logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + peer_service_health_failed + "{{" + servicename + "{{" + fmt.Sprintf("%s", err)
			continue
		}
		break
	}
	if c == 0 {
		logger.Println("[I] Give up switching to async replication!")
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + give_async_repl
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
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + give_async_repl
		UploadLog(logkey, logvalue)
		return
	}
	logger.Println("[I] " + servicename + " peer service exist in CS!")
	timestamp = time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + service_exist + "{{" + servicename
	var addr string
	var status string
	for index := range healthpair {
		if healthpair[index].Node.Address != "" {
			addr = healthpair[index].Node.Address
		}
		for checkindex := range healthpair[index].Checks {
			if healthpair[index].Checks[checkindex].Status == "passing" {
				status = healthpair[index].Checks[checkindex].Status
			} else if healthpair[index].Checks[checkindex].Status == "warning" {

				status = healthpair[index].Checks[checkindex].Status
				break
			} else if healthpair[index].Checks[checkindex].Status == "critical" {
				status = healthpair[index].Checks[checkindex].Status
				break
			} else {
				status = "invalid"
				break
			}
		}
	}
	ip := beego.AppConfig.String("ip")
	port := beego.AppConfig.String("port")
	username := beego.AppConfig.String("username")
	password := beego.AppConfig.String("password")
	if status == "passing" {
		logger.Println("[I] Service health status is passing in CS!")
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + service_passing
		logger.Println("[I] Give up switching to async replication!")
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + give_async_repl
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
			logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + give_async_repl
			UploadLog(logkey, logvalue)
			return
		} else if strings.EqualFold(Switch, "on") {
			logger.Println("[I] Current switch_async value is " + Switch)
			timestamp := time.Now().Unix()
			logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + switch_async + "{{" + Switch
			checkio_thread(ip, port, username, password, addr)
			UploadLog(logkey, logvalue)
		} else {
			logger.Println("[I] Config file switch_async format error,off or on!")
			timestamp := time.Now().Unix()
			logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + switch_async_format_err
			logger.Println("[I] Give up switching to async replication!")
			logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + give_async_repl
			UploadLog(logkey, logvalue)
			return
		}
	} else if status == "critical" {
		logger.Println("[I] Service health status is critical in CS!")
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + service_critical
		var put string
		cnt := 0
		put = "1"
		kvvalue := []byte(put)
		kvotherhostname := consulapi.KVPair{
			Key:   "cmha/service/" + servicename + "/db/" + other_hostname + "/repl_err_counter",
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
						logger.Println("[E] Create consul-api client failed!", err)
						timestamp := time.Now().Unix()
						logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + consulapi_failed + "{{" + fmt.Sprintf("%s", err)
						continue
					}
					logger.Println("[I] Create consul-api client successfully!")
					timestamp := time.Now().Unix()
					logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + consulapi_success
					ct := 0
				tr:
					_, err = kv.Put(&kvotherhostname, nil)
					if err != nil {
						count--
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
		timestamp = time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + set_counter_success
		checkio_thread(ip, port, username, password, addr)
		UploadLog(logkey, logvalue)
	} else {
		logger.Println("[E] Not passing,not waring,not critical ,is invalid state!")
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + service_invalid
		logger.Println("[I] Give up switching to async replication!")
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + give_async_repl
		UploadLog(logkey, logvalue)
		return
	}
}

func UploadLog(logkey, logvalue string) {
	config := &consulapi.Config{
		Datacenter: beego.AppConfig.String("datacenter"),
		Token:      beego.AppConfig.String("token"),
	}
	kvhostname := consulapi.KVPair{
		Key:   logkey,
		Value: []byte(logvalue),
	}
	_, err := kv.Put(&kvhostname, nil)
	if err != nil {
		logger.Println("[E] Upload log to CS failed!", err)
		count := len(service_ip)
		for i, _ := range service_ip {
			config.Address = service_ip[i] + ":" + beego.AppConfig.String("service_port")
			client, err = consulapi.NewClient(config)
			if err != nil {
				count--
				logger.Println("[E] Create consul-api client failed!", err)
				timestamp := time.Now().Unix()
				logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + consulapi_failed + "{{" + fmt.Sprintf("%s", err)
				continue
			}
			logger.Println("[I] Create consul-api client successfully!")
			timestamp := time.Now().Unix()
			logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + consulapi_success
			_, err = kv.Put(&kvhostname, nil)
			if err != nil {
				count--
				logger.Println("[E] Upload log to CS failed!", err)
				continue
			}
			break
		}
		if count == 0 {
			logger.Println("[I] Monitor Handler Completed")
			return
		}
	}
	logger.Println("[I] Upload log to CS successfully!")
	logger.Println("[I] Monitor Handler Completed")
}
