package main

import (
	"errors"
	"fmt"

	"strings"

	"github.com/upmio/cmha-cli/cliconfig"
)

func Service_Status(args ...string) error {

	if len(args) != 1 {
		err := errors.New("Need a service name as argument.")
		return err
	}

	client, err := cliconfig.Consul_Client_Init()

	if err != nil {
		fmt.Println("cluster Create consul-api client failure!", err)
		return err
	}

	var (
		leaderip string
	)

	kv := client.KV()
	if kp, _, _ := kv.Get("cmha/service/"+args[0]+"/db/leader", nil); kp.Session != "" {
		serviceleader := strings.Split(string(kp.Value), " ")

		leaderip = strings.Split(serviceleader[2], ":")[0]

	}

	servicename := args[0]

	health := client.Health()
	health_service, _, err := health.Service(servicename, "", false, nil)
	if err != nil {
		fmt.Printf("Error: Fail to query cluster health service.\n%s\n", err)
		return err
	}

	if len(health_service) < 1 {
		fmt.Printf("Not exist service such as %s.\n", servicename)
		return nil
	}

	_data_chap := make([][]string, 1)
	_data_db := make([][]string, 1)

	for _, servicedata := range health_service {

		_node_name := servicedata.Node.Node
		// prepare for color render
		node_name := _node_name

		_node_address := servicedata.Node.Address
		// prepare for color render
		node_address := _node_address

		_node_port := servicedata.Service.Port
		// prepare for color render
		node_port := fmt.Sprintf(":%d", _node_port)

		var node_is_leader string = ""
		if node_address == leaderip {
			node_name = ColorRender(node_name, COLOR_SUM)
			node_address = ColorRender(node_address, COLOR_SUM)
			node_port = ColorRender(node_port, COLOR_SUM)
			node_is_leader = ColorRender("y", COLOR_SUM)
		}

		//node type
		node_tag := servicedata.Service.Tags[0]
		var node_type = "db"

		if strings.HasPrefix(node_tag, "chap-") {

			node_type = "chap"

			chap_node_info := make([]string, 8)
			chap_node_info[0] = servicename
			chap_node_info[1] = node_name
			chap_node_info[2] = node_address
			chap_node_info[3] = node_type

			has_unknown_status := false

			for _, checkdata := range servicedata.Checks {

				node_status := checkdata.Status

				node_checkid := checkdata.CheckID

				has_unknown_status = node_status == "critical" && node_checkid == "serfHealth"

				switch node_status {
				case "passing":
					node_status = ColorRender("OK", COLOR_NORMAL)
				case "critical":
					node_status = ColorRender("Fail", COLOR_ERROR)
				default:
					node_status = ColorRender(node_status, COLOR_WARNNING)
				}

				if node_checkid == "serfHealth" {
					chap_node_info[5] = node_status
				} else {
					if chap_node_info[5] == ColorRender("Fail", COLOR_ERROR){
						chap_node_info[4] = ColorRender("UnKnown", COLOR_WARNNING)
					}else{
						chap_node_info[4] = node_status
					}
				}
				//add role
				role_key := "cmha/service/" + args[0] + "/chap/role/" + node_name
				if node_status == ColorRender("OK", COLOR_NORMAL) {
					kp, _, err := kv.Get(role_key, nil)
					if err != nil {
						fmt.Printf("%s\n", err)
					}
					if kp.Value != nil {
						chap_node_info[6] = string(kp.Value)
					} else {
						chap_node_info[6] = ColorRender("UnKnown", COLOR_WARNNING)
					}
				}else{
					chap_node_info[6] = ColorRender("UnKnown", COLOR_WARNNING)
				}
				if chap_node_info[5] == ColorRender("Fail", COLOR_ERROR){
                                        chap_node_info[6] = ColorRender("UnKnown", COLOR_WARNNING)
                                }else{
					kp, _, err := kv.Get(role_key, nil)
                                        if err != nil {
                                                fmt.Printf("%s\n", err)
                                        }
                                        if kp.Value != nil {
                                                chap_node_info[6] = string(kp.Value)
                                        } else {
                                                chap_node_info[6] = ColorRender("UnKnown", COLOR_WARNNING)
                                        }
				}
				//add vip
				if chap_node_info[6] == "master" {
					kp, _, err := kv.Get("cmha/service/"+args[0]+"/chap/VIP", nil)
					if err != nil {
						fmt.Printf("%s\n", err)
					}
					if kp.Value != nil {
						chap_node_info[7] = string(kp.Value) + node_port
					} else {
						chap_node_info[7] = ColorRender("UnKnown", COLOR_WARNNING)
					}
                                }else {
					chap_node_info[7] = ""
				}
				if chap_node_info[5] == ColorRender("Fail", COLOR_ERROR){
                                        chap_node_info[7] = ""
                                }
			}

			if has_unknown_status {
				chap_node_info[4] = ColorRender("UnKnown", COLOR_WARNNING)
			}

			_data_chap = append(_data_chap, chap_node_info)

		} else {
			// Node Type : DB
			db_node_info := make([]string, 7)
			db_node_info[0] = servicename
			db_node_info[1] = node_name
			db_node_info[2] = node_address + node_port
			db_node_info[3] = node_is_leader
			db_node_info[4] = node_type


			has_unknown_status := false

			for _, checkdata := range servicedata.Checks {

				node_status := checkdata.Status

				node_checkid := checkdata.CheckID

				has_unknown_status = node_status == "critical" && node_checkid == "serfHealth"

				switch node_status {
				case "passing":
					node_status = ColorRender("OK", COLOR_NORMAL)
				case "critical":
					node_status = ColorRender("Fail", COLOR_ERROR)
				default:
					node_status = ColorRender(node_status, COLOR_WARNNING)
				}

				if node_checkid == "serfHealth" {
					db_node_info[6] = node_status
				} else {
					if db_node_info[6] == ColorRender("Fail", COLOR_ERROR){
                                                db_node_info[5] = ColorRender("UnKnown", COLOR_WARNNING)
                                                
                                        }else{
						db_node_info[5] = node_status
					}
				}

			}

			if has_unknown_status {
				db_node_info[5] = ColorRender("UnKnown", COLOR_WARNNING)
			}

			//Get Data From DB
			data_from_db, err := ArrangeServiceData(_node_address, fmt.Sprint(_node_port),db_node_info[6])
			if err != nil {
				fmt.Printf("%s\n", err)
			}

			// repl_err_count
			repl_err_count, err := getReplErrCount(_node_name, servicename)

			if err != nil {
				fmt.Println(err.Error())
			}

			db_node_info = append(db_node_info, data_from_db...)
			db_node_info = append(db_node_info, repl_err_count)

			_data_db = append(_data_db, db_node_info)

		}

	}

	_th_chap := []string{
		"service",
		"node",
		"address",
		"type",
		"status",
		"agent",
		"role",
		"vip",
	}

	TableRender(_th_chap, _data_chap, ALIGN_CENTRE)
	fmt.Println("")

	_th_db := []string{
		"service",
		"node",
		"address",
		"leader",
		"type",
		"status",
		"agent",
		"r/w",
		"vsr",
		"repl-status",
		"repl-err-counter",
	}

	TableRender(_th_db, _data_db, ALIGN_CENTRE)
	fmt.Println("")

	return nil
}

func ArrangeServiceData(_ip string, _port string,agentstatus string) ([]string, error) {
	// Data from DB
	_sql_DB_Variables := "show variables"

	var ret_err_string string = ""
	_filter1 := "read_only"
	_filter2 := "rpl_semi_sync_master_trysyncrepl"
	_filter3 := "rpl_semi_sync_master_keepsyncrepl"

	_sql_Slave_Status := "show slave status"
	var count int
	// filter for slave status
	_filter4 := "Slave_IO_Running"
	_filter5 := "Slave_SQL_Running"
	var _slave_status map[string]string
	var err_get_slave_status error
	_service_data, err_get_service_data :=
		ServiceData(_ip,
			_port,
			_sql_DB_Variables,
			_filter1,
			_filter2,
			_filter3)
	
	if err_get_service_data != nil {
		count +=1
		ret_err_string += err_get_service_data.Error()

	}else{
		_slave_status, err_get_slave_status =
		SlaveStatus(_ip,
			_port,
			_sql_Slave_Status,
			_filter4,
			_filter5)

		if err_get_slave_status != nil {

			ret_err_string += err_get_slave_status.Error()

		}
	}
	var ret_err error = nil
	if ret_err_string != "" {
		ret_err = errors.New(ret_err_string)
	}

	var _rw, _vsr, _repl_status = "", "", ""

	// Read only
	if _service_data[_filter1] == "OFF" {
		_rw = "rw"
	} else if _service_data[_filter1] == "ON" {
		_rw = "ro"
	} else {
		_rw = _service_data[_filter1]
	}

	// VSR
	if _service_data[_filter2] == "ON" &&
		_service_data[_filter3] == "ON" {
		_vsr = "ON"
	} else {
		_vsr = "OFF"
	}

	//repl-status
	if _slave_status[_filter4] == "Yes" &&
		_slave_status[_filter5] == "Yes" {

		_repl_status = ColorRender("OK", COLOR_NORMAL)

	} else {

		_status_io := _slave_status[_filter4]
		_status_sql := _slave_status[_filter5]

		var _rt_status string = ""

		switch _status_io {
		case "Yes":
			_rt_status = ColorRender("OK", COLOR_NORMAL)
		case "No":
			_rt_status = ColorRender("No", COLOR_WARNNING)
		default:
			_rt_status = ColorRender(_status_io, COLOR_WARNNING)

		}

		_status_io = _rt_status

		switch _status_sql {
		case "Yes":
			_rt_status = ColorRender("Yes", COLOR_NORMAL)
		case "No":
			_rt_status = ColorRender("No", COLOR_WARNNING)
		default:
			_rt_status = ColorRender(_status_sql, COLOR_WARNNING)
		}

		_status_sql = _rt_status

		_repl_status = fmt.Sprintf("IO:%s;SQL:%s", _status_io, _status_sql)

	}
	if count != 0 {
		_repl_status = ColorRender("UnKnown", COLOR_WARNNING)
	}
	_ret_data := []string{_rw, _vsr, _repl_status}

	return _ret_data, ret_err

}

func getReplErrCount(_node_name, servicename string) (string, error) {

	var repl_err_count string = ""

	_key_repl_err_counter := "cmha/service/" + servicename + "/db/repl_err_counter/" + _node_name

	_client, err := cliconfig.Consul_Client_Init()

	if err != nil {
		return "", err
	}

	_kv := _client.KV()

	if _repl_err_counter, _, err := _kv.Get(_key_repl_err_counter, nil); err != nil {

		return "", err

	} else {
		if _repl_err_counter != nil {
			repl_err_count = string(_repl_err_counter.Value)

			if repl_err_count == "0" {
				repl_err_count = ColorRender(repl_err_count, COLOR_NORMAL)
			} else if repl_err_count == "1" {
				repl_err_count = ColorRender(repl_err_count, COLOR_ERROR)
			}
		}

	}

	return repl_err_count, nil

}
