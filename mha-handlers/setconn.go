package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	consulapi "github.com/hashicorp/consul/api"
)

func SetConn(ip, port, username, password string) {
	var sessionvalue string
	var timestamp int64
	c := 0
	for i := 0; i < 3; i++ {
		session := client.Session()
		sessionEntry := consulapi.SessionEntry{
			LockDelay: 10 * time.Second,
			Name:      servicename,
			Node:      hostname,
			Checks:    []string{"serfHealth", "service:" + servicename},
		}
		//Create makes a new session. Providing a session entry can customize the session. It can also be nil to use defaults.
		sessionvalue, _, err = session.Create(&sessionEntry, nil)
		if err != nil {
			c += 1
			logger.Println("[E] Session create failed!", err)
			timestamp = time.Now().Unix()
			logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + create_session_failed + "{{" + fmt.Sprintf("%s", err)
			continue
		}
		break
	}
	if c >= 3 {
		logger.Println("[I] Give up leader election")
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + give_election
		UploadLog(logkey, logvalue)
		return
	}
	//NewClient returns a new client
	logger.Println("[I] Session create successfully!")
	timestamp = time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + create_session_success
	format := beego.AppConfig.String("format")
	acquirejson := GetAcquireJson(format, hostname, ip, port, username, password)
	if acquirejson == "" {
		return
	}
	value := []byte(acquirejson)
	kv = client.KV()
	kvpair := consulapi.KVPair{
		Key:     leader,
		Value:   value,
		Session: sessionvalue,
	}
	//Acquire is used for a lock acquisiiton operation. The Key, Flags, Value and Session are respected. Returns true on success or false on failures.
	time.Sleep(15 * time.Second)
	var ok bool
	ok, _, err = kv.Acquire(&kvpair, nil)
	if err != nil {
		logger.Println("[E] Send service leader request to CS failed!", err)
		timestamp = time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + send_leader_failed + "{{" + fmt.Sprintf("%s", err)
		UploadLog(logkey, logvalue)
		return
		count := 0
		for i := 0; i < 3; i++ {
			ok, _, err = kv.Acquire(&kvpair, nil)
			if err != nil {
				count += 1
				logger.Println("[E] Send service leader request to CS failed!", err)
				timestamp = time.Now().Unix()
				logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + send_leader_failed + "{{" + fmt.Sprintf("%s", err)
				continue
			}
			break
		}
		if count >= 3 {
			UploadLog(logkey, logvalue)
			return
		}
	}
	logger.Println("[I] Send service leader request to CS successfully!")
	timestamp = time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + send_leader_success
	if !ok {
		logger.Println("[E] Becoming service leader failed! Connection string is " + ip + " " + port)
		timestamp = time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + becoming_string_failed + "{{" + ip + "{{" + port
		SetRead_only(username, password, port, 1)
		UploadLog(logkey, logvalue)
		return
	} else {
		logger.Println("[I] Becoming service leader successfully! Connection string is " + ip + " " + port)
		timestamp = time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + becoming_string_success + "{{" + ip + "{{" + port
		other_hostname := beego.AppConfig.String("otherhostname")
		SetRepl_err_counter(other_hostname)
		UploadLog(logkey, logvalue)
		return
	}
}
