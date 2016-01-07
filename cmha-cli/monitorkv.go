package main

import (
	"fmt"
	"github.com/astaxie/beego"
	consulapi "github.com/hashicorp/consul/api"
)

func MonitorKv(args ...string) error {
	if len(args) > 0 {
		config := &consulapi.Config{
			Datacenter: beego.AppConfig.String("cmha-datacenter"),
			Token:      beego.AppConfig.String("cmha-token"),
			Address:    beego.AppConfig.String("cmha-server-ip") + ":8500",
		}
		client, err := consulapi.NewClient(config)
		if err != nil {
			fmt.Println("dbleader Create consul-api client failure!", err)
			return err
		}
		kv := client.KV()
		key := "monitor/" + args[0]
		kvpair, _, err := kv.Get(key, nil)
		if err != nil {
			fmt.Println("Get key " + key + " failure!", err)
			return err
		}
		if kvpair == nil {
			fmt.Println("No key "+ key +" or key value is null")
			return nil
		}
		value := string(kvpair.Value)
		fmt.Println(value)
	}
	return nil
}
