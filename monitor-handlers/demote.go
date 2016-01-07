package main

import (
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

func CheckService() {
	beego.Info("Monitor Handler Triggered")
	time.Sleep(10000 * time.Millisecond)
	servicename = beego.AppConfig.String("servicename")
	service_ip = beego.AppConfig.Strings("service_ip")
	other_hostname = beego.AppConfig.String("otherhostname")
	tag = beego.AppConfig.String("tag")
	config := &consulapi.Config{
		Datacenter: beego.AppConfig.String("datacenter"),
		Token:      beego.AppConfig.String("token"),
	}
	var healthpair []*consulapi.ServiceEntry
	var kv *consulapi.KV
//	var kvPair *consulapi.KVPair
	var client *consulapi.Client
	var health *consulapi.Health
	var err error
	for i, _ := range service_ip {
		config.Address = service_ip[i] + ":" + beego.AppConfig.String("service_port")
		client, err = consulapi.NewClient(config)
		if err != nil {
			beego.Error("Create consul-api client failed!", err)
			beego.Info("Give up switching to async replication")
			beego.Info("Monitor Handler Completed")
			return
		}
		beego.Info("Create consul-api client successfully!")
		health = client.Health()
		healthpair, _, err = health.Service(servicename, tag, false, nil)
		if err != nil {
			beego.Error("Get peer database health status from CS failed!", err)
			beego.Info("Give up switching to async replication")
                        beego.Info("Monitor Handler Completed")
			continue
		}
		break
	}
	beego.Info("Get peer database health status from CS successfully!")
	kv = client.KV()
	if len(healthpair) <= 0 {
		beego.Error("Peer database service not exist in CS!")
		beego.Info("Give up switching to async replication")
                beego.Info("Monitor Handler Completed")
		return
	}
	beego.Info("Peer database service exist in CS!")
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
		beego.Info("Service health status is passing in CS!")
		beego.Info("Give up switching to async replication")
                beego.Info("Monitor Handler Completed")
		return
	} else if status == "warning" {
		beego.Info("Warning! Peer database " + other_hostname + " replicaton error. Service health status is warning in CS ")
		Switch := beego.AppConfig.String("switch_async")
		if strings.EqualFold(Switch, "off") {
			beego.Info("Current switch_async value is "+ Switch) 
			beego.Info("Give up switching to async replication")
                        beego.Info("Monitor Handler Completed")
			return
		} else if strings.EqualFold(Switch, "on") {
			beego.Info("Current switch_async value is "+ Switch)
        	      	checkio_thread(ip, port, username, password, addr)
		} else {
			beego.Info("Config file switch_async format error,off or on!")
			beego.Info("Give up switching to async replication")
                        beego.Info("Monitor Handler Completed")
			return
		}
	} else if status == "critical" {
			beego.Info("Service health status is critical in CS!")
			var put string
			put = "1"
			kvvalue := []byte(put)
			kvotherhostname := consulapi.KVPair{
				Key:   "monitor/" + other_hostname,
				Value: kvvalue,
			}
			_, err = kv.Put(&kvotherhostname, nil)
			if err != nil {
				beego.Error("Set peer database repl_err_counter to 1 in CS failed!", err)
				beego.Info("Give up switching to async replication")
	                        beego.Info("Monitor Handler Completed")
				return
			}
			beego.Info("Set peer database repl_err_counter to 1 in CS successfully!")
        	      	checkio_thread(ip, port, username, password, addr)
	} else {
		beego.Info("Not passing,not waring,not critical ,is invalid state!")
		beego.Info("Give up switching to async replication")
                beego.Info("Monitor Handler Completed")
		return
	}
}
