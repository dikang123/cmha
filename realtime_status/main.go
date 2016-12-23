package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/upmio/realtime_status/file"
	"github.com/upmio/realtime_status/info"
	"github.com/upmio/realtime_status/log"
	"github.com/upmio/realtime_status/time"
)

var (
	ServiceNameFlag = flag.String("servicename", "", " usage: the a group service name")
	SysInfoFlag     = flag.String("sysinfo", "", "usage: the a sys info")
	DbInfoFlag      = flag.String("dbinfo", "", "usage: the db info")
	Debug           = flag.Bool("debug", false, "debug model ,write log to std.err")
)
var Version = "version 1.1.7"

var seqpath = "/tmp/.realtime_cache/real_status_seq"
var realtime_cachedirpath = "/tmp/.realtime_cache/"

func main() {
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "-v" || arg == "--version" {
			fmt.Println(Version)
			return
		}
	}
	flag.Parse()
	if *Debug {
		log.ConfigLevel("debug")
	}
	realtimedir, err := file.PathExists(realtime_cachedirpath)
	if err != nil {
		fmt.Println("check dir exists failed!", err)
		os.Exit(2)
	}
	if !realtimedir {
		err := os.MkdirAll(realtime_cachedirpath, 0777)
		if err != nil {
			fmt.Printf("dir create filed!", err)
			os.Exit(2)
		}
	}
	_, C_time_data := time.GetNowTime()
	save_counter, _ := beego.AppConfig.Int("save_counter")
	isexist, _ := file.PathExists(seqpath)
	var newid int
	var old_id int
	if isexist {
		old_id_byte, _ := ioutil.ReadFile(seqpath)
		old_id_string := string(old_id_byte)
		old_id, _ = strconv.Atoi(old_id_string)
		if old_id == save_counter {
			newid = 1
		} else {
			newid = old_id + 1
		}
	} else {
		newid = 1
	}
	old_id = newid
	user := beego.AppConfig.String("mysql_user")
	password := beego.AppConfig.String("mysql_password")
	port := beego.AppConfig.String("mysql_port")
	hostname, err := os.Hostname()
	if err != nil {
		log.Errorf("get hostname failed!", err)
		os.Exit(2)
	}
	if *ServiceNameFlag != "" && *SysInfoFlag != "" && *DbInfoFlag != "" {
		if *SysInfoFlag == "sys_info" && *DbInfoFlag == "db_info" {
			info.InitFirstSysinfo()
			info.InitFirstDbInfo(user, password, port)
			info.CpuLoad()
			info.GetSysInfo()
			info.GetMysql(user, password, port)
			info.InsertDataToCs(*ServiceNameFlag, hostname, newid, old_id, C_time_data, "sys_and_db")
		} else {
			os.Exit(0)
		}
	} else if *ServiceNameFlag != "" && *SysInfoFlag != "" && *DbInfoFlag == "" {
		if *SysInfoFlag == "sys_info" {
			info.InitFirstSysinfo()
			info.CpuLoad()
			info.GetSysInfo()
			info.InsertDataToCs(*ServiceNameFlag, hostname, newid, old_id, C_time_data, "sys")
		}
	} else {
		os.Exit(0)
	}
}
