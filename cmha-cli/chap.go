package main

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/ryanuber/columnize"
)

func Chap(args ...string) error {
	if len(args) > 0 {
		config := &consulapi.Config{
			Datacenter: beego.AppConfig.String("cmha-datacenter"),
			Token:      beego.AppConfig.String("cmha-token"),
			Address:    beego.AppConfig.String("cmha-server-ip") + ":8500",
		}
		/*config := &consulapi.Config{
			Datacenter: beego.AppConfig.String("cmha-datacenter"),
			Token:      beego.AppConfig.String("cmha-token"),
		}
		var healthpair []*consulapi.ServiceEntry
		//var client *consulapi.Client
		var health *consulapi.Health
		for i, _ := range service_ip {
			config.Address = service_ip[i] + ":" + beego.AppConfig.String("cmha-server-ip")
			client, err := consulapi.NewClient(config)
			if err != nil {
				fmt.Println("chap Create consul-api client failure!", err)
				return err
			}
			health := client.Health()
			healthservice, _, err := health.Service(args[0], "", false, nil)
			if err != nil {
				fmt.Println("chap Query cluster health service failure!", err)
				continue
			}
			//		var d = 0
			var chaphealthy = true
			//		var dbhealthy = true
			if len(healthservice) < 1 {
				fmt.Println("not ", args[0], " service")
				return nil
			}
			break
		}*/
		client, err := consulapi.NewClient(config)
		if err != nil {
			fmt.Println("chap Create consul-api client failure!", err)
			return err
		}
		health := client.Health()
		healthservice, _, err := health.Service(args[0], "", false, nil)
		if err != nil {
			fmt.Println("chap Query cluster health service failure!", err)
			return err
		}
		//		var d = 0
		var chaphealthy = true
		//		var dbhealthy = true
		if len(healthservice) < 1 {
			fmt.Println("not ", args[0], " service")
			return nil
		}
		var c = 0
		fmt.Println("---------------------------------------------------------------------")
		//chapslice := []string{}
		chapslice := make([]string, 0, 10)
		//a := "Node\t\tAddress\t\tService\t\tStatus\t\tType"
		a := "Node|Address|Service|Status|Type"
		chapslice = append(chapslice, a)
		for index := range healthservice {
			//var issusehost string
			var chapcount = 0
			if healthservice[index].Service.Tags[0] == "chap-master" || healthservice[index].Service.Tags[0] == "chap-slave" {
				var unknown = 0
				for checkindex := range healthservice[index].Checks {
					if healthservice[index].Checks[checkindex].Status == "critical" {
						//issusehost = healthservice[index].Node.Address
						chaphealthy = false
						chapcount += 1
						if healthservice[index].Checks[checkindex].CheckID == "serfHealth" {
							//a := healthservice[index].Node.Node + "\t" + healthservice[index].Node.Address + "\thagent" + "\t\tFail" + "\t\tchap"
							a := fmt.Sprintf("%s|%s|%s|%s|%s", healthservice[index].Node.Node, healthservice[index].Node.Address, "hagent", "Fail", "chap")
							chapslice = append(chapslice, a)
							//fmt.Println(healthservice[index].Node.Node, healthservice[index].Node.Address, "HAgent", "Not OK")
						} else {
							//a := healthservice[index].Node.Node + "\t" + healthservice[index].Node.Address + "\t" + args[0] + "\t\tFail" + "\t\tchap"
							a := fmt.Sprintf("%s|%s|%s|%s|%s", healthservice[index].Node.Node, healthservice[index].Node.Address, args[0], "Fail", "chap")
							chapslice = append(chapslice, a)
							//fmt.Println(healthservice[index].Node.Node, healthservice[index].Node.Address, healthservice[index].Checks[checkindex].CheckID, "Not OK")
						}
						//fmt.Println("故障机器", issusehost)
						//break
					} else {
						if chapcount == 0 {
							chaphealthy = true
						}
						if healthservice[index].Checks[checkindex].CheckID == "serfHealth" {
							//a := healthservice[index].Node.Node + "\t" + healthservice[index].Node.Address + "\thagent" + "\t\tOK" + "\t\tchap"
							a := fmt.Sprintf("%s|%s|%s|%s|%s", healthservice[index].Node.Node, healthservice[index].Node.Address, "hagent", "OK", "chap")
							chapslice = append(chapslice, a)
							//fmt.Println(healthservice[index].Node.Node, healthservice[index].Node.Address, "HAgent", "OK")
						} else {
							for checkindex := range healthservice[index].Checks {
								if healthservice[index].Checks[checkindex].CheckID == "serfHealth" {
									if healthservice[index].Checks[checkindex].Status == "critical" {
										unknown += 1
									}
								}
							}
							if unknown != 0 {
								//a := healthservice[index].Node.Node + "\t" + healthservice[index].Node.Address + "\t" + args[0] + "\t\tUnKnown" + "\t\tchap"
								a := fmt.Sprintf("%s|%s|%s|%s|%s", healthservice[index].Node.Node, healthservice[index].Node.Address, args[0], "UnKnown", "chap")
								chapslice = append(chapslice, a)
							} else {
								if healthservice[index].Checks[checkindex].Status != "passing" {
									//a := healthservice[index].Node.Node + "\t" + healthservice[index].Node.Address + "\t" + args[0] + "\t\t" + healthservice[index].Checks[checkindex].Status + "\t\tchap"
									a := fmt.Sprintf("%s|%s|%s|%s|%s", healthservice[index].Node.Node, healthservice[index].Node.Address, args[0], healthservice[index].Checks[checkindex].Status, "chap")
									chapslice = append(chapslice, a)
									//fmt.Println(healthservice[index].Node.Node, healthservice[index].Node.Address, healthservice[index].Checks[checkindex].CheckID, healthservice[index].Checks[checkindex].Status)
								} else {
									//a := healthservice[index].Node.Node + "\t" + healthservice[index].Node.Address + "\t" + args[0] + "\t\tOK" + "\t\tchap"
									a := fmt.Sprintf("%s|%s|%s|%s|%s", healthservice[index].Node.Node, healthservice[index].Node.Address, args[0], "OK", "chap")
									chapslice = append(chapslice, a)
									//fmt.Println(healthservice[index].Node.Node, healthservice[index].Node.Address, healthservice[index].Checks[checkindex].CheckID, "OK")
								}
							}
							/*if healthservice[index].Checks[checkindex].Status != "passing" {
								a := healthservice[index].Node.Node + " " + healthservice[index].Node.Address + " " + healthservice[index].Checks[checkindex].CheckID + " " + healthservice[index].Checks[checkindex].Status
								chapslice = append(chapslice, a)
								//fmt.Println(healthservice[index].Node.Node, healthservice[index].Node.Address, healthservice[index].Checks[checkindex].CheckID, healthservice[index].Checks[checkindex].Status)
							} else {
								a := healthservice[index].Node.Node + " " + healthservice[index].Node.Address + " " + healthservice[index].Checks[checkindex].CheckID + " OK"
								chapslice = append(chapslice, a)
								//fmt.Println(healthservice[index].Node.Node, healthservice[index].Node.Address, healthservice[index].Checks[checkindex].CheckID, "OK")
							}*/
						}
					}
				}
				if chaphealthy {
					c += 1
				}
			}
		}
		fmt.Println("cmha chap  ", c, "/", 2)
		/*for chapi := range chapslice {
			fmt.Println(chapslice[chapi])
		}*/
		output := columnize.SimpleFormat(chapslice)
		fmt.Println(string(output))
		fmt.Println("---------------------------------------------------------------------")
		return nil

	}
	config := &consulapi.Config{
		Datacenter: beego.AppConfig.String("cmha-datacenter"),
		Token:      beego.AppConfig.String("cmha-token"),
		Address:    beego.AppConfig.String("cmha-server-ip") + ":8500",
	}
	/*config := &consulapi.Config{
		Datacenter: beego.AppConfig.String("cmha-datacenter"),
		Token:      beego.AppConfig.String("cmha-token"),
	}
	var catalog *consulapi.Catalog
	var service *consulapi.ServiceEntry
	for i, _ := range service_ip {
		config.Address = service_ip[i] + ":" + beego.AppConfig.String("cmha-server-ip")
		client, err := consulapi.NewClient(config)
		if err != nil {
			fmt.Println("chap Create consul-api client failure!", err)
			continue
		}
		catalog := client.Catalog()
		health := client.Health()
		services, _, err := catalog.Services(nil)
		if err != nil {
			fmt.Println("chap Query services failure!", err)
			continue
		}
		break
	}*/
	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Println("chap Create consul-api client failure!", err)
		return err
	}
	catalog := client.Catalog()
	health := client.Health()
	services, _, err := catalog.Services(nil)
	if err != nil {
		fmt.Println("chap Query services failure!", err)
		return err
	}
	var ishealthy = true
	for k, _ := range services {
		chapslice1 := []string{}
		if k != "consul" {
			service, _, err := catalog.Service(k, "", nil)
			if err != nil {
				fmt.Println("chap Query catalog service "+k, err)
				return err
			}
			var count = 0
			//var issusehost string
			fmt.Println("-------------------------")
			fmt.Println("      ", k)
			for _, value := range service {
				if value.ServiceTags[0] == "chap-master" || value.ServiceTags[0] == "chap-slave" {
					healthservice, _, err := health.Service(k, value.ServiceTags[0], false, nil)
					if err != nil {
						fmt.Println("chap Check cluster health service failure!", err)
						return err
					}
					for index := range healthservice {
						for checkindex := range healthservice[index].Checks {
							if strings.EqualFold(value.Address, healthservice[index].Node.Address) {
								if healthservice[index].Checks[checkindex].Status == "critical" {
									//issusehost = value.Address
									ishealthy = false
									a := "Fault machine: " + healthservice[index].Node.Address
									chapslice1 = append(chapslice1, a)
									//fmt.Println("故障机器", issusehost)
									break
								} else {
									ishealthy = true
								}
							}
						}
						if ishealthy {
							count += 1
						}

					}
				}
			}
			fmt.Println("cmha chap  ", count, "/", 2)
			for chapi1 := range chapslice1 {
				fmt.Println(chapslice1[chapi1])
			}
		}
	}
	fmt.Println("-------------------------")
	return nil
}
