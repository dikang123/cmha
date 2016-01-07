package main

import (
	"database/sql"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"time"
)


func Set(ip, port, username, password string, issynchronous int) {
	dsName := username + ":" + password + "@tcp(" + "localhost" + ":" + port + ")/"
	db, err := sql.Open("mysql", dsName)
	if err != nil {
		beego.Error("Connection to the database failure!", err)
		return
	}
	beego.Info("Connection to the database success!")
	defer db.Close()
	err = db.Ping()
	if err != nil {
		beego.Error("ping() database  failure!", err)
		return
	}
	beego.Info("ping() database success!")
	keepsyncrepl := "set global rpl_semi_sync_master_keepsyncrepl=" + strconv.Itoa(issynchronous)
	_, err = db.Query(keepsyncrepl)
	if err != nil {
		beego.Error(keepsyncrepl + "  failure!")
		return
	}
	beego.Info(keepsyncrepl + " success!")
	trysyncrepl := "set global rpl_semi_sync_master_trysyncrepl=" + strconv.Itoa(issynchronous)
	_, err = db.Query(trysyncrepl)
	if err != nil {
		beego.Error(trysyncrepl + " failure!")
		return
	}
	beego.Info(trysyncrepl + " success!")
}

func checkio_thread(ip, port, username, password, addr string) {
	time.Sleep(60000 * time.Millisecond)
	dsName := username + ":" + password + "@tcp(" + addr + ":" + port + ")/"
	db, err := sql.Open("mysql", dsName)
	if err != nil {
		beego.Error("Connection to the database failure!", err)
		return
	}
	beego.Info("Connection to the database success!")
	defer db.Close()
	err = db.Ping()
	if err != nil {
		beego.Error("ping() database  failure!", err)
		Set(ip, port, username, password, 0)
		return
	}
	beego.Info("ping() database success!")
	//	for index := 2; index > 0; index-- {

	row, err := db.Query("show slave status")
	if err != nil {
		beego.Error("Inquiry slave status failure!", err)
		return
	}
	beego.Info("Inquiry slave status success!")
	cols, _ := row.Columns()
	buffer := make([]interface{}, len(cols))
	data := make([]interface{}, len(cols))
	for i, _ := range buffer {
		buffer[i] = &data[i]
	}
	for row.Next() {
		err = row.Scan(buffer...)
		if err != nil {
			beego.Error("scan() traversal slave status failure!", err)
			return
		}
	}
	mapField2Data := make(map[string]interface{}, len(cols))
	for k, col := range data {
		mapField2Data[cols[k]] = col
	}
	Slave_IO_Running := mapField2Data["Slave_IO_Running"]
	if string(Slave_IO_Running.([]uint8)) != "Yes" {
		beego.Error("IO copy the thread is not normal! ")
		Set(ip, port, username, password, 0)
		return
	}
	beego.Info("IO copy the thread Normal!")
	//}
}
