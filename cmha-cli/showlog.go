package main

import (
	"fmt"

	"github.com/upmio/cmha-cli/cliconfig"
	"errors"
	"strconv"
	"strings"
	"time"
)

func ShowLog(_service string, _node string, _logtype string, _timestamps ...string) error {
	if len(_timestamps) > 0 {

		if len(_timestamps) == 1 {
			if len(_timestamps[0]) == 10 {

				show_log_detail(_service, _node, _logtype, _timestamps[0])

			} else if len(_timestamps[0]) == 21 {

				_ids := strings.Split(_timestamps[0], ",")

				if len(_ids) != 2 {

					err := errors.New(fmt.Sprintf("%s is not a valid log id.", _timestamps[0]))
					return err

				}

				show_logs(_service, _node, _logtype, _ids[0], _ids[1])

			} else {

				err := errors.New(fmt.Sprintf("%s is not a valid log id.", _timestamps[0]))
				return err

			}

		} else {

			for _, logtimestamp := range _timestamps {
				show_log_detail(_service, _node, _logtype, logtimestamp)
			}
		}

	} else {
		show_all_log(_service, _node, _logtype)
	}

	return nil

}

func show_logs(service string, node string, logtype string, _from string, _to string) {

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

	_key := fmt.Sprintf("%s/%s/%s/%s/%s/%s/", "cmha", "service", service, "log", node, logtype)

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
				show_log_detail(service, node, logtype, logid_timestamp)
			}

		}

		return
	}

	fmt.Printf("No %s log on %s/%s\n", logtype, service, node)
	fmt.Println()

}

func show_all_log(service string, node string, logtype string) {
	client, err := cliconfig.Consul_Client_Init()

	if err != nil {
		fmt.Println("cluster Create consul-api client failure!", err)
		return
	}

	kv := client.KV()

	_key := fmt.Sprintf("%s/%s/%s/%s/%s/%s/", "cmha", "service", service, "log", node, logtype)
	logs, _, _err := kv.List(_key, nil)
	if _err != nil {
		fmt.Println(_err.Error())
		return
	}

	_log_count := len(logs)
	_logs_data := make([][]string, _log_count)

	if logs != nil {
		for k, v := range logs {
			logid_items := strings.Split(v.Key, "/")
			logid_timestamp := logid_items[len(logid_items)-1]

			_logs_data[k] = make([]string, 3)
			_logs_data[k][0] = logid_timestamp
			_logs_data[k][1] = StringToTime(logid_timestamp)
			_logs_data[k][2] = v.Key

		}

		_th := []string{"ID",
			"Date & Time",
			"Connet"}
		TableRender(_th, _logs_data, ALIGN_CENTRE)
		fmt.Println()

		return
	}

	fmt.Printf("No %s log on %s/%s\n", logtype, service, node)
	fmt.Println()

}

func show_log_detail(service string, node string, logtype string, _timestamp ...string) {

	client, err := cliconfig.Consul_Client_Init()

	if err != nil {
		fmt.Println("cluster Create consul-api client failure!", err)
		return
	}

	kv := client.KV()

	var _id string = ""
	if len(_timestamp) > 0 {
		_id = _timestamp[0]
	}

	_key := fmt.Sprintf("%s/%s/%s/%s/%s/%s/%s", "cmha", "service", service, "log", node, logtype, _id)

	logs, _, _err := kv.List(_key, nil)
	if _err != nil {
		fmt.Println(_err.Error())
		return
	}

	fmt.Println()

	if logs != nil {
		for _, v := range logs {
			fmt.Printf("%s", parseLogData(string(v.Value), logtype))
		}

		fmt.Println()
		return

	}

	fmt.Printf("No %s log on %s/%s/%s\n", logtype, service, node, _timestamp[0])
	fmt.Println()

}

func parseLogData(encoded_log_data string, log_type string) string {

	if len(encoded_log_data) < 1 {
		return ""
	}

	var _log_content string = ""

	logs_encoded := strings.Split(encoded_log_data, "|")

	for _, log_encoded := range logs_encoded {

		var _args_log []string
		if len(log_encoded) > 14 {

			_arr_args_log := strings.Split(log_encoded, "{{")

			if len(_arr_args_log) > 1 {

				for _, _arg_value := range _arr_args_log[1:] {
					_args_log = append(_args_log, _arg_value)
				}

			}

		}

		switch log_type {
		case LOG_MHA:
			if len(_args_log) > 0 {
				_log_content = get_mha_log_by_index(log_encoded[10:13], _args_log)
			} else {
				_log_content = get_mha_log_by_index(log_encoded[10:13])
			}
		case LOG_MONITOR:
			if len(_args_log) > 0 {
				_log_content = get_monitor_log_by_index(log_encoded[10:13], _args_log)
			} else {
				_log_content = get_monitor_log_by_index(log_encoded[10:13])
			}
		default:
			_log_content = ""
		}

		fmt.Printf("[%s] ", StringToTime(log_encoded[:10]))
		fmt.Printf("%s", _log_content)
		fmt.Printf("\n")

	}

	return ""

}

func StringToTime(instr string) string {

	var _time string
	if t_int, err := strconv.ParseInt(instr, 10, 64); err != nil {
		return ""
	} else {
		_t := time.Unix(t_int, 0)
		_time = fmt.Sprintf("%s", _t.Local())
	}

	return _time
}

func StringNanoToTime(instr string)string{
	var _time string
        if t_int, err := strconv.ParseInt(instr, 10, 64); err != nil {
                return ""
        } else {
		_t :=time.Unix(t_int/1e9, 0)
                _time = fmt.Sprintf("%s", _t.Local())
        }

        return _time
}

func get_monitor_log_by_index(_index string, _args ...[]string) string {

	map_monitor_log := make(map[string]string)
	map_monitor_log["001"] = "[I] Monitor Handler Triggered"
	map_monitor_log["002"] = "[E] Create consul-api client failed!"
	map_monitor_log["003"] = "[I] Give up switching to async replication!"
	map_monitor_log["004"] = "[I] Monitor Handler Completed"
	map_monitor_log["005"] = "[I] Create consul-api client successfully!"

	map_monitor_log["006"] = "[E] Get peer service %s health status from CS failed!"
	map_monitor_log["007"] = "[I] Get peer service %s health status from CS successfully!"

	map_monitor_log["008"] = "[E] %s peer service not exist in CS!"
	map_monitor_log["009"] = "[I] %s peer service exist in CS!"
	map_monitor_log["010"] = "[I] Service health status is passing in CS!"

	map_monitor_log["011"] = "[W] Warning! Peer database %s replicaton error. Service health status is warning in CS!"
	map_monitor_log["012"] = "[I] Current switch_async value is %s"

	map_monitor_log["013"] = "[I] Config file switch_async format error,off or on!"
	map_monitor_log["014"] = "[I] Service health status is critical in CS!"
	map_monitor_log["015"] = "[E] Set peer database repl_err_counter to 1 in CS failed!"
	map_monitor_log["016"] = "[I] Set peer database repl_err_counter to 1 in CS successfully!"
	map_monitor_log["017"] = "[E] Not passing,not waring,not critical ,is invalid state!"
	map_monitor_log["018"] = "[E] Get and check current service leader from CS failed!"
	map_monitor_log["019"] = "[I] Get and check current service leader from CS successfully!"


	map_monitor_log["020"] = "[I] %s is not service leader!"
	map_monitor_log["021"] = "[I] %s is service leader!"
	map_monitor_log["022"] = "[E] Get %s service health status from CS failed!"
	map_monitor_log["023"] = "[I] Get %s service health status from CS successfully!"
	map_monitor_log["024"] = "[I] %s service health status is %s"

	map_monitor_log["025"] = "[E] Create connection object to local database failed!"
	map_monitor_log["026"] = "[I] Create connection object to local database successfully!"
	map_monitor_log["027"] = "[E] Connected to local database failed!"
	map_monitor_log["028"] = "[I] Connected to local database successfully!"
	map_monitor_log["029"] = "[E] Set rpl_semi_sync_master_keepsyncrepl=0 failed!"
	map_monitor_log["030"] = "[I] Set rpl_semi_sync_master_keepsyncrepl=0 successfully!"
	map_monitor_log["031"] = "[E] Set rpl_semi_sync_master_trysyncrepl=0 failed!"
	map_monitor_log["032"] = "[I] Set rpl_semi_sync_master_trysyncrepl=0 successfully!"
	map_monitor_log["033"] = "[I] Switching local database to async replication!"
	map_monitor_log["034"] = "[I] Monitor Handler Sleep 60s!"
	map_monitor_log["035"] = "[I] Connecting to peer database......"
	map_monitor_log["036"] = "[E] Create connection object to peer database failed!"
	map_monitor_log["037"] = "[I] Create connection object to peer database successfully!"
	map_monitor_log["038"] = "[E] Connected to the peer database failed!"
	map_monitor_log["039"] = "[I] Connected to the peer database successfully!"
	map_monitor_log["040"] = "[I] Checking peer database I/O thread status. Failed!"
	map_monitor_log["041"] = "[I] Checking peer database I/O thread status. Successfully!"
	map_monitor_log["042"] = "[E] Resolve slave status failed!"

	map_monitor_log["043"] = "[I] The I/O thread status is %s!"

	var ret_monitor_log string = ""
	if len(_args) > 0 {

		_log_args := stringlist_to_interface(_args[0])
		ret_monitor_log = fmt.Sprintf(map_monitor_log[_index], _log_args...)

	} else {
		ret_monitor_log = fmt.Sprintf(map_monitor_log[_index])
	}

	return ret_monitor_log

}

func get_mha_log_by_index(_index string, _args ...[]string) string {

	map_mha_log := make(map[string]string)
	map_mha_log["001"] = "[I] MHA Handler Triggered"
	map_mha_log["002"] = "[E] Create consul-api client failed!"
	map_mha_log["003"] = "[I] Give up leader election"
	map_mha_log["004"] = "[I] MHA Handler Completed"
	map_mha_log["005"] = "[I] Create consul-api client successfully!"
	map_mha_log["006"] = "[E] Get and check current service leader from CS failed!"
	map_mha_log["007"] = "[I] Get and check current service leader from CS successfully!"

	map_mha_log["008"] = "[E] Get %s repl_err_counter=%s failed!"
	map_mha_log["009"] = "[I] Get %s repl_err_counter=%s successfully!"
	map_mha_log["010"] = "[E] %s give up leader election"

	map_mha_log["011"] = "[E] Not service leader,Please create!"
	map_mha_log["012"] = "[I] Leader exist!"
	map_mha_log["013"] = "[I] Leader does not exist!"

	map_mha_log["014"] = "[E] Get and check %s service health status failed!"
	map_mha_log["015"] = "[I] Get and check %s service health status successfully!"
	map_mha_log["016"] = "[I] %s service does not exist!"
	map_mha_log["017"] = "[I] %s service exist!"
	map_mha_log["018"] = "[E] %s not is %s!"

	map_mha_log["019"] = "[E] Clean service leader value in CS failed!"
	map_mha_log["020"] = "[I] Clean service leader value in CS successfully!"
	map_mha_log["021"] = "[E] Status is critical!"
	map_mha_log["022"] = "[I] Status is not critical"
	map_mha_log["023"] = "[E] Session create failed!"
	map_mha_log["024"] = "[I] Session create successfully!"
	map_mha_log["025"] = "[E] format error,json or hap!"
	map_mha_log["026"] = "[E] Send service leader request to CS failed!"
	map_mha_log["027"] = "[I] Send service leader request to CS successfully!"

	map_mha_log["028"] = "[E] Becoming service leader failed! Connection string is %s:%s"
	map_mha_log["029"] = "[I] Becoming service leader successfully! Connection string is %s:%s"

	map_mha_log["030"] = "[E] Set peer database repl_err_counter to 1 in CS failed!"
	map_mha_log["031"] = "[I] Set peer database repl_err_counter to 1 in CS successfully!"
	map_mha_log["032"] = "[E] Create connection object to local database failed!"
	map_mha_log["033"] = "[I] Create connection object to local database successfully!"
	map_mha_log["034"] = "[E] Connected to local database failed!"
	map_mha_log["035"] = "[I] Connected to local database successfully!"
	map_mha_log["036"] = "[E] Set local database Read_only mode failed!"
	map_mha_log["037"] = "[I] Local database downgrade failed!"
	map_mha_log["038"] = "[I] Set local database Read_only mode successfully!"
	map_mha_log["039"] = "[I] Local database downgrade successfully!"
	map_mha_log["040"] = "[E] Stop local database replication I/O thread failed!"
	map_mha_log["041"] = "[I] Stop local database replication I/O thread successfully!"
	map_mha_log["042"] = "[E] Checking local database SQL thread status. Failed!"
	map_mha_log["043"] = "[I] Checking local database SQL thread status. Succeed!"
	map_mha_log["044"] = "[E] Resolve slave status failed!"

	map_mha_log["045"] = "[E] The SQL thread status is %s!"

	map_mha_log["046"] = "[I] The SQL thread status is Yes!"
	map_mha_log["047"] = "[E] Checking relay log applying status failed!"
	map_mha_log["048"] = "[I] Checking relay log applying status successfully!"
	map_mha_log["049"] = "[E] Resolve master_pos_wait failed!"
	map_mha_log["050"] = "[E] Relay log applying failed!"
	map_mha_log["051"] = "[I] Relay log applying completed!"
	map_mha_log["052"] = "[E] Set rpl_semi_sync_master_keepsyncrepl=0 failed!"
	map_mha_log["053"] = "[I] Set rpl_semi_sync_master_keepsyncrepl=0 successfully!"
	map_mha_log["054"] = "[E] Set rpl_semi_sync_master_trysyncrepl=0 failed!"
	map_mha_log["055"] = "[I] Set rpl_semi_sync_master_trysyncrepl=0 successfully!"
	map_mha_log["056"] = "[I] Switching local database to async replication!"
	map_mha_log["057"] = "[E] Set local database Read/Write mode failed!"
	map_mha_log["058"] = "[I] Set local database Read/Write mode successfully!"

	var ret_mha_log string = ""
	if len(_args) > 0 {

		_log_args := stringlist_to_interface(_args[0])
		ret_mha_log = fmt.Sprintf(map_mha_log[_index], _log_args...)

	} else {

		ret_mha_log = fmt.Sprintf(map_mha_log[_index])
	}
	return ret_mha_log
}

func stringlist_to_interface(list []string) []interface{} {
	vals := make([]interface{}, len(list))
	for k, v := range list {
		vals[k] = v
	}
	return vals
}
