#!/bin/bash
####version 0.6.0
CONSUL_IP=$1
killall -0 haproxy && killall -0 consul && killall -0 consul-template
if [ $? -ne 0 ]; then
	exit 2
fi

curl  --connect-timeout 3 -X GET http://${CONSUL_IP}:8500/v1/kv/service/innosql/leader > /dev/null 2>&1
if [ $? -ne 0 ]; then
	exit 2
fi
