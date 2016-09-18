package main

import (
	"fmt"
	"github.com/ryanuber/columnize"
)

func help(args ...string) {

	if len(args) < 1 {

		fmt.Printf("List of all %s commands:\n", Prog)

		help0 := []string{
			"? | (\\?) | Synonym for `help'.",
			"clear | (\\c) | Clear the current input statement.",
			"exit | (\\q) | Exit cmha-cli. Same as quit.",
			"help | (\\h) | Display this help.",
			"purge | (\\p) | Clear log .",
			"quit | (\\q) | Quit cmha-cli.",
			"show | (\\s) | Show cluster status or cluster detail.Takes 'status' or 'cluster' as argument.",
			"use | (\\u) | Use a service. Takes service name as argument.",
		}

		if service_name != "" {
			str := "set |  | Set database variables."
			help0 = append(help0, str)

			str = "purge | (\\p) | Clear log ."
			help0 = append(help0, str)
		}

		helpFormat(Prog, help0)

	} else if len(args) == 1 {
		switch args[0] {
		case "show":

			help_show := []string{
				"show status | | Show cluster status.",
				"show cluster | |  Show cluster detail.",
				"show alert_boot | | Show alert switch status." ,
			}
			
			if service_name == "CS"{
				str := "show alerts | [<logid>,[<logid>]] Show range alerts"
                                help_show = append(help_show, str)

                                str = "show alerts list | Show alerts information for the last 3 days"
                                help_show = append(help_show, str)
			}else {
				str := "show info | | Show service status."
				help_show = append(help_show, str)

				str = "show mhalog | <Node name> [<logid>[ <logid>]...] | Show mha log."
				help_show = append(help_show, str)

				str = "show mhalog | <Node name> [<logid>[,<logid>]...] | Show mha log."
				help_show = append(help_show, str)

				str = "show monitorlog | <Node name> [<logid>[ <logid>]...] | Show monitor log."
				help_show = append(help_show, str)

				str = "show monitorlog | <Node name> [<logid>[,<logid>]...] | Show monitor log."
				help_show = append(help_show, str)
				
				str = "show alerts | [<logid>,[<logid>]] Show range alerts"
				help_show = append(help_show, str)
				
				str = "show alerts list | Show alerts information for the last 3 days"
				help_show = append(help_show, str)
			}

			helpFormat(args[0], help_show)

		case "purge":
			help_purge := []string{
				"purge mhalog | <hostname> <all> | Drop host all mha-logs",
				"purge mhalog | <hostname> <logid> | Drop host a mha-log",
				"purge mhalog | <hostname> <logidfrom,logidto> | Drop host some mha-logs",
				"purge monitorlog | <hostname> <all> | Drop host all monitor-logs",
				"purge monitorlog | <hostname> <logid> | Drop host a monitor-log",
				"purge monitorlog | <hostname> <logidfrom,logidto> | Drop host some monitor-logs",
				"purge alerts | <all> | Drop all alerts-logs",
				"purge alerts  | <logidfrom,logidto> | Drop alerts-logs",
			}

			helpFormat(args[0], help_purge)

		case "use":

			help_use := []string{
				"use | (\\u) | Use a service. Takes service name as argument.",
			}

			helpFormat(args[0], help_use)

		case "set":
			var help_set []string
			if service_name == "CS"{
				help_set = []string{
					"set alert_boot | <on/off> | Set alert switch",
				}
                        }else {
				help_set = []string{
					"set vsr | <db_host_ipaddr> | <db_port> | on | Enable database VSR Function | e.g. set vsr 192.168.2.1 3306 on",
					"set read_only | <db_host_ipaddr> | <db_port> | <on|off> | Set database into Read_Only mode | e.g. set read_only 192.168.2.1 3306 on",
					"set slave_start | <db_host_ipaddr> | <db_port> |  | Start database slave threads | e.g. set slave_start 192.168.2.1 3306",
					"set slave_stop | <db_host_ipaddr> | <db_port> |  | Stop database slave threads | e.g. set slave_stop 192.168.2.1 3306",
					"set repl_err_counter | <db_hostname> |  |  | Reset Replication Error Counter | (1=Error; 0=Normal)",
					"set alert_boot | <on/off> | Set alert switch",
				}
			} 

			helpFormat(args[0], help_set)

		default:
			default_help()

		}

	} else if len(args) == 2 {

		if args[0] == "show" {

			help_showlog := []string{}

			if service_name == "" {
				fmt.Printf("Must to be use a service. Before use command '%s' .\n", ("show " + args[1]))
			} else {

				var str string

				switch args[1] {
				case "mhalog":

					str = "show mhalog | <Node name> [<logid>[,<logid>]...] | Show mha log."
					help_showlog = append(help_showlog, str)
					str = "show mhalog | <Node name> [<logid>[ <logid>]...] | Show mha log."
					help_showlog = append(help_showlog, str)

				case "monitorlog":

					str = "show monitorlog | <Node name> [<logid>[,<logid>]...] | Show monitor log."
					help_showlog = append(help_showlog, str)
					str = "show monitorlog | <Node name> [<logid>[ <logid>]...] | Show monitor log."
					help_showlog = append(help_showlog, str)

				}

				helpFormat((args[0] + " " + args[1]), help_showlog)

			}
		} else {
			default_help()
		}

	}
}

func default_help() {

	help_help := []string{
		"help show  | display show usage.",
		"help use  | display use usage.",
		"help set  | display set usage.",
		"help purge  | purge log.",
		"help help  | display help usage.",
	}

	helpFormat("help", help_help)

}

func helpFormat(_title string, help_content []string) {

	fmt.Println("")
	fmt.Printf("List of all '%s' commands:\n", _title)

	result := columnize.SimpleFormat(help_content)

	fmt.Println(result)

	fmt.Println("")
}
