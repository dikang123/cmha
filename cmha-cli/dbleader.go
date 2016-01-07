package main

import (
	"fmt"
	"github.com/astaxie/beego"
	consulapi "github.com/hashicorp/consul/api"
)

func DbLeader(args ...string) error {
	if len(args) > 0 {
		config := &consulapi.Config{
			Datacenter: beego.AppConfig.String("cmha-datacenter"),
			Token:      beego.AppConfig.String("cmha-token"),
			Address:    beego.AppConfig.String("cmha-server-ip") + ":8500",
		}
		/*config := &consulapi.Config{
			Datacenter: beego.AppConfig.String("cmha-datacenter"),
			Token:      beego.AppConfig.String("cmha-token"),
		}
		var kv *consulapi.KV
		var keys *consulapi.KVPair
		for i, _ := range service_ip {
			config.Address = service_ip[i] + ":" + beego.AppConfig.String("cmha-server-ip")
			client, err := consulapi.NewClient(config)
			if err != nil {
				fmt.Println("dbleader Create consul-api client failure!", err)
				continue
			}
			kv := client.KV()
			keys, _, err := kv.Keys("", "", nil)
			if err != nil {
				fmt.Println("dbleader get keys failure!", err)
				continue
			}
			if keys == nil {
				fmt.Println("dbleader not kv!")
				continue
			}
			break
		}*/
		client, err := consulapi.NewClient(config)
		if err != nil {
			fmt.Println("dbleader Create consul-api client failure!", err)
			return err
		}
		kv := client.KV()
		keys, _, err := kv.Keys("", "", nil)
		if err != nil {
			fmt.Println("dbleader get keys failure!", err)
			return err
		}
		if keys == nil {
			fmt.Println("dbleader not kv!")
			return nil
		}
		key := "service/" + args[0] + "/leader"
		var iskey = false
		for value := range keys {
			//fmt.Println(key, keys[value])
			if key == keys[value] {
				iskey = true
				break
			} else {
				continue
			}
		}
		if !iskey {
			fmt.Println("not ", args[0], " kv")
			return nil
		}
		kvpair, _, err := kv.Get(key, nil)
		if err != nil {
			fmt.Println("dbleader Get key failure!", err)
			return err
		}
		if kvpair.Session == "" {
			fmt.Println(args[0], " leader not exist!")
			return nil
		}
		fmt.Println("----------------------------------------")
		value := string(kvpair.Value)
		fmt.Println(value)
		fmt.Println("----------------------------------------")
		return nil
	}
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
	catalog := client.Catalog()
	services, _, err := catalog.Services(nil)
	if err != nil {
		fmt.Println("db Query services failure!", err)
		return err
	}
	for k, _ := range services {
		if k != "consul" {
			var key string
			fmt.Println("----------------------------------------")
			fmt.Println("      ", k)
			key = "service/" + k + "/leader"
			kvpair, _, err := kv.Get(key, nil)
			if err != nil {
				fmt.Println("dbleader Get key failure!", err)
				return err
			}
			if kvpair.Session == "" {
				fmt.Println(k, " leader not exist!")
				continue
			}
			value := string(kvpair.Value)
			fmt.Println(value)
		}
	}
	fmt.Println("----------------------------------------")
	return nil
}
