package alerts

import (
	"strconv"
	"strings"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/upmio/alerts-handler/log"
)

func GetNowTime() (int64, string) {
	C_time := time.Now().Unix()
	timeLayout := "2006-01-02 15:04:05"
	dataTimeStr := time.Unix(C_time, 0).Format(timeLayout)
	return C_time, dataTimeStr
}

func GetLastTime(last_leader_key string, kv *consulapi.KV) (string, error) {
	last_leader, err := GetKv(last_leader_key, kv)
	if err != nil {
		log.Errorf("Get last leader failed!", err)
		return "", err
	}
	last_leader_key_Value := string(last_leader.Value)
	last_leader_key_Value_split := strings.Split(last_leader_key_Value, "||")
	var last_time string
	for i := range last_leader_key_Value_split {
		if i == 1 {
			log.Infof("last_leader_key_Value_split ", last_leader_key_Value_split[i])
			last_time = last_leader_key_Value_split[i]
		}
	}
	return last_time, nil
}

func IsAlerts(last_time, servicename, last_leader, kvValue,last_leader_key string, C_time int64, kv *consulapi.KV) error {
	last_timess, _ := strconv.ParseInt(last_time, 10, 64)
	S_time := C_time - last_timess
	log.Infof("C_time:",C_time)
	log.Infof("last_time:",last_timess)
	log.Infof("S_time:",S_time)
	S_time_string := strconv.FormatInt(S_time, 10)
	S_time_int, err := strconv.Atoi(S_time_string)
	if err != nil {
		log.Errorf("s_time string to int:",err)
		return err
	}
	if S_time_int > 30 {
		log.MyLoGGer().Println("[ERROR]: [" + servicename + "] MySQL occurs failover,Leader from the [" + last_leader + "] switch to [" + kvValue + "]")
		C_time,dataTimeStr:=GetNowTime()
		C_time_string := strconv.FormatInt(C_time, 10)
        	last_leader_value := dataTimeStr + "||" + C_time_string + "||" + kvValue
        	err := PutKv(servicename, last_leader_key, last_leader_value, kv)
        	if err != nil {
                	log.Errorf("Put last leader failed!", err)
                	return err
        	}
	}
	return nil
}

func AlertsAndPutKv(servicename,kvValue, dataTimeStr, last_leader_key string, kv *consulapi.KV) error {
	log.MyLoGGer().Println("[ERROR]: [" + servicename + "] MySQL occurs failover,Leader switch to " + kvValue)
	C_time,dataTimeStr:=GetNowTime()
	C_time_string := strconv.FormatInt(C_time, 10)
	last_leader_Value := dataTimeStr + "||" + C_time_string + "||" + kvValue
	log.Info("last_leader_value:",last_leader_Value)
	err := PutKv(servicename, last_leader_key, last_leader_Value, kv)
	if err != nil {
		log.Errorf("Put last leader failed!", err)
		return err
	}
	return nil
}
