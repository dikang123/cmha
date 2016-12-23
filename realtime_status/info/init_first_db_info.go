package info

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/upmio/realtime_status/file"
)

var mysql1 []string

func InitFirstDbInfo(user, password, port string) {
	ismysql, _ := file.PathExists(mysql_statpath)
	if !ismysql {
		mysql1 = GetMysqlData(user, password, port)
	} else {
		mysql1 = DbInfo(mysql_statpath)
	}
}

func GetMysqlData(user, password, port string) []string {
	dsName := user + ":" + password + "@tcp(" + "localhost" + ":" + port + ")/"
	db, err := sql.Open("mysql", dsName)

	if err != nil {
		fmt.Println("db open failed!", err)
		os.Exit(2)
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
	var value string
	for rows.Next() {
		var Str string
		err := rows.Scan(&Str, &value)
		if err != nil {
			fmt.Println("db scan failed!", err)
			os.Exit(2)
		}
		mysql1 = append(mysql1, value)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	return mysql1
	//Slave_IO_Running_String := string(Slave_IO_Running.([]uint8))
}

func DbInfo(filename string) []string {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("open file filed!", err)
		os.Exit(2)
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)

		//if strings.Contains(line, types) {
		dbinfo := strings.Split(line, " ")
		return dbinfo
		//}
		//handler(line)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			os.Exit(2)
		}

	}
	return nil
}
