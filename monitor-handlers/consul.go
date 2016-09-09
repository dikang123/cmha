package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	consulapi "github.com/hashicorp/consul/api"
)

func GetConfig(address string) *consulapi.Config {
	config := &consulapi.Config{
		Datacenter: beego.AppConfig.String("datacenter"),
		Token:      beego.AppConfig.String("token"),
		Address:    address,
	}
	return config
}

func GetClient(config *consulapi.Config, logvalue, consulapi_failed string) (*consulapi.Client, string, error) {
	client, err := consulapi.NewClient(config)
	if err != nil {
		logger.Println("[E] Create consul-api client failed!", err)
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + consulapi_failed + "{{" + fmt.Sprintf("%s", err)
		return nil, logvalue, err
	}
	return client, logvalue, nil
}

func GetAddrAndStatus(healthpair []*consulapi.ServiceEntry) (string, string) {
	var addr string
	var status string
	for index := range healthpair {
		if healthpair[index].Node.Address != "" {
			addr = healthpair[index].Node.Address
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
	}
	return addr, status
}

func PutRepl(servicename, hostname string) consulapi.KVPair {
	var put string
	put = "1"
	kvvalue := []byte(put)
	kvotherhostname := consulapi.KVPair{
		Key:   "cmha/service/" + servicename + "/db/repl_err_counter/" + hostname,
		Value: kvvalue,
	}
	return kvotherhostname
}

func GetIsHealthAndStatus(ip string, healthpair []*consulapi.ServiceEntry) (bool, string) {
	var ishealthy = true
	var status string
	for index := range healthpair {
		if healthpair[index].Node.Address == ip {
			for checkindex := range healthpair[index].Checks {
				if healthpair[index].Checks[checkindex].Status == "critical" {
					ishealthy = false
					status = healthpair[index].Checks[checkindex].Status
					break
				} else if healthpair[index].Checks[checkindex].Status == "warning" {
					status = healthpair[index].Checks[checkindex].Status
					break
				} else {
					status = healthpair[index].Checks[checkindex].Status
					break
				}
			}
		}
	}
	return ishealthy, status

}
