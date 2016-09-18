package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/fatih/color"
	_ "github.com/go-sql-driver/mysql"
	"github.com/olekukonko/tablewriter"
	"github.com/upmio/cmha-cli/cliconfig"
)

type Color_option int

const (
	COLOR_NORMAL Color_option = iota
	COLOR_WARNNING
	COLOR_ERROR
	COLOR_SUM
)

const (
	ALIGN_DEFAULT = iota
	ALIGN_CENTRE
	ALIGN_RIGHT
	ALIGN_LEFT
)

const (
	LOG_MHA     string = "mha-handlers"
	LOG_MONITOR string = "monitor-handlers"
	LOG_ALERT string = "alerts_counter"
)

func TableRender(theader []string, _data [][]string, align int) {

	_table := tablewriter.NewWriter(os.Stdout)
	_table.SetHeader(theader)

	for _, _v := range _data {
		_table.Append(_v)
	}

	_table.SetAlignment(align)

	_table.Render()

}

func ColorRender(_input string, _color Color_option) string {

	switch _color {
	case COLOR_NORMAL:
		return color.GreenString(_input)
	case COLOR_ERROR:
		return color.RedString(_input)
	case COLOR_WARNNING:
		return color.YellowString(_input)
	case COLOR_SUM:
		return color.BlueString(_input)
	default:
		return _input

	}

}

func ServiceData(_ip string, _port string, _sql string, filters ...string) (map[string]string, error) {
	// two columns

	_user := cliconfig.GetUserName()
	_password := cliconfig.GetPassword()

	_dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", _user, _password, _ip, _port)

	_db, err := sql.Open("mysql", _dsn)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer _db.Close()

	_rows, err := _db.Query(_sql)
	if err != nil {
		return nil, err
	}
	defer _rows.Close()

	cols, _ := _rows.Columns()

	buffer := make([]interface{}, len(cols))
	data := make([]interface{}, len(cols))

	for i, _ := range buffer {
		buffer[i] = &data[i]
	}

	var _counter int = 0

	_ret := make(map[string]string)

	for _rows.Next() {

		var _colname, _colvalue string

		err = _rows.Scan(buffer...)
		_counter++

		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}

		for _k_col, _v_col := range data {

			if _k_col == 0 {

				// for filters ; It is TRUE forever,if No filter
				isFound := true

				if len(filters) > 0 {

					isFound = false

					for _, filter := range filters {

						if filter == string(_v_col.([]uint8)) {
							isFound = true

							// match filter ; break make filter loop
							break
						}
					}

					if !isFound {
						// not match filter ; break k/v loop
						break
					}
				}

				if isFound {
					_colname = string(_v_col.([]uint8))
				}
			}

			if _k_col == 1 {
				_colvalue = string(_v_col.([]uint8))
			}

		}

		if len(_colname) > 0 && len(_colvalue) > 0 {
			_ret[_colname] = _colvalue
		}

	}

	return _ret, nil

}

func SlaveStatus(_ip string, _port string, _sql string, columns ...string) (map[string]string, error) {

	if len(columns) < 1 {
		err := errors.New("At least need a column name.")
		return nil, err
	}

	_user := cliconfig.GetUserName()
	_password := cliconfig.GetPassword()

	_dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", _user, _password, _ip, _port)

	_db, err := sql.Open("mysql", _dsn)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer _db.Close()

	_rows, err := _db.Query(_sql)
	if err != nil {
		return nil, err
	}

	defer _rows.Close()

	cols, _ := _rows.Columns()

	buffer := make([]interface{}, len(cols))
	data := make([]interface{}, len(cols))

	for i, _ := range buffer {
		buffer[i] = &data[i]
	}

	for _rows.Next() {

		err = _rows.Scan(buffer...)

		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}

	}

	mapField2Data := make(map[string]interface{}, len(cols))
	for k, col := range data {
		mapField2Data[cols[k]] = col
	}

	_ret := make(map[string]string)

	for _, columsname := range columns {
		slave_status_value := mapField2Data[columsname]
		_ret[columsname] = string(slave_status_value.([]uint8))
	}

	return _ret, nil

}

//return min, max
func int64_arrange(num1 int64, num2 int64) (int64, int64) {
	if num1 > num2 {
		return num2, num1
	}
	return num1, num2
}

func StringToInt(instr string) (int64, error) {

	_ret, err := strconv.ParseInt(instr, 10, 64)

	if err != nil {
		return 0, err
	}

	return _ret, nil
}
