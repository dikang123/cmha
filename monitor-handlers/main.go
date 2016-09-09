package main

import (
	"fmt"
	"log"
	"os"
)

var logfile *os.File
var logger *log.Logger

func init() {
	var err error
	logfile, err = os.OpenFile("logs/monitor-handlers.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0)
	if err != nil {
		fmt.Printf("%s\r\n", err.Error())
		os.Exit(-1)
	}
}

func main() {
	defer logfile.Close()
	logger = log.New(logfile, "", log.Ldate|log.Ltime|log.Lshortfile)
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "-v" || arg == "--version" {
			fmt.Println("1.1.5-Beta.5")
			return
		} else {
			return
		}
	}
	CheckService()
}
