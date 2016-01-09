#!/bin/bash
####version 1.1.2
config_dir="/usr/local/cmha/consul-template.d"
configure_file="${config_dir}/haproxy.ctmpl"
bin_file="/usr/local/bin/consul-template"
haproxy_file="/etc/haproxy/haproxy.cfg"
ip=$2
case $1 in
     stop)
#	id=`ps -ef | grep  consul | grep agent | awk '{print $2}'`
#	kill -9 $id
#	count=`ps -ef | grep  consul | grep agent |wc -l`
	kill -15 `pidof consul-template`
        if [ $? -eq 0 ]; then
                echo "consul-template stop success"
                exit 0
        else
                echo "consul-template stop failure!"
		exit 1
        fi
	;;
     start)
	if [ ! -f ${configure_file} ] && [ ! -x ${bin_file} ]; then
		echo "Consul-template configure file or binary not exit,Please configure!"
		exit 1
	fi
	echo ${ip} ${configure_file} ${haproxy_file}
	consul-template -consul ${ip}:8500 -syslog=true -syslog-facility=LOCAL5 -template "${configure_file}:${haproxy_file}:service haproxy restart" & >/dev/null 2>&1
	
#	consul agent -config-dir $config_dir & >/dev/null 2>&1
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
	echo $"Usage: $0 {start|stop|reload}"
	exit 1
	;;
esac
