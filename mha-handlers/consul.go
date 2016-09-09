package main

import(
	"fmt"
	"time"
	"strconv"
	"github.com/astaxie/beego"
	consulapi "github.com/hashicorp/consul/api"
)
func GetConfig(address string)*consulapi.Config{
	config := &consulapi.Config{
                Datacenter: beego.AppConfig.String("datacenter"),
                Token:      beego.AppConfig.String("token"),
                Address:    address,
        }
	return config
}

func GetClient(config *consulapi.Config,logvalue,consulapi_failed string)(*consulapi.Client,string,error){
	client, err := consulapi.NewClient(config)
	if err != nil {
                logger.Println("[E] Create consul-api client failed!", err)
                timestamp := time.Now().Unix()
                logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + consulapi_failed + "{{" + fmt.Sprintf("%s", err)
		return nil,logvalue,err	
        }
	return client,logvalue,nil
}
