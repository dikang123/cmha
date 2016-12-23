package info

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var one_m, five_m, fifteen_m string

func CpuLoad() {
	line := Load("/proc/loadavg")

	for i := range line {
		if i == 0 {
			one_m = line[i]
		} else if i == 1 {
			five_m = line[i]
		} else if i == 2 {
			fifteen_m = line[i]
		}
	}
}

func Load(filename string) []string {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("open file failed!", err, filename)
		os.Exit(2)
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)

		//if strings.Contains(line, types) {
		load := strings.Split(line, " ")
		return load
		//}
		//handler(line)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			os.Exit(2)
		}
	}
	return nil
}
