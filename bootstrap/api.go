package main

import (
	"fmt"
	//	"reflect"
	"time"

	"github.com/astaxie/beego"
	consulapi "github.com/hashicorp/consul/api"
)

var hostname string
var other_hostname string
var ip string
var port string
var username string
var password string
var servicename string
var service_ip []string

func SetConn() {
	hostname = beego.AppConfig.String("hostname")
	other_hostname = beego.AppConfig.String("otherhostname")
	ip = beego.AppConfig.String("ip")
	port = beego.AppConfig.String("port")
	username = beego.AppConfig.String("username")
	password = beego.AppConfig.String("password")
	servicename = beego.AppConfig.String("servicename")
	service_ip = beego.AppConfig.Strings("service_ip")
	//Config is used to configure the creation of a client
	config := &consulapi.Config{
		Datacenter: beego.AppConfig.String("datacenter"),
		Token:      beego.AppConfig.String("token"),
	}
	var kvPair *consulapi.KVPair
	var client *consulapi.Client
	var kv *consulapi.KV
	var err error
	for i, _ := range service_ip {
		config.Address = service_ip[i] + ":" + beego.AppConfig.String("service_port")
		client, err = consulapi.NewClient(config)
		if err != nil {
			beego.Error("Create a consul-api client failure!", err)
			return
		}
		beego.Info("Create a consul-api client success!")
		//KV is used to return a handle to the K/V apis
		kv = client.KV()
		//Get is used to lookup a single key
		kvPair, _, err = kv.Get("service/"+servicename+"/leader", nil)
		if err != nil {
			beego.Error("Get a service/"+servicename+"/leader key failure!", err)
			continue
		}
		break
	}
	var put string
	put = "0"
	kvvalue := []byte(put)
	kvhostname := consulapi.KVPair{
		Key:   "monitor/" + hostname,
		Value: kvvalue,
	}
	kvotherhostname := consulapi.KVPair{
		Key:   "monitor/" + other_hostname,
		Value: kvvalue,
	}
	_, err = kv.Put(&kvhostname, nil)
	if err != nil {
		beego.Error("monitor/"+hostname+" put failure!", err)
		return
	}
	beego.Info("monitor/" + hostname + " put success!")
	_, err = kv.Put(&kvotherhostname, nil)
	if err != nil {
		beego.Error("monitor/"+other_hostname+" put failure!", err)
		return
	}
	beego.Info("monitor/" + other_hostname + " put success!")
	//NewClient returns a new client
	beego.Info("Get a service/" + servicename + "/leader key success!")
	if kvPair == nil {
		beego.Error("service/" + servicename + "/leader not found, please create the key!")
		return
	}
	//Are there external connection string provided
	if kvPair.Session != "" {
		beego.Info("There are external connection string provided!")
		time.Sleep(1000)
		return
	}
	//Health returns a handle to the health endpoints
	health := client.Health()
	//Checks is used to return the checks associated with a service
	healthvalue, _, err := health.Checks(servicename, nil)
	if err != nil {
		beego.Error("Return to service-related checks failure!", err)
		return
	}
	if len(healthvalue) <= 0 {
		beego.Info("Without this service, or service is not a healthy state!")
		return
	}
	var islocal bool
	for index := range healthvalue {
		fmt.Println(healthvalue[index].Node, hostname)
		if healthvalue[index].Node == hostname {
			islocal = true
			beego.Info("Native " + servicename + " service is healthy!")
			break
		}

	}
	if !islocal {
		beego.Info("Native " + servicename + " service unhealthy or the service does not exist!")
		return
	} else {
		//Session returns a handle to the session endpoints
		session := client.Session()
		sessionEntry := consulapi.SessionEntry{
			LockDelay: 10 * time.Second,
			Name:      servicename,
			Node:      hostname,
			Checks:    []string{"serfHealth", "service:" + servicename},
		}
		//Create makes a new session. Providing a session entry can customize the session. It can also be nil to use defaults.
		sessionvalue, _, err := session.Create(&sessionEntry, nil)
		if err != nil {
			beego.Error("Session creation failure!", err)
			return
		}
		format := beego.AppConfig.String("format")
		var acquirejson string
		if format == "json" {
			acquirejson = `{"Node":"` + hostname + `","Ip":"` + ip + `","Port":` + port + `,"Username":"` + username + `","Password":"` + password + `"}`
		} else if format == "hap" {
			acquirejson = "server" + " " + hostname + " " + ip + ":" + port
		} else {
			beego.Error("format error!")
			return
		}
		value := []byte(acquirejson)
		kvpair := consulapi.KVPair{
			Key:     "service/" + servicename + "/leader",
			Value:   value,
			Session: sessionvalue,
		}
		//Acquire is used for a lock acquisiiton operation. The Key, Flags, Value and Session are respected. Returns true on success or false on failures.
		ok, _, err := kv.Acquire(&kvpair, nil)
		if err != nil {
			beego.Error("Set the connection string master failure!", err)
			return
		}
		if !ok {
			beego.Error("kv acquire failure!")
			return
		} else {
			beego.Info("kv acquire success!")
		}
	}
}
