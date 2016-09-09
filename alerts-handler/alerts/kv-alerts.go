package alerts

import (
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/upmio/alerts-handler/log"
)

var leader = "leader"
var monitor = "monitor"
var role = "role"
var alert_counter_key string
func IsLerderorMonitorOrRole(servicename, kvname string) {
	if kvname == leader {
		err := IsLeader(servicename, kvname)
		if err != nil {
			log.Errorf("leader kv alerts failed!", err)
			return
		}
	} else if kvname == monitor {
		log.Info("monitor")
		err := IsMonitor(servicename, kvname)
		if err != nil {
			log.Errorf("monitor kv alerts failed!", err)
			return
		}
	}
}

func IsLeader(servicename, kvname string) error {
	time.Sleep(10 * time.Second)
	C_time := time.Now().Unix()
        var timeLayout = "2006-01-02 15:04:05"
        var dataTimeStr = time.Unix(C_time, 0).Format(timeLayout)
	_, _, adress := GetConf()
	config := GetConfig(adress)
	var err error
	client, err := GetClient(config)
	if err != nil {
		log.Error(err)
		return err
	}
	kv := ReturnKv(client)
	leader_key := GetLeaderKey(servicename)
	last_leader_key :=GetLastLeaterKey(servicename)
	_alert := "cmha/service/" + servicename +"/alerts/alerts_counter/"
	kvValue, last_leaderValue, _ := ConsulGetAndPutKv(leader_key,last_leader_key, servicename, kv)
	var last_leader string
	if last_leaderValue !=""{
		last_leader =GetLastLeader(last_leaderValue)
	}
	if kvValue == "" {
		alertid :=time.Now().UnixNano()
		id,alertdate := GetNowTime()
		log.Infof("leader does not exist:",alertid,alertdate)	
		alert_counter_key := _alert + strconv.FormatInt(alertid,10)
		str := "[ERROR]: ["  + servicename + "] MySQL occurs failover,The current leader does not exist"
	        alter_counter_value := "[" + alertdate + "]" + "@" + str
                err = PutKv(servicename, alert_counter_key, alter_counter_value, kv)
                if err != nil {
                	log.Errorf("Put %s failed!", alert_counter_key, err)
                        return err
                }
		log.MyLoGGer(id).Println("[ERROR]: [" + servicename + "] MySQL occurs failover,The current leader does not exist")
		
	} else if last_leader == "" {
		err = AlertsAndPutKv(servicename,kvValue, dataTimeStr, last_leader_key, kv,_alert,C_time)
		if err != nil {
			log.Error(err)
			return err
		}
	} else {
		last_time, err := GetLastTime(last_leader_key, kv)
		if err != nil {
			log.Error(err)
			return err
		}
		err = IsAlerts(last_time, servicename, last_leader, kvValue,last_leader_key,C_time, kv,dataTimeStr,_alert)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

func IsMonitor(servicename, kvname string) error {
	C_time := time.Now().Unix()
	var timeLayout = "2006-01-02 15:04:05"
	_ = time.Unix(C_time, 0).Format(timeLayout)
	var err error
	ip := beego.AppConfig.String("ip")
	frequency := beego.AppConfig.String("frequency")
	frequency_int, err := strconv.Atoi(frequency)
	if err != nil {
		log.Errorf("frequency string to int error: %s", err)
		return err
	}
	port := beego.AppConfig.String("port")
	adress := ip + ":" + port
	config := &consulapi.Config{
		Datacenter: beego.AppConfig.String("datacenter"),
		Token:      beego.AppConfig.String("token"),
		Address:    adress,
	}
	var kvPair *consulapi.KVPair
	client, err := consulapi.NewClient(config)
	if err != nil {
		log.Errorf("[E] Create consul-api client failed!", err)
		return err
	}
	kv := client.KV()
	var hosts []string
	_repl_err_counter :="cmha/service/" + servicename + "/db/repl_err_counter/"
	logs, _, _err := kv.List(_repl_err_counter, nil)
        if _err != nil {
		log.Errorf("query %s list failed!",_repl_err_counter,err)
                return err
        }
        if logs != nil {
		for _, v := range logs {
			key_items := strings.Split(v.Key, "/")
                        key_hostname := key_items[len(key_items)-1]
			hosts = append(hosts,key_hostname)
		}
	}	
	
	for host := range hosts {
		var addr string
		repl_err_counter := "cmha/service/" + servicename + "/db/repl_err_counter/" + hosts[host]
		kvPair, err = GetKv(repl_err_counter, kv)
		if err != nil {
			log.Errorf("Get repl_err_counter failed!", err)
			continue
		}
		kvValue := string(kvPair.Value)
		health := client.Health()
		healthpair, _, err := health.Service(servicename, "", false, nil)
		if err != nil {
			log.Errorf("err", err)
			continue
		}
		for index := range healthpair {
			if healthpair[index].Node.Node != hosts[host] {
				continue
			}
			addr = healthpair[index].Node.Address
			_keys := "cmha/service/" + servicename +"/alerts/alerts_counter/"
			_alerts_key := "cmha/service/"+servicename +"/alerts/" +hosts[host] + "/"
			_key := "cmha/service/"+servicename +"/alerts/" +hosts[host] + "/"
			now_alerts := "[ERROR]: [" + servicename + "] " + hosts[host] + " " + addr + " repl_err_counter is 1"
			log.Infof("kvValue",kvValue)
			if kvValue == "1" {
				logs, _, _err := kv.List(_key, nil)
			        if _err != nil {
                			log.Errorf("query %s list failed!",_key,err)
             				break
        			}
				if logs != nil {
					c := 0
					var alerts_frequency []string
					for _, v := range logs {
						key_items := strings.Split(v.Key, "/")
			                        key_timestamp := key_items[len(key_items)-1]
						L_time, _err := StringToInt(key_timestamp)
                        			if _err != nil {
         			                       log.Errorf("string to int failed!",err)
                                			continue
                        			}
						S_time := C_time-L_time
						S_time_string := strconv.FormatInt(S_time, 10)
                                                S_time_int, err := strconv.Atoi(S_time_string)
						if err != nil {
                                        	        log.Errorf("S_time string to int error: %s", err)
                                                        break
                                                }
						if S_time_int <= frequency_int {
							alert_kvPair, err := GetKv(v.Key, kv)
							if err != nil {
                                         			log.Errorf("Get %s failed!",v.Key,err)
                                         			break
                                 			}
							if alert_kvPair != nil {
								alert_kvPair_value := string(alert_kvPair.Value)
								if alert_kvPair_value == now_alerts {
									c +=1
								}
								
							}else{
								continue
							}
						}else{
							alerts_frequency = append(alerts_frequency,v.Key)
						}			
					}
					if c == 0 {
						alertid :=time.Now().UnixNano()
						id, alertdate := GetNowTime()
						log.Infof("repl_err_counter is 1 c==0:",alertid,alertdate)
					       	alter_counter_value := "[" + alertdate + "]" + "@" + "[ERROR]: [" + servicename + "] " + hosts[host] + " " + addr + " repl_err_counter is 1"
                           			keys := _keys + strconv.FormatInt(alertid,10)
						alerts_key := _alerts_key + strconv.FormatInt(alertid,10)
					  	alerts_key_value := now_alerts
                                              	err = PutKv(servicename, keys, alter_counter_value, kv)
                                              	if err != nil {
                                                      	log.Errorf("Put %s failed!", keys, err)
                                               	  	break
                                               	}
						log.MyLoGGer(id).Println("[ERROR]: [" + servicename + "] " + hosts[host] + " " + addr + " repl_err_counter is 1")
                                             	err = PutKv(servicename,alerts_key,alerts_key_value,kv)
                                              	if err != nil {
    	                                          	log.Errorf("Put %s failed!",alerts_key,err)
                                                      	break
                           	           	}	
						for i:=0;i<len(alerts_frequency);i++ {
							err = DeleteKv(alerts_frequency[i],kv)
							if err != nil {
								log.Errorf("DeleteKv:",err)
								continue
							}
						}
					}
				}else{
					alertid :=time.Now().UnixNano()
					id, alertdate := GetNowTime()
					log.Infof("repl_err_counter is 1:",alertid,alertdate)
					keys := _keys + strconv.FormatInt(alertid,10)
					alerts_key := _alerts_key + strconv.FormatInt(alertid,10)
                                        alert_counter_value := "[" + alertdate + "]" + "@" + "[ERROR]: [" + servicename + "] " + hosts[host] + " " + addr + " repl_err_counter is 1"
                                        alerts_key_value := now_alerts
                                        err = PutKv(servicename, keys, alert_counter_value, kv)
                                        if err != nil {
                                                log.Errorf("Put %s failed!", keys, err)
                                                break
                                        }
					log.MyLoGGer(id).Println("[ERROR]: [" + servicename + "] " + hosts[host] + " " + addr + " repl_err_counter is 1")
                                        err = PutKv(servicename,alerts_key,alerts_key_value,kv)
                                        if err != nil {
                                                log.Errorf("Put %s failed!",alerts_key,err)
                                                break
                                        }

				}
			}
		}
	}
	return nil 
}
