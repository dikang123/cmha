#!/bin/bash
<<<<<<< HEAD
#version 1.1.5-Beta
=======
#version 1.1.5-Beta.1
>>>>>>> 126d33b0306a2de4f2f5445489f9e46636c7c67e
VIP=$1
HOSTNAME=`hostname`
LOCAL_HOST="127.0.0.1"
STATUS=$2
SERVICE_NAME=$3

backup_status(){
	curl -X PUT -d 'backup' http://${LOCAL_HOST}:8500/v1/kv/cmha/service/${SERVICE_NAME}/chap/role/${HOSTNAME} >>/dev/null
	ip address | grep $VIP 
	if [ $? = 0 ]; then
		ip -s -s a f to $VIP/32
	fi
}

master_status(){
	curl -X PUT -d 'master' http://${LOCAL_HOST}:8500/v1/kv/cmha/service/${SERVICE_NAME}/chap/role/${HOSTNAME} >>/dev/null
}

fault_status(){
	curl -X PUT -d 'fault' http://${LOCAL_HOST}:8500/v1/kv/cmha/service/${SERVICE_NAME}/chap/role/${HOSTNAME} >>/dev/null
	ip -s -s a f to $VIP/32
}


case $STATUS in 

	 master)
		master_status
	     ;;

	backup)
		backup_status	
	;;

	fault)
		fault_status
	;;
esac
