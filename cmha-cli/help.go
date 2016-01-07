package main

import (
	"fmt"
)

func help(args ...string) {
	if len(args) < 1 {
		fmt.Println("Usage:")
		fmt.Println("\b\b\tCommand\tDescription\t")
		fmt.Println("*\t\b\bshow\tDisplay Current Status\t")
		fmt.Println("\tshow cluster\t\t\t\tDisplay cluster status")
		fmt.Println("\tshow chap [service_name]\t\t\t\tDisplay proxy status(e.g.show chap;show chap innosql)")
		fmt.Println("\tshow db [service_name]\t\t\tDisplay database group and database status(e.g.show db;show db innosql)")
		fmt.Println("\tshow dbleader [service_name]\t\t\t\tDisplay current database leader(e.g.show dbleader;show dbleader innosql)")
		fmt.Println("\tshow cmhaleader\t\t\t\tDisplay current CMHA cluster leader")
		fmt.Println("\tshow repl_err_counter <db_hostname>\tDisplay replication err counter(e.g.show repl_err_counter repl64-mysql1)")
		fmt.Println("*\t\b\bset\tSet database variables")
		fmt.Println("\tset vsr <db_host_ipaddr> <db_port> on\t\tEnable database VSR Function(e.g.set vsr 192.168.2.1 3306 on)")
                fmt.Println("\tset read_only <db_host_ipaddr> <db_port> on\tSet database into Read_Only mode(e.g.set read_only 192.168.2.1 3306 on)")
                fmt.Println("\tset slave_start <db_host_ipaddr> <db_port>\tStart database slave threads(e.g.set slave_start 192.168.2.1 3306)")
                fmt.Println("\tset slave_stop <db_host_ipaddr> <db_port>\tStop database slave threads(e.g.set slave_stop 192.168.2.1 3306)")
                fmt.Println("\tset repl_err_counter <db_hostname>\t\tReset Replication Error Counter(1=Error;0=Normal)")
	} else if len(args) == 1 {
		switch args[0] {
		case "show":
			fmt.Println("Usage:")
			fmt.Println("\b\b\tCommand\tDescription\t")
			fmt.Println("*\t\b\bshow\tDisplay Current Status\t")
			fmt.Println("\tshow cluster\t\t\t\tDisplay cluster status")
			fmt.Println("\tshow chap\t\t\t\tDisplay proxy status")
			fmt.Println("\tshow db [service_name]\t\t\tDisplay database group and database status(e.g.show db;show db innosql)")
			fmt.Println("\tshow dbleader\t\t\t\tDisplay current database leader")
			fmt.Println("\tshow cmhaleader\t\t\t\tDisplay current CMHA cluster leader")
			fmt.Println("\tshow repl_err_counter <db_hostname>\tDisplay replication err counter(e.g.show repl_err_counter repl64-mysql1)")
		case "set":
			fmt.Println("Usage:")
			fmt.Println("\b\b\tCommand\tDescription\t")
			fmt.Println("*\t\b\bset\t\tSet database variables")
			fmt.Println("\tset vsr <db_host_ipaddr> <db_port> on\t\tEnable database VSR Function(e.g.set vsr 192.168.2.1 3306 on)")
			fmt.Println("\tset read_only <db_host_ipaddr> <db_port> on\tSet database into Read_Only mode(e.g.set read_only 192.168.2.1 3306 on)")
			fmt.Println("\tset slave_start <db_host_ipaddr> <db_port>\tStart database slave threads(e.g.set slave_start 192.168.2.1 3306)")
			fmt.Println("\tset slave_stop <db_host_ipaddr> <db_port>\tStop database slave threads(e.g.set slave_stop 192.168.2.1 3306)")
			fmt.Println("\tset repl_err_counter <db_hostname>\t\tReset Replication Error Counter(1=Error;0=Normal)")
		}

	}
}
