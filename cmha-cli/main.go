package main

import (
	"fmt"

	"os"
	"strings"
	"unsafe"

	"github.com/0-T-0/go.linenoise"
	"github.com/upmio/cmha-cli/cliconfig"
)

var (
	service_ip              []string
	VERSION                 = "1.1.8"
	Prog                    = "CMHA CLI"
	history_file            = ".cmha-cli.history"
	max_number_history_line = 300
	service_name            = "CS"
	cmdline                 = ""
	nodes_data              [2][3]string
)

func welcomebanner() {
	banner := `
Welcome to the %s.
Version: %s


Copyright [2016] [BSG China]

Licensed under the Apache License, Version 2.0 (the "License");

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

`
	fmt.Printf(banner, ColorRender(Prog, COLOR_SUM), VERSION)

}

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
			fmt.Println("version", VERSION)
			return
		}
	}

	client, err := cliconfig.Consul_Client_Init()
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

	linenoise.SetMultiLine(true) //with multi line set,
	//the input will display in multi line when you type more than a line long
	linenoise.SetCompletionCallback(cb)
	linenoise.HistoryLoad(history_file) //load from disk
	linenoise.HistorySetMaxLen(max_number_history_line)

	welcomebanner()

	for {

		cmdline = ""

		if str, end := linenoise.Scan(renderprompt(service_name)); end {

			fmt.Printf("Unexpected error\n")
			cliquit()

		} else {
			cmdline = str
		}

		linenoise.HistoryAdd(cmdline) //add to memory
		linenoise.HistorySave(history_file)

		fields := strings.Fields(cmdline)

		if len(fields) == 0 {
			sayUnrecognized()
			continue
		}

		switch fields[0] {
		case "exit", "quit", "\\q":
			cliquit()

		case "clear", "\\c":
			linenoise.ClearScreen()

		case "help", "\\h", "\\?":
			helpCMD(fields[1:])
		case "show", "\\s":

			if len(fields) < 2 {
				help("show")
				continue
			}

			switch fields[1] {
			case "cluster":
				Cluster()

			case "status":
				Status()
			case "alert_boot":
				if service_name == "" {
					fmt.Printf("Must to be use a service. Before use command 'show %s' .\n", fields[1])
					help("use")
					continue
				}
				AlertBoot(service_name)
			case "alerts":
				if service_name == "" {
					fmt.Println("Welcome see CS alerts info.")
					if len(fields) == 2 {
						if err := show_all_alerts("", ""); err != nil {
							fmt.Println(err.Error())
						}
					} else if len(fields) == 3 {
						if fields[2] == "list" {
							if err := show_all_alerts("", "list"); err != nil {
								fmt.Println(err.Error())
							}
						} else if len(fields[2]) == 19 {
							if err := show_alerts_detail("", fields[2]); err != nil {
								fmt.Println(err.Error())
							}
						} else if len(fields[2]) == 39 {
							_ids := strings.Split(fields[2], ",")

							if len(_ids) != 2 {
								fmt.Println("%s is not a valid alerts id.", fields[2])
								continue

							}
							if err := show_alerts("", _ids[0], _ids[1]); err != nil {
								fmt.Println(err.Error())
							}
						}
					}
				} else {
					fmt.Println("Welcome see " + service_name + " alerts info.")
					if len(fields) == 2 {
						if err := show_all_alerts(service_name, ""); err != nil {
							fmt.Println(err.Error())
						}
					} else if len(fields) == 3 {
						if fields[2] == "list" {
							if err := show_all_alerts(service_name, "list"); err != nil {
								fmt.Println(err.Error())
							}
						} else if len(fields[2]) == 19 {
							if err := show_alerts_detail(service_name, fields[2]); err != nil {
								fmt.Println(err.Error())
							}
						} else if len(fields[2]) == 39 {
							_ids := strings.Split(fields[2], ",")

							if len(_ids) != 2 {

								fmt.Println("%s is not a valid alerts id.", fields[2])
								continue

							}
							if err := show_alerts(service_name, _ids[0], _ids[1]); err != nil {
								fmt.Println(err.Error())
							}
						}
					}
				}
			case "mhalog":
				if service_name == "" {
					fmt.Printf("Must to be use a service. Before use command 'show %s' .\n", fields[1])
					help("use")
					continue
				}

				if len(fields) <= 2 {

					help("show", fields[1])
					continue
				}

				if len(fields) == 3 {
					node_name := fields[2]
					if err := ShowLog(service_name, node_name, LOG_MHA); err != nil {
						fmt.Println(err.Error())
					}
				}

				if len(fields) > 3 {
					node_name := fields[2]
					if err := ShowLog(service_name, node_name, LOG_MHA, fields[3:]...); err != nil {
						fmt.Println(err.Error())
					}
				}

			case "monitorlog":
				if service_name == "" {
					fmt.Printf("Must to be use a service. Before use command 'show %s' .\n", fields[1])
					help("use")
					continue
				}

				if len(fields) <= 2 {
					help("show", fields[1])
					continue
				}

				if len(fields) == 3 {
					node_name := fields[2]
					if err := ShowLog(service_name, node_name, LOG_MONITOR); err != nil {
						fmt.Println(err.Error())
					}
				}

				if len(fields) > 3 {
					node_name := fields[2]
					if err := ShowLog(service_name, node_name, LOG_MONITOR, fields[3:]...); err != nil {
						fmt.Println(err.Error())
					}
				}

			case "info":

				if service_name == "" {
					fmt.Printf("Must to be use a service. Before use command '%s' .\n", "show info")
					help("use")
					continue
				} else if service_name == "CS" {
					help("show")
				} else {
					Service_Status(service_name)
				}

			default:
				help("show")
			}

		case "use", "\\u":

			if len(fields) != 2 {
				help("use")
				continue
			}

			//serive name
			if fields[1] != "CS" {
				if ok, err := Use_service(fields[1]); ok {

					service_name = fields[1]
					fmt.Printf("Service has changed to %s\n", ColorRender(service_name, COLOR_SUM))

					if _nodes_data, err := get_Nodes_Info(service_name); err != nil {
						fmt.Printf("%s\n", err.Error())
					} else {

						for i, _ := range _nodes_data {
							for j, v := range _nodes_data[i] {
								nodes_data[i][j] = v
							}

						}
					}

				} else {
					fmt.Printf("%s\n", err)
					help("use")
					continue
				}
			} else {
				service_name = fields[1]
				fmt.Printf("Service has changed to %s\n", ColorRender(service_name, COLOR_SUM))
				continue
			}

		case "purge", "\\p":

			if service_name == "" {
				fmt.Printf("Must to be use a service. Before use command '%s' .\n", fields[0])
				help("use")
				continue
			}
			if len(fields) == 1 {
				help(fields[0])
				continue

			}
			if fields[1] == "alerts" {
				if len(fields) < 3 {

					help(fields[0])
					continue

				}
			} else {
				if len(fields) < 4 {

					help(fields[0])
					continue

				}
			}

			var logtype string

			if fields[1] == "mhalog" {
				logtype = LOG_MHA

			} else if fields[1] == "monitorlog" {

				logtype = LOG_MONITOR
			} else if fields[1] == "alerts" {
				logtype = LOG_ALERT
			} else {
				help(fields[0])
				continue
			}
			var _timestamps []string
			var node_name string
			if logtype == LOG_ALERT {
				_timestamps = fields[2:]
			} else {
				node_name = fields[2]

				_timestamps = fields[3:]
			}

			if _timestamps[0] == "all" {

				if logtype == LOG_ALERT {
					if err := PurgeAlertLog(service_name, logtype); err != nil {
						fmt.Printf("\n%s\n", err.Error())
						help(fields[0])
						continue
					}

				} else {
					if err := PurgeLog(service_name, node_name, logtype); err != nil {

						fmt.Printf("\n%s\n", err.Error())
						help(fields[0])
						continue
					}
				}

			} else {
				if logtype == LOG_ALERT {
					if len(_timestamps[0]) == 39 {
						if err := PurgeAlertLog(service_name, logtype, _timestamps...); err != nil {
							fmt.Printf("\n%s\n", err.Error())
							help(fields[0])
							continue
						}
					} else {
						help(fields[0])
						continue
					}
				} else {
					if err := PurgeLog(service_name, node_name, logtype, _timestamps...); err != nil {

						fmt.Printf("\n%s\n", err.Error())
						help(fields[0])
						continue
					}
				}

			}
		case "set":

			if service_name == "" {
				fmt.Printf("Must to be use a service. Before use command '%s' .\n", fields[0])
				help("use")
				continue
			} else if len(fields) < 3 {
				help("set")
			} else {
				setCMD(fields[1:], service_name)
			}
		default:
			sayUnrecognized()
		}

	}
}

func show_monitor_log() error {

	return nil

}

func show_mha_log() error {

	return nil

}

func helpCMD(args []string) {
	if len(args) < 1 {
		help()
	} else {
		help(args[0])
	}

}

func setCMD(args []string, service_name string) {
	if service_name == "CS" {
		if len(args) == 2 {
			switch args[0] {
			case "alert_boot":
				if args[1] == "on" || args[1] == "off" {
					Setmonitorkv(service_name, "alert", args[1])
				} else {
					fmt.Println("Parameter error,eg:set alert on/off")
				}
			default:
				sayUnrecognized()
			}
		}
	} else {

		if len(args) == 1 {
			switch args[0] {
			case "vsr":
				fmt.Println("Parameter error,eg:set vsr 192.168.2.1 3306 on")
			case "read_only":
				fmt.Println("Parameter error,eg:set read_only 192.168.2.1 3306 on")
			case "slave_start":
				fmt.Println("Parameter error,eg:set slave_start 192.168.2.1 3306")
			case "slave_stop":
				fmt.Println("Parameter error,eg:set slave_stop 192.168.2.1 3306")
			default:
				sayUnrecognized()
			}
		} else if len(args) == 2 {
			switch args[0] {
			case "repl_err_counter":

				_nodename := args[1]
				if is_valid_node(_nodename) {
					Setmonitorkv(_nodename, service_name, "repl")

				} else {
					say_invalid_node(_nodename)
				}
			case "alert_boot":
				if args[1] == "on" || args[1] == "off" {
					Setmonitorkv(service_name, "alert", args[1])
				} else {
					fmt.Println("Parameter error,eg:set alert on/off")
				}

			case "vsr":
				fmt.Println("Parameter error,eg:set vsr 192.168.2.1 3306 on")
			case "read_only":
				fmt.Println("Parameter error,eg:set read_only 192.168.2.1 3306 on")
			case "slave_start":
				fmt.Println("Parameter error,eg:set slave_start 192.168.2.1 3306")
			case "slave_stop":
				fmt.Println("Parameter error,eg:set slave_stop 192.168.2.1 3306")

			default:
				sayUnrecognized()
			}
		} else if len(args) == 3 {
			switch args[0] {
			case "slave_start":

				_ip := args[1]
				_port := args[2]
				if is_valid_node(_ip, _port) {

					SetSlaveStart(_ip, _port)

				} else {
					say_invalid_node(_ip, _port)
				}

			case "slave_stop":

				_ip := args[1]
				_port := args[2]
				if is_valid_node(_ip, _port) {

					SetSlaveStop(_ip, _port)

				} else {
					say_invalid_node(_ip, _port)
				}

			case "vsr":
				fmt.Println("Parameter error,eg:set vsr 192.168.2.1 3306 on")
			case "read_only":
				fmt.Println("Parameter error,eg:set read_only 192.168.2.1 3306 on")
			default:
				sayUnrecognized()
			}
		} else if len(args) == 4 {

			switch args[0] {
			case "vsr":

				_ip := args[1]
				_port := args[2]
				if is_valid_node(_ip, _port) {

					SetVsr(_ip, _port, args[3])

				} else {
					say_invalid_node(_ip, _port)
				}

			case "read_only":

				_ip := args[1]
				_port := args[2]
				if is_valid_node(_ip, _port) {

					SetReadOnly(_ip, _port, args[3])

				} else {
					say_invalid_node(_ip, _port)
				}

			case "slave_start":
				fmt.Println("Parameter error,eg:set slave_start 192.168.2.1 3306")
			case "slave_stop":
				fmt.Println("Parameter error,eg:set slave_stop 192.168.2.1 3306")
			default:
				sayUnrecognized()

			}

		} else {
			fmt.Println("Missing Parameters")
			help("set")
		}
	}
}

func say_invalid_node(_nodedata ...string) {

	var node_info string
	if len(_nodedata) == 1 {
		node_info = _nodedata[0]
	}

	if len(_nodedata) == 2 {
		node_info = fmt.Sprintf("%s:%s", _nodedata[0], _nodedata[1])
	}

	fmt.Println()
	fmt.Printf("There is NOT a node in %s, such as %s\n", service_name, node_info)
	fmt.Println()

}

// verify a node in service by node name, or node ip and  port.
// is_valid_node(nodename)
// is_valid_node(ip,port)
func is_valid_node(_nodedata ...string) bool {
	if len(_nodedata) < 1 {
		return false
	}

	var ret bool = false

	// verify node name
	if len(_nodedata) == 1 {
		for _, node_data := range nodes_data {

			if _nodedata[0] == node_data[0] {
				ret = true
				continue
			}
		}
	}

	// verify ip and port
	if len(_nodedata) == 2 {
		for _, node_data := range nodes_data {

			if _nodedata[0] == node_data[1] &&
				_nodedata[1] == node_data[2] {
				ret = true
				continue
			}
		}
	}

	return ret

}

func cliquit() {

	fmt.Printf("Bye\n")
	os.Exit(0)

}

func cleanCommand(cmd string) ([]string, error) {
	cmd_args := strings.Fields(cmd)
	return cmd_args, nil
}

func sayUnrecognized() {

	fmt.Println("Unknown command. Use 'help'.")
	fmt.Println("")

}

func renderprompt(servicename string) string {

	sname := servicename

	if len(sname) == 0 {
		sname = "CS"
	}

	return fmt.Sprintf("CMHA [%s]> ", sname)
}
