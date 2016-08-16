package alerts

import (
	"bufio"
	"crypto/md5"
	"fmt"
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
		fmt.Println("string", a_string)
		status, checkid, node, output := PareInfo(a_string)
		fmt.Println(status, checkid, node, output)
		nodes, _, err := catalog.Nodes(nil)
		if err != nil {
			log.Errorf("query nodes error:", err)
			continue
		}
		var address string
		for index := range nodes {
			if nodes[index].Node == node {
				address = nodes[index].Address
				break
			}
		}
		md5s := status + checkid + node + address + output
		ms := []byte(md5s)
		m := md5.Sum(ms)
		var alert_counter_key string
		if iscs {
			alert_counter_key = node + "_cs_alerts_counter"
			log.Info(alert_counter_key)
		} else {
			alert_counter_key = node + "_service_alerts_counter"
			log.Info(alert_counter_key)
		}
		log.Infof("alert_counter_key:%s", alert_counter_key)
		alert_kvPair, err := GetKv(alert_counter_key, kv)
		if err != nil {
			log.Errorf("Get leader failed!", err)
			//			return
			break
		}
		if alert_kvPair != nil {
			log.Info("111cs")
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
					log.Info("222cs")
					log.Infof("alerts")
					ServiceAlerts(output, servicename, node, address, status, checkid)
					//					log.MyLoGGer().Println("[ERROR]: [" + servicename + "] " + hosts[host] + " " + addr + " repl_err_counter is 1")
					alter_counter_value := dataTimeStr + "@" + string(m[:])
					err = PutKv(servicename, alert_counter_key, alter_counter_value, kv)
					if err != nil {
						log.Errorf("Put %s failed!", alert_counter_key, err)
						break
					}
				}
			} else {
				log.Info("333cs")
				ServiceAlerts(output, servicename, node, address, status, checkid)
				//				log.MyLoGGer().Println("[ERROR]: [" + servicename + "] " + hosts[host] + " " + addr + " repl_err_counter is 1")
				alter_counter_value := dataTimeStr + "@" + string(m[:])
				err = PutKv(servicename, alert_counter_key, alter_counter_value, kv)
				if err != nil {
					log.Errorf("Put %s failed!", alert_counter_key, err)
					break
				}
				log.Infof("put key success")
			}
		} else {
			ServiceAlerts(output, servicename, node, address, status, checkid)
			log.Infof("get alert_counter_key is empary")
			//			log.MyLoGGer().Println("[ERROR]: [" + servicename + "] " + hosts[host] + " " + addr + " repl_err_counter is 1")
			alert_counter_value := dataTimeStr + "@" + string(m[:])
			err = PutKv(servicename, alert_counter_key, alert_counter_value, kv)
			if err != nil {
				log.Errorf("Put %s failed!", alert_counter_key, err)
				break
			}
		}
	}
}

func ServiceAlerts(output, servicename, node, address, status, checkid string) {
	if strings.Contains(output, "haproxy") {
		log.MyLoGGer().Println("[ERROR]: [" + servicename + "] " + node + " " + address + " haproxy service " + status)
	} else if strings.Contains(output, "keepalived") {
		log.MyLoGGer().Println("[ERROR]: [" + servicename + "] " + node + " " + address + " keepalived service " + status)
	} else if strings.Contains(output, "consul-template") {
		log.MyLoGGer().Println("[ERROR]: [" + servicename + "] " + node + " " + address + " consul-template service " + status)
	} else if status == "warning" {
		log.MyLoGGer().Println("[ERROR]: [" + servicename + "] " + node + " " + address + " replication IO_thread " + status)
	} else if checkid == "serfHealth" {
		log.MyLoGGer().Println("[ERROR]: [" + servicename + "] " + node + " " + address + " consul service " + output + " " + status)
	} else {
		log.MyLoGGer().Println("[ERROR]: [" + servicename + "] " + node + " " + address + " " + status)
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
