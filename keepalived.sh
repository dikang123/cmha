#!/bin/bash
####version 1.1.2
CONSUL_IP1=$1
CONSUL_IP2=$2
CONSUL_IP3=$3
SERVICE_NAME=$4
killall -0 haproxy && killall -0 consul && killall -0 consul-template
if [ $? -ne 0 ]; then
	exit 2
fi

curl  --connect-timeout 3 -X GET http://${CONSUL_IP1}:8500/v1/kv/service/${SERVICE_NAME}/leader > /dev/null 2>&1
if [ $? -ne 0 ]; then
	curl  --connect-timeout 3 -X GET http://${CONSUL_IP2}:8500/v1/kv/service/${SERVICE_NAME}/leader > /dev/null 2>&1
	if [ $? -ne 0 ]; then
        	curl  --connect-timeout 3 -X GET http://${CONSUL_IP3}:8500/v1/kv/service/${SERVICE_NAME}/leader > /dev/null 2>&1
		if [ $? -ne 0 ]; then
        		exit 2
		fi
	fi
fi
