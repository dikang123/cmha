package main

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/ryanuber/columnize"
)

func Cluster(args ...string) error {
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
		var health *consulapi.Health
		var healthservicer []*consulapi.ServiceEntry
		for i, _ := range service_ip {
			config.Address = service_ip[i] + ":" + beego.AppConfig.String("cmha-server-ip")
			client, err := consulapi.NewClient(config)
			if err != nil {
				fmt.Println("cluster Create consul-api client failure!", err)
				return err
			}
			health := client.Health()
			healthservice, _, err := health.Service(args[0], "", false, nil)
			if err != nil {
				fmt.Println("cluster Query cluster health service failure!", err)
				continue
			}
			break
		}
		*/
		client, err := consulapi.NewClient(config)
		if err != nil {
			fmt.Println("cluster Create consul-api client failure!", err)
			return err
		}
		health := client.Health()
		healthservice, _, err := health.Service(args[0], "", false, nil)
		if err != nil {
			fmt.Println("cluster Query cluster health service failure!", err)
			return err
		}
		var c = 0
		var d = 0
		var chaphealthy = true
		var dbhealthy = true
		if len(healthservice) < 1 {
			fmt.Println("not ", args[0], " service")
			return nil
		}
		fmt.Println("---------------------------------------------------------------------")
		//chapslice := []string{}
		chapslice := make([]string, 0, 10)
		//a := "Node\t\tAddress\t\tService\t\tStatus\t\tType"
		a := "Node|Address|Service|Status|Type"
		chapslice = append(chapslice, a)
		//dbslice := []string{}
		for index := range healthservice {
			var chapcount = 0
			var dbcount = 0
			if healthservice[index].Service.Tags[0] == "chap-master" || healthservice[index].Service.Tags[0] == "chap-slave" {
				var unknown = 0
				for checkindex := range healthservice[index].Checks {
					//fmt.Println(healthservice[index].Checks[checkindex].Status)
					if healthservice[index].Checks[checkindex].Status == "critical" {
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
						//	break
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
								//a := healthservice[index].Node.Node + "\t" + healthservice[index].Node.Address + "\t" + args[0] + "\t\tUnknown" + "\t\tchap"
								a := fmt.Sprintf("%s|%s|%s|%s|%s", healthservice[index].Node.Node, healthservice[index].Node.Address, args[0], "Unknown", "chap")
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
			if healthservice[index].Service.Tags[0] == "master" || healthservice[index].Service.Tags[0] == "slave" {
				var unknown = 0
				for checkindex := range healthservice[index].Checks {
					if healthservice[index].Checks[checkindex].Status == "critical" {
						dbhealthy = false
						dbcount += 1
						if healthservice[index].Checks[checkindex].CheckID == "serfHealth" {
							//a := healthservice[index].Node.Node + "\t" + healthservice[index].Node.Address + "\tmagent" + "\t\tFail" + "\t\tdb"
							a := fmt.Sprintf("%s|%s|%s|%s|%s", healthservice[index].Node.Node, healthservice[index].Node.Address, "magent", "Fail", "db")
							chapslice = append(chapslice, a)
							//fmt.Println(healthservice[index].Node.Node, healthservice[index].Node.Address, "MAgent", "Not OK")
						} else {
							//a := healthservice[index].Node.Node + "\t" + healthservice[index].Node.Address + "\t" + args[0] + "\t\tFail" + "\t\tdb"
							a := fmt.Sprintf("%s|%s|%s|%s|%s", healthservice[index].Node.Node, healthservice[index].Node.Address, args[0], "Fail", "db")
							chapslice = append(chapslice, a)
							//fmt.Println(healthservice[index].Node.Node, healthservice[index].Node.Address, healthservice[index].Checks[checkindex].CheckID, "Not OK")
						}
						//break
					} else {
						if dbcount == 0 {
							dbhealthy = true
						}

						if healthservice[index].Checks[checkindex].CheckID == "serfHealth" {
							//a := healthservice[index].Node.Node + "\t" + healthservice[index].Node.Address + "\tmagent" + "\t\tOK" + "\t\tdb"
							a := fmt.Sprintf("%s|%s|%s|%s|%s", healthservice[index].Node.Node, healthservice[index].Node.Address, "magent", "OK", "db")
							chapslice = append(chapslice, a)
							//fmt.Println(healthservice[index].Node.Node, healthservice[index].Node.Address, "MAgent", "OK")
						} else {
							for checkindex := range healthservice[index].Checks {
								if healthservice[index].Checks[checkindex].CheckID == "serfHealth" {
									if healthservice[index].Checks[checkindex].Status == "critical" {
										unknown += 1
									}
								}
							}
							if unknown != 0 {
								//a := healthservice[index].Node.Node + "\t" + healthservice[index].Node.Address + "\t" + args[0] + "\t\tUnKnown" + "\t\tdb"
								a := fmt.Sprintf("%s|%s|%s|%s|%s", healthservice[index].Node.Node, healthservice[index].Node.Address, args[0], "UnKnown", "db")
								chapslice = append(chapslice, a)
							} else {
								if healthservice[index].Checks[checkindex].Status != "passing" {
									//a := healthservice[index].Node.Node + "\t" + healthservice[index].Node.Address + "\t" + args[0] + "\t\t" + healthservice[index].Checks[checkindex].Status + "\t\tdb"
									a := fmt.Sprintf("%s|%s|%s|%s|%s", healthservice[index].Node.Node, healthservice[index].Node.Address, args[0], healthservice[index].Checks[checkindex].Status, "db")
									chapslice = append(chapslice, a)
									//fmt.Println(healthservice[index].Node.Node, healthservice[index].Node.Address, healthservice[index].Checks[checkindex].CheckID, healthservice[index].Checks[checkindex].Status)
								} else {
									//a := healthservice[index].Node.Node + "\t" + healthservice[index].Node.Address + "\t" + args[0] + "\t\tOK" + "\t\tdb"
									a := fmt.Sprintf("%s|%s|%s|%s|%s", healthservice[index].Node.Node, healthservice[index].Node.Address, args[0], "OK", "db")
									chapslice = append(chapslice, a)
									//fmt.Println(healthservice[index].Node.Node, healthservice[index].Node.Address, healthservice[index].Checks[checkindex].CheckID, "OK")
								}
							}
							/*if healthservice[index].Checks[checkindex].Status != "passing" {
								a := healthservice[index].Node.Node + " " + healthservice[index].Node.Address + " " + healthservice[index].Checks[checkindex].CheckID + " " + healthservice[index].Checks[checkindex].Status
								dbslice = append(dbslice, a)
								//fmt.Println(healthservice[index].Node.Node, healthservice[index].Node.Address, healthservice[index].Checks[checkindex].CheckID, healthservice[index].Checks[checkindex].Status)
							} else {
								a := healthservice[index].Node.Node + " " + healthservice[index].Node.Address + " " + healthservice[index].Checks[checkindex].CheckID + " OK"
								dbslice = append(dbslice, a)
								//fmt.Println(healthservice[index].Node.Node, healthservice[index].Node.Address, healthservice[index].Checks[checkindex].CheckID, "OK")
							}*/

						}
					}
				}
				if dbhealthy {
					d += 1
				}
			}
		}
		fmt.Println("cmha chap  ", c, "/", 2)
		fmt.Println("cmha db    ", d, "/", 2)
		/*for chapi := range chapslice {
			fmt.Println(chapslice[chapi])
		}*/
		output := columnize.SimpleFormat(chapslice)
		fmt.Println(string(output))
		//fmt.Println("cmha db    ", d, "/", 2)
		//for dbi := range dbslice {
		//	fmt.Println(dbslice[dbi])
		//}
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
	var status *consulapi.Status
	for i, _ := range service_ip {
		config.Address = service_ip[i] + ":" + beego.AppConfig.String("cmha-server-ip")
		client, err := consulapi.NewClient(config)
		if err != nil {
			fmt.Println("cluster Create consul-api client failure!", err)
			return err
		}
		status := client.Status()
		peers, err := status.Peers()
		if err != nil {
			fmt.Println("cluster Query consul service faiiure!", err)
			continue
		}
		break
	}*/

	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Println("cluster Create consul-api client failure!", err)
		return err
	}
	status := client.Status()
	peers, err := status.Peers()
	if err != nil {
		fmt.Println("cluster Query consul service faiiure!", err)
		return err
	}
	agent := client.Agent()
	members, err := agent.Members(false)
	if err != nil {
		fmt.Println("cluster Query cluster members failure!", err)
		return err
	}
	peerssize := len(peers)
	var count = peerssize
	slice := []string{}
	for _, peersvalue := range peers {
		a := strings.Split(peersvalue, ":")
		for _, membersvalue := range members {
			if strings.EqualFold(a[0], membersvalue.Addr) {
				if membersvalue.Status != 1 {
					a := "Fault machine: " + membersvalue.Addr
					slice = append(slice, a)
					count -= 1
				}
			}
		}
	}
	fmt.Println("-------------------------")
	fmt.Println("       total")
	fmt.Println("cmha server", count, "/", peerssize)
	for i := range slice {
		fmt.Println(slice[i])
	}
	catalog := client.Catalog()
	health := client.Health()
	services, _, err := catalog.Services(nil)
	if err != nil {
		fmt.Println("cluster Query services failure!", err)
		return err
	}
	var chap = 0
	var db = 0
	s := len(services)
	s -= 1
	var isdbhealthy = true
	var ischaphealthy = true
	for k, _ := range services {
		if k != "consul" {
			service, _, err := catalog.Service(k, "", nil)
			if err != nil {
				fmt.Println("cluster Query catalog service "+k, err)
				return err
			}
			for _, chapvalue := range service {
				if chapvalue.ServiceTags[0] == "chap-master" || chapvalue.ServiceTags[0] == "chap-slave" {
					chaphealthservice, _, err := health.Service(k, chapvalue.ServiceTags[0], false, nil)
					if err != nil {
						fmt.Println("cluster Query cluster health service failure!", err)
						return err
					}
					for chapindex := range chaphealthservice {
						for chapcheckindex := range chaphealthservice[chapindex].Checks {
							if strings.EqualFold(chapvalue.Address, chaphealthservice[chapindex].Node.Address) {
								if chaphealthservice[chapindex].Checks[chapcheckindex].Status == "critical" {
									ischaphealthy = false

									break
								} else {
									ischaphealthy = true
								}
							}
						}
						if ischaphealthy {
							chap += 1
						}
					}
				}
			}
			for _, dbvalue := range service {
				if dbvalue.ServiceTags[0] == "master" || dbvalue.ServiceTags[0] == "slave" {
					dbhealthservice, _, err := health.Service(k, dbvalue.ServiceTags[0], false, nil)
					if err != nil {
						fmt.Println("cluster Query cluster health service failure!", err)
						return err
					}
					for dbindex := range dbhealthservice {
						for dbcheckindex := range dbhealthservice[dbindex].Checks {
							if strings.EqualFold(dbvalue.Address, dbhealthservice[dbindex].Node.Address) {
								if dbhealthservice[dbindex].Checks[dbcheckindex].Status == "critical" {
									isdbhealthy = false

									break
								} else {
									isdbhealthy = true
								}
							}
						}
						if isdbhealthy {
							db += 1
						}
					}
				}
			}
		}
	}
	s *= 2
	fmt.Println("cmha chap  ", chap, "/", s)
	fmt.Println("cmha db    ", db, "/", s)
	fmt.Println("-------------------------")
	isdbhealthy = true
	ischaphealthy = true
	for k, _ := range services {
		chapslice1 := []string{}
		dbslice1 := []string{}
		if k != "consul" {
			service, _, err := catalog.Service(k, "", nil)
			if err != nil {
				fmt.Println("cluster Query catalog service "+k, err)
				return err
			}
			chap = 0
			db = 0
			fmt.Println("      ", k)
			for _, chapvalue := range service {
				if chapvalue.ServiceTags[0] == "chap-master" || chapvalue.ServiceTags[0] == "chap-slave" {
					chaphealthservice, _, err := health.Service(k, chapvalue.ServiceTags[0], false, nil)
					if err != nil {
						fmt.Println("cluster Query cluster health service failure!", err)
						return err
					}
					for chapindex := range chaphealthservice {
						for chapcheckindex := range chaphealthservice[chapindex].Checks {
							if strings.EqualFold(chapvalue.Address, chaphealthservice[chapindex].Node.Address) {
								if chaphealthservice[chapindex].Checks[chapcheckindex].Status == "critical" {
									ischaphealthy = false
									a := "Fault machine: " + chaphealthservice[chapindex].Node.Address
									chapslice1 = append(chapslice1, a)
									//fmt.Println("故障机器:", chaphealthservice[chapindex].Node.Address)
									break
								} else {
									ischaphealthy = true
								}
							}
						}
						if ischaphealthy {
							chap += 1
						}
					}
				}
			}
			for _, dbvalue := range service {
				if dbvalue.ServiceTags[0] == "master" || dbvalue.ServiceTags[0] == "slave" {
					dbhealthservice, _, err := health.Service(k, dbvalue.ServiceTags[0], false, nil)
					if err != nil {
						fmt.Println("cluster Query cluster health service failure!", err)
						return err
					}
					for dbindex := range dbhealthservice {
						for dbcheckindex := range dbhealthservice[dbindex].Checks {
							if strings.EqualFold(dbvalue.Address, dbhealthservice[dbindex].Node.Address) {
								if dbhealthservice[dbindex].Checks[dbcheckindex].Status == "critical" {
									isdbhealthy = false
									a := "Fault machine: " + dbhealthservice[dbindex].Node.Address
									dbslice1 = append(dbslice1, a)
									//fmt.Println("故障机器:", dbhealthservice[dbindex].Node.Address)
									break
								} else {
									isdbhealthy = true
								}
							}
						}
						if isdbhealthy {
							db += 1
						}
					}
				}

			}
			fmt.Println("cmha chap  ", chap, "/", 2)
			for chapi1 := range chapslice1 {
				fmt.Println(chapslice1[chapi1])
			}
			fmt.Println("cmha db    ", db, "/", 2)
			for dbi1 := range dbslice1 {
				fmt.Println(dbslice1[dbi1])
			}
			fmt.Println("-------------------------")
		}

	}
	return nil
}
