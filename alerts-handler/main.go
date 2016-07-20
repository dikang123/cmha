package main

import (
	"flag"
	"fmt"
	//	"log"
	"os"

	"github.com/upmio/cmha/alerts-handler/alerts"
	"github.com/upmio/cmha/alerts-handler/log"
)

var logfile *os.File

//var logger *log.Logger

/*func init() {
	var err error
	logfile, err = os.OpenFile("logs/alerts-handler.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm|os.ModeTemporary)
	if err != nil {
		fmt.Printf("%s\r\n", err.Error())
		os.Exit(-1)
	}
}*/

var (
	ServiceFlag = flag.String("service", "", " usage: the horus-server manager ip")
	KvFlag      = flag.String("kv", "", "usage: the horus-server manager port")
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
			return
		}
		defer logfd.Close()
		if err := log.LogInit(*Loglvl, logfd); err != nil {
			fmt.Println("log init failed", err)
			return
		}
	}
	if *ServiceFlag != "" && *KvFlag == "" {
		servicename := *ServiceFlag
		alerts.IsDbOrChap(servicename)
	} else if *ServiceFlag != "" && *KvFlag != "" {
		log.Info("kv")
	} else {
		log.Info("error")
	}
}
