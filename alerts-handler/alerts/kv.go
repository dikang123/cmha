package alerts

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/upmio/alerts-handler/log"
)

func PutKv(servicename, key, kvValue string, kv *consulapi.KV) error {
	value := consulapi.KVPair{
		Key:   key,
		Value: []byte(kvValue),
	}
	_, err := kv.Put(&value, nil)
	if err != nil {
		if kvValue == "" {
			log.Errorf("Create %s failed! %s", key, servicename)
		} else {
			log.Errorf("Put %s failed! %s", key, servicename)
		}
		return err
	}
	return nil
}

func GetKv(key string, kv *consulapi.KV) (*consulapi.KVPair, error) {
	kvPair, _, err := kv.Get(key, nil)
	if err != nil {
		log.Errorf("Get %s failed!", key, err)
		return nil, err
	}
	return kvPair, nil
}


func DeleteKv(key string, kv *consulapi.KV)error{
	_, err := kv.Delete(key, nil)
	if err != nil {
		log.Errorf("Delete failed:",key,err)
		return err
	}
	return nil
}

func ReadFile(catalog *consulapi.Catalog, kv *consulapi.KV, C_time int64, dataTimeStr, servicename string, iscs bool) {
	frequency := beego.AppConfig.String("frequency")
	frequency_int, err := strconv.Atoi(frequency)
	if err != nil {
		log.Errorf("frequency string to int error:", err)
	}
	var fi *os.File
	if iscs {
		fi, err = os.Open("/usr/local/cmha/alerts-handler/logs/template_cs_node_alerts.log")
		if err != nil {
			log.Errorf("Error: %s\n", err)
			return
		}
	} else {
		fi, err = os.Open("/usr/local/cmha/alerts-handler/logs/template_service_node_alerts.log")
		if err != nil {
			log.Errorf("Error: %s\n", err)
			return
		}
	}

	defer fi.Close()
	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		a_string := string(a)
		status, checkid, node, output := PareInfo(a_string)
		nodes, _, err := catalog.Nodes(nil)
		if err != nil {
			log.Errorf("query nodes error:", err)
			continue
		}
		var address string
		var nodename string
		for index := range nodes {
			if nodes[index].Node == node {
				address = nodes[index].Address
				nodename = nodes[index].Node
				break
			}
		}
		var _alert_counter_key,alert_counter_key string
		var _alerts_key,alerts_key string
		var _key string
		if iscs {
			_alert_counter_key = "cmha/service/CS/alerts/alerts_counter/"
			_alerts_key = "cmha/service/CS/alerts/" + nodename + "/"
			_key = "cmha/service/CS/alerts/" + nodename + "/"
		} else {
			_alert_counter_key = "cmha/service/" + servicename + "/alerts/alerts_counter/"
			_alerts_key = "cmha/service/" +servicename + "/alerts/" +nodename +"/"
			_key = "cmha/service/"+servicename +"/alerts/" + nodename + "/"
		}
		now_alerts :=NowAlerts(output, servicename, node, address, status, checkid)
		logs, _, _err := kv.List(_key, nil)
                if _err != nil {
                	log.Errorf("query %s list failed!",_key,err)
                        break
                }
		if logs != nil {
			c := 0
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
                                        		c += 1
							break
                                        	}
					}else{
                                        	continue
                                	}
                                }
			}
			if c == 0 {
				alertid :=time.Now().UnixNano()
				id, alertdate := GetNowTime()
				log.Infof("service or cs:",alertid,alertdate,_alert_counter_key)
				str := ServiceAlerts(output, servicename, node, address, status, checkid,id)
                                alter_counter_value := "[" + alertdate + "]" + "@" + str
                                alerts_key_value := now_alerts
				if iscs{
					alert_counter_key = _alert_counter_key + strconv.FormatInt(alertid,10)
		                	alerts_key = _alerts_key + strconv.FormatInt(alertid,10)
				}else{
					alert_counter_key = _alert_counter_key + strconv.FormatInt(alertid,10)
		                	alerts_key = _alerts_key + strconv.FormatInt(alertid,10)
				}
                                err = PutKv(servicename, alert_counter_key, alter_counter_value, kv)
                                if err != nil {
                                	log.Errorf("Put %s failed!", alert_counter_key, err)
                                        break
                                }
				
                                err = PutKv(servicename,alerts_key,alerts_key_value,kv)
                                if err != nil {
                                	log.Errorf("Put %s failed!",alerts_key,err)
                                        break
                               	}	
			}
		}else {
			alertid :=time.Now().UnixNano()
			id, alertdate := GetNowTime()
			log.Infof("cs or service:",alertid,alertdate,_alert_counter_key)
			str := ServiceAlerts(output, servicename, node, address, status, checkid,id)
                        alert_counter_value := "[" + alertdate + "]" + "@" + str
                        alerts_key_value := now_alerts
			if iscs{
                                        alert_counter_key = _alert_counter_key + strconv.FormatInt(alertid,10)
                                        alerts_key = _alerts_key + strconv.FormatInt(alertid,10)
                        }else{                                        
					alert_counter_key = _alert_counter_key + strconv.FormatInt(alertid,10)
                                        alerts_key = _alerts_key + strconv.FormatInt(alertid,10)
                        }
                        err = PutKv(servicename, alert_counter_key, alert_counter_value, kv)
                        if err != nil {
                                log.Errorf("Put %s failed!", alert_counter_key, err)
                                break
                        }

                        err = PutKv(servicename,alerts_key,alerts_key_value,kv)
                        if err != nil {
                                log.Errorf("Put %s failed!",alerts_key,err)
                                break
                        }
		}
	}
}

func ServiceAlerts(output, servicename, node, address, status, checkid string,alertid int64) string{
	if strings.Contains(output, "haproxy") {
		log.MyLoGGer(alertid).Println("[ERROR]: [" + servicename + "] " + node + " " + address + " haproxy service " + status)
		str := "[ERROR]: ["  + servicename + "] " + node + " " + address + "  haproxy service " + status
		return str
	} else if strings.Contains(output, "keepalived") {
		log.MyLoGGer(alertid).Println("[ERROR]: [" + servicename + "] " + node + " " + address + " keepalived service " + status)
		str := "[ERROR]: [" + servicename + "] " + node + " " + address + " keepalived service " + status
		return str
	} else if strings.Contains(output, "consul-template") {
		log.MyLoGGer(alertid).Println("[ERROR]: [" + servicename + "] " + node + " " + address + " consul-template service " + status)
		str := "[ERROR]: [" + servicename + "] " + node + " " + address + " consul-template service " + status
		return str
	} else if status == "warning" {
		log.MyLoGGer(alertid).Println("[ERROR]: [" + servicename + "] " + node + " " + address + " replication IO_thread " + output + " " + status)
		str := "[ERROR]: [" + servicename + "] " + node + " " + address + " " + output + " " + status
		return str
	} else {
		log.MyLoGGer(alertid).Println("[ERROR]: [" + servicename + "] " + node + " " + address + " " + output + " " + status)
		str := "[ERROR]: [" + servicename + "] " + node + " " + address + " " + output + " " + status
		return str
	}
}

func NowAlerts(output, servicename, node, address, status, checkid string) string {
        if strings.Contains(output, "haproxy") {
                str := "[ERROR]: ["  + servicename + "] " + node + " " + address + "  haproxy service " + status
                return str
        } else if strings.Contains(output, "keepalived") {
                str := "[ERROR]: [" + servicename + "] " + node + " " + address + " keepalived service " + status
                return str
        } else if strings.Contains(output, "consul-template") {
                str := "[ERROR]: [" + servicename + "] " + node + " " + address + " consul-template service " + status
                return str
        } else if status == "warning" {
                str := "[ERROR]: [" + servicename + "] " + node + " " + address + " " + output + " " + status
                return str
        } else {
                str := "[ERROR]: [" + servicename + "] " + node + " " + address + " " + output + " " + status
                return str
        }
}

func PareInfo(a_string string) (string, string, string, string) {
	var status, checkid, node, output string
	if strings.Contains(a_string, "%%") {
		a_split := strings.Split(a_string, "%%")
		for split := range a_split {
			if strings.Contains(a_split[split], "||") && strings.Contains(a_split[split], "Status") {
				a_status := strings.Split(a_split[split], "||")
				for i := range a_status {
					if i == 1 {
						status = a_status[i]
					}
				}
			} else if strings.Contains(a_split[split], "||") && strings.Contains(a_split[split], "CheckID") {
				a_checkid := strings.Split(a_split[split], "||")
				for i := range a_checkid {
					if i == 1 {
						checkid = a_checkid[i]
					}
				}
			} else if strings.Contains(a_split[split], "||") && strings.Contains(a_split[split], "Node") {
				a_node := strings.Split(a_split[split], "||")
				for i := range a_node {
					if i == 1 {
						node = a_node[i]
					}
				}
			} else if strings.Contains(a_split[split], "||") && strings.Contains(a_split[split], "Output") {
				a_output := strings.Split(a_split[split], "||")
				for i := range a_output {
					if i == 1 {
						output = a_output[i]
					}
				}
			}
		}
	}
	return status, checkid, node, output
}
