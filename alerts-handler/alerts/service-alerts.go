package alerts

import (
	"fmt"

	"github.com/astaxie/beego"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/upmio/cmha/alerts-handler/log"
)

func IsDbOrChap(servicename string) {
	fmt.Println("servicename", servicename)
	ip := beego.AppConfig.String("ip")
	port := beego.AppConfig.String("port")
	address := ip + ":" + port
	config := &consulapi.Config{
		Datacenter: beego.AppConfig.String("datacenter"),
		Token:      beego.AppConfig.String("token"),
		Address:    address,
	}
	var healthpair []*consulapi.ServiceEntry
	var health *consulapi.Health
	var err error
	client, err := consulapi.NewClient(config)
	if err != nil {
		log.Errorf("[E] Create consul-api client failed!", err)
		return
	}
	health = client.Health()
	healthpair, _, err = health.Service(servicename, "", false, nil)
	if err != nil {
		return
	}
	var addr, types, hostname string
	var status string
	for index := range healthpair {
		if len(healthpair[index].Service.Tags) == 1 {
			for tagindex := range healthpair[index].Service.Tags {
				if healthpair[index].Service.Tags[tagindex] == "master" && healthpair[index].Service.Tags[tagindex] == "slave" {
					types = "db"
				} else {
					types = "chap"
				}
			}
		}
		if healthpair[index].Node.Address != "" {
			addr = healthpair[index].Node.Address
		}
		if healthpair[index].Node.Node != "" {
			hostname = healthpair[index].Node.Node
		}
		for checkindex := range healthpair[index].Checks {
			if healthpair[index].Checks[checkindex].Status == "passing" {
				status = healthpair[index].Checks[checkindex].Status
			} else if healthpair[index].Checks[checkindex].Status == "warning" {
				status = healthpair[index].Checks[checkindex].Status
				break
			} else if healthpair[index].Checks[checkindex].Status == "critical" {
				status = healthpair[index].Checks[checkindex].Status
				break
			} else {
				status = "invalid"
				break
			}
		}
		if status == "warning" && types == "db" {
			log.Warnf("%s system of %s node %s %s MySQL replication critical", servicename, types, hostname, addr)
			continue
		} else if status == "critical" && types == "db" {
			log.Warnf("%s system of %s node %s %s status critical", servicename, types, hostname, addr)
			continue
		} else if status == "critical" && types == "chap" {
			log.Warnf("%s system of %s node %s %s status critical", servicename, types, hostname, addr)
		}
	}
}
