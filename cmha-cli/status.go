package main

import (
	"fmt"

	"github.com/upmio/cmha-cli/cliconfig"
	"sort"
	"strings"
)

func Status(args ...string) error {

	client, err := cliconfig.Consul_Client_Init()

	if err != nil {
		fmt.Println("Error: status.go Create consul-api client is failure!", err)
		return err
	}

	status := client.Status()

	leader, err := status.Leader()
	if err != nil {
		fmt.Println("Error: [status] Query for a known leader is failure!", err)
		return err
	}

	servers, err := status.Peers()
	if err != nil {
		fmt.Println("Error: [status] Query for a known raft peers is failure!", err)
		return err
	}

	agent := client.Agent()
	members, err := agent.Members(false)
	if err != nil {
		fmt.Println("cluster Query cluster members failure!", err)
		return err
	}

	// server info map["addr:port"]:[name,status]
	servers_info := make(map[string][2]string)

	for _, _v := range members {

		// only use consul server
		if _v.Tags["role"] != "consul" {
			continue
		}

		_hostname := _v.Name
		var _status string = ColorRender("OK", COLOR_NORMAL)
		if _v.Status != 1 {
			_status = ColorRender("Fail", COLOR_ERROR)
		}

		_server_info := [2]string{_hostname, _status}

		server_addr := _v.Addr
		servers_info[server_addr] = _server_info

	}

	_len := len(servers)

	_data := make([][]string, _len)

	sort.Strings(servers)

	for i := range _data {

		_server_item := servers[i]

		_server_ip := strings.Split(_server_item, ":")[0]

		_server_info := servers_info[_server_ip]

		_data[i] = make([]string, 5)

		if _server_item == leader {

			//service
                        _data[i][0] = "CS"
			// name
			_data[i][1] = ColorRender(_server_info[0], COLOR_SUM)
			// ip & port
			_data[i][2] = ColorRender(_server_item, COLOR_SUM)

			// leader
			_data[i][3] = ColorRender("y", COLOR_SUM)

		} else {
			//service
			_data[i][0] = "CS" 
			// name
			_data[i][1] = _server_info[0]
			// ip & port
			_data[i][2] = _server_item
			_data[i][3] = ""
		}

		// status
		_data[i][4] = _server_info[1]

	}

	_th := []string{
		"Service",
		"Name",
		"Consul Server (ip:port)",
		"Leader",
		"Status",
	}

	TableRender(_th, _data, ALIGN_CENTRE)

	fmt.Println("")

	return nil
}
