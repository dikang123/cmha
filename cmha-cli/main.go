package main

import (
	"fmt"
	"github.com/HardySimpson/linenoise"
	"github.com/astaxie/beego"
	consulapi "github.com/hashicorp/consul/api"
	"os"
	"strings"
	"unsafe"
)

var service_ip []string

// this func is called every time when people hit the tab button
func cb(buf string, lc unsafe.Pointer) {
	if !strings.HasPrefix(buf, "cd") {
		return
	}
	words := strings.Fields(buf)
	if len(words) == 1 {
		return
	}
	baseDir, _ := os.Open(".")
	names, _ := baseDir.Readdirnames(0)
	for _, n := range names {
		if strings.HasPrefix(n, words[1]) {
			linenoise.AddCompletion(lc, "cd "+n)
		}
	}
	return
}

func main() {
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "-v" || arg == "--version" {
			fmt.Println("version 1.1.4")
			return
		} else {
			return
		}
	}
	config := &consulapi.Config{
		Datacenter: beego.AppConfig.String("cmha-datacenter"),
		Token:      beego.AppConfig.String("cmha-token"),
		Address:    beego.AppConfig.String("cmha-server-ip") + ":8500",
	}
	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Println(" Create consul-api client failure!", err)
		return
	}
	status := client.Status()
	leader, err := status.Leader()
	if err != nil {
		fmt.Println("Please reset cmha tool configure  file cmha-server-ip or no cluster leader please create cluster!")
		return
	}
	if leader == "" {
		fmt.Println("No cluster leader!")
		return
	}
	/*service_ip = beego.AppConfig.Strings("cmha-server-ip")
	//fmt.Println(service_ip)
	config := &consulapi.Config{
		Datacenter: beego.AppConfig.String("datacenter"),
		Token:      beego.AppConfig.String("token"),
		//	Address: beego.AppConfig.String("service_ip") + ":" + beego.AppConfig.String("service_port"),
	}
	//	var kvPair *consulapi.KVPair
	var client *consulapi.Client
	var err error
	var a int
	for i, _ := range service_ip {
		//fmt.Println("aaaaaa", service_ip[i])
		config.Address = service_ip[i] + ":" + "8500"
		client, err = consulapi.NewClient(config)
		if err != nil {
			//beego.Error("Create consul-api client failure!", err)
			return
		}
		//beego.Info(" Create consul-api client success!")
		status := client.Status()
		leader, err := status.Leader()
		if err != nil {
			//fmt.Println("Please reset cmha tool configure  file cmha-server-ip or no cluster leader please create cluster!", service_ip[i])
			//fmt.Println(service_ip[i])
			//return
			fmt.Println("configure file cmha-server-ip err", service_ip[i])
			a += 1
			continue
		}
		if leader == "" {
			fmt.Println("No cluster leader!", service_ip[i])
			//return
			a += 1
			continue
		}
		break
	}
	if a == len(service_ip) {
		fmt.Println("Please reset cmha tool configure  file cmha-server-ip or no cluster leader please create cluster!")
		return
	}*/
	linenoise.SetMultiLine(true) //with multi line set,
	//the input will display in multi line when you type more than a line long
	linenoise.SetCompletionCallback(cb)
	linenoise.HistoryLoad("history.txt") //load from disk
	linenoise.HistorySetMaxLen(10)
	help()
MainLoop: //max line keep in memory and disk
	for {
		line, end := linenoise.Scan("cmha > ")
		if end {
			return
		}

		linenoise.HistoryAdd(line) //add to memory
		linenoise.HistorySave("history.txt")
		arr := strings.Trim(line, " ") //add to disk
		cmds, _ := cleanCommand(arr)
		//fmt.Println(len(cmds))
		if len(cmds) == 1 {
			//		Rout:
			switch arr {
			case "exit", "quit":
				fmt.Println("Exiting...")
				break MainLoop
			case "clear":
				linenoise.ClearScreen()
			case "help":
				help()
			case "show":
				help("show")
			case "set":
				help("set")
			/*case "check":
				help("check")
			case "monitor":
				help("monitor")*/
			default:
				fmt.Println("Unknown command")
			}
		} else if len(cmds) == 2 {
			switch cmds[0] {
			case "show":
				switch cmds[1] {
				case "cluster":
					Cluster()
				case "chap":
					Chap()
				case "db":
					Db()
				case "dbleader":
					DbLeader()
				case "cmhaleader":
					CmhaLeader()
				/*case "monitorkv":
					MonitorKv()*/
				default:
					fmt.Println("Unknown command")
				}
			case "set":
				switch cmds[1] {
				case "vsr":
					fmt.Println("Parameter error,eg:set vsr 192.168.2.1 3306 on")
				case "read_only":
					fmt.Println("Parameter error,eg:set read_only 192.168.2.1 3306 on")
				case "slave_start":
					fmt.Println("Parameter error,eg:set slave_start 192.168.2.1 3306")
				case "slave_stop":
					fmt.Println("Parameter error,eg:set slave_stop 192.168.2.1 3306")
				default:
					fmt.Println("Unknown command")
				}
			default:
				fmt.Println("Unknown command")
			}
		} else if len(cmds) == 3 {
			switch cmds[0] {
			case "show":
				switch cmds[1] {
				case "cluster":
					Cluster(cmds[2])
				case "chap":
					Chap(cmds[2])
				case "db":
					Db(cmds[2])
				case "dbleader":
					DbLeader(cmds[2])
				case "repl_err_counter":
					MonitorKv(cmds[2])
				default:
					fmt.Println("Unknown command")
				}
			case "set":
				switch cmds[1] {
				case "repl_err_counter":
					Setmonitorkv(cmds[2])
				case "vsr":
					fmt.Println("Parameter error,eg:set vsr 192.168.2.1 3306 on")
				case "read_only":
					fmt.Println("Parameter error,eg:set read_only 192.168.2.1 3306 on")
				case "slave_start":
					fmt.Println("Parameter error,eg:set slave_start 192.168.2.1 3306")
				case "slave_stop":
					fmt.Println("Parameter error,eg:set slave_stop 192.168.2.1 3306")
				default:
					fmt.Println("Unknown command")
				}
			default:
				fmt.Println("Unknown command")
			}
		} else if len(cmds) == 4 {
			switch cmds[0] {
			case "set":
				switch cmds[1] {
				case "slave_start":
					SetSlaveStart(cmds[2], cmds[3])
				case "slave_stop":
					SetSlaveStop(cmds[2], cmds[3])
				case "vsr":
					fmt.Println("Parameter error,eg:set vsr 192.168.2.1 3306 on")
				case "read_only":
					fmt.Println("Parameter error,eg:set read_only 192.168.2.1 3306 on")
				}
			}
		} else if len(cmds) == 5 {
			switch cmds[0] {
			case "set":
				switch cmds[1] {
				case "vsr":
					SetVsr(cmds[2], cmds[3], cmds[4])
				case "read_only":
					SetReadOnly(cmds[2], cmds[3], cmds[4])
				case "slave_start":
					fmt.Println("Parameter error,eg:set slave_start 192.168.2.1 3306")
				case "slave_stop":
					fmt.Println("Parameter error,eg:set slave_stop 192.168.2.1 3306")
				}
			}
		} else {
			fmt.Println("Missing Parameters")
		}
	}
}
func cleanCommand(cmd string) ([]string, error) {
	//cmd_args := strings.Split(strings.Trim(cmd, " \n"), " ")
	cmd_args := strings.Fields(cmd)
	return cmd_args, nil
}
