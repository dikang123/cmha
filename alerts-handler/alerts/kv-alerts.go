package alerts

import (
	"crypto/md5"
	//	"fmt"
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

func IsLerderorMonitorOrRole(servicename, kvname string) {
	if kvname == leader {
		err := IsLeader(servicename, kvname)
		if err != nil {
			log.Errorf("leader kv alerts failed!", err)
			return
		}
	} else if kvname == monitor {
		err := IsMonitor(servicename, kvname)
		if err != nil {
			log.Errorf("monitor kv alerts failed!", err)
			return
		}
	}
}

func IsLeader(servicename, kvname string) error {
	time.Sleep(10 * time.Second)
	log.Info("leader")
	C_time, dataTimeStr := GetNowTime()
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
//	key := GetServiceNameLeader(servicename)
	last_leader_key :=GetLastLeaterKey(servicename)
	log.Infof("last_leader_key:",last_leader_key)
	kvValue, last_leaderValue, _ := ConsulGetAndPutKv(leader_key,last_leader_key, servicename, kv)
	log.Infof("last_leaderValue:",last_leaderValue)
	var last_leader string
	if last_leaderValue !=""{
		last_leader =GetLastLeader(last_leaderValue)
	}
	log.Infof("last_leader:",last_leader)
	if kvValue == "" {
		log.Info("111leader")
		log.MyLoGGer().Println("[ERROR]: [" + servicename + "] MySQL occurs failover,The current leader does not exist")
	} else if last_leader == "" {
		log.Info("222leader")
		err = AlertsAndPutKv(servicename,kvValue, dataTimeStr, last_leader_key, kv)
		if err != nil {
			log.Error(err)
			return err
		}
	} else {
		log.Info("333leader")
		last_time, err := GetLastTime(last_leader_key, kv)
		if err != nil {
			log.Error(err)
			return err
		}
		log.Infof("last_time:",last_time)
		log.Info("444leader")
		err = IsAlerts(last_time, servicename, last_leader, kvValue,last_leader_key,C_time, kv)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

func IsMonitor(servicename, kvname string) error {
	log.Info("monitor")
	var C_time = time.Now().Unix()

	log.Infof("C_time:%d test", C_time)
	var timeLayout = "2006-01-02 15:04:05"
	var dataTimeStr = time.Unix(C_time, 0).Format(timeLayout)
	log.Infof("dataTimeStr:%s", dataTimeStr)
	var err error
	log.Infof("servicename: %s,kvname: %s", servicename, kvname)
	ip := beego.AppConfig.String("ip")
	frequency := beego.AppConfig.String("frequency")
	frequency_int, err := strconv.Atoi(frequency)
	if err != nil {
		log.Errorf("frequency string to int error: %s", err)
		return err
	}
	port := beego.AppConfig.String("port")
	hostname := beego.AppConfig.Strings("hostname")
	otherhostname := beego.AppConfig.Strings("otherhostname")
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
	for host := range hostname {
		hosts = append(hosts, hostname[host])
	}
	for otherhost := range otherhostname {
		hosts = append(hosts, otherhostname[otherhost])
	}
	for host := range hosts {
		var addr string
		log.Infof("host:%s", hosts[host])
		monitor_key := "monitor/" + hosts[host]
		kvPair, err = GetKv(monitor_key, kv)
		if err != nil {
			log.Errorf("Get leader failed!", err)
			//	return
			continue
		}
		log.Info("get monitor_key success!")
		kvValue := string(kvPair.Value)
		health := client.Health()
		healthpair, _, err := health.Service(servicename, "", false, nil)
		if err != nil {
			log.Errorf("err", err)
			//	return
			continue
		}
		log.Infof("get servicename:%s success!", servicename)
		log.Infof("size:%s", len(healthpair))
		for index := range healthpair {
			if healthpair[index].Node.Node != hosts[host] {
				log.Infof("not local node:%s,%s", healthpair[index].Node.Node, hosts[host])
				continue
			}
			log.Infof("local node:%s,%s", healthpair[index].Node.Node, hosts[host])
			addr = healthpair[index].Node.Address
			md5s := hosts[host] + addr + kvValue
			ms := []byte(md5s)
			m := md5.Sum(ms)
			key := hosts[host] + "_alerts_counter"
			if kvValue == "1" {
				log.Info("111monitor")
				alert_counter_key := hosts[host] + "_alerts_counter"
				log.Infof("alert_counter_key:%s", alert_counter_key)
				alert_kvPair, err := GetKv(alert_counter_key, kv)
				if err != nil {
					log.Errorf("Get leader failed!", err)
					//			return
					break
				}
				if alert_kvPair != nil {
					log.Info("222monitor")
					log.Infof("get alert_counter_key is not empry")
					alert_kvPair_value := string(alert_kvPair.Value)
					if strings.Contains(alert_kvPair_value, string(m[:])) {
						alert_kvPair_values := strings.Split(alert_kvPair_value, "@")
						var L_time string
						for i := range alert_kvPair_values {
							if i == 0 {
								L_time = alert_kvPair_values[i]
								break
							}
						}
						L_times, _ := time.ParseInLocation("2006-01-02 15:04:05", L_time, time.Local)
						L_timess := L_times.Unix()
						S_time := C_time - L_timess
						S_time_string := strconv.FormatInt(S_time, 10)
						S_time_int, err := strconv.Atoi(S_time_string)
						if err != nil {
							log.Errorf("S_time string to int error: %s", err)
							break
						}
						if S_time_int > frequency_int {
							log.Info("333monitor")
							log.Infof("alerts")
							log.MyLoGGer().Println("[ERROR]: [" + servicename + "] " + hosts[host] + " " + addr + " repl_err_counter is 1")
							alter_counter_value := dataTimeStr + "@" + string(m[:])
							err = PutKv(servicename, key, alter_counter_value, kv)
							if err != nil {
								log.Errorf("Put %s failed!", key, err)
								break
							}
						}
					} else {
						log.Info("444monitor")
						log.MyLoGGer().Println("[ERROR]: [" + servicename + "] " + hosts[host] + " " + addr + " repl_err_counter is 1")
						alter_counter_value := dataTimeStr + "@" + string(m[:])
						err = PutKv(servicename, key, alter_counter_value, kv)
						if err != nil {
							log.Errorf("Put %s failed!", key, err)
							break
						}
						log.Infof("put key success")
					}
				} else {
					log.Info("555monitor")
					log.Infof("get alert_counter_key is empary")
					log.MyLoGGer().Println("[ERROR]: [" + servicename + "] " + hosts[host] + " " + addr + " repl_err_counter is 1")
					alert_counter_value := dataTimeStr + "@" + string(m[:])
					err = PutKv(servicename, key, alert_counter_value, kv)
					if err != nil {
						log.Errorf("Put %s failed!", key, err)
						break
					}
				}
			}
		}
	}
	return nil
}
