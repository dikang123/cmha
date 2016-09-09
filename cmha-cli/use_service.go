package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/upmio/cmha-cli/cliconfig"
)

func Use_service(args ...string) (bool, error) {
	if len(args) != 1 {
		err := errors.New("Error: Only accept one argument; But receive argument number : " + fmt.Sprint(len(args)))
		return false, err
	}

	client, err := cliconfig.Consul_Client_Init()

	if err != nil {
		fmt.Println("cluster Create consul-api client failure!", err)
		return false, err
	}

	status := client.Status()
	peers, err := status.Peers()
	if err != nil {
		fmt.Println("cluster Query consul service failure!", err)
		return false, err
	}

	agent := client.Agent()
	members, err := agent.Members(false)
	if err != nil {
		fmt.Println("cluster Query cluster members failure!", err)
		return false, err
	}

	peerssize := len(peers)
	var count = peerssize
	slice := []string{}
	for _, peersvalue := range peers {
		a := strings.Split(peersvalue, ":")
		for _, membersvalue := range members {
			if strings.EqualFold(a[0], membersvalue.Addr) {
				if membersvalue.Status != 1 {
					a := "Fault machine: " + membersvalue.Addr
					slice = append(slice, a)
					count -= 1
				}
			}
		}
	}

	catalog := client.Catalog()
	services, _, err := catalog.Services(nil)

	if err != nil {
		fmt.Println("cluster Query services failure!", err)
		return false, err
	}

	var IsServiceName bool = false

	for _k, _ := range services {
		if _k == "Statistics" {
			continue
		}
		if _k == args[0] {
			IsServiceName = true
			break
		}
	}

	if IsServiceName {
		return true, nil
	} else {
		err := errors.New("Error: No such service by name - " + args[0])
		return false, err
	}

}

func get_Nodes_Info(servicename string) ([][]string, error) {

	client, err := cliconfig.Consul_Client_Init()

	if err != nil {
		return nil, err
	}

	health := client.Health()

	health_service, _, err := health.Service(servicename, "", false, nil)
	if err != nil {
		return nil, err
	}

	if len(health_service) < 1 {
		err := errors.New(fmt.Sprintf("Not exist service such as %s.\n", servicename))
		return nil, err
	}

	_data := make([][]string, 2)

	var i int = 0

	for _, servicedata := range health_service {

		node_name := servicedata.Node.Node
		node_address := servicedata.Node.Address
		node_port := servicedata.Service.Port

		// skip chap node
		node_tag := servicedata.Service.Tags[0]
		if strings.HasPrefix(node_tag, "chap-") {
			continue

		}

		node_info := make([]string, 3)
		node_info[0] = node_name
		node_info[1] = node_address
		node_info[2] = fmt.Sprintf("%d", node_port)

		_data[i] = node_info

		i++

	}

	return _data, nil

}
