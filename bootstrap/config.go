package main

import(
	"github.com/astaxie/beego"
	consulapi "github.com/hashicorp/consul/api"	
)

func ReadHostConf()(string,string){
	hostname := beego.AppConfig.String("hostname")
        other_hostname := beego.AppConfig.String("otherhostname")
	return hostname,other_hostname
}

func ReadDatabaseConf()(string,string,string,string){
	ip := beego.AppConfig.String("ip")
        port := beego.AppConfig.String("port")
        username := beego.AppConfig.String("username")
        password := beego.AppConfig.String("password")
	return ip,port,username,password
}

func ReadServiceConf()string{
	servicename := beego.AppConfig.String("servicename")
	return servicename
}

func ReadCaConf()string {
	consul_agent_ip := beego.AppConfig.String("consul_agent_ip")
        consul_agent_port := beego.AppConfig.String("consul_agent_port")
        address = consul_agent_ip + ":" + consul_agent_port
	return address
}

func ReturnLeaderAndLastLeader()(string,string){
	leader := "cmha/service/" + servicename + "/db/leader"
        last_leader :="cmha/service/" + servicename + "/db/last_leader"
	return leader,last_leader
}

func ReturnConsulConfig(address string)*consulapi.Config{
	 config := &consulapi.Config{
                Datacenter: beego.AppConfig.String("datacenter"),
                Token:      beego.AppConfig.String("token"),
                Address:    address,
        }
	return config
}

func ReadFormat()string{
	format := beego.AppConfig.String("format")
	return format
}
