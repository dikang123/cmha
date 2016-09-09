package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/upmio/cmha-cli/cliconfig"
)

var username = cliconfig.GetUserName()
var password = cliconfig.GetPassword()

func SetVsr(args ...string) error {
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
			switch args[2] {
			case "on":
				fmt.Printf("Set VSR %s success;Replication switch to VSR.\n", args[2])
			case "off":
				fmt.Printf("Set VSR %s success;Replication switch to Async.\n", args[2])
			}
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

			var rw_mode string = ""
			switch args[2] {
			case "on":
				rw_mode = "Read-Only"
			case "off":
				rw_mode = "Read-Write"
			}

			if rw_mode == "" {
				fmt.Printf("%s:%s set readyonly to %s mode\n", args[0], args[1], args[2])
			} else {
				fmt.Printf("%s:%s switch to %s mode\n", args[0], args[1], rw_mode)
			}
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
	if len(args) >= 3 {

		client, err := cliconfig.Consul_Client_Init()

		if err != nil {
			fmt.Println("Setmonitorkv Create consul-api client failure!", err)
			return err
		}
		kv := client.KV()
		var put string
		var key string
		if args[2] == "repl" {
			key = "cmha/service/" + args[1] + "/db/repl_err_counter/" + args[0]
			put = "0"
		} else if args[1] == "alert" {
			key = "cmha/service/" + args[0] + "/alerts/alert_boot"
			if args[2] == "on" {
				put = "enable"
			} else if args[2] == "off" {
				put = "disable"
			}
		}

		kvvalue := []byte(put)

		kvotherhostname := consulapi.KVPair{
			Key:   key,
			Value: kvvalue,
		}
		_, err = kv.Put(&kvotherhostname, nil)
		if err != nil {
			if args[2] == "repl" {
				fmt.Println("reset "+args[0]+" repl_err_counter failure", err)
			} else if args[1] == "alert" {
				fmt.Println("reset "+args[0]+" alert switch failure", err)
			}

			return err
		}
		if args[2] == "repl" {
			fmt.Println("reset " + args[0] + " repl_err_counter success")
		} else if args[1] == "alert" {
			fmt.Println("reset " + args[0] + " alert switch success")
		}

	}
	return nil
}
