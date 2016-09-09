package main

import (
	"github.com/astaxie/beego"
)

func ReadCaAddress() string {
	consul_agent_ip := beego.AppConfig.String("consul_agent_ip")
	consul_agent_port := beego.AppConfig.String("consul_agent_port")
	address := consul_agent_ip + ":" + consul_agent_port
	return address
}

func ReadServiceName() string {
	servicename = beego.AppConfig.String("servicename")
	return servicename
}

func ReadHostName() (string, string) {
	hostname = beego.AppConfig.String("hostname")
	other_hostname = beego.AppConfig.String("otherhostname")
	return hostname, other_hostname
}

func ReadTag() string {
	tag = beego.AppConfig.String("tag")
	return tag
}

func ReadDatabaseCfg() (string, string, string, string) {
	ip := beego.AppConfig.String("ip")
	port := beego.AppConfig.String("port")
	username := beego.AppConfig.String("username")
	password := beego.AppConfig.String("password")
	return ip, port, username, password
}
