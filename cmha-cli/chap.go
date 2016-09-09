package main

import (
	"fmt"
	"strings"

	"github.com/ryanuber/columnize"
	"github.com/upmio/cmha-cli/cliconfig"
)

func Chap(args ...string) error {
	if len(args) > 0 {
		client, err := cliconfig.Consul_Client_Init()

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
		var chaphealthy = true
		if len(healthservice) < 1 {
			fmt.Println("not ", args[0], " service")
			return nil
		}
		var c = 0
		fmt.Println("---------------------------------------------------------------------")
		chapslice := make([]string, 0, 10)
		a := "Node|Address|Service|Status|Type"
		chapslice = append(chapslice, a)
		for index := range healthservice {
			var chapcount = 0
			if healthservice[index].Service.Tags[0] == "chap-master" || healthservice[index].Service.Tags[0] == "chap-slave" {
				var unknown = 0
				for checkindex := range healthservice[index].Checks {
					if healthservice[index].Checks[checkindex].Status == "critical" {
						chaphealthy = false
						chapcount += 1
						if healthservice[index].Checks[checkindex].CheckID == "serfHealth" {
							a := fmt.Sprintf("%s|%s|%s|%s|%s", healthservice[index].Node.Node, healthservice[index].Node.Address, "hagent", "Fail", "chap")
							chapslice = append(chapslice, a)
						} else {
							a := fmt.Sprintf("%s|%s|%s|%s|%s", healthservice[index].Node.Node, healthservice[index].Node.Address, args[0], "Fail", "chap")
							chapslice = append(chapslice, a)
						}
					} else {
						if chapcount == 0 {
							chaphealthy = true
						}
						if healthservice[index].Checks[checkindex].CheckID == "serfHealth" {
							a := fmt.Sprintf("%s|%s|%s|%s|%s", healthservice[index].Node.Node, healthservice[index].Node.Address, "hagent", "OK", "chap")
							chapslice = append(chapslice, a)
						} else {
							for checkindex := range healthservice[index].Checks {
								if healthservice[index].Checks[checkindex].CheckID == "serfHealth" {
									if healthservice[index].Checks[checkindex].Status == "critical" {
										unknown += 1
									}
								}
							}
							if unknown != 0 {
								a := fmt.Sprintf("%s|%s|%s|%s|%s", healthservice[index].Node.Node, healthservice[index].Node.Address, args[0], "UnKnown", "chap")
								chapslice = append(chapslice, a)
							} else {
								if healthservice[index].Checks[checkindex].Status != "passing" {
									a := fmt.Sprintf("%s|%s|%s|%s|%s", healthservice[index].Node.Node, healthservice[index].Node.Address, args[0], healthservice[index].Checks[checkindex].Status, "chap")
									chapslice = append(chapslice, a)
								} else {
									a := fmt.Sprintf("%s|%s|%s|%s|%s", healthservice[index].Node.Node, healthservice[index].Node.Address, args[0], "OK", "chap")
									chapslice = append(chapslice, a)
								}
							}
						}
					}
				}
				if chaphealthy {
					c += 1
				}
			}
		}
		fmt.Println("cmha chap  ", c, "/", 2)
		output := columnize.SimpleFormat(chapslice)
		fmt.Println(string(output))
		fmt.Println("---------------------------------------------------------------------")
		return nil

	}

	client, err := cliconfig.Consul_Client_Init()
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
									ishealthy = false
									a := "Fault machine: " + healthservice[index].Node.Address
									chapslice1 = append(chapslice1, a)
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
