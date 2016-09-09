package main

import (
	consulapi "github.com/hashicorp/consul/api"
)

func UploadLog(logkey, logvalue string) {
	kvhostname := consulapi.KVPair{
		Key:   logkey,
		Value: []byte(logvalue),
	}
	_, err = kv.Put(&kvhostname, nil)
	if err != nil {
		logger.Println("[E] Upload log to CS failed!", err)
		count := 0
		for i := 0; i < 3; i++ {
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
	logger.Println("[I] MHA Handler Completed")
}
