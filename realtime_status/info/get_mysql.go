package info

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var (
	Com_insert                       int
	Com_update                       int
	Com_delete                       int
	Com_select                       int
	TPS                              int
	Innodb_buffer_pool_read_requests int
	hit                              float64
	Threads_running                  int
	Threads_connected                int
	Threads_created                  int
	Threads_cached                   int
)
var mysql_statpath = "/tmp/.realtime_cache/mysql_stat"

func GetMysql(user, password, port string) {
	dsName := user + ":" + password + "@tcp(" + "localhost" + ":" + port + ")/"
	db, err := sql.Open("mysql", dsName)
	if err != nil {
		fmt.Println("db open failed!",err)
		os.Exit(0)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println("db ping failed!", err)
		os.Exit(2)
	}
	rows, err := db.Query("show global status where Variable_name in ('Com_select','Com_insert','Com_update','Com_delete','Threads_running','Threads_connected','Threads_cached','Threads_created','Bytes_received','Bytes_sent','Innodb_rows_inserted','Innodb_rows_updated','Innodb_rows_deleted','Innodb_rows_read','Innodb_buffer_pool_read_requests','Innodb_buffer_pool_reads','Innodb_buffer_pool_pages_data','Innodb_buffer_pool_pages_free','Innodb_buffer_pool_pages_dirty','Innodb_buffer_pool_pages_flushed','Innodb_data_reads','Innodb_data_writes','Innodb_data_read','Innodb_data_written','Innodb_os_log_fsyncs','Innodb_os_log_written')")
	if err != nil {
		fmt.Println("db query failed!", err)
		os.Exit(2)
	}
	defer rows.Close()
	var mysql2 []string
	var value string
	for rows.Next() {
		var Str string
		err := rows.Scan(&Str, &value)
		if err != nil {
			fmt.Println("db scan failed!", err)
			os.Exit(2)

		}
		mysql2 = append(mysql2, value)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println("err!", err)
		os.Exit(2)
	}
	//	Bytes_received := RdIos(mysql1, mysql2, 0)
	_ = RdIos(mysql1, mysql2, 0)
	//	Bytes_sent := RdIos(mysql1, mysql2, 1)
	_ = RdIos(mysql1, mysql2, 1)
	Com_delete = RdIos(mysql1, mysql2, 2)
	Com_insert = RdIos(mysql1, mysql2, 3)
	Com_select = RdIos(mysql1, mysql2, 4)
	Com_update = RdIos(mysql1, mysql2, 5)
	//Innodb_buffer_pool_pages_data := RdIos(mysql1, mysql2, 6)
	_ = RdIos(mysql1, mysql2, 6)
	//	Innodb_buffer_pool_pages_dirty := RdIos(mysql1, mysql2, 7)
	_ = RdIos(mysql1, mysql2, 7)
	//	Innodb_buffer_pool_pages_flushed := RdIos(mysql1, mysql2, 8)
	_ = RdIos(mysql1, mysql2, 8)
	//	Innodb_buffer_pool_pages_free := RdIos(mysql1, mysql2, 9)
	_ = RdIos(mysql1, mysql2, 9)
	Innodb_buffer_pool_read_requests = RdIos(mysql1, mysql2, 10)
	Innodb_buffer_pool_reads := RdIos(mysql1, mysql2, 11)
	//	Innodb_data_read := RdIos(mysql1, mysql2, 12)
	_ = RdIos(mysql1, mysql2, 12)
	//	Innodb_data_reads := RdIos(mysql1, mysql2, 13)
	_ = RdIos(mysql1, mysql2, 13)
	//	Innodb_data_writes := RdIos(mysql1, mysql2, 14)
	_ = RdIos(mysql1, mysql2, 14)
	//	Innodb_data_written := RdIos(mysql1, mysql2, 15)
	_ = RdIos(mysql1, mysql2, 15)
	//	Innodb_os_log_fsyncs := RdIos(mysql1, mysql2, 16)
	_ = RdIos(mysql1, mysql2, 16)
	//	Innodb_os_log_written := RdIos(mysql1, mysql2, 17)
	_ = RdIos(mysql1, mysql2, 17)
	//	Innodb_rows_deleted := RdIos(mysql1, mysql2, 18)
	_ = RdIos(mysql1, mysql2, 18)
	//	Innodb_rows_inserted := RdIos(mysql1, mysql2, 19)
	_ = RdIos(mysql1, mysql2, 19)
	//	Innodb_rows_read := RdIos(mysql1, mysql2, 20)
	_ = RdIos(mysql1, mysql2, 20)
	//	Innodb_rows_updated := RdIos(mysql1, mysql2, 21)
	_ = RdIos(mysql1, mysql2, 21)
	Threads_cached = ReturnValue(mysql2, 22)
	Threads_connected = ReturnValue(mysql2, 23)
	Threads_created = RdIos(mysql1, mysql2, 24)
	Threads_running = ReturnValue(mysql2, 25)
	TPS = Com_delete + Com_insert + Com_update
	if Innodb_buffer_pool_read_requests == 0 {
		hit = 100.00
	} else {
		hit = (float64(Innodb_buffer_pool_read_requests) - float64(Innodb_buffer_pool_reads))/float64(Innodb_buffer_pool_read_requests)*float64(100)

	}
	WriteProcStat(mysql2, mysql_statpath)
}

func ReturnValue(mysql2 []string, index int) int {
	var value int
	for i := range mysql2 {
		if i == index {
			value = StringToInt(mysql2[i])
		}
	}
	return value
}
