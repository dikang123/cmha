package main

import (
	"fmt"
	"github.com/astaxie/beego"
	consulapi "github.com/hashicorp/consul/api"
)

func CmhaLeader(args ...string) error {
	config := &consulapi.Config{
		Datacenter: beego.AppConfig.String("cmha-datacenter"),
		Token:      beego.AppConfig.String("cmha-token"),
		Address:    beego.AppConfig.String("cmha-server-ip") + ":8500",
	}
	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Println("leader.go Create consul-api client failure!", err)
		return err
	}
	status := client.Status()
	leader, err := status.Leader()
	if err != nil {
		fmt.Println("leader.go Query leader status failure!", err)
		return err
	}
	fmt.Println("---------------------------------")
	fmt.Println("cmha leader " + leader)
	fmt.Println("---------------------------------")
	return nil
}
