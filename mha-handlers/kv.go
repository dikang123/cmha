package main

import (
	"strconv"
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

func GetLeaderKey(servicename string) string {
	leader = "cmha/service/" + servicename + "/db/leader"
	return leader
}

func GetAcquireJson(format, hostname, ip, port, username, password string) string {
	acquirejson := ""
	if format == "json" {
		acquirejson = `{"Node":"` + hostname + `","Ip":"` + ip + `","Port":` + port + `,"Username":"` + username + `","Password":"` + password + `"}`
	} else if format == "hap" {
		acquirejson = "server" + " " + hostname + " " + ip + ":" + port
	} else {
		logger.Println("[E] format error,json or hap!")
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + format_json_hap_failed
		logger.Println("[I] Give up leader election")
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + give_election
		UploadLog(logkey, logvalue)
		//return
	}
	return acquirejson

}

func IsStatus(healthpair []*consulapi.ServiceEntry, ip string) (string, int64) {
	var ishealthy = true
	for index := range healthpair {
		for checkindex := range healthpair[index].Checks {
			if healthpair[index].Checks[checkindex].Node == hostname {
				if healthpair[index].Checks[checkindex].Status == "critical" {
					ishealthy = false
					break
				}
			}
		}
	}
	if !ishealthy {
		logger.Println("[E] Status is critical!")
		logvalue, timestamp := AddLogValue(status_critical, logvalue)
		logger.Println("[I] Give up leader election")
		logvalue, timestamp = AddLogValue(give_election, logvalue)
		UploadLog(logkey, logvalue)
		return logvalue, timestamp
	} else {
		logger.Println("[I] Status is not critical")
		logvalue, timestamp := AddLogValue(status_nocritical, logvalue)
		slave(ip, port, username, password)
		return logvalue, timestamp
	}
}
