package cliconfig

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	consulapi "github.com/hashicorp/consul/api"
	"net"
)

func GetUserName() string {
	rt := beego.AppConfig.String("cmha-to-tool-user")
	return rt
}

func GetPassword() string {
	rt := beego.AppConfig.String("cmha-to-tool-password")
	return rt
}

func Consul_Client_Init() (*consulapi.Client, error) {


	ip := beego.AppConfig.String("cmha-server-ip")
	port := beego.AppConfig.String("cmha-server-port")
	datacenter := beego.AppConfig.String("cmha-datacenter")
	token := beego.AppConfig.String("cmha-token")

	if ip == "" {
		errMsg := "Can not link to Consul Server."
                return nil, errors.New(errMsg)
	}

	config := &consulapi.Config{
		Datacenter: datacenter,
		Token:      token,
		Address:    net.JoinHostPort(ip, port),
	}

	client, err := consulapi.NewClient(config)
	if err != nil {
		errMsg := fmt.Sprintf("Create consul-api client failure!%s", err)
		return nil, errors.New(errMsg)
	}

	return client, nil

}

func isActiveLink(ip string, port string) (bool, error) {

	ip_with_port := net.JoinHostPort(ip, port)

	conn, err := net.Dial("tcp", ip_with_port)
	if err != nil {
		return false, err
	}

	defer conn.Close()
	return true, nil

}
