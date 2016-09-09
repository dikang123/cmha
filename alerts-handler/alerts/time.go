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

func IsAlerts(last_time, servicename, last_leader, kvValue,last_leader_key string, C_time int64, kv *consulapi.KV,dataTimeStr string,_alert string) error {
	last_timess, _ := strconv.ParseInt(last_time, 10, 64)
	S_time := C_time - last_timess
	S_time_string := strconv.FormatInt(S_time, 10)
	S_time_int, err := strconv.Atoi(S_time_string)
	if err != nil {
		log.Errorf("s_time string to int:",err)
		return err
	}
	if S_time_int > 30 {
		alertid :=time.Now().UnixNano()
		id, alertdate := GetNowTime()
		log.Infof("leader from the switch to:",alertid,alertdate)
		alert_counter_key := _alert + strconv.FormatInt(alertid,10)
		str := "[ERROR]: ["  + servicename + "]  MySQL occurs failover,Leader from the [" + last_leader + "] switch to [" + kvValue + "]"
                alter_counter_value := "[" + alertdate + "]" + "@" + str
                err = PutKv(servicename, alert_counter_key, alter_counter_value, kv)
                if err != nil {
                        log.Errorf("Put %s failed!", alert_counter_key, err)
                        return err
                }
		log.MyLoGGer(id).Println("[ERROR]: [" + servicename + "] MySQL occurs failover,Leader from the [" + last_leader + "] switch to [" + kvValue + "]")

		C_time_now,data:=GetNowTime()
		C_time_string := strconv.FormatInt(C_time_now, 10)
        	last_leader_value := data + "||" + C_time_string + "||" + kvValue
        	err := PutKv(servicename, last_leader_key, last_leader_value, kv)
        	if err != nil {
                	log.Errorf("Put last leader failed!", err)
                	return err
        	}

	}
	return nil
}

func AlertsAndPutKv(servicename,kvValue, dataTimeStr, last_leader_key string, kv *consulapi.KV,_alert string,C_time int64) error {
	alertid :=time.Now().UnixNano()
	id, alertdate := GetNowTime()
	log.Infof("leader switch to:",alertid,alertdate)
	str := "[ERROR]: [" + servicename + "] MySQL occurs failover,Leader switch to [" + kvValue + "]"
        alter_counter_value := "[" + alertdate + "]" + "@" + str
	alert_counter_key := _alert + strconv.FormatInt(alertid,10)
        err := PutKv(servicename, alert_counter_key, alter_counter_value, kv)
        if err != nil {
         	log.Errorf("Put %s failed!", alert_counter_key, err)
                return err
        }
	log.MyLoGGer(id).Println("[ERROR]: [" + servicename + "] MySQL occurs failover,switch to [" + kvValue + "]")

	C_time,data:=GetNowTime()
	C_time_string := strconv.FormatInt(C_time, 10)
	last_leader_Value := data + "||" + C_time_string + "||" + kvValue
	err = PutKv(servicename, last_leader_key, last_leader_Value, kv)
	if err != nil {
		log.Errorf("Put last leader failed!", err)
		return err
	}
	return nil
}

func StringToInt(instr string) (int64, error) {

        _ret, err := strconv.ParseInt(instr, 10, 64)

        if err != nil {
                return 0, err
        }

        return _ret, nil
}
