package log

import (
	"os"

	"github.com/astaxie/beego"
)

func AlertFile() (*os.File, error) {
	logfile_path := beego.AppConfig.String("logfile_path")
	logfile, err := os.OpenFile(logfile_path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0)
	if err != nil {
		Errorf("%s\r\n", err.Error())
		return nil, err
		//os.Exit(-1)
	}
	defer logfile.Close()
	return logfile, nil
}
