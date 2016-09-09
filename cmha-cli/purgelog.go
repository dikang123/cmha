package main

import (
	"fmt"

	"github.com/upmio/cmha-cli/cliconfig"
	"strings"
	"errors"
)

func PurgeLog(_service string, _node string, _logtype string, _timestamps ...string) error {

	if len(_timestamps) > 0 {
		if len(_timestamps) == 1 {
			if len(_timestamps[0]) == 10 {

				purge_a_log(_service, _node, _logtype, _timestamps[0])

			} else if len(_timestamps[0]) == 21 {

				_ids := strings.Split(_timestamps[0], ",")

				if len(_ids) != 2 {

					err := errors.New(fmt.Sprintf("%s is not a valid log id.", _timestamps[0]))
					return err

				}

				purge_logs(_service, _node, _logtype, _ids[0], _ids[1])

			} else {

				err := errors.New(fmt.Sprintf("%s is not a valid log id.", _timestamps[0]))
				return err

			}

		} else {

			for _, logtimestamp := range _timestamps {
				purge_a_log(_service, _node, _logtype, logtimestamp)
			}
		}

	} else {
		purge_all_logs(_service, _node, _logtype)
	}

	return nil

}

func purge_logs(service string, node string, logtype string, _from string, _to string) {

	time_from, err_from := StringToInt(_from)
	if err_from != nil {
		fmt.Println(err_from)
		return
	}

	time_to, err_to := StringToInt(_to)
	if err_to != nil {
		fmt.Println(err_to)
		return
	}

	time_from, time_to = int64_arrange(time_from, time_to)

	client, err := cliconfig.Consul_Client_Init()

	if err != nil {
		fmt.Println("cluster Create consul-api client failure!", err)
		return
	}

	kv := client.KV()

	_key := fmt.Sprintf("%s/%s/%s/", service, node, logtype)

	logs, _, _err := kv.List(_key, nil)
	if _err != nil {
		fmt.Println(_err.Error())
		return
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
				purge_a_log(service, node, logtype, logid_timestamp)
			}

		}

		return
	}

	fmt.Printf("No %s log on %s/%s\n", logtype, service, node)
	fmt.Println()

}

func purge_all_logs(service string, node string, logtype string) {

	client, err := cliconfig.Consul_Client_Init()

	if err != nil {
		fmt.Println("cluster Create consul-api client failure!", err)
		return
	}

	kv := client.KV()

	_key := fmt.Sprintf("%s/%s/%s/", service, node, logtype)

	logs, _, _err := kv.List(_key, nil)
	if _err != nil {
		fmt.Println(_err.Error())
		return
	}

	if logs != nil {
		for _, v := range logs {
			logid_items := strings.Split(v.Key, "/")
			logid_timestamp := logid_items[len(logid_items)-1]

			purge_a_log(service, node, logtype, logid_timestamp)

		}

		return
	}

	fmt.Printf("No %s log on %s/%s\n", logtype, service, node)
	fmt.Println()

}

func purge_a_log(service string, node string, logtype string, _timestamp string) {

	client, err := cliconfig.Consul_Client_Init()

	if err != nil {
		fmt.Println("cluster Create consul-api client failure!", err)
		return
	}

	kv := client.KV()

	_id := _timestamp

	_key := fmt.Sprintf("%s/%s/%s/%s", service, node, logtype, _id)

	if _, err := kv.Delete(_key, nil); err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		fmt.Printf("Log %s has deleted.\n", _key)
		return
	}

}
