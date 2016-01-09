package main

import (
	"database/sql"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	consulapi "github.com/hashicorp/consul/api"
	"strconv"
	"time"
	"strings"
)

var kv *consulapi.KV

func Set(ip, port, username, password string, issynchronous int) {
	beego.Info("Set peer database repl_err_counter to 1 in CS")
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
	for i, _ := range service_ip {
                config.Address = service_ip[i] + ":" + beego.AppConfig.String("service_port")
                client, err = consulapi.NewClient(config)
                if err != nil {
                        beego.Error("Create consul-api client failed!", err)
                        continue
                }
                beego.Info("Create consul-api client successfully!")
                break
        }
        kv = client.KV()
	health = client.Health()
/*	var put string
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
        beego.Info("Set peer database repl_err_counter to 1 in CS successfully!")*/
	SetRepl_err_counter(other_hostname)
	kvPair, _, err = kv.Get("service/"+servicename+"/leader", nil)
        if err != nil {
		beego.Error("Get and check current service leader from CS failed!", err)
		beego.Info("Give up switching to async replication")
                beego.Info("Monitor Handler Completed")
		return
	}
	beego.Info("Get and check current service leader from CS successfully!")
        content :=string(kvPair.Value)
        istrue := strings.Contains(content, ip)
        if  !istrue {
		beego.Info(ip +" is not service leader!")
		beego.Info("Give up switching to async replication!")
                beego.Info("Monitor Handler Completed")
            	return
        }
	beego.Info(ip +" is service leader!")
	healthpair, _, err = health.Service(servicename, "", false, nil)
	if err != nil {
		beego.Error("Get peer database health status from CS failed!", err)
		beego.Info("Give up switching to async replication")
                beego.Info("Monitor Handler Completed")
		return
	}
	beego.Info("Get peer database health status from CS successfully!")
	var ishealthy = true
	var status string
	for index := range healthpair {
		if healthpair[index].Node.Address == ip {
			for checkindex := range healthpair[index].Checks {
					if healthpair[index].Checks[checkindex].Status == "critical" {
						ishealthy = false
						status = healthpair[index].Checks[checkindex].Status
						break
					}else if healthpair[index].Checks[checkindex].Status == "warning"{
						status = healthpair[index].Checks[checkindex].Status
						break
					}else{
						status = healthpair[index].Checks[checkindex].Status
                                                break	
					}
			}
		}
	}
	if !ishealthy {
		beego.Info(ip + " service health status is "+ status)
		beego.Info("Give up switching to async replication")
                beego.Info("Monitor Handler Completed")
		return
	}	
	beego.Info(ip + " service health status is "+ status)
	dsName := username + ":" + password + "@tcp(" + "localhost" + ":" + port + ")/"
	db, err := sql.Open("mysql", dsName)
	if err != nil {
		beego.Error("Create connection object to peer database failed!", err)
		beego.Info("Give up switching to async replication")
                beego.Info("Monitor Handler Completed")
		return
	}
	beego.Info("Create connection object to peer database successfully!")
	defer db.Close()
	err = db.Ping()
	if err != nil {
		beego.Error("Connected to the peer database failed!", err)
		beego.Info("Give up switching to async replication")
                beego.Info("Monitor Handler Completed")
		return
	}
	beego.Info("Connected to the peer database successfully!")
	keepsyncrepl := "set global rpl_semi_sync_master_keepsyncrepl=" + strconv.Itoa(issynchronous)
	_, err = db.Query(keepsyncrepl)
	if err != nil {
		beego.Error("Set rpl_semi_sync_master_keepsyncrepl=0 failed!")
                beego.Info("Monitor Handler Completed")
		return
	}
	beego.Info("Set rpl_semi_sync_master_keepsyncrepl=0 successfully!")
	trysyncrepl := "set global rpl_semi_sync_master_trysyncrepl=" + strconv.Itoa(issynchronous)
	_, err = db.Query(trysyncrepl)
	if err != nil {
		beego.Error("Set rpl_semi_sync_master_trysyncrepl=0 failed!")
                beego.Info("Monitor Handler Completed")
		return
	}
	beego.Info("Set rpl_semi_sync_master_trysyncrepl=0 successfully!")
	beego.Info("Switching local database to async replication!")
	beego.Info("Monitor Handler Completed")
}

func SetRepl_err_counter(hostname string){
        count := 0
        var put string
//        other_hostname := beego.AppConfig.String("otherhostname")
        put = "1"
        kvvalue := []byte(put)
        kvotherhostname := consulapi.KVPair{
                Key:   "monitor/" + hostname,
                Value: kvvalue,
        }
   try:  _, err := kv.Put(&kvotherhostname, nil)
        if err != nil {
                beego.Error("Set peer database repl_err_counter to 1 in CS failed!", err)
                if count ==2 {
                        beego.Info("Monitor Handler Completed")
                        return
                }
                count++
                goto try
        }
        beego.Info("Set peer database repl_err_counter to 1 in CS successfully!")
        beego.Info("MHA Handler Completed")
}

func checkio_thread(ip, port, username, password, addr string) {
	beego.Info("Monitor Handler Sleep 60s!")
	beego.Info("Connecting to peer database......")
	time.Sleep(60000 * time.Millisecond)
	dsName := username + ":" + password + "@tcp(" + addr + ":" + port + ")/"
	db, err := sql.Open("mysql", dsName)
	if err != nil {
		beego.Error("Create connection object to peer database failed!", err)
		beego.Info("Give up switching to async replication")
                beego.Info("Monitor Handler Completed")
		return
	}
	beego.Info("Create connection object to peer database successfully!")
	defer db.Close()
	err = db.Ping()
	if err != nil {
		beego.Error("Connected to the peer database failed!", err)
		Set(ip, port, username, password, 0)
		return
	}
	beego.Info("Connected to the peer database successfully!")
	row, err := db.Query("show slave status")
	if err != nil {
		beego.Error("Checking peer database I/O thread status. Failed!", err)
		Set(ip, port, username, password, 0)
		return
	}
	beego.Info("Checking peer database I/O thread status. Successfully!")
	cols, _ := row.Columns()
	buffer := make([]interface{}, len(cols))
	data := make([]interface{}, len(cols))
	for i, _ := range buffer {
		buffer[i] = &data[i]
	}
	for row.Next() {
		err = row.Scan(buffer...)
		if err != nil {
			beego.Error("Resolve slave status failed!", err)
			beego.Info("Give up switching to async replication")
                	beego.Info("Monitor Handler Completed")
			return
		}
	}
	mapField2Data := make(map[string]interface{}, len(cols))
	for k, col := range data {
		mapField2Data[cols[k]] = col
	}
	Slave_IO_Running := mapField2Data["Slave_IO_Running"]
	if string(Slave_IO_Running.([]uint8)) != "Yes" {
		beego.Info("The I/O thread status is "+ string(Slave_IO_Running.([]uint8)) + "!")
		Set(ip, port, username, password, 0)
		return
	}
	beego.Info("The I/O thread status is " + string(Slave_IO_Running.([]uint8)) + "!")
	beego.Info("Give up switching to async replication.")
	beego.Info("Monitor Handler Completed")
}
