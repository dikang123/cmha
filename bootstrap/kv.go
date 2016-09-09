package main

import (
	"time"
	"errors"
	"strconv"
	"github.com/astaxie/beego"
	consulapi "github.com/hashicorp/consul/api"
)

func GetRepl(servicename,host string,kvvalue []byte)(consulapi.KVPair,string){
	key := "cmha/service/" + servicename + "/db/repl_err_counter/" + host
	kvhost := consulapi.KVPair{
                Key:   key,
                Value: kvvalue,
        }
	return kvhost,key
}

func PutValue()[]byte{
	var put string
        put = "0"
        kvvalue := []byte(put)
	return kvvalue
}
func Put(kvpair *consulapi.KVPair,key string,kv *consulapi.KV)error{
	_, err := kv.Put(kvpair, nil)
        if err != nil {
                beego.Error(key + " put failure!", err)
                return err
        }
	return nil
}

func SessionAcquireLeader(format,hostname,ip,port,username,password,leader,sessionvalue string,kv *consulapi.KV)(bool,error){
	var acquirejson string
        if format == "json" {
        	acquirejson = `{"Node":"` + hostname + `","Ip":"` + ip + `","Port":` + port + `,"Username":"` + username + `","Password":"` + password + `"}`
        } else if format == "hap" {
        	acquirejson = "server" + " " + hostname + " " + ip + ":" + port
        } else {
               	beego.Error("format error!")
                return false,errors.New("format error!")
        }
        value := []byte(acquirejson)
        kvpair := consulapi.KVPair{
        	Key:     leader,
                Value:   value,
                Session: sessionvalue,
        }
        //Acquire is used for a lock acquisiiton operation. The Key, Flags, Value and Session are respected. Returns true on success or false on failures.
        ok, _, err := kv.Acquire(&kvpair, nil)
        if err != nil {
       		beego.Error("Set the connection string master failure!", err)
                return false,err
        }
	return ok,nil
}

func PutLastLeader(leadervalue,last_leader string,kv *consulapi.KV)error{
	C_time,dataTimeStr :=GetNowTime()
	C_time_string := strconv.FormatInt(C_time, 10)
	last_leader_value := dataTimeStr + "||" + C_time_string + "||" +leadervalue
	kvlastleader := consulapi.KVPair{
		Key:   last_leader,
		Value: []byte(last_leader_value),
	}
	err := Put(&kvlastleader,last_leader,kv)
	if err != nil {
		beego.Error("Put failed:",last_leader,err)
		return err
	}
	return nil
		
}

func GetNowTime() (int64, string) {
        C_time := time.Now().Unix()
        timeLayout := "2006-01-02 15:04:05"
        dataTimeStr := time.Unix(C_time, 0).Format(timeLayout)
        return C_time, dataTimeStr
}
