package main

import (
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
var address string

func SetConn(client *consulapi.Client) {
	hostname,other_hostname = ReadHostConf()
	ip,port,username,password = ReadDatabaseConf()
	servicename = ReadServiceConf()
	address =ReadCaConf()
	leader,last_leader := ReturnLeaderAndLastLeader()
	var kvPair *consulapi.KVPair
	var kv *consulapi.KV
	var err error
	for i := 0; i < 3; i++ {
		//KV is used to return a handle to the K/V apis
		kv = client.KV()
		//Get is used to lookup a single key
		kvPair, _, err = kv.Get(leader, nil)
		if err != nil {
			beego.Error("Get a "+leader+" key failure!", err)
			continue
		}
		break
	}
	kvvalue := PutValue()
	kvhostname,kvhostnamekey :=GetRepl(servicename,hostname,kvvalue)
	kvotherhostname,kvotherhostnamekey := GetRepl(servicename,other_hostname,kvvalue)
	err =Put(&kvhostname,kvhostnamekey,kv)
	if err != nil {
		beego.Error("Put failed:",kvhostnamekey,err)
		return
	}
	err =Put(&kvotherhostname,kvotherhostnamekey,kv)
        if err != nil {
                beego.Error("Put failed:",kvotherhostnamekey,err)
                return
        }
	//NewClient returns a new client
	if kvPair == nil {
		beego.Error(leader + "not found, please create the key!")
		return
	}
	//Are there external connection string provided
	if kvPair.Session != "" {
		beego.Info("There are external connection string provided!")
		time.Sleep(1000)
		return
	}
	islocal,err :=HealthCheck(servicename,hostname,client)
	if err != nil {
		beego.Error("HealthCheck failed:",err)
		return
	}
	if !islocal {
		beego.Info("Native " + servicename + " service unhealthy or the service does not exist!")
		return
	} else {
		sessionvalue,err :=CreateSession(servicename,hostname,client)
		if err != nil {
			beego.Error("CreateSessiob failed:",err)
			return
		}
		format := ReadFormat()
		ok,err :=SessionAcquireLeader(format,hostname,ip,port,username,password,leader,sessionvalue,kv)
		if err != nil {
			beego.Error("SessionAcquireLeader failed:",err)
			return
		}
		if !ok {
			beego.Error("kv acquire failure!")
			return
		} else {
			leadervalue := "server" + " " + hostname + " " + ip + ":" + port
			err = PutLastLeader(leadervalue,last_leader,kv)		
			if err != nil {
				beego.Error("PutLastLeader failed:",last_leader,err)
				return
			}
		}
	}
}
