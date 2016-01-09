#!/bin/bash
####version 1.1.2
config_dir="/usr/local/cmha/consul.d"
configure_file="${config_dir}/config.json"
bin_file="/usr/local/bin/consul"
addr=`grep -w "rpc" ${configure_file}  |awk '{print $2}' |sed 's/\"//g'`
case $1 in
     stop)
#	id=`ps -ef | grep  consul | grep agent | awk '{print $2}'`
#	kill -9 $id
#	count=`ps -ef | grep  consul | grep agent |wc -l`
	kill -15 `pidof consul`
        if [ $? -eq 0 ]; then
                echo "consul stop success"
                exit 0
        else
                echo "consul stop failure!"
		exit 1
        fi
	;;
     start)
	if [ ! -f ${configure_file} ] && [ ! -x ${bin_file} ]; then
		echo "Consul configure file or binary not exit,Please configure!"
		exit 1
	fi
	consul agent -config-dir $config_dir & >/dev/null 2>&1
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
	
#	if ! $(echo $2 | grep -q '^[0-9".]\+$' >/dev/null); then
#		echo "rpc addr error,format reload ip"
#		exit 1
#	fi
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
