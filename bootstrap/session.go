package main

import(
	"time"
	"github.com/astaxie/beego"
	consulapi "github.com/hashicorp/consul/api"	
)

func CreateSession(servicename,hostname string,client *consulapi.Client)(string,error){
	//Session returns a handle to the session endpoints
        session := client.Session()
        sessionEntry := consulapi.SessionEntry{
        	LockDelay: 10 * time.Second,
                Name:      servicename,
               	Node:      hostname,
                Checks:    []string{"serfHealth", "service:" + servicename},
        }
	//Create makes a new session. Providing a session entry can customize the session. It can also be nil to use defaults.
        sessionvalue, _, err := session.Create(&sessionEntry, nil)
        if err != nil {
        	beego.Error("Session creation failure!", err)
                return "",err
        }
	return sessionvalue,nil
}
