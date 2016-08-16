package alerts

import (
	"os"
	"strings"
	"time"

	"github.com/astaxie/beego"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/upmio/alerts-handler/log"
)

var client *consulapi.Client

func CsAlerts(servicename string) {
	log.Info("cs")
	var catalog *consulapi.Catalog
	var kv *consulapi.KV
	log.Infof("CsAlerts")
	var C_time = time.Now().Unix()
	var timeLayout = "2006-01-02 15:04:05"
	var dataTimeStr = time.Unix(C_time, 0).Format(timeLayout)
	ips := beego.AppConfig.String("ip")
	service_ip := beego.AppConfig.Strings("service_ip")
	config := &consulapi.Config{
		Datacenter: beego.AppConfig.String("datacenter"),
		Token:      beego.AppConfig.String("token"),
	}
	logfiles, err := os.OpenFile("/usr/local/cmha/alerts-handler/logs/template_cs_node_alerts.log", os.O_RDWR|os.O_APPEND|os.O_CREATE|os.O_EXCL, 0)
	if err != nil {
		log.Errorf("%s", err.Error())
		os.Exit(-1)
	}
	defer logfiles.Close()
	for ip := range service_ip {
		config.Address = ips + ":" + beego.AppConfig.String("port")
		client, err := consulapi.NewClient(config)
		if err != nil {
			log.Errorf("consul api error:", err)
			continue
		}
		kv = client.KV()
		catalog = client.Catalog()
		nodes, _, err := catalog.Nodes(nil)
		if err != nil {
			log.Errorf("query nodes error:", err)
			continue
		}
		var hostname string
		for index := range nodes {
			if nodes[index].Address == service_ip[ip] {
				hostname = nodes[index].Node
				log.Infof("hostname:%s,service_ip:%s",hostname,service_ip[ip])
				break
			}
		}
		health := client.Health()
		health_check, _, err := health.State("critical", nil)
		for healthindex := range health_check {
			if health_check[healthindex].Node == hostname {
				log.Infof("Node:%s,hostname:%s",health_check[healthindex].Node,hostname)
				output_str := strings.Replace(health_check[healthindex].Output, "\n", "", -1)
				str := "{Status||" + health_check[healthindex].Status + "%%CheckID||" + health_check[healthindex].CheckID + "%%Node||" + health_check[healthindex].Node + "%%Output||" + output_str + "\n"
				log.Infof("str:%s",str)
				logfiles.WriteString(str)
			}
		}
		//catalog := client.Catalog()
	}
	ReadFile(catalog, kv, C_time, dataTimeStr, servicename, true)
	err = os.Remove("/usr/local/cmha/alerts-handler/logs/template_cs_node_alerts.log")
	if err != nil {
		log.Errorf("remove template_cs_node_alerts.log file error:", err)
		return
	}
}
