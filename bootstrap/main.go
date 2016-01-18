package main

import (
	"fmt"
	"os"
	"time"

	"github.com/astaxie/beego"
)

func init() {
	beego.SetLogger("file", `{"filename":"logs/bootstrap.log"}`)
	beego.SetLogFuncCall(true)
}

func main() {
	defer beego.BeeLogger.Close()
	defer time.Sleep(100 * time.Millisecond)
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "-v" || arg == "--version" {
			fmt.Println("version 1.1.4")
			return
		} else {
			return
		}
	}
	SetConn()
}
