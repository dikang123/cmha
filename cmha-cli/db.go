package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/ryanuber/columnize"
	"github.com/upmio/cmha-cli/cliconfig"
)

var kv *consulapi.KV

func Db(args ...string) error {
	if len(args) > 0 {

		client, err := cliconfig.Consul_Client_Init()
		if err != nil {
			fmt.Println("db Create consul-api client failure!", err)
			return err
		}
		health := client.Health()
		healthservice, _, err := health.Service(args[0], "", false, nil)
		if err != nil {
			fmt.Println("db Query cluster health service failure!", err)
			return err
		}
		var dbhealthy = true
		if len(healthservice) < 1 {
			fmt.Println("not ", args[0], " service")
			return nil
		}
		var d = 0
		fmt.Println("---------------------------------------------------------------------------------------------------------------------")
		dbslice := make([]string, 0, 10)
		a := "Node|Address|Service|Status|Vsr|Read-Only|Repl-Status|Type|repl_err_counter"
		dbslice = append(dbslice, a)

		username := cliconfig.GetUserName()
		password := cliconfig.GetPassword()

		for index := range healthservice {
			var dbcount = 0
			if healthservice[index].Service.Tags[0] == "master" ||
				healthservice[index].Service.Tags[0] == "slave" {

				var critical = 0
				var VSR = "err"
				var READ_ONLY = "err"
				var REPL_STATUS = "err"
				for checkindex := range healthservice[index].Checks {
					if healthservice[index].Checks[checkindex].Status == "critical" {
						critical += 1
					}
				}

				if critical == 0 {
					dsName := username + ":" +
						password +
						"@tcp(" + healthservice[index].Node.Address + ":" +
						strconv.Itoa(healthservice[index].Service.Port) +
						")/" + "?timeout=3"

					db, err := sql.Open("mysql", dsName)

					if err != nil {
						fmt.Println("open database failure", err)
						return err
					}
					defer db.Close()
					err = db.Ping()
					if err != nil {
						fmt.Println("connection to the database failure", err)
						goto Here
					}
					sql := "show variables like" + "'" + "rpl_semi_sync_master_keepsyncrepl" + "'"
					sql1 := "show variables like" + "'" + "rpl_semi_sync_master_trysyncrepl" + "'"
					row, err := db.Query(sql)
					if err != nil {
						fmt.Println("query rpl_semi_sync_master_keepsyncrepl failure!", err)
						return err
					}
					rows, err := db.Query(sql1)
					if err != nil {
						fmt.Println("query rpl_semi_sync_master_trysyncrepl failure!", err)
						return err
					}
					colss, _ := rows.Columns()
					cols, _ := row.Columns()
					buffer1 := make([]interface{}, len(colss))
					data1 := make([]interface{}, len(colss))
					buffer := make([]interface{}, len(cols))
					data := make([]interface{}, len(cols))
					for i, _ := range buffer {
						buffer[i] = &data[i]
					}
					for i1, _ := range buffer1 {
						buffer1[i1] = &data1[i1]
					}

					for row.Next() {
						err = row.Scan(buffer...)
						if err != nil {
							fmt.Println("scan() traversal rpl_semi_sync_master_keepsyncrepl failure!", err)
							return err
						}
					}
					mapField2Data := make(map[string]interface{}, len(cols))
					for k, col := range data {
						mapField2Data[cols[k]] = col
					}
					rpl_semi_sync_master_keepsyncrepl := mapField2Data["Value"]
					for rows.Next() {
						err = rows.Scan(buffer1...)
						if err != nil {
							fmt.Println("scan() traversal rpl_semi_sync_master_trysyncrepl failure!", err)
							return err
						}
					}
					mapField2Data1 := make(map[string]interface{}, len(colss))
					for k1, co := range data1 {
						mapField2Data1[colss[k1]] = co
					}
					rpl_semi_sync_master_trysyncrepl := mapField2Data1["Value"]

					if string(rpl_semi_sync_master_keepsyncrepl.([]uint8)) == "ON" &&
						string(rpl_semi_sync_master_trysyncrepl.([]uint8)) == "ON" {

						VSR = "ON"
					} else {
						VSR = "OFF"
					}
					rowss, err := db.Query("show variables like 'read_only'")
					if err != nil {
						fmt.Println("query read_only failure!", err)
						return err
					}
					colsss, _ := rowss.Columns()
					buffer2 := make([]interface{}, len(colsss))
					data2 := make([]interface{}, len(colsss))
					for i2, _ := range buffer2 {
						buffer2[i2] = &data2[i2]
					}
					for rowss.Next() {
						err = rowss.Scan(buffer2...)
						if err != nil {
							fmt.Println("scan() traversal read_only failure!", err)
							return err
						}
					}
					mapField2Data2 := make(map[string]interface{}, len(colsss))
					for k2, co1 := range data2 {
						mapField2Data2[colsss[k2]] = co1
					}
					read_only := mapField2Data2["Value"]

					if string(read_only.([]uint8)) == "ON" {
						READ_ONLY = "readonly"
					} else {
						READ_ONLY = "readwrite"
					}
					r, err := db.Query("show slave status")
					if err != nil {
						fmt.Println("query slave status failure!", err)
						return err
					}
					c, _ := r.Columns()
					b := make([]interface{}, len(c))
					d1 := make([]interface{}, len(c))
					for i3, _ := range b {
						b[i3] = &d1[i3]
					}
					for r.Next() {
						err = r.Scan(b...)
						if err != nil {
							fmt.Println("scan() traversal slave status failure!", err)
							return err
						}
					}
					mapField2Data3 := make(map[string]interface{}, len(c))
					for k3, co2 := range d1 {
						mapField2Data3[c[k3]] = co2
					}
					Slave_SQL_Running := mapField2Data3["Slave_SQL_Running"]
					Slave_IO_Running := mapField2Data3["Slave_IO_Running"]

					if string(Slave_IO_Running.([]uint8)) == "Yes" &&
						string(Slave_SQL_Running.([]uint8)) == "Yes" {

						REPL_STATUS = "OK"
					} else {
						REPL_STATUS = "IO:" + string(Slave_IO_Running.([]uint8)) +
							";SQL:" + string(Slave_SQL_Running.([]uint8))
					}
				}
				kv = client.KV()

			Here:
				var unknown = 0
				for checkindex := range healthservice[index].Checks {
					var Value string
					key := "cmha/service/" + args[0] + "/db/" + healthservice[index].Checks[checkindex].Node + "/repl_err_counter"
					kvpair, _, err := kv.Get(key, nil)
					if err != nil {
						fmt.Println("Get key "+key+" failure!", err)
						return err
					}
					if kvpair == nil {
						fmt.Println("No key " + key + " or key value is null")
						return nil
					}
					Value = string(kvpair.Value)
					if healthservice[index].Checks[checkindex].Status == "critical" {
						dbhealthy = false
						dbcount += 1
						if healthservice[index].Checks[checkindex].CheckID == "serfHealth" {
							a := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s",
								healthservice[index].Node.Node,
								healthservice[index].Node.Address,
								"magent",
								"Fail",
								" ",
								" ",
								" ",
								"agent",
								" ")

							dbslice = append(dbslice, a)
						} else {
							a := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s",
								healthservice[index].Node.Node,
								healthservice[index].Node.Address,
								args[0],
								"Fail",
								VSR,
								READ_ONLY,
								REPL_STATUS,
								"db",
								Value)

							dbslice = append(dbslice, a)
						}
					} else {
						if dbcount == 0 {
							dbhealthy = true
						}
						if healthservice[index].Checks[checkindex].CheckID == "serfHealth" {
							a := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s",
								healthservice[index].Node.Node,
								healthservice[index].Node.Address,
								"magent",
								"OK",
								" ",     //vsr
								" ",     //read-only
								" ",     //repl_status
								"agent", //type
								" ")

							dbslice = append(dbslice, a)
						} else {
							for checkindex := range healthservice[index].Checks {
								if healthservice[index].Checks[checkindex].CheckID == "serfHealth" {
									if healthservice[index].Checks[checkindex].Status == "critical" {
										unknown += 1
									}
								}
							}
							if unknown != 0 {
								a := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s",
									healthservice[index].Node.Node,
									healthservice[index].Node.Address,
									args[0],
									"UnKnown",
									VSR,
									READ_ONLY,
									REPL_STATUS,
									"db",
									Value)

								dbslice = append(dbslice, a)
							} else {
								if healthservice[index].Checks[checkindex].Status != "passing" {
									a := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s",
										healthservice[index].Node.Node,
										healthservice[index].Node.Address,
										args[0],
										healthservice[index].Checks[checkindex].Status,
										VSR,
										READ_ONLY,
										REPL_STATUS,
										"db",
										Value)
									dbslice = append(dbslice, a)

								} else {
									a := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s",
										healthservice[index].Node.Node,
										healthservice[index].Node.Address,
										args[0],
										"OK",
										VSR,
										READ_ONLY,
										REPL_STATUS,
										"db",
										Value)

									dbslice = append(dbslice, a)
								}
							}
						}
					}
				}
				if dbhealthy {
					d += 1
				}
			}
		}
		keys, _, err := kv.Keys("", "", nil)
		if err != nil {
			fmt.Println("dbleader get keys failure!", err)
			return err
		}
		if keys == nil {
			fmt.Println("dbleader not kv!")
			return nil
		}
		key := "cmha/service/" + args[0] + "/db/leader"
		var iskey = false
		for value := range keys {
			if key == keys[value] {
				iskey = true
				break
			} else {
				continue
			}
		}
		if !iskey {
			fmt.Println("not ", args[0], " kv")
			return nil
		}
		kvpair, _, err := kv.Get(key, nil)
		if err != nil {
			fmt.Println("dbleader Get key failure!", err)
			return err
		}
		if kvpair.Session == "" {
			fmt.Println(args[0], " leader not exist!")
			return nil
		}
		fmt.Println("cmha db  ", d, "/", 2, "Leader", string(kvpair.Value))
		output := columnize.SimpleFormat(dbslice)
		fmt.Println(string(output))
		fmt.Println("---------------------------------------------------------------------------------------------------------------------")
		return nil
	}

	client, err := cliconfig.Consul_Client_Init()

	if err != nil {
		fmt.Println("db Create consul-api client failure!", err)
		return err
	}
	catalog := client.Catalog()
	health := client.Health()
	services, _, err := catalog.Services(nil)
	if err != nil {
		fmt.Println("db Query services failure!", err)
		return err
	}
	var ishealthy = true
	for k, _ := range services {
		dbslice1 := []string{}
		if k != "consul" {
			service, _, err := catalog.Service(k, "", nil)
			if err != nil {
				fmt.Println("db Query catalog service "+k, err)
				return err
			}
			var count = 0
			//var issusehost string
			fmt.Println("-------------------------")
			fmt.Println("      ", k)
			for _, value := range service {
				if value.ServiceTags[0] == "master" || value.ServiceTags[0] == "slave" {
					healthservice, _, err := health.Service(k, value.ServiceTags[0], false, nil)
					if err != nil {
						fmt.Println("db Check cluster health service failure!", err)
						return err
					}
					for index := range healthservice {
						for checkindex := range healthservice[index].Checks {
							if strings.EqualFold(value.Address, healthservice[index].Node.Address) {
								if healthservice[index].Checks[checkindex].Status == "critical" {
									ishealthy = false
									a := "Fault machine: " + healthservice[index].Node.Address
									dbslice1 = append(dbslice1, a)
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
			fmt.Println("cmha db  ", count, "/", 2)
			for dbi1 := range dbslice1 {
				fmt.Println(dbslice1[dbi1])
			}
		}
	}
	fmt.Println("-------------------------")
	return nil
}
