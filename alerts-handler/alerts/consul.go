package alerts

import (
	"strings"
	"github.com/astaxie/beego"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/upmio/alerts-handler/log"
)

func GetConfig(adress string) *consulapi.Config {
	config := &consulapi.Config{
		Datacenter: beego.AppConfig.String("datacenter"),
		Token:      beego.AppConfig.String("token"),
		Address:    adress,
	}
	return config
}

func GetClient(config *consulapi.Config) (*consulapi.Client, error) {
	client, err := consulapi.NewClient(config)
	if err != nil {
		log.Errorf("[E] Create consul-api client failed!", err)
		return nil, err
	}
	return client, nil
}

func ReturnKv(client *consulapi.Client) *consulapi.KV {
	kv := client.KV()
	return kv
}

func GetServiceNameLeader(servicename string) string {
	key := servicename + "_leader"
	return key
}

func GetLeaderKey(servicename string) string {
	leader_key := "service/" + servicename + "/leader"
	return leader_key

}

func GetLastLeaterKey(servicename string) string {
	last_leader := "cmha/" +servicename +"/last_leader"
	return last_leader
}

func GetLastLeader(last_leaderValue string) string{
	var last_leader string
	if strings.Contains(last_leaderValue,"||"){
		last_leader_split :=strings.Split(last_leaderValue,"||")
		if len(last_leader_split) == 3{
			for i:=range last_leader_split{
				if i ==2{
					last_leader = last_leader_split[i]
					break	
				}
			}
		}
	}
	return last_leader
}

func ConsulGetAndPutKv(leader_key,last_leader, servicename string, kv *consulapi.KV) (string, string, error) {
	kvPair, err := GetKv(leader_key, kv)
	if err != nil {
		log.Errorf("Get leader failed!", err)
		return "", "", err
	}
/*	err = PutKv(servicename, key, "", kv)
	if err != nil {
		log.Errorf("Put last leader failed!", err)
		return "", "", err
	}*/
	kvValue := string(kvPair.Value)
	last_leaders, err := GetKv(last_leader, kv)
	if err != nil {
		log.Errorf("Get last leader failed!", err)
		return "", "", err
	}
	var last_leaderValue string
	if last_leaders != nil {
		last_leaderValue = string(last_leaders.Value)
	}else{
		last_leaderValue =""
	}
	log.Infof("kvValue:%s,last_leaderValue:%s",kvValue,last_leaderValue)
	return kvValue, last_leaderValue, nil
}
