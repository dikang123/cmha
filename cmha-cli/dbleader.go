package main

import (
	"fmt"

	"github.com/upmio/cmha-cli/cliconfig"
)

func DbLeader(args ...string) error {
	if len(args) > 0 {

		client, err := cliconfig.Consul_Client_Init()

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
		key := "cmha/service/" + args[0] + "/db/leader"
		var iskey = false
		for value := range keys {
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

	client, err := cliconfig.Consul_Client_Init()

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
			key = "cmha/service/" + k + "/db/leader"
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
