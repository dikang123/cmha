package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/upmio/mysqlcheck/check"
)

var (
	UserFlag     = flag.String("user", "", "usage: the mysql user")
	PasswordFlag = flag.String("password", "", "usage: the mysql password")
)

func main() {
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "-v" || arg == "--version" {
			fmt.Println("1.1.5-Beta")
			return
		}
	}
	flag.Parse()
	if *UserFlag == "" || *PasswordFlag == "" {
		os.Exit(2)
	}
	user := *UserFlag
	password := *PasswordFlag
	check.IsPingType(user, password)
}
