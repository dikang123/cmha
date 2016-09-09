package main

import (
	"fmt"

	"github.com/upmio/cmha-cli/cliconfig"
)

func MonitorKv(args ...string) error {
	if len(args) > 0 {

		client, err := cliconfig.Consul_Client_Init()

		if err != nil {
			fmt.Println("dbleader Create consul-api client failure!", err)
			return err
		}

		kv := client.KV()
		key := "monitor/" + args[0]
		kvpair, _, err := kv.Get(key, nil)

		if err != nil {
			fmt.Println("Get key "+key+" failure!", err)
			return err
		}

		if kvpair == nil {
			fmt.Println("No key " + key + " or key value is null")
			return nil
		}

		value := string(kvpair.Value)
		fmt.Println(value)
	}
	return nil
}
