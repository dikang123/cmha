package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/astaxie/beego"
	"github.com/upmio/alerts-handler/alerts"
	"github.com/upmio/alerts-handler/log"
)

var (
	ServiceFlag = flag.String("service", "", " usage: the a group service name")
	KvFlag      = flag.String("kv", "", "usage: the alerts type leader or monitor")
	Loglvl      = flag.String("loglevel", "info", " the log level")
	Logfile     = flag.String("logfile", "./logs/alerts-handler.log", "the log file ")
	Debug       = flag.Bool("debug", false, "debug model ,write log to std.err")
)

func main() {

	flag.Parse()
	if *Debug {
		log.ConfigLevel("debug")
	} else {
		logfd, err := os.OpenFile(*Logfile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("open logfile failed", err)
			os.Exit(-1)
		}
		defer logfd.Close()
		if err := log.LogInit(*Loglvl, logfd); err != nil {
			fmt.Println("log init failed", err)
			os.Exit(-1)
		}
	}
	alert_boot := beego.AppConfig.String("alert_boot")

	if alert_boot != "enable" {
		os.Exit(0)
	}
	time.Sleep(5 * time.Second)

	logfile_path := beego.AppConfig.String("logfile_path")
	isleader := alerts.IsCsLeader()
	if !isleader {
		os.Exit(0)
	}
	log.Infof("[INFO]: %s alerts Handler Triggered", *ServiceFlag)
/*	logfile, err := log.AlertFile()
	if err != nil {
		log.Errorf("%s\r\n", err.Error())
		os.Exit(-1)
	}*/

	logfile, err := os.OpenFile(logfile_path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0)
	if err != nil {
		log.Errorf("%s\r\n", err.Error())
		os.Exit(-1)
	}
	defer logfile.Close()
	log.MyLogInit(logfile)

	if *ServiceFlag != "" && *KvFlag == "" {
		servicename := *ServiceFlag
		if servicename == "CS" {
			alerts.CsAlerts(servicename)
		} else {
			alerts.IsDbOrChap(servicename)
		}
	} else if *ServiceFlag != "" && *KvFlag != "" {
		servicename := *ServiceFlag
		kvname := *KvFlag
		alerts.IsLerderorMonitorOrRole(servicename, kvname)
	} else {
		log.Info("error")
	}
	log.Info("[INFO]: %s alerts Handler Completed", *ServiceFlag)
}
