#!/bin/bash
<<<<<<< HEAD
####version 1.1.5-Beta
=======
####version 1.1.5-Beta.1
>>>>>>> 126d33b0306a2de4f2f5445489f9e46636c7c67e
LOCAL_HOST="127.0.0.1"
SERVICE_NAME=$1
killall -0 haproxy && killall -0 consul && killall -0 consul-template && killall -0 keepalived
if [ $? -ne 0 ]; then
	exit 2
fi

curl  --connect-timeout 1 -X GET http://${LOCAL_HOST}:8500/v1/kv/cmha/service/${SERVICE_NAME}/db/leader > /dev/null 2>&1
if [ $? -ne 0 ]; then
	echo "GET DB LEADER FAIL"
        exit 2
else
	echo "check ok"
fi

