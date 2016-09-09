package check

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func CheckMysqlHealth(user, password, ip, port, defaultDb, timeout string,checktime int) int {
	var MYSQL_OK int
	db, err := Conn(user, password, ip, port, defaultDb, timeout)
	if err != nil {
		if checktime == 1{
			fmt.Print(err)
		}
		MYSQL_OK = 1
		return MYSQL_OK
	}
	err = SetMysql(db, "set sql_log_bin=0;")
	if err != nil {
		if checktime == 1{
			fmt.Print(err)
		}
		MYSQL_OK = 1
		return MYSQL_OK
	}
	err = SetMysql(db, "set innodb_lock_wait_timeout=3;")
	if err != nil {
		if checktime == 1 {
			fmt.Print(err)
		}
		MYSQL_OK = 1
		return MYSQL_OK
	}
	err = SetMysql(db, "set lock_wait_timeout=3;")
	if err != nil {
		if checktime == 1{
			fmt.Print(err)
		}
		MYSQL_OK = 1
		return MYSQL_OK
	}
	tx, err := Tx(db)
	if err != nil {
		if checktime == 1{ 
			fmt.Print(err)
		}
		MYSQL_OK = 1
		return MYSQL_OK
	}
	err = MysqlOperation(tx)
	if err != nil {
		if checktime == 1{
			fmt.Print(err)
		}
		MYSQL_OK = 1
		return MYSQL_OK
	}
	MYSQL_OK = 0
	return MYSQL_OK
}

func SelectCheckMysqlHealth(user, password, ip, port, defaultDb, timeout string,checktime int) int {
	var MYSQL_OK int
	db, err := Conn(user, password, ip, port, defaultDb, timeout)
	if err != nil {
		if checktime == 1{
			fmt.Print(err)
		}
		MYSQL_OK = 1
		return MYSQL_OK
	}
	defer db.Close()
	err = ExecSelect(db)
	if err != nil {
		if checktime == 1 {
			fmt.Print(err)
		}
		MYSQL_OK = 1
		return MYSQL_OK
	}
	MYSQL_OK = 0
	return MYSQL_OK
}

func ExecSelect(db *sql.DB) error {
	row, err := db.Query("select cmha_name from cmha_check;")
	if err != nil {
		return err
	}
	defer row.Close()
	err = row.Err()
	if err != nil {
		return err
	}
	return nil
}

func Conn(user, password, ip, port, defaultDb, t string) (*sql.DB, error) {
	time_int, _ := strconv.Atoi(t)
	time_duration := time.Duration(time_int) * time.Second
	_t := time_duration
	openstr := user + ":" + password + "@tcp(" + ip + ":" + port + ")/" + defaultDb + "?timeout=" + _t.String()
	db, err := sql.Open("mysql", openstr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Tx(db *sql.DB) (*sql.Tx, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func SetMysql(db *sql.DB, sqlstr string) error {
	log_bin, err := db.Query(sqlstr)
	if err != nil {
		return err
	}
	defer log_bin.Close()
	err = log_bin.Err()
	if err != nil {
		return err
	}
	return nil
}

func MysqlOperation(tx *sql.Tx) error {
	err := MysqlExec(tx, "select cmha_name from cmha_check;")
	if err != nil {
		tx.Rollback()
		return err
	}
	err = MysqlExec(tx, "delete from cmha_check;")
	if err != nil {
		tx.Rollback()
		return err
	}
	nowtime_string := GetNowTime()
	insertstr := "insert into cmha_check(id,cmha_name,create_time) values(1,'cmha_check_insert','" + nowtime_string + "');"
	err = MysqlExec(tx, insertstr)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = MysqlExec(tx, "update cmha_check set cmha_name='cmha_check_update' where id=1;")
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func GetNowTime() string {
	C_time := time.Now().Unix()
	timeLayout := "2006-01-02 15:04:05"
	dataTimeStr := time.Unix(C_time, 0).Format(timeLayout)
	return dataTimeStr
}

func MysqlExec(tx *sql.Tx, sqlstr string) error {
	stmt, err := tx.Query(sqlstr)
	if err != nil {
		return err
	}
	err = stmt.Err()
	if err != nil {
		return err
	}
	defer stmt.Close()
	return nil
}
