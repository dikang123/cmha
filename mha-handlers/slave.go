package main

import (
	"database/sql"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
)

func slave(ip, port, username, password string) {
	dsName := username + ":" + password + "@tcp(" + "localhost" + ":" + port + ")/"
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
	_, err = db.Query("stop slave io_thread")
	if err != nil {
		beego.Error("Stop local database replication I/O thread failed!", err)
		beego.Info("Give up leader election")
                beego.Info("MHA Handler Completed")
		return
	}
	beego.Info("Stop local database replication I/O thread successfully!")
	row, err := db.Query("show slave status")
	if err != nil {
		beego.Error("Checking local database I/O thread status. Failed!", err)
		beego.Info("Give up leader election")
                beego.Info("MHA Handler Completed")
		return
	}
	beego.Info("Checking local database I/O thread status. Succeed!")
	cols, _ := row.Columns()
	buffer := make([]interface{}, len(cols))
	data := make([]interface{}, len(cols))
	for i, _ := range buffer {
		buffer[i] = &data[i]
	}
	for row.Next() {
		err = row.Scan(buffer...)
		if err != nil {
			beego.Error("Resolve slave status failed!", err)
			beego.Info("Give up leader election")
	                beego.Info("MHA Handler Completed")
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
		beego.Error("The SQL thread status is " + string(Slave_SQL_Running.([]uint8)) + "!")
		beego.Info("Give up leader election")
                beego.Info("MHA Handler Completed")
		return
	}
	beego.Info("The SQL thread status is Yes!")
	sqlstr := "select master_pos_wait(?,?)"
	rowss, err := db.Query(sqlstr, Master_Log_File, Read_Master_Log_Pos)
	if err != nil {
		beego.Error("Checking relay log applying status failed!", err)
		beego.Info("Give up leader election")
                beego.Info("MHA Handler Completed")
		return
	}
	beego.Info("Checking relay log applying status successfully!")
	var master_pos_wait string
	for rowss.Next() {
		err = rowss.Scan(&master_pos_wait)
		if err != nil {
			beego.Error("Resolve master_pos_wait failed!", err)
			beego.Info("Give up leader election")
	                beego.Info("MHA Handler Completed")
			return
		}
		if master_pos_wait < "0" && master_pos_wait == "null" {
			beego.Error("Relay log applying failed!", err)
			beego.Info("Give up leader election")
	                beego.Info("MHA Handler Completed")
			return
		}
		beego.Info("Relay log applying completed!")
		_, err := db.Query("set global rpl_semi_sync_master_keepsyncrepl=0")
		if err != nil {
			beego.Error("Set rpl_semi_sync_master_keepsyncrepl=0 failed!")
			beego.Info("Give up leader election")
	                beego.Info("MHA Handler Completed")
			return
		}
		beego.Info("Set rpl_semi_sync_master_keepsyncrepl=0 successfully!")
		_, err = db.Query("set global rpl_semi_sync_master_trysyncrepl=0")
		if err != nil {
			beego.Error("Set rpl_semi_sync_master_trysyncrepl=0 failed!")
			beego.Info("Give up leader election")
	                beego.Info("MHA Handler Completed")
			return
		}
		beego.Info("Switching local database to async replication!")
		_, err = db.Query("set global read_only=0")
		if err != nil {
			beego.Error("Set local database Read/Write mode failed!", err)
			beego.Info("Give up leader election")
	                beego.Info("MHA Handler Completed")
			return
		}
		beego.Info("Set local database Read/Write mode successfully!")
	}
	SetConn(ip, port, username, password)
}
