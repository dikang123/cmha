package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func checkio_thread(ip, port, username, password, addr string,logvalue string) string{
	logger.Println("[I] Monitor Handler Sleep 60s!")
	timestamp := time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + handlers_sleep
	time.Sleep(60000 * time.Millisecond)
	logger.Println("[I] Connecting to peer database......")
	timestamp = time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + connecting_database
	dsName := username + ":" + password + "@tcp(" + addr + ":" + port + ")/"
	db, err := sql.Open("mysql", dsName)
	if err != nil {
		logger.Println("[E] Create connection object to peer database failed!", err)
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + create_peer_database_object_failed + "{{" + fmt.Sprintf("%s", err)
		logger.Println("[I] Give up switching to async replication!")
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + give_async_repl
		return logvalue
	}
	logger.Println("[I] Create connection object to peer database successfully!")
	timestamp = time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + create_peer_database_object_success
	defer db.Close()
	err = db.Ping()
	if err != nil {
		logger.Println("[E] Connected to the peer database failed!", err)
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + connected_peer_database_failed + "{{" + fmt.Sprintf("%s", err)
		logvalue =Set(ip, port, username, password, 0,logvalue)
		return logvalue
	}
	logger.Println("[I] Connected to the peer database successfully!")
	timestamp = time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + connected_peer_database_success
	row, err := db.Query("show slave status")
	if err != nil {
		logger.Println("[I] Checking peer database I/O thread status. Failed!", err)
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + check_perr_dabases_io_failed + "{{" + fmt.Sprintf("%s", err)
		logvalue = Set(ip, port, username, password, 0,logvalue)
		return logvalue
	}
	logger.Println("[I] Checking peer database I/O thread status. Successfully!")
	timestamp = time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + check_perr_dabases_io_success
	cols, _ := row.Columns()
	buffer := make([]interface{}, len(cols))
	data := make([]interface{}, len(cols))
	for i, _ := range buffer {
		buffer[i] = &data[i]
	}
	for row.Next() {
		err = row.Scan(buffer...)
		if err != nil {
			logger.Println("[E] Resolve slave status failed!", err)
			timestamp := time.Now().Unix()
			logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + resolve_slave_status_failed + "{{" + fmt.Sprintf("%s", err)
			logger.Println("[I] Give up switching to async replication!")
			logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + give_async_repl
			return logvalue
		}
	}
	mapField2Data := make(map[string]interface{}, len(cols))
	for k, col := range data {
		mapField2Data[cols[k]] = col
	}
	Slave_IO_Running := mapField2Data["Slave_IO_Running"]
	if string(Slave_IO_Running.([]uint8)) != "Yes" {
		logger.Println("[I] The I/O thread status is " + string(Slave_IO_Running.([]uint8)) + "!")
		timestamp := time.Now().Unix()
		logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + io_status + "{{" + string(Slave_IO_Running.([]uint8))
		logvalue = Set(ip, port, username, password, 0,logvalue)
		return logvalue
	}
	logger.Println("[I] The I/O thread status is " + string(Slave_IO_Running.([]uint8)) + "!")
	timestamp = time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + io_status + "{{" + string(Slave_IO_Running.([]uint8))
	logger.Println("[I] Give up switching to async replication!")
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + give_async_repl
	return logvalue
}
