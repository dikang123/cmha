#!/bin/bash
####version 1.1.5
config_dir="/usr/local/cmha/consul-template.d"
configure_file="${config_dir}/haproxy.ctmpl"
bin_file="/usr/local/bin/consul-template"
haproxy_file="/etc/haproxy/haproxy.cfg"
ip=`grep -w "rpc" /usr/local/cmha/consul.d/config.json|awk '{print $2}' |sed 's/\"//g'`
pid=`pidof consul-template`
case $1 in
     stop)
	if [ "$pid" != "" ];then
		kill -15 $pid
        		if [ $? -eq 0 ]; then
        		        echo "consul-template stop success"
        		        exit 0
        		else
        		        echo "consul-template stop failure!"
				exit 1
        		fi
	else
		echo "consul-template not running!"
	fi
	;;
     start)
	if [ "$pid" != "" ];then
		echo "consul-template is running!"
		exit 0
	fi
	if [ ! -f ${configure_file} ] && [ ! -x ${bin_file} ]; then
		echo "Consul-template configure file or binary not exit,Please configure!"
		exit 1
	fi
	consul-template -consul ${ip}:8500 -syslog=true -syslog-facility=LOCAL5 -template "${configure_file}:${haproxy_file}:service haproxy restart" & >/dev/null 2>&1
	count=`ps -ef | grep  consul-template |grep haproxy |wc -l`
	if [ ${count} -eq 0 ]; then
		echo "consul start failure"
		exit 1
	else
		echo "consul start success!"
		exit 0
	fi
	;;
     *)
	echo $"Usage: $0 {start|stop}"
	exit 1
	;;
esac
