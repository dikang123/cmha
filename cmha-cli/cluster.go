package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/upmio/cmha-cli/cliconfig"
)

func Cluster(args ...string) error {

	client, err := cliconfig.Consul_Client_Init()

	if err != nil {
		fmt.Println("cluster Create consul-api client failure!", err)
		return err
	}

	status := client.Status()
	peers, err := status.Peers()
	if err != nil {
		fmt.Println("cluster Query consul service failure!", err)
		return err
	}

	agent := client.Agent()
	members, err := agent.Members(false)
	if err != nil {
		fmt.Println("cluster Query cluster members failure!", err)
		return err
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
	health := client.Health()
	services, _, err := catalog.Services(nil)
	if err != nil {
		fmt.Println("cluster Query services failure!", err)
		return err
	}

	_counter_services := len(services) - 1

	_service_names := make([]string, _counter_services)

	var j = 0
	for _k, _ := range services {
		if _k != "consul" {
			_service_names[j] = _k
			j++
		}
	}

	sort.Strings(_service_names)

	_data := make([][]string, _counter_services)

	for _, servicename := range _service_names {
		if servicename == "Statistics" {
			continue
		}

		health_service, _, err := health.Service(servicename, "", false, nil)
		if err != nil {
			fmt.Printf("Error: Fail to query cluster health service.\n%s\n", err)
			return err
		}

		if len(health_service) < 1 {
			fmt.Printf("Not exist service such as %s.\n", servicename)
			return nil
		}

		for _, servicedata := range health_service {

			node_name := servicedata.Node.Node
			node_address := servicedata.Node.Address

			//node type
			node_tag := servicedata.Service.Tags[0]
			var node_type = "db"
			if strings.HasPrefix(node_tag, "chap-") {
				node_type = "chap"

			}

			node_info := make([]string, 6)
			node_info[0] = servicename
			node_info[1] = node_name
			node_info[2] = node_address
			node_info[3] = node_type

			has_unknown_status := false

			for _, checkdata := range servicedata.Checks {

				node_status := checkdata.Status

				node_checkid := checkdata.CheckID

				has_unknown_status = node_status == "critical" && node_checkid == "serfHealth"

				switch node_status {
				case "passing":
					node_status = ColorRender("OK", COLOR_NORMAL)
				case "critical":
					node_status = ColorRender("Fail", COLOR_ERROR)
				default:
					node_status = ColorRender(node_status, COLOR_WARNNING)
				}

				if node_checkid == "serfHealth" {
					node_info[5] = node_status
				} else {
					node_info[4] = node_status
				}

			}

			if has_unknown_status {
				node_info[4] = ColorRender("UnKnown", COLOR_WARNNING)
			}

			_data = append(_data, node_info)
		}

	}

	_th := []string{
		"service",
		"node",
		"address",
		"type",
		"status",
		"agent status",
	}

	TableRender(_th, _data, ALIGN_CENTRE)
	fmt.Println("")

	return nil
}
