package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"

	consulapi "github.com/hashicorp/consul/api"
)

func UploadLog(logkey, logvalue string) {
	config := &consulapi.Config{
		Datacenter: beego.AppConfig.String("datacenter"),
		Token:      beego.AppConfig.String("token"),
		Address:    address,
	}
	kvhostname := consulapi.KVPair{
		Key:   logkey,
		Value: []byte(logvalue),
	}
	_, err := kv.Put(&kvhostname, nil)
	if err != nil {
		logger.Println("[E] Upload log to CS failed!", err)
		count := 0
		for i := 0; i < 3; i++ {
			client, err = consulapi.NewClient(config)
			if err != nil {
				count += 1
				logger.Println("[E] Create consul-api client failed!", err)
				timestamp := time.Now().Unix()
				logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + consulapi_failed + "{{" + fmt.Sprintf("%s", err)
				continue
			}
			logger.Println("[I] Create consul-api client successfully!")
			timestamp := time.Now().Unix()
			logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + consulapi_success
			_, err = kv.Put(&kvhostname, nil)
			if err != nil {
				count += 1
				logger.Println("[E] Upload log to CS failed!", err)
				continue
			}
			break
		}
		if count >= 3 {
			logger.Println("[I] Monitor Handler Completed")
			return
		}
	}
	logger.Println("[I] Upload log to CS successfully!")
	logger.Println("[I] Monitor Handler Completed")
}
