package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

func SetRepl_err_counter(hostname string, client *consulapi.Client,logvalue string) string{
	cnt := 0
	kvotherhostname := PutRepl(servicename, hostname)
try:
	_, err := kv.Put(&kvotherhostname, nil)
	if err != nil {
		logger.Println("[E] Set peer database repl_err_counter to 1 in CS failed!", err)
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + set_counter_failed + "{{" + fmt.Sprintf("%s", err)
		if cnt == 2 {
			count := 0
			for i := 0; i < 3; i++ {
				ct := 0
			tr:
				_, err = kv.Put(&kvotherhostname, nil)
				if err != nil {
					count += 1
					logger.Println("[E] Set peer database repl_err_counter to 1 in CS failed! CS ip = "+consul_agent_ip, err)
					timestamp := time.Now().Unix()
					logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + set_counter_failed + "{{" + consul_agent_ip + "{{" + fmt.Sprintf("%s", err)
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
				UploadLog(logkey, logvalue)
				os.Exit(1)
			}
		} else {
			cnt++
			goto try
		}
	}
	logger.Println("[I] Set peer database repl_err_counter to 1 in CS successfully!")
	timestamp := time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + set_counter_success
	return logvalue
}
