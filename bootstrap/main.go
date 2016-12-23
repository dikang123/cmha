package main

import (
	"fmt"
	"os"
	"time"

	"github.com/astaxie/beego"
	consulapi "github.com/hashicorp/consul/api"
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
			fmt.Println("version 1.1.7")
			return
		} else {
			return
		}
	}
	address := ReadCaConf()	
	//Config is used to configure the creation of a client
        config := ReturnConsulConfig(address)
	client, err := consulapi.NewClient(config)
        if err != nil {
         	beego.Error("Create a consul-api client failure!", err)
                return
        }	
	SetConn(client)
}
