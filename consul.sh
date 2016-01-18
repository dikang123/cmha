#!/bin/bash
####version 1.1.5
config_dir="/usr/local/cmha/consul.d"
configure_file="${config_dir}/config.json"
bin_file="/usr/local/bin/consul"
addr=`grep -w "rpc" ${configure_file}  |awk '{print $2}' |sed 's/\"//g'`
pid=`pidof consul`
case $1 in
     stop)
	if [ "$pid" != "" ];then
	kill -15 $pid
        if [ $? -eq 0 ]; then
                echo "consul stop success"
                exit 0
        else
                echo "consul stop failure!"
		exit 1
        fi
	else
		echo "consul not running!"
		exit 0
	fi
	;;
     start)
	if [ "$pid" != "" ];then
		echo "consul is running!"
		exit 0
	fi
	if [ ! -f ${configure_file} ] && [ ! -x ${bin_file} ]; then
		echo "Consul configure file or binary not exit,Please configure!"
		exit 1
	fi
	consul agent -config-dir $config_dir --bind=$addr & >/dev/null 2>&1
	count=`ps -ef | grep  consul | grep agent |wc -l`
	if [ ${count} -eq 0 ]; then
		echo "consul start failure"
		exit 1
	else
		echo "consul start success!"
		exit 0
	fi
	;;
     reload)
	if [ "$pid" = "" ];then
		 echo "consul not running!"
                exit 0
	fi
	if [ ! -f ${configure_file} ] && [ ! -x ${bin_file} ]; then
                echo "Consul configure file or binary not exit,Please configure!"
                exit 1
        fi
	consul $1 -rpc-addr=${addr}:8400 >/dev/null 2>&1
	if [ $? -ne 0 ]; then
		echo "consul reload failure"
		exit 1
	else
		echo "consul reload success"	
		exit 0
	fi
	;;
     *)
	echo $"Usage: $0 {start|stop|reload}"
	exit 1
	;;
esac
