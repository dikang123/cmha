package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/upmio/cmha-cli/cliconfig"
)

//show all alerts log
func show_all_alerts(_service string, _list string) error {
	client, err := cliconfig.Consul_Client_Init()

	if err != nil {
		fmt.Println("cluster Create consul-api client failure!", err)
		return err
	}

	kv := client.KV()
	var _key string
	if _service == "" {
		_key = fmt.Sprintf("%s/%s/%s/%s/%s/", "cmha", "service", "CS", "alerts", "alerts_counter")
	} else {
		_key = fmt.Sprintf("%s/%s/%s/%s/%s/", "cmha", "service", _service, "alerts", "alerts_counter")
	}

	logs, _, _err := kv.List(_key, nil)
	if _err != nil {
		fmt.Println(_err.Error())
		return err
	}

	_log_count := len(logs)
	_logs_data := make([][]string, _log_count)

	if logs != nil {
		islist := true
		for k, v := range logs {
			if _list == "" {
				logid_items := strings.Split(v.Key, "/")
				logid_timestamp := logid_items[len(logid_items)-1]
				keys := strings.Split(v.Key, "service/")
				var key string
				for i := range keys {
					if i == 1 {
						key = keys[i]
					}

				}

				_logs_data[k] = make([]string, 3)
				_logs_data[k][0] = logid_timestamp
				_logs_data[k][1] = StringNanoToTime(logid_timestamp)
				_logs_data[k][2] = key

			} else {
				logid_items := strings.Split(v.Key, "/")
				logid_timestamp := logid_items[len(logid_items)-1]
				logid_timestamp_int64, _err := strconv.ParseInt(logid_timestamp, 10, 64)
				if _err != nil {
					fmt.Println(_err.Error())
					return err
				}
				nowtime := time.Now()
				d, _ := time.ParseDuration("-24h")
				timestamp := nowtime.Add(d * 3).Unix()
				if logid_timestamp_int64 >= timestamp {
					logs, _, _err := kv.List(v.Key, nil)
				        if _err != nil {
				                fmt.Println(_err.Error())
				                return err
				        }

       					fmt.Println()

				        if logs != nil {
			        	        for _, v := range logs {
				                        fmt.Printf("%s", strings.Replace(string(v.Value), "@", " ", -1))
                				}

			        	}
				}
				islist = false
			}

		}
		if islist {
			_th := []string{"ID",
				"Date & Time",
				"Connet"}
			TableRender(_th, _logs_data, ALIGN_CENTRE)
		}
		fmt.Println()

		return nil
	}

	fmt.Printf("No alerts on %s\n", _service)
	fmt.Println()
	return nil
}

func show_alerts_detail(service string, _timestamp ...string) error {

	client, err := cliconfig.Consul_Client_Init()

	if err != nil {
		fmt.Println("cluster Create consul-api client failure!", err)
		return err
	}

	kv := client.KV()

	var _id string = ""
	if len(_timestamp) > 0 {
		_id = _timestamp[0]
	}
	var _key string
	if service == "" {
		_key = fmt.Sprintf("%s/%s/%s/%s/%s/%s", "cmha", "service", "CS", "alerts", "alerts_counter", _id)
	} else {
		_key = fmt.Sprintf("%s/%s/%s/%s/%s/%s", "cmha", "service", service, "alerts", "alerts_counter", _id)
	}

	logs, _, _err := kv.List(_key, nil)
	if _err != nil {
		fmt.Println(_err.Error())
		return err
	}

	fmt.Println()

	if logs != nil {
		for _, v := range logs {
			fmt.Printf("%s", strings.Replace(string(v.Value), "@", " ", -1))
		}

		fmt.Println()
		return nil

	}

	fmt.Printf("No alerts on %s,%s\n", service, _timestamp[0])
	fmt.Println()
	return nil

}

func show_alerts(service string, _from string, _to string) error {

	time_from, err_from := StringToInt(_from)
	if err_from != nil {
		fmt.Println(err_from)
		return err_from
	}

	time_to, err_to := StringToInt(_to)
	if err_to != nil {
		fmt.Println(err_to)
		return err_to
	}

	time_from, time_to = int64_arrange(time_from, time_to)

	client, err := cliconfig.Consul_Client_Init()

	if err != nil {
		fmt.Println("cluster Create consul-api client failure!", err)
		return err
	}

	kv := client.KV()
	var _key string
	if service == "" {
		_key = fmt.Sprintf("%s/%s/%s/%s/%s/", "cmha", "service", "CS", "alerts", "alerts_counter")
	} else {
		_key = fmt.Sprintf("%s/%s/%s/%s/%s/", "cmha", "service", service, "alerts", "alerts_counter")
	}

	logs, _, _err := kv.List(_key, nil)
	if _err != nil {
		fmt.Println(_err.Error())
		return err
	}

	if logs != nil {
		for _, v := range logs {
			logid_items := strings.Split(v.Key, "/")
			logid_timestamp := logid_items[len(logid_items)-1]

			_logid, _err := StringToInt(logid_timestamp)
			if _err != nil {
				fmt.Printf("%s\n", err.Error())
				continue
			}

			if _logid >= time_from && _logid <= time_to {
				show_alerts_detail(service, logid_timestamp)
			}
		}

		return nil
	}

	fmt.Printf("No alerts on %s\n", service)
	fmt.Println()
	return nil

}

func AlertBoot(service string) error {
	client, err := cliconfig.Consul_Client_Init()

	if err != nil {
		fmt.Println("cluster Create consul-api client failure!", err)
		return err
	}

	kv := client.KV()
	key := "cmha/service/" + service + "/alerts/alert_boot"
	kvpair, _,err := kv.Get(key, nil)
	if string(kvpair.Value) != "" {
		fmt.Println(string(kvpair.Value))
	} else {
		fmt.Println("alert_boot is empty")
	}
	return nil
}
