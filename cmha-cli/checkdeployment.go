package main

import (
	"fmt"
	"github.com/astaxie/beego"
	consulapi "github.com/hashicorp/consul/api"
)

func CheckDeployment(args ...string) error {
	config := &consulapi.Config{
		Datacenter: beego.AppConfig.String("cmha-datacenter"),
		Token:      beego.AppConfig.String("cmha-token"),
		Address:    beego.AppConfig.String("cmha-server-ip") + ":8500",
	}
	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Println("deploymnet.go Create consul-api client failure!", err)
		return err
	}
	catalog := client.Catalog()
	services, _, err := catalog.Services(nil)
	kv := client.KV()
	health := client.Health()
	if err != nil {
		fmt.Println("deployment.go Query services failure!", err)
		return err
	}

	var isdeployment = true
	for k, _ := range services {
		if k != "consul" {
			var ishealth = true
			var key string
			key = "service/" + k + "/leader"
			kvpair, _, err := kv.Get(key, nil)
			if err != nil {
				fmt.Println("deployment.go Get key failure!", err)
				return err
			}
			if kvpair == nil {
				fmt.Println("deployment.go Not found " + key + ", Please create!")
				continue
			}
			if kvpair.Value == nil && kvpair.Session == "" {
				fmt.Println("deployment.go Auto deployment failure!")
				continue
			}
			healthservice, _, err := health.Service(k, "", false, nil)
			if err != nil {
				fmt.Println("deployment.go Check cluster health service failure!", err)
				continue
			}
			for index := range healthservice {
				for checkindex := range healthservice[index].Checks {
					if healthservice[index].Checks[checkindex].Status != "passing" {
						ishealth = false
					}
				}
				if !ishealth {
					isdeployment = false
					break

				}
			}
			if !isdeployment {
				fmt.Println("Check auto deployment error,Please redeployment!")
				break
			}
		}
	}
	if isdeployment {
		fmt.Println("Check auto deployment OK !")
	}
	return nil
}
