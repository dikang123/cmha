package main

import (
//	"fmt"
	"strings"
//	"bytes"
//	"encoding/binary"
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
	var kvPair *consulapi.KVPair
	var client *consulapi.Client
	var health *consulapi.Health
	var err error
	for i, _ := range service_ip {
		config.Address = service_ip[i] + ":" + beego.AppConfig.String("service_port")
		client, err = consulapi.NewClient(config)
		if err != nil {
			beego.Error("Create consul-api client failure!", err)
			return
		}
		beego.Info(" Create consul-api client success!")
		health = client.Health()
		healthpair, _, err = health.Service(servicename, tag, false, nil)
		if err != nil {
			beego.Error("Health check execution nodes and services /v1/health/service/"+servicename+"?tag="+tag+" failure!", err)
			continue
		}
		break
	}
	kv = client.KV()
	beego.Info("Health check execution nodes and services /v1/health/service/" + servicename + "?tag=" + tag + " success!")
	if len(healthpair) <= 0 {
		beego.Error("tag=" + tag + " of " + servicename + " service does not exist!")
		return
	}
	beego.Info("tag=" + tag + " of " + servicename + " service exist!")
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
		beego.Info(tag + " status is passing!")
		return
	} else if status == "warning" {
		beego.Info(tag + " status is warning!")
		Switch := beego.AppConfig.String("switch_async")
		if strings.EqualFold(Switch, "off") {
			beego.Info("Not set asynchronous!")
			return
		} else if strings.EqualFold(Switch, "on") {
			var put string
			put = "1"
			kvvalue := []byte(put)
			kvotherhostname := consulapi.KVPair{
				Key:   "monitor/" + other_hostname,
				Value: kvvalue,
			}
			_, err = kv.Put(&kvotherhostname, nil)
			if err != nil {
				beego.Error("monitor/"+other_hostname+" put failure!", err)
				return
			}
			beego.Info("monitor/" + other_hostname + " put success!")
			kvPair, _, err = kv.Get("service/"+servicename+"/leader", nil)
        	//      fmt.Println(len(kvPair))
       //         	beego.Info(kvPair.Value)
                	if err != nil {
                        	beego.Error("Get a service/"+servicename+"/leader key failure!", err)
                        	return
               	 	}
			content :=string(kvPair.Value)
        		istrue := strings.Contains(content, ip)
  //      		fmt.Println(b)
       			if  istrue {
        	      		checkio_thread(ip, port, username, password, addr)
      	                	return
        		}else{
              			beego.Info(ip +" not leader")
                		return
        		}
		} else {
			beego.Info("Config file switch format error,off or on!")
			return
		}
	} else if status == "critical" {
			beego.Info(tag +" status is critical!")
			var put string
			put = "1"
			kvvalue := []byte(put)
			kvotherhostname := consulapi.KVPair{
				Key:   "monitor/" + other_hostname,
				Value: kvvalue,
			}
			_, err = kv.Put(&kvotherhostname, nil)
			if err != nil {
				beego.Error("monitor/"+other_hostname+" put failure!", err)
				return
			}
			beego.Info("monitor/" + other_hostname + " put success!")
			kvPair, _, err := kv.Get("service/"+servicename+"/leader",nil)
                //	fmt.Println(kvPa.Value)
                	if err != nil {
                        	beego.Error("Get a service/"+servicename+"/leader key failure!", err)
                        	return
               	 	}
			content :=string(kvPair.Value)
        		istrue := strings.Contains(content, ip)
  //      		fmt.Println(b)
//      		fmt.Println(content)
       			if  istrue {
        	      		checkio_thread(ip, port, username, password, addr)
      	                	return
        		}else{
              			beego.Info(ip +" not leader")
                		return
        		}
	//	checkio_thread(ip, port, username, password, addr)
	} else {
		beego.Info("Not passing,not waring,not critical ,is invalid state!")
		return
	}
}
