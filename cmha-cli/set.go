package main

import (
	"database/sql"
	"fmt"
	//	"fmt"
	"strconv"
	"strings"
//	"reflect"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	consulapi "github.com/hashicorp/consul/api"
)

var username = beego.AppConfig.String("cmha-to-tool-user")
var password = beego.AppConfig.String("cmha-to-tool-password")

func SetVsr(args ...string) error {
	if len(args) == 3 {
		dsName := username + ":" + password + "@tcp(" + args[0] + ":" + args[1] + ")/"
		db, err := sql.Open("mysql", dsName)
		if err != nil {
			fmt.Println("open database failure")
			return err
		}
		fmt.Println(username, password, args[0], args[1])
		defer db.Close()
		err = db.Ping()
		if err != nil {
			fmt.Println("connection to the database failure")
			return err
		}
		var keep int
		if strings.EqualFold(args[2], "on") {
			keep = 1
		} else if strings.EqualFold(args[2], "on") {
			keep = 0
		}
		keepsyncrepl := "set global rpl_semi_sync_master_keepsyncrepl=" + strconv.Itoa(keep)
		_, err = db.Query(keepsyncrepl)
		if err != nil {
			fmt.Println("set rpl_semi_sync_master_keepsyncrepl failure")
			return err
		}
		trysyncrepl := "set global rpl_semi_sync_master_trysyncrepl=" + strconv.Itoa(keep)
		_, err = db.Query(trysyncrepl)
		if err != nil {
			fmt.Println("set rpl_semi_sync_master_trysyncrepl failure")
			return err
		} else {
			fmt.Println("set rpl_semi_sync_master_keepsyncrepl and rpl_semi_sync_master_trysyncrepl success")
		}
	}
	return nil
}

func SetReadOnly(args ...string) error {
	if len(args) == 3 {
		dsName := username + ":" + password + "@tcp(" + args[0] + ":" + args[1] + ")/"
		db, err := sql.Open("mysql", dsName)
		if err != nil {
			fmt.Println("open database failure")
			return err
		}
		defer db.Close()
		err = db.Ping()
		if err != nil {
			fmt.Println("connection to the database failure")
			return err
		}
		read_only := "set global read_only=" + args[2]
		_, err = db.Query(read_only)
		if err != nil {
			fmt.Println("set read_only failure")
			return err
		} else {
			fmt.Println("set read_only success")
		}
	}
	return nil
}

func SetSlaveStart(args ...string) error {
	if len(args) == 2 {
		dsName := username + ":" + password + "@tcp(" + args[0] + ":" + args[1] + ")/"
		db, err := sql.Open("mysql", dsName)
		if err != nil {
			fmt.Println("open database failure")
			return err
		}
		defer db.Close()
		err = db.Ping()
		if err != nil {
			fmt.Println("connection to the database failure")
			return err
		}
		slave := "start slave"
		_, err = db.Query(slave)
		if err != nil {
			fmt.Println("start slave thread failure")
			return err
		} else {
			fmt.Println("start slave thread success")
		}
	}
	return nil
}

func SetSlaveStop(args ...string) error {
	if len(args) == 2 {
		dsName := username + ":" + password + "@tcp(" + args[0] + ":" + args[1] + ")/"
		db, err := sql.Open("mysql", dsName)
		if err != nil {
			fmt.Println("open database failure")
			return err
		}
		defer db.Close()
		err = db.Ping()
		if err != nil {
			fmt.Println("connection to the database failure")
			return err
		}
		slave := "stop slave"
		_, err = db.Query(slave)
		if err != nil {
			fmt.Println("stop slave thread failure")
			return err
		} else {
			fmt.Println("stop slave thread success")
		}
	}
	return nil
}

func Setmonitorkv(args ...string) error {
	if len(args) == 1 {
		config := &consulapi.Config{
			Datacenter: beego.AppConfig.String("cmha-datacenter"),
			Token:      beego.AppConfig.String("cmha-token"),
			Address:    beego.AppConfig.String("cmha-server-ip") + ":8500",
		}
		client, err := consulapi.NewClient(config)
		if err != nil {
			fmt.Println("Setmonitorkv Create consul-api client failure!", err)
			return err
		}
		kv := client.KV()
		var put string
		put = "0"
		kvvalue := []byte(put)
//		fmt.Println(kvvalue)
		kvotherhostname := consulapi.KVPair{
			Key:   "monitor/" + args[0],
			Value: kvvalue,
		}
		_, err = kv.Put(&kvotherhostname, nil)
		if err != nil {
			fmt.Println("monitor/"+args[0]+" put failure",err)
			return err
		}
		fmt.Println("monitor/"+args[0]+" put success")
	}	
	return nil
}
