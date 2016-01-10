package main

import (
	//"strings"
	"strconv"
	"time"
	//"fmt"
	"github.com/astaxie/beego"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	consulapi "github.com/hashicorp/consul/api"
)

var service_ip []string
var servicename string
var hostname string
var port string
var username string
var password string
var kv *consulapi.KV
var client *consulapi.Client
var logvalue string
var logkey string
const (
        triggered = "001"
        consulapi_failed = "002"
        give_election = "003"
        completed = "004"
        consulapi_success = "005"
        current_check_failed = "006"
        current_check_success = "007"
        get_counter_failed = "008"
        get_counter_success = "009"
        ip_election = "010"
        create_counter = "011"
        leader_exist = "012"
        leader_noexist = "013"
        get_health_failed = "014"
        get_health_success = "015"
        service_noexist = "016"
        service_exist = "017"
        ip_noservice = "018"
        clean_kv_failed = "019"
        clean_kv_success = "020"
        status_critical = "021"
        status_nocritical = "022"
        create_session_failed = "023"
        create_session_success ="024"
        format_json_hap_failed = "025"
        send_leader_failed = "026"
        send_leader_success = "027"
        becoming_string_failed = "028"
        becoming_string_success = "029"
        set_counter_failed = "030"
        set_counter_success = "031"
        create_database_object_failed = "032"
        create_database_object_success = "033"
        connected_database_failed = "034"
        connected_database_success = "035"
        set_read_only_failed = "036"
        database_downgrade_failed = "037"
        set_read_only_success = "038"
        database_downgrade_success = "039"
        stop_repl_io_failed = "040"
        stop_repl_io_success = "041"
        check_io_failed = "042"
        check_io_success = "043"
        resolve_slave_status_failed = "044"
        sql_status_noyes = "045"
        sql_status_yes = "046"
        exec_master_pos_wait_failed = "047"
        exec_master_pos_wait_success = "048"
        resolve_master_pos_wait_failed = "049"
        relay_log_apply_failed = "050"
        relay_log_apply_completed = "051"
        set_keepsyncrepl_failed = "052"
        set_keepsyncrepl_success = "053"
        set_trysyncrepl_failed = "054"
        set_trysyncrepl_success = "055"
        switch_local_async_repl = "056"
        set_read_write_failed = "057"
        set_read_write_success = "058"
 )
func SessionAndChecks() {
//	beego.Info("MHA Handler Triggered")
	logger.Println("[I] MHA Handler Triggered")
	timestamp := time.Now().Unix()
	logvalue = strconv.FormatInt(timestamp, 10) + triggered
	ip := beego.AppConfig.String("ip")
//	Switch := beego.AppConfig.String("switch")
//	if strings.EqualFold(Switch, "off") {
//		beego.Info(ip +" switch="+Switch+",give up leader election!")
//		beego.Info("MHA Handler Completed")
//		return
//	} else if strings.EqualFold(Switch, "on") {
//		beego.Info(ip + " switch=" +Switch)
//		beego.Info("Begin leader election!")
//	} else {
//		beego.Info("Config file switch format error,switch="+Switch+",Should off or on!")
//		beego.Info("Give up leader election")
//		beego.Info("MHA Handler Completed")
//		return
//	}
	service_ip = beego.AppConfig.Strings("service_ip")
	servicename = beego.AppConfig.String("servicename")
	hostname = beego.AppConfig.String("hostname")
	port = beego.AppConfig.String("port")
	username = beego.AppConfig.String("username")
	password = beego.AppConfig.String("password")
	logkey = "mha-handlers/"+ hostname + "/" + strconv.FormatInt(timestamp, 10)
	config := &consulapi.Config{
		Datacenter: beego.AppConfig.String("datacenter"),
		Token:      beego.AppConfig.String("token"),
	}
	var kvPair *consulapi.KVPair
	var kvMonitor *consulapi.KVPair
//	var client *consulapi.Client
//	var kv *consulapi.KV
	var err error
	for i, _ := range service_ip {
		config.Address = service_ip[i] + ":" + beego.AppConfig.String("service_port")
		client, err = consulapi.NewClient(config)
		if err != nil {
//			beego.Error("Create consul-api client failed!", err)
//			beego.Info("Give up leader election")
//	                beego.Info("MHA Handler Completed")
			logger.Println("[E] Create consul-api client failed!",err)
			timestamp := time.Now().Unix()
        		logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + consulapi_failed
			logger.Println("[I] Give up leader election")
			logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + give_election
			logger.Println("[I] MHA Handler Completed")
			logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + completed
			UploadLog(logkey,logvalue)
			return
		}
	//	beego.Info("Create consul-api client successfully!")
		logger.Println("[I] Create consul-api client successfully!")
		timestamp := time.Now().Unix()
               	logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + consulapi_success
//		fmt.Println(logvalue)
		//KV is used to return a handle to the K/V apis
		kv = client.KV()
		//Get is used to lookup a single key
		kvPair, _, err = kv.Get("service/"+servicename+"/leader", nil)
		if err != nil {
//			beego.Error("Get and check current service leader from CS failed!", err)
			logger.Println("[E] Get and check current service leader from CS failed!",err)
			timestamp := time.Now().Unix()
	                logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + current_check_failed
			continue
		}
		break
	}
//	beego.Info("Get and check current service leader from CS successfully!")
	logger.Println("[I] Get and check current service leader from CS successfully!")
	timestamp = time.Now().Unix()
        logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + current_check_success
	kvMonitor, _, err = kv.Get("monitor/"+hostname, nil)
	if err != nil {
//		beego.Error("Get " + ip + "repl_err_counter=0/1 failed", err)
		logger.Println("[E] Get " + ip + " repl_err_counter=0/1 failed",err)
		timestamp := time.Now().Unix()
        	logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + get_counter_failed
//		beego.Info("Give up leader election")
		logger.Println("[I] Give up leader election")
		logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + give_election
//                beego.Info("MHA Handler Completed")
		logger.Println("[I] MHA Handler Completed")
		logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + completed
		UploadLog(logkey,logvalue)
		return
	}
	logger.Println("[E] Get " + ip + " repl_err_counter=0/1 successfully!")
	timestamp = time.Now().Unix()
        logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + get_counter_success
	kvValue :=string(kvMonitor.Value)
	if kvValue != "0" {
//		beego.Error(ip + " give up leader election")	
		logger.Println("[E] " + ip + " give up leader election")
		timestamp = time.Now().Unix()
	        logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + ip_election
//              beego.Info("MHA Handler Completed")
		logger.Println("[I] MHA Handler Completed")
		logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + completed
		UploadLog(logkey,logvalue)
		return
	}
	//NewClient returns a new client
	if kvPair == nil {
//		beego.Error("Not replication counter,Please create!")
		logger.Println("[E] Not replication counter,Please create!")
		timestamp = time.Now().Unix()
                logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + create_counter
//		beego.Info("Give up leader election")
		logger.Println("[I] Give up leader election")
		logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + give_election
//                beego.Info("MHA Handler Completed")
		logger.Println("[I] MHA Handler Completed")
		logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + completed
		UploadLog(logkey,logvalue)
		return
	}
	//Are there external connection string provided
	if kvPair.Session != "" {
//		beego.Info("Leader exist!")
		logger.Println("[I] Leader exist!")
		timestamp = time.Now().Unix()
                logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + leader_exist
//		time.Sleep(1 * time.Second)
//		beego.Info("Give up leader election")
//                beego.Info("MHA Handler Completed")
		logger.Println("[I] Give up leader election")
		logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + give_election
		logger.Println("[I] MHA Handler Completed")
		logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + completed
		UploadLog(logkey,logvalue)
		return
	}
//	beego.Info("Leader does not exist!")
	logger.Println("[I] Leader does not exist!")
	timestamp = time.Now().Unix()
        logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + leader_noexist
/*	dsName := username + ":" + password + "@tcp(" + "localhost" + ":" + port + ")/"
        db, err := sql.Open("mysql", dsName)
        if err != nil {
                beego.Error("Create connection object to local database failed!", err)
                beego.Info("Give up leader election")
                beego.Info("MHA Handler Completed")
                return
        }
        beego.Info("Create connection object to local database successfully!")
        defer db.Close()
        err = db.Ping()
        if err != nil {
                beego.Error("Connected to local database failed!", err)
                beego.Info("Give up leader election")
                beego.Info("MHA Handler Completed")
                return
        }
        beego.Info("Connected to local database successfully!")
	_, err = db.Query("set global read_only=0")
        if err != nil {
 	       beego.Error("Set local database Read/Write mode failed!", err)
               beego.Info("Give up leader election")
               beego.Info("MHA Handler Completed")
               return
        }
        beego.Info("Set local database Read/Write mode successfully!")*/
	SetRead_only(username,password,port,1)
	//Health returns a handle to the health endpoints
	health := client.Health()
	//Checks is used to return the checks associated with a service
	healthvalue, _, err := health.Checks(servicename, nil)
	if err != nil {
//		beego.Error("Get and check "+ ip + " service health status failed!", err)
//		beego.Info("Give up leader election")
//                beego.Info("MHA Handler Completed")
		logger.Println("[E] Get and check " + ip + " service health status failed!",err)
		timestamp = time.Now().Unix()
                logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + get_health_failed
		logger.Println("[I] Give up leader election")
		logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + give_election
		logger.Println("[I] MHA Handler Completed")
		logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + completed
		UploadLog(logkey,logvalue)
		return
	}
//	beego.Info("Get and check "+ ip + " service health status successfully!")
	logger.Println("[I] Get and check "+ ip + " service health status successfully!")
	timestamp = time.Now().Unix()
        logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + get_health_success
	if len(healthvalue) <= 0 {
//		beego.Info(servicename + " service does not exist!")
//		beego.Info("Give up leader election")
//                beego.Info("MHA Handler Completed")
		logger.Println("[I] " + servicename + " service does not exist!")
		timestamp = time.Now().Unix()
		logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + service_noexist
		logger.Println("[I] Give up leader election")
		logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + give_election
		logger.Println("[I] MHA Handler Completed")
		logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + completed
		UploadLog(logkey,logvalue)
		return
	}
	var islocal bool
	for index := range healthvalue {
		if healthvalue[index].Node == hostname {
			islocal = true
	//		beego.Info(servicename + " service exist!")
			logger.Println("[I] " + servicename + " service exist!")
			timestamp = time.Now().Unix()
	                logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + service_exist
			break
		}

	}
	if !islocal {
//		beego.Info(ip + " not is " +servicename + "!")
		logger.Println("[E] " + ip + " not is " +servicename + "!")
		timestamp = time.Now().Unix()
                logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + ip_noservice
//		beego.Info("Give up leader election")
		logger.Println("[I] Give up leader election")
		logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + give_election
                logger.Println("[I] MHA Handler Completed")
		logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + completed
		UploadLog(logkey,logvalue)
		return
	} else {
		updatevalue := consulapi.KVPair{
			Key:   "service/" + servicename + "/leader",
			Value: []byte(""),
		}
		_, err = kv.Put(&updatevalue, nil)
		if err != nil {
//			beego.Error("Clean service leader value in CS failed!", err)
			logger.Println("[E] Clean service leader value in CS failed!",err)
			timestamp = time.Now().Unix()
	                logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + clean_kv_failed
//			beego.Info("Give up leader election")
//	                beego.Info("MHA Handler Completed")
			logger.Println("[I] Give up leader election")
			logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + give_election
	                logger.Println("[I] MHA Handler Completed")
			logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + completed
			UploadLog(logkey,logvalue)
			return
		}
//		beego.Info("Clean service leader value in CS successfully!")
		logger.Println("[I] Clean service leader value in CS successfully!")
		timestamp = time.Now().Unix()
                logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + clean_kv_success
		healthpair, _, err := health.Service(servicename, "", false, nil)
		if err != nil {
//			beego.Error("Get and check " + ip + " service health status failed!", err)
//			beego.Info("Give up leader election")
//	                beego.Info("MHA Handler Completed")
			logger.Println("[E] Get and check " + ip + " service health status failed!",err)
			timestamp = time.Now().Unix()
	                logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + get_health_failed
			logger.Println("[I] Give up leader election")
			logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + give_election
                        logger.Println("[I] MHA Handler Completed")
			logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + completed
			UploadLog(logkey,logvalue)
			return
		}
//		beego.Info("Get and check " + ip + " service health status successfully!")
		logger.Println("[I] Get and check " + ip + " service health status successfully!")
		timestamp = time.Now().Unix()
                logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + get_health_success
		var ishealthy = true
		hostname := beego.AppConfig.String("hostname")
		for index := range healthpair {
			for checkindex := range healthpair[index].Checks {
				if healthpair[index].Checks[checkindex].Node == hostname {
					if healthpair[index].Checks[checkindex].Status == "critical" {
						ishealthy = false
						break
					}
				}
			}
		}
		if !ishealthy {
//			beego.Info("Status is critical!")
//			beego.Info("Give up leader election")
//	                beego.Info("MHA Handler Completed")
			logger.Println("[E] Status is critical!")
			timestamp = time.Now().Unix()
	                logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + status_critical
			logger.Println("[I] Give up leader election")
			logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + give_election
                        logger.Println("[I] MHA Handler Completed")
			logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + completed
			UploadLog(logkey,logvalue)
			return
		} else {
//			beego.Info("Status is not critical")
			logger.Println("[I] Status is not critical")
			timestamp = time.Now().Unix()
                        logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + status_nocritical
			slave(ip, port, username, password)
		}
	}
}
func SetConn(ip, port, username, password string) {
	config := &consulapi.Config{
		Datacenter: beego.AppConfig.String("datacenter"),
		Token:      beego.AppConfig.String("token"),
	}
	var client *consulapi.Client
	var err error
	var sessionvalue string
	var timestamp int64
	for i, _ := range service_ip {
		config.Address = service_ip[i] + ":" + beego.AppConfig.String("service_port")
		client, err = consulapi.NewClient(config)
		if err != nil {
//			beego.Error("Create consul-api client failed! ", err)
			logger.Println("[E] Create consul-api client failed!",err)
			timestamp = time.Now().Unix()
                        logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + consulapi_failed
//			beego.Info("Give up leader election")
//                	beego.Info("MHA Handler Completed")
			logger.Println("[I] Give up leader election")
			logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + give_election
                        logger.Println("[I] MHA Handler Completed")
			logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + completed
			UploadLog(logkey,logvalue)
			return
		}
//		beego.Info("Create  consul-api client successfully!")
		logger.Println("[I] Create consul-api client successfully!")
		timestamp = time.Now().Unix()
                logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + consulapi_success
		session := client.Session()
		sessionEntry := consulapi.SessionEntry{
			LockDelay: 10 * time.Second,
			Name:      servicename,
			Node:      hostname,
			Checks:    []string{"serfHealth", "service:" + servicename},
		}
		//Create makes a new session. Providing a session entry can customize the session. It can also be nil to use defaults.
		sessionvalue, _, err = session.Create(&sessionEntry, nil)
		if err != nil {
//			beego.Error("Session create failed!", err)
			logger.Println("[E] Session create failed!",err)
			timestamp = time.Now().Unix()
	                logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + create_session_failed
			continue
		}
		break
	}
	//NewClient returns a new client
	beego.Info("Session create successfully!")
	logger.Println("[I] Session create successfully!")
	timestamp = time.Now().Unix()
        logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + create_session_success
	format := beego.AppConfig.String("format")
	var acquirejson string
	if format == "json" {
		acquirejson = `{"Node":"` + hostname + `","Ip":"` + ip + `","Port":` + port + `,"Username":"` + username + `","Password":"` + password + `"}`
	} else if format == "hap" {
		acquirejson = "server" + " " + hostname + " " + ip + ":" + port
	} else {
//		beego.Error("format error,json or hap!")
//		beego.Info("Give up leader election")
//                beego.Info("MHA Handler Completed")
		logger.Println("[E] format error,json or hap!")
		timestamp = time.Now().Unix()
	        logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + format_json_hap_failed
		logger.Println("[I] Give up leader election")
		logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + give_election
              	logger.Println("[I] MHA Handler Completed")	
		logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + completed
		UploadLog(logkey,logvalue)
		return
	}
	value := []byte(acquirejson)
	kv = client.KV()
	kvpair := consulapi.KVPair{
		Key:     "service/" + servicename + "/leader",
		Value:   value,
		Session: sessionvalue,
	}
	//Acquire is used for a lock acquisiiton operation. The Key, Flags, Value and Session are respected. Returns true on success or false on failures.
	time.Sleep(15 * time.Second)
	var ok bool
	ok, _, err = kv.Acquire(&kvpair, nil)
	if err != nil {
//		beego.Error("Send service leader request to CS failed! ", err)
//                beego.Info("MHA Handler Completed")
		logger.Println("[E] Send service leader request to CS failed!",err)
		timestamp = time.Now().Unix()
                logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + send_leader_failed
		logger.Println("[I] MHA Handler Completed")
		logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + completed
		UploadLog(logkey,logvalue)
		return
	}
//	beego.Info("Send service leader request to CS successfully!")
	logger.Println("[I] Send service leader request to CS successfully!")
	timestamp = time.Now().Unix()
        logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + send_leader_success
	if !ok {
//		time.Sleep(5 * time.Second)
		beego.Info("Becoming service leader failed! Connection string is " + ip + " " + port)
		logger.Println("[E] Becoming service leader failed! Connection string is " + ip + " " + port)
		timestamp = time.Now().Unix()
	        logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + becoming_string_failed
		SetRead_only(username,password,port,1)
//                beego.Info("Monitor Handler Completed")
		logger.Println("[I] MHA Handler Completed")
		logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + completed
		UploadLog(logkey,logvalue)
                return
	} else {
//		beego.Info("Becoming service leader successfully! Connection string is " + ip + " " + port)
		logger.Println("[I] Becoming service leader successfully! Connection string is " + ip + " " + port)
		timestamp = time.Now().Unix()
                logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + becoming_string_success
	//	var put string
		other_hostname := beego.AppConfig.String("otherhostname")
		SetRepl_err_counter(other_hostname)
        /*        put = "1"
                kvvalue := []byte(put)
                kvotherhostname := consulapi.KVPair{
      	        	Key:   "monitor/" + other_hostname,
                        Value: kvvalue,
                }
                _, err = kv.Put(&kvotherhostname, nil)
                if err != nil {
//                	beego.Error("Set peer database repl_err_counter to 1 in CS failed!", err)
//                        beego.Info("Monitor Handler Completed")
			logger.Println("[E] Set peer database repl_err_counter to 1 in CS failed!",err)
			timestamp = time.Now().Unix()
	                logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + set_counter_failed
			logger.Println("[I] MHA Handler Completed")
			logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + completed
			UploadLog(logkey,logvalue)
                        return
                }
//                beego.Info("Set peer database repl_err_counter to 1 in CS successfully!")
		logger.Println("[I] Set peer database repl_err_counter to 1 in CS successfully!")
		timestamp = time.Now().Unix()
                logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + set_counter_success
                beego.Info("MHA Handler Completed")*/
	}
}

func SetRepl_err_counter(hostname string){
	count := 0
	var put string
	var timestamp int64
//        other_hostname := beego.AppConfig.String("otherhostname")
        put = "1"
        kvvalue := []byte(put)
      	kvotherhostname := consulapi.KVPair{
        	Key:   "monitor/" + hostname,
                Value: kvvalue,
        }
   try:  _, err := kv.Put(&kvotherhostname, nil)
        if err != nil {
//         	beego.Error("Set peer database repl_err_counter to 1 in CS failed!", err)
		logger.Println("[E] Set peer database repl_err_counter to 1 in CS failed!",err)
		timestamp = time.Now().Unix()
                logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + set_counter_failed
		if count ==2 {
//			beego.Info("Monitor Handler Completed")
			logger.Println("[I] MHA Handler Completed")
			timestamp = time.Now().Unix()
	                logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + completed
			UploadLog(logkey,logvalue)
			return
		}
		count++
                goto try
	}
//        beego.Info("Set peer database repl_err_counter to 1 in CS successfully!")
	logger.Println("[I] Set peer database repl_err_counter to 1 in CS successfully!")
	timestamp = time.Now().Unix()
        logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + set_counter_success
//        beego.Info("MHA Handler Completed")
	logger.Println("[I] MHA Handler Completed")
	logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + completed
}

func SetRead_only(username,password,port string,value int){
	var timestamp int64
	dsName := username + ":" + password + "@tcp(" + "localhost" + ":" + port + ")/"
        db, err := sql.Open("mysql", dsName)
        if err != nil {
//                beego.Error("Create connection object to local database failed!", err)
//                beego.Info("Give up leader election")
//                beego.Info("MHA Handler Completed")
		logger.Println("[E] Create connection object to local database failed!",err)
		timestamp = time.Now().Unix()
	        logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + create_database_object_failed
		logger.Println("[I] Give up leader election")
		logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + give_election
                logger.Println("[I] MHA Handler Completed")
		logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + completed
		UploadLog(logkey,logvalue)	
                return
        }
//        beego.Info("Create connection object to local database successfully!")
	logger.Println("[I] Create connection object to local database successfully!")
	timestamp = time.Now().Unix()
        logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + create_database_object_success		
        defer db.Close()
        err = db.Ping()
        if err != nil {
//                beego.Error("Connected to local database failed!", err)
//                beego.Info("Give up leader election")
//                beego.Info("MHA Handler Completed")
		logger.Println("[E] Connected to local database failed!",err)
		timestamp = time.Now().Unix()                 
		logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + connected_database_failed
		logger.Println("[I] Give up leader election")
		logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + give_election
                logger.Println("[I] MHA Handler Completed")
		logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + completed
		UploadLog(logkey,logvalue)
                return
        }
//        beego.Info("Connected to local database successfully!")
	logger.Println("[I] Connected to local database successfully!")
	timestamp = time.Now().Unix()                                 
	logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + connected_database_success
	read_only := "set global read_only=" + strconv.Itoa(value)
        _, err = db.Query(read_only)
        if err != nil {
//               	beego.Error("Set local database Read_only mode failed!", err)
//		beego.Info("Local database downgrade failed!")
//               	beego.Info("Give up leader election")
//               	beego.Info("MHA Handler Completed")
		logger.Println("[E] Set local database Read_only mode failed!",err)
		timestamp = time.Now().Unix()                                 
		logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + set_read_only_failed
		logger.Println("[I] Local database downgrade failed!")
		logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + database_downgrade_failed
		logger.Println("[I] Give up leader election")
		logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + give_election
                logger.Println("[I] MHA Handler Completed")
		logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + completed
		UploadLog(logkey,logvalue)
               	return
        }
//        beego.Info("Set local database Read_only mode successfully!")
//	beego.Info("Local database downgrade successfully!")
	logger.Println("[I] Set local database Read_only mode successfully!")
	timestamp = time.Now().Unix()                                           
	logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + set_read_only_success
	logger.Println("[I] Local database downgrade successfully!")
	logvalue = logvalue +"|"+ strconv.FormatInt(timestamp, 10) + database_downgrade_success
}

func UploadLog(logkey,logvalue string){
	kvhostname := consulapi.KVPair{
        	Key:   logkey,
                Value: []byte(logvalue),
        }
        _, err := kv.Put(&kvhostname, nil)
     	if err != nil {
     		logger.Println("[E] Set peer database repl_err_counter to 1 in CS failed!",err)
                logger.Println("[I] MHA Handler Completed")
                return
       	}
	logger.Println("[E] Set peer database repl_err_counter to 1 in CS success!")
}
