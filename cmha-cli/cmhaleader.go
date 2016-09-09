package main

import (
	"fmt"

	"github.com/upmio/cmha-cli/cliconfig"
)

func CmhaLeader(args ...string) error {

	client, err := cliconfig.Consul_Client_Init()

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
