package info

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/astaxie/beego"
	consulapi "github.com/hashicorp/consul/api"
)

var real_status_cachepath = "/tmp/.realtime_cache/real_status_cache"
var real_status_seqpath ="/tmp/.realtime_cache/real_status_seq"

type SysDb struct {
	Sys *SysInfo `json:"sys"`
	Db  *Dbinfo  `json:"db"`
}

type SyS struct {
	Sys *SysInfo `json:"sys"`
	Db struct{} `json:"db"`
}


type SysInfo struct {
	Id            string `json:"id"`
	Time          string `json:"time"`
	OneSecond     string `json:"one_m"`
	FiveSecond    string `json:"five_m"`
	FifteenSecond string `json:"fifteen_m"`
	Usr           string `json:"usr"`
	Sys           string `json:"sys"`
	Idle          string `json:"idle"`
	Iow           string `json:"iow"`
	Si            string `json:"si"`
	So            string `json:"so"`
	Recv          string `json:"recv"`
	Send          string `json:"send"`
	Rs            string `json:"r_s"`
	Ws            string `json:"w_s"`
	RkBs          string `json:"rkB_s"`
	WkBs          string `json:"wkB_s"`
	Queue         string `json:"queue"`
	Await         string `json:"await"`
	Svctm         string `json:"svctm"`
	Util          string `json:"util"`
	Ncpu	      string `json:"ncpu"`
}

type Dbinfo struct {
	Id   string `json:"id"`
	Time string `json:"time"`
	Ins  string `json:"ins"`
	Upd  string `json:"upd"`
	Del  string `json:"del"`
	Sel  string `json:"sel"`
	Qps  string `json:"qps"`
	Tps  string `json:"tps"`
	Lor  string `json:"lor"`
	Hit  string `json:"hit"`
	Run  string `json:"run"`
	Con  string `json:"con"`
	Cre  string `json:"cre"`
	Cac  string `json:"cac"`
}

func InsertDataToCs(servicename, hostname string, new_id, old_id int, C_time,types string) {
	rs := []*SysDb{}
	var sysdb SysDb
	var sys SyS
	sy := []*SyS{}
	if types == "sys_and_db" {
		sysdb = SysDb{
			Sys: &SysInfo{
				Id:            strconv.Itoa(new_id),
				Time:          C_time,
				OneSecond:     one_m,
				FiveSecond:    five_m,
				FifteenSecond: fifteen_m,
				Usr:           strconv.Itoa(user_diff_1),
				Sys:           strconv.Itoa(system_diff_1),
				Idle:          strconv.Itoa(idle_diff_1),
				Iow:           strconv.Itoa(iowait_diff_1),
				Si:            strconv.Itoa(pswpin_diff),
				So:            strconv.Itoa(pswpout_diff),
				Recv:          diff_recvs,
				Send:          diff_sends,
				Rs:            strconv.FormatFloat(rd_ios_s, 'f', 1, 64),
				Ws:            strconv.FormatFloat(wr_ios_s, 'f', 1, 64),
				RkBs:          strconv.FormatFloat(rkbs, 'f', 1, 64),
				WkBs:          strconv.FormatFloat(wkbs, 'f', 1, 64),
				Queue:         strconv.FormatFloat(queue, 'f', 1, 64),
				Await:         strconv.FormatFloat(await, 'f', 1, 64),
				Svctm:         strconv.FormatFloat(svc_t, 'f', 1, 64),
				Util:          strconv.FormatFloat(busy, 'f', 1, 64),
				Ncpu:	       strconv.Itoa(ncpu),
			},
			Db: &Dbinfo{
				Id:   strconv.Itoa(new_id),
				Time: C_time,
				Ins:  strconv.Itoa(Com_insert),
				Upd:  strconv.Itoa(Com_update),
				Del:  strconv.Itoa(Com_delete),
				Sel:  strconv.Itoa(Com_select),
				Qps:  strconv.Itoa(Com_select),
				Tps:  strconv.Itoa(TPS),
				Lor:  strconv.Itoa(Innodb_buffer_pool_read_requests),
				Hit:  strconv.FormatFloat(hit, 'f', 2, 64),
				Run:  strconv.Itoa(Threads_running),
				Con:  strconv.Itoa(Threads_connected),
				Cre:  strconv.Itoa(Threads_created),
				Cac:  strconv.Itoa(Threads_cached),
			},
		}
		rs = append(rs,&sysdb)
		sysdbinfo, _ := json.Marshal(rs)
		WriteJson(sysdbinfo, real_status_cachepath)
		PutCs(servicename, hostname, sysdbinfo)
	} else if types == "sys" {
		sys = SyS{
			Sys: &SysInfo{
				Id:            strconv.Itoa(new_id),
				Time:          C_time,
				OneSecond:     one_m,
				FiveSecond:    five_m,
				FifteenSecond: fifteen_m,
				Usr:           strconv.Itoa(user_diff_1),
				Sys:           strconv.Itoa(system_diff_1),
				Idle:          strconv.Itoa(idle_diff_1),
				Iow:           strconv.Itoa(iowait_diff_1),
				Si:            strconv.Itoa(pswpin_diff),
				So:            strconv.Itoa(pswpout_diff),
				Recv:          diff_recvs,
				Send:          diff_sends,
				Rs:            strconv.FormatFloat(rd_ios_s, 'f', 1, 64),
				Ws:            strconv.FormatFloat(wr_ios_s, 'f', 1, 64),
				RkBs:          strconv.FormatFloat(rkbs, 'f', 1, 64),
				WkBs:          strconv.FormatFloat(wkbs, 'f', 1, 64),
				Queue:         strconv.FormatFloat(queue, 'f', 1, 64),
				Await:         strconv.FormatFloat(await, 'f', 1, 64),
				Svctm:         strconv.FormatFloat(svc_t, 'f', 1, 64),
				Util:          strconv.FormatFloat(busy, 'f', 1, 64),
				Ncpu:          strconv.Itoa(ncpu),
			},
		}
		sy = append(sy,&sys)
		sysinfo, _ := json.Marshal(sy)
		WriteJson(sysinfo, real_status_cachepath)
		PutCs(servicename, hostname, sysinfo)
	}
	WriteOldId(old_id, real_status_seqpath)
	os.Exit(0)
}

func WriteJson(byteinfo []byte, filepath string) {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("openfile error:", err, filepath)
		return
	}
	defer file.Close()
	err = ioutil.WriteFile(filepath, byteinfo, os.ModePerm)
	if err != nil {
		fmt.Println("write file error:", err)
		return
	}
}

func PutCs(servicename, hostname string, byteinfo []byte) {
	config := &consulapi.Config{
		Datacenter: beego.AppConfig.String("datacenter"),
		Token:      beego.AppConfig.String("token"),
		Address:    beego.AppConfig.String("LOCAL_HOST") + ":8500",
	}
	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Println("consulapi create failed!", err)
		os.Exit(2)
	}
	kv := client.KV()
	put := consulapi.KVPair{
		Key:   "cmha/service/" + servicename + "/real_status/" + hostname + "/1",
		Value: byteinfo,
	}
	_, err = kv.Put(&put, nil)
	if err != nil {
		fmt.Println("put kv failed:", err)
		os.Exit(2)
	}
}

func WriteOldId(writedata int, filepath string) {
	var writestring = strconv.Itoa(writedata)
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("openfile error:", err, filepath)
		os.Exit(2)
	}
	defer file.Close()
	err = ioutil.WriteFile(filepath, []byte(writestring), os.ModePerm)
	if err != nil {
		fmt.Println("write file error:", err)
		os.Exit(2)
	}
}
