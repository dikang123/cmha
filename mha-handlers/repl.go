package main

import (
	"fmt"
	"strconv"
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

func SetRepl_err_counter(hostname string) {
	cnt := 0
	var put string
	var timestamp int64
	put = "1"
	kvvalue := []byte(put)
	kvotherhostname := consulapi.KVPair{
		Key:   "cmha/service/" + servicename + "/db/repl_err_counter/" + hostname,
		Value: kvvalue,
	}
try:
	_, err := kv.Put(&kvotherhostname, nil)
	if err != nil {
		logger.Println("[E] Set peer database repl_err_counter to 1 in CS failed!", err)
		timestamp = time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + set_counter_failed + "{{" + fmt.Sprintf("%s", err)
		UploadLog(logkey, logvalue)
		return
		if cnt == 2 {
			count := 0
			for i := 0; i < 3; i++ {
				ct := 0
			tr:
				_, err := kv.Put(&kvotherhostname, nil)
				if err != nil {
					count += 1
					logger.Println("[E] Set peer database repl_err_counter to 1 in CS failed!", err)
					timestamp = time.Now().Unix()
					logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + set_counter_failed + "{{" + fmt.Sprintf("%s", err)
					if ct == 2 {
						continue
					} else {
						ct++
						goto tr
					}
				}
				break
			}
			if count >= 3 {
				return
			}
		} else {
			cnt++
			goto try
		}
	}
	logger.Println("[I] Set peer database repl_err_counter to 1 in CS successfully!")
	timestamp = time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + set_counter_success
}
