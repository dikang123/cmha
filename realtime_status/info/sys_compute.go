package info

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)


func RdIos(sys_io1, sys_io2 []string, index int) int {
	var sys_io1_total int
	var sys_io2_total int
	for io1 := range sys_io1 {
		if io1 == index {
			sys_io1_total = sys_io1_total + StringToInt(sys_io1[io1])
		}
	}
	for io2 := range sys_io2 {
		if io2 == index {
			sys_io2_total = sys_io2_total + StringToInt(sys_io2[io2])
		}
	}
	return int(sys_io2_total - sys_io1_total)
}

func ComputeFloat(sys_io1, sys_io2 []string, index int) string {
	var sys_io1_total float64
	var sys_io2_total float64
	for io1 := range sys_io1 {
		if io1 == index {
			sys_io1_total = sys_io1_total + StringToFloat64(sys_io1[io1])
		}
	}
	for io2 := range sys_io2 {
		if io2 == index {
			sys_io2_total = sys_io2_total + StringToFloat64(sys_io2[io2])
		}
	}
	floatvalue := float64(sys_io2_total - sys_io1_total)
	s := fmt.Sprintf("%.1f", floatvalue)
	return s
}

func ComputeFloat2(sys_io1, sys_io2 []string, index int) string {
	var sys_io1_total float64
	var sys_io2_total float64
	for io1 := range sys_io1 {
		if io1 == index {
			sys_io1_total = sys_io1_total + StringToFloat64(sys_io1[io1])
		}
	}
	for io2 := range sys_io2 {
		if io2 == index {
			sys_io2_total = sys_io2_total + StringToFloat64(sys_io2[io2])
		}
	}
	floatvalue := float64(sys_io2_total - sys_io1_total)
	s := fmt.Sprintf("%.2f", floatvalue)
	return s
}

func GetSwap(filename string) []string {
	f, err := os.Open(filename)
	if err != nil {
		return nil
	}
	buf := bufio.NewReader(f)
	var swaps []string
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if strings.Contains(line, "pswpin") || strings.Contains(line, "pswpout") {
			swap := strings.Split(line, " ")
			for i := range swap {
				if i == 1 {
					swaps = append(swaps, swap[i])
				}
			}

		}
		//handler(line)
		if err != nil {
			if err == io.EOF {
				return swaps
			}
			return nil
		}
	}
	return swaps
}
func StringToFloat64(cpu string) float64 {

	f, _ := strconv.ParseFloat(cpu, 64)
	return f
}
