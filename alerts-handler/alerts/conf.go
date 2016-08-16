package alerts

import (
	"github.com/astaxie/beego"
)

func GetConf() (string, string, string) {
	ip := beego.AppConfig.String("ip")
	port := beego.AppConfig.String("port")
	adress := ip + ":" + port
	return ip, port, adress
}
