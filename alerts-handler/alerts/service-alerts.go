package alerts

import (
	"os"
	"strings"
	"time"

	"github.com/astaxie/beego"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/upmio/alerts-handler/log"
)

func IsDbOrChap(servicename string) {
	C_time := time.Now().Unix()

	var timeLayout = "2006-01-02 15:04:05"
	var dataTimeStr = time.Unix(C_time, 0).Format(timeLayout)
	var err error
	ip := beego.AppConfig.String("ip")
	port := beego.AppConfig.String("port")
	address := ip + ":" + port
	config := &consulapi.Config{
		Datacenter: beego.AppConfig.String("datacenter"),
		Token:      beego.AppConfig.String("token"),
		Address:    address,
	}
	logfiles, err := os.OpenFile("/usr/local/cmha/alerts-handler/logs/template_service_node_alerts.log", os.O_RDWR|os.O_APPEND|os.O_CREATE|os.O_EXCL, 0)
	if err != nil {
		log.Errorf("%s", err.Error())
		os.Exit(-1)
	}
	defer logfiles.Close()
	var healthpair []*consulapi.ServiceEntry
	var health *consulapi.Health
	client, err := consulapi.NewClient(config)
	if err != nil {
		log.Errorf("Create consul-api client failed!", err)
		return
	}
	health = client.Health()
	kv := client.KV()
	healthpair, _, err = health.Service(servicename, "", false, nil)
	if err != nil {
		log.Errorf("err", err)
		return
	}
	for index := range healthpair {
		for checkindex := range healthpair[index].Checks {
			if healthpair[index].Checks[checkindex].Status == "warning" {
				output_str := strings.Replace(healthpair[index].Checks[checkindex].Output, "\n", "", -1)
				str := "{Status||" + healthpair[index].Checks[checkindex].Status + "%%CheckID||" + healthpair[index].Checks[checkindex].CheckID + "%%Node||" + healthpair[index].Checks[checkindex].Node + "%%Output||" + output_str + "\n"
				logfiles.WriteString(str)
			} else if healthpair[index].Checks[checkindex].Status == "critical" {
				output_str := strings.Replace(healthpair[index].Checks[checkindex].Output, "\n", "", -1)
				str := "{Status||" + healthpair[index].Checks[checkindex].Status + "%%CheckID||" + healthpair[index].Checks[checkindex].CheckID + "%%Node||" + healthpair[index].Checks[checkindex].Node + "%%Output||" + output_str + "\n"
				logfiles.WriteString(str)
			}
		}
	}
	catalog := client.Catalog()
	ReadFile(catalog, kv, C_time, dataTimeStr, servicename, false)
	err = os.Remove("/usr/local/cmha/alerts-handler/logs/template_service_node_alerts.log")
	if err != nil {
		log.Errorf("remove template_service_node_alerts.log file error:", err)
		return
	}
}
