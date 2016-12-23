package info

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
)

var ncpu int
var HZ, _ = beego.AppConfig.Int("HZ")
var (
	user_diff_1   int
	system_diff_1 int
	idle_diff_1   int
	iowait_diff_1 int
	pswpin_diff   int
	pswpout_diff  int
	diff_recv     string
	diff_recvs    string
	diff_send     string
	diff_sends    string
	rd_ios_s      float64
	wr_ios_s      float64
	rkbs          float64
	wkbs          float64
	queue         float64
	await         float64
	svc_t         float64
	busy          float64
)
var proc_statpath = "/tmp/.realtime_cache/proc_stat"
var  proc_diskstatspath = "/tmp/.realtime_cache/proc_diskstats"
var proc_net_devpath = "/tmp/.realtime_cache/proc_net_dev"
var proc_vmstatpath = "/tmp/.realtime_cache/proc_vmstat"

func IntToFloat(a,b int)(float64,float64){
	c := strconv.FormatFloat(float64(a),'f',1,64)
        d := strconv.FormatFloat(float64(b),'f',1,64)
	e, _ := strconv.ParseFloat(c, 64)
        f, _ := strconv.ParseFloat(d, 64)
	return e,f

}

func GetSysInfo() {
	GetNcpu()
	total_2, sys_cpu2 := SysCpu1Total1("/proc/stat", "cpu")
	user_diff := UserDiff(sys_cpu1, sys_cpu2)
	system_diff := SystemDiff(sys_cpu1, sys_cpu2)
	idle_diff := IdleDiffAndIowaitDiff(sys_cpu1, sys_cpu2, 4)
	iowait_diff := IdleDiffAndIowaitDiff(sys_cpu1, sys_cpu2, 5)
	total_diff := total_2 - total_1
	user_diff_f,total_diff_f := IntToFloat(user_diff,total_diff)
	user_diff_1 = int(user_diff_f / total_diff_f * 100 + 0.5)
	system_diff_f,_ := IntToFloat(system_diff,total_diff)
	system_diff_1 = int(system_diff_f / total_diff_f * 100 + 0.5)
	idle_diff_f,_ := IntToFloat(idle_diff,total_diff)
	idle_diff_1 = int(idle_diff_f / total_diff_f * 100 + 0.5)
	iowait_diff_f,_ := IntToFloat(iowait_diff,total_diff)
	iowait_diff_1 = int(iowait_diff_f / total_diff_f * 100 + 0.5)
	WriteProcStat(sys_cpu2, proc_statpath)
	deltams := (float64(user_diff + system_diff + idle_diff + iowait_diff)) * 1000.0 / float64(ncpu) / float64(HZ)
	sys_io2 := SysDisk("/proc/diskstats", "dm-0")
	rd_ios := ComputeFloat(sys_io1, sys_io2, 3)
	//	rd_merges := RdIos(sys_io1, sys_io2, 4)
	_ = RdIos(sys_io1, sys_io2, 4)
	rd_sectors := ComputeFloat(sys_io1, sys_io2, 5)
	rd_ticks := ComputeFloat(sys_io1, sys_io2, 6)
	wr_ios := ComputeFloat(sys_io1, sys_io2, 7)
	//	wr_merges := RdIos(sys_io1, sys_io2, 8)
	_ = RdIos(sys_io1, sys_io2, 8)
	wr_sectors := ComputeFloat(sys_io1, sys_io2, 9)
	wr_ticks := ComputeFloat(sys_io1, sys_io2, 10)
	ticks := ComputeFloat(sys_io1, sys_io2, 12)
	aveq := RdIos(sys_io1, sys_io2, 13)
	n_ios := StringToFloat64(rd_ios) + StringToFloat64(wr_ios)
	n_ticks := StringToFloat64(rd_ticks) + StringToFloat64(wr_ticks)
	n_kbytes := StringToFloat64(rd_sectors) + StringToFloat64(wr_sectors)*2.0
	queue = float64(aveq) / deltams
	if n_ios != 0 {
		//size := n_kbytes / n_ios
		_ = n_kbytes / n_ios
		await = n_ticks / n_ios
		svc_t = StringToFloat64(ticks) / n_ios
	} else {
		//size := 0.0
		_ = 0.0
		await = 0.0
		svc_t = 0.0
	}
	busy = StringToFloat64(ticks) / deltams * float64(100)
	if busy > 100 {
		busy = 100.0
	}
	if StringToFloat64(rd_sectors) == 0 {
		rkbs = 0.0
	} else {
		rkbs = StringToFloat64(rd_sectors) / float64(2) / deltams * float64(1000)
	}
	if StringToFloat64(wr_sectors) == 0 {
		wkbs = 0.0
	} else {
		wkbs = StringToFloat64(wr_sectors) / float64(2) / deltams * float64(1000)
	}
	if StringToFloat64(rd_ios) == 0 {
		rd_ios_s = 0.0
	} else {
		rd_ios_s = StringToFloat64(rd_ios) / deltams * float64(1000)
	}
	if StringToFloat64(wr_ios) == 0 {
		wr_ios_s = 0.0
	} else {
		wr_ios_s = StringToFloat64(wr_ios) / deltams * float64(1000)
	}
	WriteProcStat(sys_io2, proc_diskstatspath)
	net_se := SysNetDevFirst("/proc/net/dev", net_dev)
	diff_recv = ComputeFloat(net_first, net_se, 0)
	diff_send = ComputeFloat(net_first, net_se, 1)

	if StringToFloat64(diff_recv) < 1024 {
		diff_recvs = diff_recv + "B/s"
	} else if StringToFloat64(diff_recv) > 1048576 {
		diff_recvs = strconv.FormatFloat(StringToFloat64(diff_recv)/1048576, 'f', 1, 64) + "MB/s"
	} else {
		diff_recvs = strconv.FormatFloat(StringToFloat64(diff_recv)/1024, 'f', 1, 64) + "KB/s"
	}

	if StringToFloat64(diff_send) < 1024 {
		diff_sends = diff_send + "B/s"
	} else if StringToFloat64(diff_send) > 1048576 {
		diff_sends = strconv.FormatFloat(StringToFloat64(diff_send)/1048576, 'f', 1, 64) + "MB/s"
	} else {
		diff_sends = strconv.FormatFloat(StringToFloat64(diff_send)/1024, 'f', 1, 64) + "KB/s"
	}
	WriteProcStat(net_se, proc_net_devpath)
	swap2 := GetSwap("/proc/vmstat")
	pswpin_diff = RdIos(swap1, swap2, 0)
	pswpout_diff = RdIos(swap1, swap2, 1)
	WriteProcStat(swap2, proc_vmstatpath)
}

func UserDiff(sys_cpu1, sys_cpu2 []string) int {
	var sys_cpu1_total int
	var sys_cpu2_total int
	for cpu1 := range sys_cpu1 {
		if cpu1 == 1 || cpu1 == 2 {
			sys_cpu1_total = sys_cpu1_total + StringToInt(sys_cpu1[cpu1])
		}
	}
	for cpu2 := range sys_cpu2 {
		if cpu2 == 1 || cpu2 == 2 {
			sys_cpu2_total = sys_cpu2_total + StringToInt(sys_cpu2[cpu2])
		}
	}
	return sys_cpu2_total - sys_cpu1_total
}

func SystemDiff(sys_cpu1, sys_cpu2 []string) int {
	var sys_cpu1_total int
	var sys_cpu2_total int
	for cpu1 := range sys_cpu1 {
		if cpu1 == 3 || cpu1 == 6 || cpu1 == 7 {
			sys_cpu1_total = sys_cpu1_total + StringToInt(sys_cpu1[cpu1])
		}
	}
	for cpu2 := range sys_cpu2 {
		if cpu2 == 3 || cpu2 == 6 || cpu2 == 7 {
			sys_cpu2_total = sys_cpu2_total + StringToInt(sys_cpu2[cpu2])
		}
	}
	return sys_cpu2_total - sys_cpu1_total
}

func IdleDiffAndIowaitDiff(sys_cpu1, sys_cpu2 []string, index int) int {
	var sys_cpu1_total int
	var sys_cpu2_total int
	for cpu1 := range sys_cpu1 {
		if cpu1 == index {
			sys_cpu1_total = sys_cpu1_total + StringToInt(sys_cpu1[cpu1])
		}
	}
	for cpu2 := range sys_cpu2 {
		if cpu2 == index {
			sys_cpu2_total = sys_cpu2_total + StringToInt(sys_cpu2[cpu2])
		}
	}
	return sys_cpu2_total - sys_cpu1_total
}

func WriteProcStat(writedata []string, filepath string) {
	var writestring string
	for i := range writedata {
		if writedata[i] != " " {
			writestring = writestring + writedata[i] + " "
		}
	}
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

func GetNcpu() {
	f, err := os.Open("/proc/cpuinfo")
	if err != nil {
		fmt.Println("open file failed!", err)
		os.Exit(2)
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	var counter int
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)

		if strings.Contains(line, "processor") {
			counter += 1
		}
		//handler(line)
		if err != nil {
			if err == io.EOF {
				ncpu = counter
				return
			}
			os.Exit(2)
		}
	}
}
