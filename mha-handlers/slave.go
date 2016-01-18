package main

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"strconv"
)

func slave(ip, port, username, password string) {
	var timestamp int64
	dsName := username + ":" + password + "@tcp(" + "localhost" + ":" + port + ")/"
	db, err := sql.Open("mysql", dsName)
	if err != nil {
		logger.Println("[E] Create connection object to local database failed!",err)
		timestamp = time.Now().Unix()
		logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + create_database_object_failed + "{{" + fmt.Sprintf("%s",err)
		logger.Println("[I] Give up leader election")
		logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + give_election
		UploadLog(logkey,logvalue)
		return
	}
	logger.Println("[I] Create connection object to local database successfully!")
	timestamp = time.Now().Unix()
        logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + create_database_object_success
	defer db.Close()
	err = db.Ping()
	if err != nil {
		logger.Println("[E] Connected to local database failed!",err)
		timestamp = time.Now().Unix()
                logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + connected_database_failed + "{{" + fmt.Sprintf("%s",err)
		logger.Println("[I] Give up leader election")
		logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + give_election
		UploadLog(logkey,logvalue)
		return
	}
	logger.Println("[I] Connected to local database successfully!")
	timestamp = time.Now().Unix()
        logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + connected_database_success
	_, err = db.Query("stop slave io_thread")
	if err != nil {
		logger.Println("[E] Stop local database replication I/O thread failed!",err)
		timestamp = time.Now().Unix()
	        logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + stop_repl_io_failed + "{{" + fmt.Sprintf("%s",err)
		logger.Println("[I] Give up leader election")
		logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + give_election
		UploadLog(logkey,logvalue)
		return
	}
	logger.Println("[I] Stop local database replication I/O thread successfully!")
	timestamp = time.Now().Unix()
        logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + stop_repl_io_success
	row, err := db.Query("show slave status")
	if err != nil {
		logger.Println("[E] Checking local database SQL thread status. Failed!",err)
		timestamp = time.Now().Unix()
                logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + check_io_failed + "{{" + fmt.Sprintf("%s",err)
		logger.Println("[I] Give up leader election")
		logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + give_election
		UploadLog(logkey,logvalue)
		return
	}
	logger.Println("[I] Checking local database SQL thread status. Succeed!")
	timestamp = time.Now().Unix()
        logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + check_io_success
	cols, _ := row.Columns()
	buffer := make([]interface{}, len(cols))
	data := make([]interface{}, len(cols))
	for i, _ := range buffer {
		buffer[i] = &data[i]
	}
	for row.Next() {
		err = row.Scan(buffer...)
		if err != nil {
			logger.Println("[E] Resolve slave status failed!",err)
			timestamp = time.Now().Unix()
	                logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + resolve_slave_status_failed + "{{" + fmt.Sprintf("%s",err)
			logger.Println("[I] Give up leader election")
			logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + give_election
			UploadLog(logkey,logvalue)
			return
		}
	}
	mapField2Data := make(map[string]interface{}, len(cols))
	for k, col := range data {
		mapField2Data[cols[k]] = col
	}
	Master_Log_File := mapField2Data["Master_Log_File"]
	Read_Master_Log_Pos := mapField2Data["Read_Master_Log_Pos"]
	Slave_SQL_Running := mapField2Data["Slave_SQL_Running"]
	if string(Slave_SQL_Running.([]uint8)) != "Yes" {
		logger.Println("[E] The SQL thread status is " + string(Slave_SQL_Running.([]uint8)) + "!")
		timestamp = time.Now().Unix()
                logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + sql_status_noyes + "{{" + string(Slave_SQL_Running.([]uint8))
		logger.Println("[I] Give up leader election")
		logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + give_election
		UploadLog(logkey,logvalue)
		return
	}
	logger.Println("[I] The SQL thread status is Yes!")
	timestamp = time.Now().Unix()
        logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + sql_status_yes
	sqlstr := "select master_pos_wait(?,?)"
	rowss, err := db.Query(sqlstr, Master_Log_File, Read_Master_Log_Pos)
	if err != nil {
		logger.Println("[E] Checking relay log applying status failed!",err)
		timestamp = time.Now().Unix()
	        logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + exec_master_pos_wait_failed + "{{" + fmt.Sprintf("%s",err)
		logger.Println("[I] Give up leader election")
		logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + give_election
		UploadLog(logkey,logvalue)
		return
	}
	logger.Println("[I] Checking relay log applying status successfully!")
	timestamp = time.Now().Unix()
        logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + exec_master_pos_wait_success
	var master_pos_wait string
	for rowss.Next() {
		err = rowss.Scan(&master_pos_wait)
		if err != nil {
			logger.Println("[E] Resolve master_pos_wait failed!",err)
			timestamp = time.Now().Unix()
		        logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + resolve_master_pos_wait_failed + "{{" + fmt.Sprintf("%s",err)
			logger.Println("[I] Give up leader election")
			logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + give_election
			UploadLog(logkey,logvalue)
			return
		}
		if master_pos_wait < "0" && master_pos_wait == "null" {
			logger.Println("[E] Relay log applying failed!")
			timestamp = time.Now().Unix()
                        logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + relay_log_apply_failed
			logger.Println("[I] Give up leader election")
			logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + give_election
			UploadLog(logkey,logvalue)
			return
		}
		logger.Println("[I] Relay log applying completed!")
		timestamp = time.Now().Unix()
                logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + relay_log_apply_completed
		_, err := db.Query("set global rpl_semi_sync_master_keepsyncrepl=0")
		if err != nil {
			logger.Println("[E] Set rpl_semi_sync_master_keepsyncrepl=0 failed!",err)
			timestamp = time.Now().Unix()
	                logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + set_keepsyncrepl_failed + "{{" + fmt.Sprintf("%s",err)
			logger.Println("[I] Give up leader election")
			logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + give_election
			UploadLog(logkey,logvalue)
			return
		}
		logger.Println("[I] Set rpl_semi_sync_master_keepsyncrepl=0 successfully!")
		timestamp = time.Now().Unix()
                logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + set_keepsyncrepl_success
		_, err = db.Query("set global rpl_semi_sync_master_trysyncrepl=0")
		if err != nil {
			logger.Println("[E] Set rpl_semi_sync_master_trysyncrepl=0 failed!",err)
			timestamp = time.Now().Unix()
	                logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + set_trysyncrepl_failed + "{{" + fmt.Sprintf("%s",err)
			logger.Println("[I] Give up leader election")
			logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + give_election
			UploadLog(logkey,logvalue)
			return
		}
		logger.Println("[I] Set rpl_semi_sync_master_trysyncrepl=0 successfully!")
		timestamp = time.Now().Unix()
                logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + set_trysyncrepl_success
		logger.Println("[I] Switching local database to async replication!")
		logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + switch_local_async_repl
		_, err = db.Query("set global read_only=0")
		if err != nil {
			logger.Println("[E] Set local database Read/Write mode failed!",err)
			timestamp = time.Now().Unix()
                        logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + set_read_write_failed + "{{" + fmt.Sprintf("%s",err)
			logger.Println("[I] Give up leader election")
			logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + give_election
			UploadLog(logkey,logvalue)
			return
		}
		logger.Println("[I] Set local database Read/Write mode successfully!")
		timestamp = time.Now().Unix()
                logvalue = logvalue +"|" + strconv.FormatInt(timestamp, 10) + set_read_write_success
	}
	SetConn(ip, port, username, password)
}
