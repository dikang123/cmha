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

func SessionAndChecks() {
	beego.Info("MHA Handler Triggered")
	ip := beego.AppConfig.String("ip")
	Switch := beego.AppConfig.String("switch")
	if strings.EqualFold(Switch, "off") {
		beego.Info(ip +" switch="+Switch+",give up leader election!")
		beego.Info("MHA Handler Completed")
		return
	} else if strings.EqualFold(Switch, "on") {
		beego.Info(ip + " switch=" +Switch)
		beego.Info("Begin leader election!")
	} else {
		beego.Info("Config file switch format error,switch="+Switch+",Should off or on!")
		beego.Info("Give up leader election")
		beego.Info("MHA Handler Completed")
		return
	}
	service_ip = beego.AppConfig.Strings("service_ip")
	servicename = beego.AppConfig.String("servicename")
	hostname = beego.AppConfig.String("hostname")
	config := &consulapi.Config{
		Datacenter: beego.AppConfig.String("datacenter"),
		Token:      beego.AppConfig.String("token"),
	}
	var kvPair *consulapi.KVPair
	var kvMonitor *consulapi.KVPair
	var client *consulapi.Client
	var kv *consulapi.KV
	var err error
	for i, _ := range service_ip {
		config.Address = service_ip[i] + ":" + beego.AppConfig.String("service_port")
		client, err = consulapi.NewClient(config)
		if err != nil {
			beego.Error("Create consul-api client failed!", err)
			beego.Info("Give up leader election")
	                beego.Info("MHA Handler Completed")
			return
		}
		beego.Info("Create consul-api client successfully!")
		//KV is used to return a handle to the K/V apis
		kv = client.KV()
		//Get is used to lookup a single key
		kvPair, _, err = kv.Get("service/"+servicename+"/leader", nil)
		if err != nil {
			beego.Error("Get and check current service leader from CS failed!", err)
			continue
		}
		break
	}
	beego.Info("Get and check current service leader from CS successfully!")
	kvMonitor, _, err = kv.Get("monitor/"+hostname, nil)
	if err != nil {
		beego.Error("Get " + ip + "repl_err_counter=0/1 failed", err)
		beego.Info("Give up leader election")
                beego.Info("MHA Handler Completed")
		return
	}
	kvValue :=string(kvMonitor.Value)
	if kvValue != "0" {
		beego.Error(ip + " give up leader election")
                beego.Info("MHA Handler Completed")
		return
	}
	//NewClient returns a new client
	if kvPair == nil {
		beego.Error("Not replication counter,Please create!")
		beego.Info("Give up leader election")
                beego.Info("MHA Handler Completed")
		return
	}
	//Are there external connection string provided
	if kvPair.Session != "" {
		beego.Info("Leader exist!")
		time.Sleep(1 * time.Second)
		beego.Info("Give up leader election")
                beego.Info("MHA Handler Completed")
		return
	}
	beego.Info("Leader does not exist!")
	//Health returns a handle to the health endpoints
	health := client.Health()
	//Checks is used to return the checks associated with a service
	healthvalue, _, err := health.Checks(servicename, nil)
	if err != nil {
		beego.Error("Get and check "+ ip + " service health status failed!", err)
		beego.Info("Give up leader election")
                beego.Info("MHA Handler Completed")
		return
	}
	beego.Info("Get and check "+ ip + " service health status successfully!")
	if len(healthvalue) <= 0 {
		beego.Info(servicename + " service does not exist!")
		beego.Info("Give up leader election")
                beego.Info("MHA Handler Completed")
		return
	}
	var islocal bool
	for index := range healthvalue {
		if healthvalue[index].Node == hostname {
			islocal = true
			beego.Info(servicename + " service exist!")
			break
		}

	}
	if !islocal {
		beego.Info(ip + " not is " +servicename + "!")
		beego.Info("Give up leader election")
                beego.Info("MHA Handler Completed")
		return
	} else {
		updatevalue := consulapi.KVPair{
			Key:   "service/" + servicename + "/leader",
			Value: []byte(""),
		}
		_, err = kv.Put(&updatevalue, nil)
		if err != nil {
			beego.Error("Clean service leader value in CS failed!", err)
			beego.Info("Give up leader election")
	                beego.Info("MHA Handler Completed")
			return
		}
		beego.Info("Clean service leader value in CS successfully!")
		healthpair, _, err := health.Service(servicename, "", false, nil)
		if err != nil {
			beego.Error("Get and check " + ip + " service health status failed!", err)
			beego.Info("Give up leader election")
	                beego.Info("MHA Handler Completed")
			return
		}
		beego.Info("Get and check " + ip + " service health status successfully!")
		var ishealthy = true
		hostname := beego.AppConfig.String("hostname")
		for index := range healthpair {
			for checkindex := range healthpair[index].Checks {
				if healthpair[index].Checks[checkindex].Node == hostname {
					if healthpair[index].Checks[checkindex].Status == "critical" {
						ishealthy = false
						break
					}
				}
			}
		}
		if !ishealthy {
			beego.Info("Status is critical!")
			beego.Info("Give up leader election")
	                beego.Info("MHA Handler Completed")
			return
		} else {
			beego.Info("Status is not critical")
			port := beego.AppConfig.String("port")
			username := beego.AppConfig.String("username")
			password := beego.AppConfig.String("password")
			slave(ip, port, username, password)
		}
	}
}
func SetConn(ip, port, username, password string) {
	config := &consulapi.Config{
		Datacenter: beego.AppConfig.String("datacenter"),
		Token:      beego.AppConfig.String("token"),
	}
	var client *consulapi.Client
	var kv *consulapi.KV
	var err error
	var sessionvalue string
	for i, _ := range service_ip {
		config.Address = service_ip[i] + ":" + beego.AppConfig.String("service_port")
		client, err = consulapi.NewClient(config)
		if err != nil {
			beego.Error("Create  consul-api client failed! ", err)
			beego.Info("Give up leader election")
                	beego.Info("MHA Handler Completed")
			return
		}
		beego.Info("Create  consul-api client successfully!")
		session := client.Session()
		sessionEntry := consulapi.SessionEntry{
			LockDelay: 10 * time.Second,
			Name:      servicename,
			Node:      hostname,
			Checks:    []string{"serfHealth", "service:" + servicename},
		}
		//Create makes a new session. Providing a session entry can customize the session. It can also be nil to use defaults.
		sessionvalue, _, err = session.Create(&sessionEntry, nil)
		if err != nil {
			beego.Error("Session create failed!", err)
			continue
		}
		break
	}
	//NewClient returns a new client
	beego.Info("Session create successfully!")
	format := beego.AppConfig.String("format")
	var acquirejson string
	if format == "json" {
		acquirejson = `{"Node":"` + hostname + `","Ip":"` + ip + `","Port":` + port + `,"Username":"` + username + `","Password":"` + password + `"}`
	} else if format == "hap" {
		acquirejson = "server" + " " + hostname + " " + ip + ":" + port
	} else {
		beego.Error("format error,json or hap!")
		beego.Info("Give up leader election")
                beego.Info("MHA Handler Completed")
		return
	}
	value := []byte(acquirejson)
	kv = client.KV()
	kvpair := consulapi.KVPair{
		Key:     "service/" + servicename + "/leader",
		Value:   value,
		Session: sessionvalue,
	}
	//Acquire is used for a lock acquisiiton operation. The Key, Flags, Value and Session are respected. Returns true on success or false on failures.
	time.Sleep(15 * time.Second)
	var ok bool
	ok, _, err = kv.Acquire(&kvpair, nil)
	if err != nil {
		beego.Error("Send service leader request to CS failed! ", err)
                beego.Info("MHA Handler Completed")
		return
	}
	beego.Info("Send service leader request to CS successfully!")
	if !ok {
		time.Sleep(5 * time.Second)
		beego.Warn("Becoming service leader failed! Connection string is " + ip + " " + port)
                beego.Info("MHA Handler Completed")
		return
	} else {
		beego.Info("Becoming service leader successfully! Connection string is " + ip + " " + port)
                beego.Info("MHA Handler Completed")
	}
}
