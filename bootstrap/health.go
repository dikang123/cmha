package main

import(
	"github.com/astaxie/beego"
	consulapi "github.com/hashicorp/consul/api"	
)

func HealthCheck(servicename,hostname string,client *consulapi.Client)(bool,error){
	islocal :=false
	 //Health returns a handle to the health endpoints
        health := client.Health()
        //Checks is used to return the checks associated with a service
        healthvalue, _, err := health.Checks(servicename, nil)
        if err != nil {
                beego.Error("Return to service-related checks failure!", err)
                return islocal,err
        }
        if len(healthvalue) <= 0 {
                beego.Info("Without this service, or service is not a healthy state!")
                return islocal,err
        }
        for index := range healthvalue {
                if healthvalue[index].Node == hostname {
                        islocal = true
                        beego.Info("Native " + servicename + " service is healthy!")
                        break
                }

        }
	return islocal,nil
}
