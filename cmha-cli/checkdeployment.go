package main

import (
	"fmt"

	"github.com/upmio/cmha-cli/cliconfig"
)

func CheckDeployment(args ...string) error {

	client, err := cliconfig.Consul_Client_Init()

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
			key = "cmha/service/" + k + "/db/leader"
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
