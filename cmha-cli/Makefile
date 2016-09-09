.PHONY: build pull  build-deps

projectname := cmha-cli

GOROOT := /usr/lib/golang

wsp := $(shell echo ${WORKSPACE})
ifneq (${wsp} , )
	GOPATH := ${wsp}
else 
	GOPATH := $(shell echo ${GOPATH})
endif


build: pull build-deps
	
	GOPATH=${GOPATH} go fmt 
	GOPATH=${GOPATH} go build 
	GOPATH=${GOPATH} go install 

build-deps:
	#GOPATH=${GOPATH} go get github.com/GeertJohan/go.linenoise
	GOPATH=${GOPATH} go get github.com/0-T-0/go.linenoise
	GOPATH=${GOPATH} go get github.com/astaxie/beego
	GOPATH=${GOPATH} go get github.com/fatih/color
	GOPATH=${GOPATH} go get github.com/go-sql-driver/mysql
	GOPATH=${GOPATH} go get github.com/hashicorp/consul/api
	GOPATH=${GOPATH} go get github.com/olekukonko/tablewriter
	GOPATH=${GOPATH} go get github.com/ryanuber/columnize

pull:
	git pull
