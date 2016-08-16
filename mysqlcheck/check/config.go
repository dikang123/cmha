package check

import (
	"github.com/astaxie/beego"
)

func GetConfig() (string, string, string, string, string, string) {
	host := beego.AppConfig.String("HOST")
	port := beego.AppConfig.String("PORT")
	checktime_string := beego.AppConfig.String("CHECK_TIME")
	timeout := beego.AppConfig.String("TIMEOUT")
	defaultDb := beego.AppConfig.String("DATABASE")
	ping_type := beego.AppConfig.String("PING_TYPE")
	return host, port, checktime_string, timeout, defaultDb, ping_type
}
