package main

import (
	"fmt"
//	"github.com/astaxie/beego"
	"os"
//	"time"
	"log"
)

var logfile *os.File
var logger *log.Logger
func init() {
	var err error
	logfile, err = os.OpenFile("logs/mha-handlers.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0)
	if err != nil {
		fmt.Printf("%s\r\n", err.Error())
		os.Exit(-1)
	}
//	beego.SetLogger("file", `{"filename":"logs/mha-handlers.log"}`)
//	beego.SetLogFuncCall(true)
}

func main() {
	defer logfile.Close()	
	logger = log.New(logfile, "", log.Ldate|log.Ltime|log.Lshortfile)
//	defer beego.BeeLogger.Close()
//	defer time.Sleep(100 * time.Millisecond)
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "-v" || arg == "--version" {
			fmt.Println("version 1.1.3")
			return
		} else {
			return
		}
	}
	SessionAndChecks()
}
