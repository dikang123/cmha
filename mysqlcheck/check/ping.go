package check

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func IsPingType(user, password string) {
	host, port, checktime_string, timeout, defaultDb, ping_type := GetConfig()
	if ping_type == "select,replication" || ping_type == "select" {
		TrySelectCheckTime(user, password, host, port, defaultDb, checktime_string, ping_type, timeout)
	} else if ping_type == "update,replication" || ping_type == "update" {
		TryUpdateCheckTime(user, password, host, port, defaultDb, checktime_string, ping_type, timeout)
	}else{
		fmt.Print("Configuration error")
		os.Exit(2)
	}

}

func TrySelectCheckTime(user, password, host, port, defaultDb, checktime_string, ping_type, timeout string) {
	checktime, _ := strconv.Atoi(checktime_string)
	for {
		if checktime == 0{
			break
		}else{
			checktime--
			MYSQL_OK := SelectCheckMysqlHealth(user, password, host, port, defaultDb, timeout,checktime)
			if MYSQL_OK == 0 {
				if strings.Contains(ping_type, "replication") {
					isyes, err := ShowSlave(user, password, host, port, defaultDb, timeout)
					if err != nil {
						fmt.Print(err)
						os.Exit(2)
					}
					fmt.Println("isyes:",isyes)
					if isyes == "Yes" {
						fmt.Print("check ok")
						os.Exit(0)
					} else if isyes == "noreplication"{
						fmt.Print("replication is not configured")
						os.Exit(1)
					}else {
						fmt.Print("check replication io_thread fail:", isyes)
						os.Exit(1)
					}
				} else {
					fmt.Print("check ok")
					os.Exit(0)
				}
			}
			if MYSQL_OK == 1 && checktime == 0 {
				os.Exit(2)
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func TryUpdateCheckTime(user, password, host, port, defaultDb, checktime_string, ping_type, timeout string) {
	checktime, _ := strconv.Atoi(checktime_string)
	for {
		if checktime == 0{
			break
		}else{
			checktime--
			MYSQL_OK := CheckMysqlHealth(user, password, host, port, defaultDb, timeout,checktime)
			if MYSQL_OK == 0 {
				if strings.Contains(ping_type, "replication") {
					isyes, err := ShowSlave(user, password, host, port, defaultDb, timeout)
					if err != nil {
						fmt.Print(err)
						os.Exit(2)
					}
					if isyes == "Yes" {
						fmt.Print("check ok")
						os.Exit(0)
					}else if isyes == "noreplication"{
                                                fmt.Print("replication is not configured")
                                                os.Exit(1)
                                        }else {
						fmt.Print("check replication io_thread fail:", isyes)
						os.Exit(1)
					}
				} else {
					fmt.Print("check ok")
					os.Exit(0)
				}
			}
			if MYSQL_OK == 1 && checktime == 0 {
				os.Exit(2)
			}
		}
		time.Sleep(1 * time.Second)
	}
}
