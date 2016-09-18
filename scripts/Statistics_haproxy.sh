#!/bin/bash
####version 1.1.5-Beta
LOCAL_HOST="127.0.0.1"
SERVICE_NAME=$1
hostname=`hostname`
user=admin
password=1111
ha_port=9600
cache_file=/tmp/.haproxy_cache_status
	
curl -s -X GET -u $user:$password "http://$LOCAL_HOST:$ha_port/haproxy?stats;csv" 2>>/dev/null|tail -n3|/usr/local/bin/jq --slurp --raw-input --raw-output 'split("\n") | .[0:] | map(split(",")) |map({"name": .[1], "Queue_cur": .[2],"Queue_max": .[3],"Session_cur": .[4],"Session_max": .[5],"Session_limit": .[6],"Session_Total": .[7],"net_input_Bytes": .[8] ,"net_output_Bytes": .[9] ,"Denied_req": .[10] ,"Denied_resp": .[11] ,"error_req": .[12], "error_con": .[13] ,"error_resp": .[14],"Warnings_retr": .[15],"Warnings_redis": .[16],"status": .[17]})' 1>$cache_file 2>>/dev/null
	curl --connect-timeout 1 -s -X PUT -d @$cache_file http://$LOCAL_HOST:8500/v1/kv/cmha/service/$SERVICE_NAME/chap/status/$hostname >>/dev/null 2>&1

