package alerts

import (
	"strings"

	"github.com/astaxie/beego"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/upmio/alerts-handler/log"
)

func IsCsLeader() bool {
	ip := beego.AppConfig.String("ip")
	port := beego.AppConfig.String("port")
	address := ip + ":" + port
	config := &consulapi.Config{
		Datacenter: beego.AppConfig.String("datacenter"),
		Token:      beego.AppConfig.String("token"),
		Address:    address,
	}
	client, err := consulapi.NewClient(config)
	if err != nil {
		log.MyLoGGer().Println("[E] Create consul-api client failed!", err)
		return false
	}
	state := client.Status()
	leader, err := state.Leader()
	if strings.Contains(leader, ip) {
		return true
	} else {
		return false
	}
}
