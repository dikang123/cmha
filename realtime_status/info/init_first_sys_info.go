package info

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/upmio/realtime_status/file"
)

var sys_cpu1 []string
var sys_io1 []string
var total_1 int
var net_first []string
var swap1 []string
var net_dev = beego.AppConfig.String("net_dev")

func InitFirstSysinfo() {
	isproc, _ := file.PathExists(proc_statpath)

	if isproc {
		total_1, sys_cpu1 = SysCpu1Total1(proc_statpath, "cpu")
	} else {
		total_1, sys_cpu1 = SysCpu1Total1("/proc/stat", "cpu")
		
	}
	isswap, _ := file.PathExists(proc_vmstatpath)
	if isswap {
		swap1 = ReadTmpFile(proc_vmstatpath)
	} else {
		swap1 = append(swap1, "0")
		swap1 = append(swap1, "0")
	}
	isnet, _ := file.PathExists(proc_net_devpath)
	if isnet {
		net_first = ReadTmpFile(proc_net_devpath)
	} else {
		net_first = SysNetDevFirst("/proc/net/dev", net_dev)

	}
	isdisk, _ := file.PathExists(proc_diskstatspath)
	if isdisk {
		sys_io1 = ReadTmpFile(proc_diskstatspath)
	} else {
		sys_io1 = SysDisk("/proc/diskstats", "dm-0")
	}
}

func SysCpu1Total1(filename, types string) (int, []string) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("open file failed!", err)
		return 0, nil
	}
	defer f.Close()
	var cpus []string 
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)

		if strings.Contains(line, types) {
			cpu := strings.Split(line, " ")
			var total_1 int
			for i := range cpu {
				if cpu[i] != ""{
					cpus = append(cpus,cpu[i])
					if i != 0 || i != 1 {
						total_1 = total_1 + StringToInt(cpu[i])
					}
				}
			}
			return total_1, cpus
		}
		//handler(line)
		if err != nil {
			if err == io.EOF {
				return 0, nil
			}
			return 0, nil
		}
	}
	return 0, nil
}

func SysNetDevFirst(filename, types string) []string {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("open file failed!", err)
		return nil
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		var net_first []string
		var slicemap []string
		if strings.Contains(line, types) {
			net := strings.Split(line, " ")
			//var total_1 int
			for i := range net {
				net[i] = strings.Replace(net[i], " ", "", -1)
				if net[i] != "" {
					slicemap = append(slicemap, net[i])
				}
			}
			for s := range slicemap {
				if s == 1 || s == 9 {
					net_first = append(net_first, slicemap[s])
				}
			}
			return net_first
		}
		//handler(line)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return nil
		}
	}
	return nil
}

func ReadTmpFile(filename string) []string {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("open file failed!", err)
		return nil
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)

		//if strings.Contains(line, types) {
		lines := strings.Split(line, " ")
		return lines
		//}
		//handler(line)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return nil
		}
	}
	return nil
}

func SysDisk(filename, types string) []string {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("open file failed!", err)
		return nil
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if strings.Contains(line, types) {
			disk := strings.Split(line, " ")
			var disks []string
			for i := range disk {
				if disk[i] != "" {
					disks = append(disks, disk[i])
				}
			}
			return disks
		}
		//handler(line)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return nil
		}
	}
	return nil
}

func StringToInt(cpu string) int {
	total, _ := strconv.Atoi(cpu)
	return total
}
