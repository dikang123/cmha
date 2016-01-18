#!/bin/bash
###version 1.1.4
scp_error(){
	if [ $? -ne 0 ]; then
        	echo "Scp $1 fail!"
		exit 113
        fi
}

error(){
if [ $? -ne 0 ]; then
	echo "$2 Error: $3: $4 error"
fi
if [ "$1" = "" ]; then
                echo "$2 Error: $3: $4 is empty or categories error"
                exit 111
fi
}

sed_error(){
if [ $1 -ne 0 ]; then
	echo "$2 Error: $3: replace $4 failure"
	exit 112
fi
}

bootstrap_expect=`grep "bootstrap_expect" auto-deployment.ini |awk -F '=' '{print $2}'`
error "$bootstrap_expect" "bootstrap_expect" "auto-deployment.ini" "bootstrap_expect"
db_number=`grep "db_number" auto-deployment.ini |awk -F '=' '{print $2}'`
error "$db_number" "db_number" "auto-deployment.ini" "db_number"
if [ ${bootstrap_expect} -eq 3 ]; then
	server1_ip=`awk -F '=' '/\[cmha_server\]/{a=1}a==1&&$1~/server1-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $1}'`
	error "$server1_ip" "server1_ip" "auto-deployment.ini" "server1-ip-hostname"
	server2_ip=`awk -F '=' '/\[cmha_server\]/{a=1}a==1&&$1~/server2-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $1}'`
	error "$server2_ip" "server2_ip" "auto-deployment.ini" "server2-ip-hostname"
	server3_ip=`awk -F '=' '/\[cmha_server\]/{a=1}a==1&&$1~/server3-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $1}'`
	error "$server3_ip" "server3_ip" "auto-deployment.ini" "server3-ip-hostname"
elif [ ${bootstrap_expect} -eq 5 ]; then
	server1_ip=`awk -F '=' '/\[cmha_server\]/{a=1}a==1&&$1~/server1-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $1}'`
	error "$server1_ip" "server1_ip" "auto-deployment.ini" "server1-ip-hostname"

	server2_ip=`awk -F '=' '/\[cmha_server\]/{a=1}a==1&&$1~/server2-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $1}'`
	error "$server2_ip" "server2_ip" "auto-deployment.ini" "server2-ip-hostname"
	server3_ip=`awk -F '=' '/\[cmha_server\]/{a=1}a==1&&$1~/server3-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $1}'`
	error "$server3_ip" "server3_ip" "auto-deployment.ini" "server3-ip-hostname"
	server4_ip=`awk -F '=' '/\[cmha_server\]/{a=1}a==1&&$1~/server4-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $1}'`
	error "$server4_ip" "server4_ip" "auto-deployment.ini" "server4-ip-hostname"
	server5_ip=`awk -F '=' '/\[cmha_server\]/{a=1}a==1&&$1~/server5-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $1}'`
	error "$server5_ip" "server5_ip" "auto-deployment.ini" "server5-ip-hostname"
fi
if [ ${bootstrap_expect} -eq 3 ]; then
	s_ip=${server1_ip}";"${server2_ip}";"${server3_ip}
elif [ ${bootstrap_expect} -eq 5 ]; then
	s_ip=${server1_ip}";"${server2_ip}";"${server3_ip}";"${server4_ip}";"${server5_ip}
fi

datacenter=`grep -w "datacenter" auto-deployment.ini |awk -F '=' '{print $2}'`
error "$datacenter" "datacenter" "auto-deployment.ini" "datacenter"
enable_rsyslog=`grep -w "enable-rsyslog" auto-deployment.ini |awk -F '=' '{print $2}'`
error "$enable_rsyslog" "enable-rsyslog" "auto-deployment.ini" "enable-rsyslog"
rsyslog_facility=`grep -w "rsyslog-facility" auto-deployment.ini |awk -F '=' '{print $2}'`
error "$rsyslog_facility" "rsyslog-facility" "auto-deployment.ini" "rsyslog-facility"




read_consul_server_config(){

for i in $(seq 1 ${bootstrap_expect})
do 
	local server_hostname=`awk -F '=' '/\[cmha_server\]/{a=1}a==1&&$1~/server'${i}'-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $2}'`
	error "$server_hostname" "Configure Consul Server" "auto-deployment.ini" "server${i}-ip-hostname"
	local server_ip=`awk -F '=' '/\[cmha_server\]/{a=1}a==1&&\$1~/server'${i}'-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $1}'`
	error "$server_ip" "Configure Consul Server" "auto-deployment.ini" "server${i}-ip-hostname"
	local server_password=`awk -F '=' '/\[cmha_server\]/{a=1}a==1&&$1~/server'${i}'-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $4}'`
	error "$server_password" "Configure Consul Server" "auto-deployment.ini" "server${i}-ip-hostname"
	local server_username=`awk -F '=' '/\[cmha_server\]/{a=1}a==1&&$1~/server'${i}'-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $3}'`
	error "$server_username" "Configure Consul Server" "auto-deployment.ini" "server${i}-ip-hostname"
	expect expect/expect.exp ${server_username} ${server_ip} ${server_password} "mkdir -p" "/usr/local/cmha/log/" >/dev/null 2>&1
	if [ $? -ne 0 ]; then
		echo "${server_ip} mkdir /usr/local/cmha/log/ failure"
		exit 146
	fi
	expect expect/syslog.exp ${server_username} ${server_ip} ${server_password} "/usr/local/cmha/log/consul.log" ${rsyslog_facility} >/dev/null 2>&1
	local count=0
	while read LINE
	do
        	if $(echo ${LINE} | grep -w 'bootstrap_expect' >/dev/null); then
                	local server_autodeployment_bootstrap_expect=`grep "bootstrap_expect" auto-deployment.ini |awk -F '=' '{print $2}'`
			error "$server_autodeployment_bootstrap_expect" "Configure Consul Server" "auto-deployment.ini" "bootstrap_expect"
	                local server_config_bootstrap_expect=`grep "bootstrap_expect" conf/config.json |sed 's/.$//'`
			error "$server_config_bootstrap_expect" "Configure Consul Server" "conf/config.json" "bootstrap_expect"
			sed -i "s/${server_config_bootstrap_expect}/  \"bootstrap_expect\": ${server_autodeployment_bootstrap_expect}/" conf/config.json
			sed_error "$?" "Configure Consul Server" "conf/config.json" "bootstrap_expect"
	        elif $(echo ${LINE} | grep -w 'server' >/dev/null); then
        	        sed -i "s/${LINE}/\"server\": true,/g" conf/config.json
			sed_error "$?" "Configure Consul Server" "conf/config.json" "server"
	        elif $(echo ${LINE} | grep -w 'datacenter' >/dev/null); then
                	local server_config_datacenter=`grep "datacenter" conf/config.json |awk '{print $2}'|sed 's/.$//' |sed 's/\"//g'`
			error "server_config_datacenter" "Configure Consul Server" "conf/config.json" "datacenter"
	                sed -i "s/${server_config_datacenter}/${datacenter}/g" conf/config.json
			sed_error "$?" "Configure Consul Server" "conf/config.json" "datacenter"
		elif $(echo ${LINE} | grep -w 'data_dir' >/dev/null); then
			local server_config_data_dir=`grep "data_dir" conf/config.json |awk '{print $2}'|sed 's/.$//' |sed 's/\"//g'`
                        error "$server_config_data_dir" "Configure Consul Server" "conf/config.json" "node_name"
#                        sed -i "s/  'data_dir': '/tmp/consul'/  'data_dir': '/usr/local/cmha/consul_data'/g" conf/config.json
			sed -i "s/  \"data_dir\": \"\/tmp\/consul\"/  \"data_dir\": \"\/usr\/local\/cmha\/consul_data\"/g" conf/config.json
                        sed_error "$?" "Configure Consul Server" "conf/config.json" "data_dir"
        	elif $(echo ${LINE} | grep -w 'node_name' >/dev/null); then
                	local server_config_node_name=`grep "node_name" conf/config.json |awk '{print $2}'|sed 's/.$//' |sed 's/\"//g'`
			error "$server_config_node_name" "Configure Consul Server" "conf/config.json" "node_name"
			sed -i "s/${server_config_node_name}/${server_hostname}/g" conf/config.json
			sed_error "$?" "Configure Consul Server" "conf/config.json" "node_name"
        	elif $(echo ${LINE} | grep -w 'log_level' >/dev/null); then
                	local server_autodeployment_log_level=`grep "log_level" auto-deployment.ini |awk -F '=' '{print $2}'`
			error "$server_autodeployment_log_level" "Configure Consul  Server" "auto-deployment.ini" "log_level"
	                local server_config_log_level=`grep "log_level" conf/config.json |awk '{print $2}'|sed 's/.$//' |sed 's/\"//g'`
			error "$server_config_log_level" "Configure Consul Server" "conf/config.json" "log_level"
			sed -i "s/${server_config_log_level}/${server_autodeployment_log_level}/g" conf/config.json
			sed_error "$?" "Configure Consul Server" "conf/config.json" "log_level"
		elif $(echo ${LINE} | grep -w 'enable_syslog' >/dev/null); then
			sed -i "s/${LINE}/\"enable_syslog\": ${enable_rsyslog},/g" conf/config.json
			sed_error "$?" "Configure Consul Server" "conf/config.json" "enable_syslog"
		elif $(echo ${LINE} | grep -w 'syslog_facility' >/dev/null); then
			sed -i "s/${LINE}/\"syslog_facility\": \"${rsyslog_facility}\",/g" conf/config.json
			sed_error "$?" "Configure Consul Server" "conf/config.json" "syslog_facility"
	        elif $(echo ${LINE} | grep -w 'http' >/dev/null); then
                	local server_config_http=` grep "http" conf/config.json |sed 's/.$//'`
			error "$server_config_http" "Configure Consul Server" "conf/config.json" "http"
	                sed -i "s/${server_config_http}/    \"http\": \"${server_ip}\"/" conf/config.json
			sed_error "$?" "Configure Consul Server" "conf/config.json" "http"
        	elif $(echo ${LINE} | grep -w 'rpc' >/dev/null); then
	                local server_config_rpc=`grep "rpc" conf/config.json |sed 's/.$//'`
			error "$server_config_rpc" "Configure Consul Server" "conf/config.json" "rpc"
        	        sed -i "s/${server_config_rpc}/    \"rpc\": \"${server_ip}/" conf/config.json
			sed_error "$?" "Configure Consul Server" "conf/config.json" "rpc"
	        elif $(echo ${LINE} | grep -q '^[0-9".,]\+$' >/dev/null) && [ ${bootstrap_expect} -eq 3 ]; then
        	       	if [ ${count} -eq 0 ]; then
                       		local server_config_start_join1=`echo ${LINE} |sed 's/.$//' |sed 's/\"//g'|sed 's/ //g'`
				error "$server_config_start_join1" "Configure Consul Server" "conf/config.json" "start_join 1"
	                       	sed -i "s/${server_config_start_join1}/${server1_ip}/" conf/config.json
				sed_error "$?" "Configure Consul Server" "conf/config.json" "start_join 1"
        	               	let "count+=1"
	                elif [ ${count} -eq 1 ]; then
       	                	local server_config_start_join2=`echo ${LINE} | sed 's/.$//' |sed 's/\"//g'|sed 's/ //g'`
				error "$server_config_start_join2" "Configure Consul Server" "conf/config.json" "start_join 2"
               	        	sed -i "s/${server_config_start_join2}/${server2_ip}/" conf/config.json
				sed_error "$?" "Configure Consul Server" "conf/config.json" "start_join 2"
                       		let "count+=1"
                	elif [ ${count} -eq 2 ]; then
                        	local server_config_start_join3=`echo ${LINE} |sed 's/\"//g'|sed 's/ //g'`
				error "$server_config_start_join3" "Configure Consul Server" "conf/config.json" "start_join 3"
                        	sed -i "s/${server_config_start_join3}/${server3_ip}/" conf/config.json
				sed_error "$?" "Configure Consul Server" "conf/config.json" "start_join 3"
                	fi
		elif $(echo ${LINE} | grep -q '^[0-9".,]\+$' >/dev/null) && [ ${bootstrap_expect} -eq 5 ]; then
                        if [ ${count} -eq 0 ]; then 
                                local server_config_start_join1=`echo ${LINE} |sed 's/.$//' |sed 's/\"//g'|sed 's/ //g'`
				error "$server_config_start_join1" "Configure Consul Server" "conf/config.json" "start_join 1"
                                sed -i "s/${server_config_start_join1}/${server1_ip}/" conf/config.json
				sed_error "$?" "Configure Consul Server" "conf/config.json" "start_join 1"
                                let "count+=1"
                        elif [ ${count} -eq 1 ]; then 
                                local server_config_start_join2=`echo ${LINE} | sed 's/.$//' |sed 's/\"//g'|sed 's/ //g'`
				error "$server_config_start_join2" "Configure Consul Server" "conf/config.json" "start_join 2"
                                sed -i "s/${server_config_start_join2}/${server2_ip}/" conf/config.json
				sed_error "$?" "Configure Consul Server" "conf/config.json" "start_join 2"
                                let "count+=1"
                        elif [ ${count} -eq 2 ]; then 
                                local server_config_start_join3=`echo ${LINE} |sed 's/\"//g'|sed 's/ //g' |sed 's/.$//'`
				error "$server_config_start_join3" "Configure Consul  Server" "conf/config.json" "start_join 3"
                                sed -i "s/${server_config_start_join3}/${server3_ip}/" conf/config.json
				sed_error "$?" "Configure Consul Server" "conf/config.json" "start_join 3"
				let "count+=1"
			elif [ ${count} -eq 3 ]; then
				local server_config_start_join4=`echo ${LINE} |sed 's/\"//g'|sed 's/ //g'| sed 's/.$//'`
				error "$server_config_start_join4" "Configure Consul Server" "conf/config.json" "start_join 4"
                                sed -i "s/${server_config_start_join4}/${server4_ip}/" conf/config.json
				sed_error "$?" "Configure Consul Server" "conf/config.json" "start_join 4"
                                let "count+=1"
			elif [ ${count} -eq 4 ]; then
				local server_config_start_join5=`echo ${LINE} |sed 's/\"//g'|sed 's/ //g'`
				error "$server_config_start_join5" "Configure Consul Server" "conf/config.json" "start_join 5"
                                sed -i "s/${server_config_start_join5}/${server5_ip}/" conf/config.json
				sed_error "$?" "Configure Consul Server" "conf/config.json" "start_join 5"
                        fi      

        	fi
	done < conf/config.json
	expect expect/expect.exp ${server_username} ${server_ip} ${server_password} "mkdir" "/usr/local/cmha/consul.d/" >/dev/null 2>&1
	if [ $? -ne 0 ]; then
		echo "${server_ip} mkdir /usr/local/cmha/consul.d/ failure"
		exit 146
	fi
	expect expect/scp.exp ${server_username} ${server_ip} ${server_password} "conf/config.json" /usr/local/cmha/consul.d/  > /dev/null 2>&1
	scp_error "conf/config.json"
	expect expect/scp.exp ${server_username} ${server_ip} ${server_password} "server-remove.sh" /usr/local/cmha/  > /dev/null 2>&1
	scp_error "server-remove.sh"
	echo ${server_ip} consul server configure success
	expect expect/scp.exp ${server_username} ${server_ip} ${server_password} "bin/consul" /usr/local/bin/  > /dev/null 2>&1 
	scp_error "consul"
	expect expect/consul.exp ${server_username} ${server_ip} ${server_password} "consul agent -config-dir /usr/local/cmha/consul.d/ -bind=${server_ip} &" >/dev/null 2>&1
	expect expect/scp.exp ${server_username} ${server_ip} ${server_password} "consul.sh" /usr/local/cmha/ >/dev/null 2>&1 
	scp_error "consul.sh"
	expect expect/pross.exp ${server_username} ${server_ip} ${server_password} "pgrep consul" >/dev/null 2>&1
	if [ $? -ne 0 ];then
		echo "${server_ip} consul server start fail!"
		exit 120 
	else
		echo "${server_ip} consul server start success"
	fi
done
}

read_consul_agent_config(){
for i in $(seq 1 ${db_number})
do
	for a in chap mysql
	do
		for j in master slave
		do 
			local hostname=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-'$j'-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $2}'`
			error "${hostname}" "Configure Consul  Agent" "auto-deployment.ini" "$a-$j-ip-hostname"
       	 		local ip=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-'$j'-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $1}'`
			error "$ip" "Configure Consul  Agent" "auto-deployment.ini" "$a-$j-ip-hostname"
        		local password=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-'$j'-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $4}'`
			error "$password" "Configure Consul Agent" "auto-deployment.ini" "$a-$j-ip-hostname"
        		local username=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-'$j'-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $3}'`
			error "$username" "Configure Consul Agent" "auto-deployment.ini" "$a-$j-ip-hostname"
			expect expect/expect.exp ${username} ${ip} ${password} "mkdir -p" "/usr/local/cmha/log/" >/dev/null 2>&1
			if [ $? -ne 0 ]; then
                		echo "${ip} mkdir /usr/local/cmha/log/ failure"
                        	exit 146
                	fi
			expect expect/syslog.exp ${username} ${ip} ${password} "/usr/local/cmha/log/consul.log" ${rsyslog_facility} >/dev/null 2>&1
        		local count=0
        		while read LINE
        		do
                		if $(echo ${LINE} | grep -w 'bootstrap_expect' >/dev/null); then
                        		local agent_config_bootstrap_expect=`grep "bootstrap_expect" conf/config.json |awk '{print $2}'|sed 's/.$//' |sed 's/\"//g'`
					error "$agent_config_bootstrap_expect" "Configure Consul Agent" "conf/config.json" "bootstrap_expect"
                        		sed -i "/\"bootstrap_expect\": ${agent_config_bootstrap_expect},/d" conf/config.json
					sed_error "$?" "Configure Consul Agent" "conf/config.json" "bootstrap_expect"
                		elif $(echo ${LINE} | grep -w 'server' >/dev/null); then
                        		sed -i "s/${LINE}/\"server\": false,/g" conf/config.json
					sed_error "$?" "Configure Consul Agent" "conf/config.json" "server"
                		elif $(echo ${LINE} | grep -w 'datacenter' >/dev/null); then
                        		local agent_autodeployment_datacenter=`grep "datacenter" auto-deployment.ini |awk -F '=' '{print $2}'`
					error "$agent_autodeployment_datacenter" "Configure Consul Agent" "auto-deployment.ini" "datacenter"
                        		local agent_config_datacenter=`grep "datacenter" conf/config.json |awk '{print $2}'|sed 's/.$//' |sed 's/\"//g'`
					error "$agent_config_datacenter" "Configure Consul Agent" "conf/config.json" "datacenter"
                        		sed -i "s/${agent_config_datacenter}/${agent_autodeployment_datacenter}/g" conf/config.json
					sed_error "$?" "Configure Consul Agent" "conf/config.json" "datacenter"
                		elif $(echo ${LINE} | grep -w 'node_name' >/dev/null); then
                       			local agent_config_node_name=`grep "node_name" conf/config.json |awk '{print $2}'|sed 's/.$//' |sed 's/\"//g'`
					error "$agent_config_node_name" "Configure Consul Agent" "conf/config.json" "node_name"
                        		sed -i "s/${agent_config_node_name}/${hostname}/g" conf/config.json
					sed_error "$?" "Configure Consul Agent" "conf/config.json" "node_name"
                		elif $(echo ${LINE} | grep -w 'log_level' >/dev/null); then
                        		local agent_autodeployment_log_level=`grep "log_level" auto-deployment.ini |awk -F '=' '{print $2}'`
					error "$agent_autodeployment_log_level" "Configure Consul Agent" "auto-deployment.ini" "log_level"
                        		local agent_config_log_level=`grep "log_level" conf/config.json |awk '{print $2}'|sed 's/.$//' |sed 's/\"//g'`
					error "$agent_config_log_level" "Configure Consul Agent" "conf/config.json" "log_level"
                        		sed -i "s/${agent_config_log_level}/${agent_autodeployment_log_level}/g" conf/config.json
					sed_error "$?" "Configure Consul Agent" "conf/config.json" "log_level"
				elif $(echo ${LINE} | grep -w 'enable_syslog' >/dev/null); then
                        		sed -i "s/${LINE}/\"enable_syslog\": ${enable_rsyslog},/g" conf/config.json
                        		sed_error "$?" "Configure Consul Server" "conf/config.json" "enable_syslog"
                		elif $(echo ${LINE} | grep -w 'syslog_facility' >/dev/null); then
                        		sed -i "s/${LINE}/\"syslog_facility\": \"${rsyslog_facility}\",/g" conf/config.json
                        		sed_error "$?" "Configure Consul Server" "conf/config.json" "syslog_facility"
                		elif $(echo ${LINE} | grep -w 'http' >/dev/null); then
                        		local agent_config_http=` grep "http" conf/config.json |sed 's/.$//'`
					error "$agent_config_http" "Configure Consul Agent" "conf/config.json" "http"
                        		sed -i "s/${agent_config_http}/    \"http\": \"${ip}\"/" conf/config.json
					sed_error "$?" "Configure Consul Agent" "conf/config.json" "http"
                		elif $(echo ${LINE} | grep -w 'rpc' >/dev/null); then
                        		local agent_config_rpc=`grep "rpc" conf/config.json |sed 's/.$//'`
					error "$agent_config_rpc" "Configure Consul Agent" "conf/config.json" "rpc"
                        		sed -i "s/${agent_config_rpc}/    \"rpc\": \"${ip}/" conf/config.json
					sed_error "$?" "Configure Consul Agent" "conf/config.json" "rpc"
                		elif $(echo ${LINE} | grep -q '^[0-9".,]\+$' >/dev/null) && [ ${bootstrap_expect} -eq 3 ]; then
                        		if [ ${count} -eq 0 ]; then
                                		local agent_config_start_join1=`echo ${LINE} |sed 's/.$//' |sed 's/\"//g'|sed 's/ //g'`
						error "$agent_config_start_join1" "Configure Consul Agent" "conf/config.json" "start_join 1"
                                		sed -i "s/${agent_config_start_join1}/${server1_ip}/" conf/config.json
						sed_error "$?" "Configure Consul Agent" "conf/config.json" "start_join 1"
                                		let "count+=1"
                        		elif [ ${count} -eq 1 ]; then
                                		local agent_config_start_join2=`echo ${LINE} | sed 's/.$//' |sed 's/\"//g'|sed 's/ //g'`
						error "$agent_config_start_join2" "Configure Consul Agent" "conf/config.json" "start_join 2"
                                		sed -i "s/${agent_config_start_join2}/${server2_ip}/" conf/config.json
						sed_error "$?" "Configure Consul Agent" "conf/config.json" "start_join 2"
                                		let "count+=1"
                        		elif [ ${count} -eq 2 ]; then
                                		local agent_config_start_join3=`echo ${LINE} |sed 's/\"//g'|sed 's/ //g'`
						error "$agent_config_start_join3" "Configure Consul Agent" "conf/config.json" "start_join 3"
                                		sed -i "s/${agent_config_start_join3}/${server3_ip}/" conf/config.json
						sed_error "$?" "Configure Consul Agent" "conf/config.json" "start_join 3"
                        		fi
				elif $(echo ${LINE} | grep -q '^[0-9".,]\+$' >/dev/null) && [ ${bootstrap_expect} -eq 5 ]; then
                        		if [ ${count} -eq 0 ]; then
                                		local server_config_start_join1=`echo ${LINE} |sed 's/.$//' |sed 's/\"//g'|sed 's/ //g'`
						error "$server_config_start_join1" "Configure Consul Agent" "conf/config.json" "start_join 1"
                                		sed -i "s/${server_config_start_join1}/${server1_ip}/" conf/config.json
						sed_error "$?" "Configure Consul Agent" "conf/config.json" "start_join 1"
                                		let "count+=1"
                        		elif [ ${count} -eq 1 ]; then
                                		local server_config_start_join2=`echo ${LINE} | sed 's/.$//' |sed 's/\"//g'|sed 's/ //g'`
						error "$server_config_start_join2" "Configure Consul Agent" "conf/config.json" "start_join 2"
                                		sed -i "s/${server_config_start_join2}/${server2_ip}/" conf/config.json
						sed_error "$?" "Configure Consul Agent" "conf/config.json" "start_join 2"
                                		let "count+=1"
                        		elif [ ${count} -eq 2 ]; then
                                		local server_config_start_join3=`echo ${LINE} |sed 's/\"//g'|sed 's/ //g' |sed 's/.$//'`
						error "$server_config_start_join3" "Configure Consul Agent" "conf/config.json" "start_join 3"
                                		sed -i "s/${server_config_start_join3}/${server3_ip}/" conf/config.json
						sed_error "$?" "Configure Consul Agent" "conf/config.json" "start_join 3"
                                		let "count+=1"
                        		elif [ ${count} -eq 3 ]; then
                                		local server_config_start_join4=`echo ${LINE} |sed 's/\"//g'|sed 's/ //g'| sed 's/.$//'`
						error "$server_config_start_join4" "Configure Consul Agent" "conf/config.json" "start_join 4"
                                		sed -i "s/${server_config_start_join4}/${server4_ip}/" conf/config.json
						sed_error "$?" "Configure Consul Agent" "conf/config.json" "start_join 4"
                                		let "count+=1"
                        		elif [ ${count} -eq 4 ]; then
                                		local server_config_start_join5=`echo ${LINE} |sed 's/\"//g'|sed 's/ //g'`
						error "$server_config_start_join5" "Configure Consul Agent" "conf/config.json" "start_join 5"
                                		sed -i "s/${server_config_start_join5}/${server5_ip}/" conf/config.json
						sed_error "$?" "Configure C0nsul Agent" "conf/config.json" "start_join 5"
                        		fi

                		fi
			done < conf/config.json
			expect expect/expect.exp ${username} ${ip} ${password} "mkdir -p" "/usr/local/cmha/consul.d/" >/dev/null 2>&1
        		if [ $? -ne 0 ]; then
                		echo "${ip} mkdir /usr/local/cmha/consul.d/ failure"
                		exit 146
        		fi
			expect expect/scp.exp ${username} ${ip} ${password} "conf/config.json" /usr/local/cmha/consul.d/  >/dev/null 2>&1 
		        scp_error "Agent: conf/config.json"
			echo "${ip} consul agent configure success"
        		expect expect/scp.exp ${username} ${ip} ${password} "bin/consul" /usr/local/bin/  >/dev/null 2>&1 
        		scp_error "Agent: consul"
        		expect expect/consul.exp ${username} ${ip} ${password} "consul agent -config-dir /usr/local/cmha/consul.d/ -bind=${ip} &" >/dev/null 2>&1 >/tmp/consul.log
			expect expect/pross.exp ${username} ${ip} ${password} "pgrep consul" >/dev/null 2>&1
        		if [ $? -ne 0 ];then
                		echo "${ip} consul agent start fail!"
                		exit 120
        		else
				echo "${ip} consul agent start success"
			fi
		done
	done
done
}

install_mysql(){
local dir="/usr/local/cmha/mysql/"
local destdir="/usr/local/cmha/"
local mysqlversion=`grep -w "mysqlversion" auto-deployment.ini |awk -F '=' '{print $2}'`
local mysqldir=${mysqlversion%.tar*}
echo ${mysqldir}
local mysql="package/${mysqlversion}"
local count=1
local offset=1
for i in $(seq 1 ${db_number})
do
	for j in master slave
	do
		local ip=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&\$1~/mysql-'$j'-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $1}'`
		error "$ip" "Install Mysql" "auto-deployment.ini" "mysql-$j-ip-hostname"
		local hostname=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-'$j'-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $2}'`
		error "$hostname" "Install Mysql" "auto-deployment.ini" "mysql-$j-ip-hostname"
		local password=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-'$j'-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $4}'`
		error "$password" "Install Mysql" "auto-deployment.ini" "mysql-$j-ip-hostname"
		local username=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-'$j'-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $3}'`
		error "$username" "Install Mysql" "auto-deployment.ini" "mysql-$j-ip-hostname"
		local mysql_port=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-'$j'-port/{print $2;exit}' auto-deployment.ini`
		error "$mysql_port" "Install Mysql" "auto-deployment.ini" "mysql-$j-port"
#		local mysql_rpl_reverse_recover_enabled=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-'$j'-rpl_reverse_recover_enabled/{print $2;exit}' auto-deployment.ini`
#		error "$mysql_rpl_reverse_recover_enabled" "Install Mysql" "auto-deployment.ini" "mysql-$j-rpl_reverse_recover_enabled"
		local mysql_ha_partner_host=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-'$j'-ha_partner_host/{print $2;exit}' auto-deployment.ini`
		error "$mysql_ha_partner_host" "Install Mysql" "auto-deployment.ini" "mysql-$j-ha_partner_host"
		local mysql_ha_partner_port=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-'$j'-ha_partner_port/{print $2;exit}' auto-deployment.ini`
		error "$mysql_ha_partner_port" "Install Mysql" "auto-deployment.ini" "mysql-$j-ha_partner_port"
		local mysql_ha_partner_user=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-'$j'-ha_partner_user/{print $2;exit}' auto-deployment.ini`
		error "$mysql_ha_partner_user" "Install Mysql" "auto-deployment.ini" "mysql-$j-ha_partner_user"
		local mysql_ha_partner_password=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-'$j'-ha_partner_password/{print $2;exit}' auto-deployment.ini`
		error "$mysql_ha_partner_password" "Install Mysql" "auto-deployment.ini" "mysql-$j-ha_partner_password"
		local mysql_password=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-'$j'-password/{print $2;exit}' auto-deployment.ini`
		error "$mysql_password" "Install Mysql" "auto-deployment.ini" "mysql-$j-password"
		local mysql_username=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-'$j'-username/{print $2;exit}' auto-deployment.ini`
		error "$mysql_username" "Install Mysql" "auto-deployment.ini" "mysql-$j-username"
		local cmha_check_username=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-master-cmha-check-username/{print $2;exit}' auto-deployment.ini`
	#	echo $cmha_check_username
		error "$cmha_check_username" "Install Mysql" "auto-deployment.ini" "mysql-master-cmha-check-username"
		local cmha_check_password=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-master-cmha-check-password/{print $2;exit}' auto-deployment.ini`
#		echo $cmha_check_password
		error "$cmha_check_password" "Install Mysql" "auto-deployment.ini" "mysql-master-cmha-check-password"
		local basedir=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-master-basedir/{print $2;exit}' auto-deployment.ini`
		error "$basedir" "Install Mysql" "auto-deployment.ini" "mysql-master-basedir"
		local datadir=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-master-datadir/{print $2;exit}' auto-deployment.ini`
		error "$datadir" "Install Mysql" "auto-deployment.ini" "mysql-master-datadir"
		expect expect/scp.exp ${username} ${ip} ${password} "mysql-remove.sh" /usr/local/cmha/ >/dev/null 2>&1 
		scp_error "${ip} Install mysql:　mysql-remove.sh"
		expect expect/scp.exp ${username} ${ip} ${password} ${mysql} /usr/local/cmha/ >/dev/null 2>&1 
		scp_error "${ip} Install mysql:　${mysqlversion}"
		expect expect/scp.exp ${username} ${ip} ${password} "consul.sh" /usr/local/cmha/ >/dev/null 2>&1
		scp_error "${ip} Install mysql: consul.sh"
		expect expect/expect.exp ${username} ${ip} ${password} "mkdir -p" "/usr/local/cmha/scripts/" >/dev/null 2>&1
        	if [ $? -ne 0 ]; then
                	echo "${ip} mkdir /usr/local/cmha/scripts/ failure"
                	exit 146
        	fi
		expect expect/scp.exp ${username} ${ip} ${password} "mysqlcheck.sh" /usr/local/cmha/scripts/  >/dev/null 2>&1
		scp_error "${ip} Install mysql: mysqlcheck.sh"
		expect expect/mysql.exp ${username} ${ip} ${password} ${mysqlversion} ${mysqldir} >/dev/null 2>&1
#>/dev/null 2>&1 
		expect expect/scp1.exp ${username} ${ip} ${password} /tmp/mysql.log /tmp/ >/dev/null 2>&1
		scp_error "${ip} Install mysql: /tmp/mysql.log"
		local coun=`grep -w "0" /tmp/mysql.log |wc -l`
		if [ ${coun} -ne 9 ]; then
			echo "${ip} Install mysql fail!"
			exit 115
		fi
		while read LINE
		do
			if $(echo ${LINE} | grep -w "port" | grep -w "=" >/dev/null); then
				sed -i "s/${LINE}/port		= ${mysql_port}/" conf/my.cnf
				sed_error "$?" "Install Mysql" "conf/my.cnf" "port"
			elif $(echo ${LINE} | grep -w "basedir" | grep -w "=" >/dev/null); then
				sed -i "s:${LINE}:basedir = ${basedir}:" conf/my.cnf
				sed_error "$?" "Install Mysql" "conf/my.cnf" "basedir"
			elif $(echo ${LINE} | grep -w "datadir" | grep -w "=" >/dev/null); then
                                sed -i "s:${LINE}:datadir = ${datadir}:" conf/my.cnf
                                sed_error "$?" "Install Mysql" "conf/my.cnf" "datadir"
			elif $(echo ${LINE} | grep -w "server-id" >/dev/null); then
				sed -i "s/${LINE}/server-id       = ${count}/" conf/my.cnf
				sed_error "$?" "Install Mysql" "conf/my.cnf" "server-id"
				let "count+=1"
			elif $(echo ${LINE} | grep -w "auto_increment_offset" >/dev/null); then
                                sed -i "s/${LINE}/auto_increment_offset       = ${offset}/" conf/my.cnf
                                sed_error "$?" "Install Mysql" "conf/my.cnf" "auto_increment_offset"
                                let "offset+=1"	
#			elif $(echo ${LINE} | grep -w "loose-rpl_reverse_recover_enabled" >/dev/null); then
#				sed -i "s/${LINE}/loose-rpl_reverse_recover_enabled = ${mysql_rpl_reverse_recover_enabled}/" conf/my.cnf
#				sed_error "$?" "Install Mysql" "conf/my.cnf" "loose-rpl_reverse_recover_enabled"
			elif $(echo ${LINE} | grep -w "ha_partner_host" >/dev/null); then
				sed -i "s/${LINE}/ha_partner_host = ${mysql_ha_partner_host}/" conf/my.cnf
				sed_error "$?" "Install Mysql" "conf/my.cnf" "ha_partner_host"
			elif $(echo ${LINE} | grep -w "ha_partner_port" >/dev/null); then
				sed -i "s/${LINE}/ha_partner_port = ${mysql_ha_partner_port}/" conf/my.cnf
				sed_error "$?" "Install Mysql" "conf/my.cnf" "ha_partner_port"
			elif $(echo ${LINE} | grep -w "ha_partner_user" >/dev/null); then
				sed -i "s/${LINE}/ha_partner_user = ${mysql_ha_partner_user}/" conf/my.cnf
				sed_error "$?" "Install Mysql" "conf/my.cnf" "ha_partner_user"
			elif $(echo ${LINE}| grep -w "ha_partner_password" >/dev/null); then
				sed -i "s/${LINE}/ha_partner_password = ${mysql_ha_partner_password}/" conf/my.cnf
				sed_error "$?" "Install Mysql" "conf/my.cnf" "ha_partner_password"
			fi 		
		done < conf/my.cnf
		expect expect/scp.exp ${username} ${ip} ${password} conf/my.cnf /etc/ >/dev/null 2>&1
		scp_error "${ip} Install mysql: conf/my.cnf"
		echo "${ip} mysql configure success"
		expect expect/mysql-install-plugin.exp ${username} ${ip} ${password} ${mysql_password} ${mysql_username} ${mysql_ha_partner_host} ${mysql_ha_partner_password} ${mysql_ha_partner_user} ${hostname} ${cmha_check_username} ${cmha_check_password} ${basedir} ${datadir}>/dev/null 2>&1
		echo ${username} ${ip} ${password} ${mysql_password} ${mysql_username} ${mysql_ha_partner_host} ${mysql_ha_partner_password} ${mysql_ha_partner_user} ${hostname} ${cmha_check_username} ${cmha_check_password} ${basedir} ${datadir}
		expect expect/consul.exp ${username} ${ip} ${password} "/usr/local/cmha/mysql/bin/mysql -u${mysql_username} -p${mysql_password} -h${hostname} -e \"show status\"" >/dev/null 2>&1
		if [ $? -ne 0 ]; then
			echo "${ip} install mysql fail!"
			exit 121
		else
			echo "${ip} install mysql success"
		fi
		expect expect/consul.exp ${username} ${ip} ${password} "service mysql restart" >/dev/null 2>&1
		if [ $? -ne 0 ]; then
			echo "${ip} mysql restart fail!"
			exit 124
		else
			echo "${ip} mysql restart success"
		fi
	done
done	
}


scp_expect(){
expect -c "
set timeout 100;
spawn scp -r $4 $1@$2:$5
expect {
\"*yes/no*\" {send \"yes\r\"; exp_continue}
\"*password*\" {send \"$3\r\";}
}
expect eof;"
}

configure_master_slave_relationship(){
for i in $(seq 1 ${db_number})
do
	for j in master slave
	do
		if [ $j == "master" ]; then
			local slave_host=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&\$1~/mysql-slave-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $1}'`
			error "$slave_host" "Master Slave Relationship" "auto-deployment.ini" "mysql-slave-ip-hostname"
			local slave_user=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-slave-ha_partner_user/{print $2;exit}' auto-deployment.ini`
			error "$slave_user" "Master Slave Relationship" "auto-deployment.ini" "mysql-slave-ha_partner_user"
			local slave_pass=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-slave-ha_partner_password/{print $2;exit}' auto-deployment.ini`
			error "$slave_pass" "Master Slave Relationship" "auto-deployment.ini" "mysql-slave-ha_partner_password"
			local slave_port=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-slave-ha_partner_port/{print $2;exit}' auto-deployment.ini`
			error "$slave_port" "Master Slave Relationship" "auto-deployment.ini" "mysql-slave-ha_partner_port"
			local slave_username=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-slave-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $3}'`
			error "$slave_username" "Master Slave Relationship" "auto-deployment.ini" "mysql-slave-ip-hostname"
			local slave_password=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-slave-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $4}'`
			error "slave_password" "Master Slave Relationship" "auto-deployment.ini" "mysql-slave-ip-hostname"
			local master_username=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-master-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $3}'`
			error "$master_username" "Master Slave Relationship" "auto-deployment.ini" "mysql-master-ip-hostname"
			local master_host=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&\$1~/mysql-master-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $1}'`
			error "$master_host" "Master Slave Relationship" "auto-deployment.ini" "mysql-master-ip-hostname"
			local master_password=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-master-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $4}'`
			error "$master_password" "Master Slave Relationship" "auto-deployment.ini" "mysql-master-ip-hostname"
			local mysql_user=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-slave-username/{print $2;exit}' auto-deployment.ini`
			error "$mysql_user" "Master Slave Relationship" "auto-deployment.ini" "mysql-slave-username"
			local mysql_pass=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-slave-password/{print $2;exit}' auto-deployment.ini`
			error "$mysql_pass" "Master Slave Relationship" "auto-deployment.ini" "mysql-slave-password"
			local user=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-master-username/{print $2;exit}' auto-deployment.ini`
			error "$user" "Master Slave Relationship" "auto-deployment.ini" "mysql-master-username"
			local pass=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-master-password/{print $2;exit}' auto-deployment.ini`
			error "$user" "Master Slave Relationship" "auto-deployment.ini" "mysql-master-password"
			expect expect/mysql-relationship.exp ${slave_username} ${slave_host} ${slave_password} ${mysql_user} ${mysql_pass} >/dev/null 2>&1
			expect expect/scp1.exp ${slave_username} ${slave_host} ${slave_password} /usr/local/cmha/mysql-relationship.txt /tmp/ >/dev/null 2>&1
			scp_error "Relationship: /usr/local/cmha/mysql-relationship.txt"
			while read LINE
			do
				if $(echo ${LINE} | grep -w "mysql-bin" >/dev/null); then
					master_log_file=`echo ${LINE}`
				elif $(echo ${LINE} | grep "^[0-9]*$" >/dev/null); then
					master_log_pos=`echo ${LINE}`
				fi
			done < /tmp/mysql-relationship.txt
			expect expect/mysql-change.exp ${master_username} ${master_host} ${master_password} ${slave_host} ${slave_user} ${slave_pass} ${master_log_file} ${master_log_pos} ${user} ${pass} ${slave_port} >/dev/null 2>&1
		elif [ $j == "slave" ]; then	
			local master_host=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&\$1~/mysql-master-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $1}'`
			error "$master_host" "Maste Slave Relationship" "auto-deployment.ini" "mysql-master-ip-hostname"
                        local master_user=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-master-ha_partner_user/{print $2;exit}' auto-deployment.ini`
			error "$master_user" "Master Slave Relationship" "auto-deployment.ini" "mysql-master-ha_partner_user"
			local master_port=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-master-ha_partner_port/{print $2;exit}' auto-deployment.ini`
			error "$master_port" "Master Slave Relationship" "auto-deployment.ini" "mysql-master-ha_partner_port"
                        local master_pass=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-master-ha_partner_password/{print $2;exit}' auto-deployment.ini`
			error "$master_pass" "Master Slave Relationship" "auto-deployment.ini" "mysql-master-ha_partner_password"
                        local master_username=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-master-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $3}'`
			error "$master_username" "Master Slave Relationship" "auto-deployment.ini" "mysql-master-ip-hostname"
                        local master_password=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-master-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $4}'`
			error "$master_password" "Master Slave Relationship" "auto-deployment.ini" "mysql-master-ip-hostname"
                        local slave_username=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-slave-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $3}'`
			error "$slave_username" "Master Slave Relationship" "auto-deployment.ini" "mysql-slave-ip-hostname"
                        local slave_host=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&\$1~/mysql-slave-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $1}'`
			error "$slave_host" "Master Slave Relationship" "auto-deployment.ini" "mysql-slave-ip-hostname"
                        local slave_password=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-slave-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $4}'`
			error "$slave_password" "Master slave Relationship" "auto-deployment.ini" "mysql-slave-ip-hostname"
			local mysql_user=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-master-username/{print $2;exit}' auto-deployment.ini`
			error "$mysql_user" "Master slave Relationship" "auto-deployment.ini" "mysql-master-username"
			local mysql_pass=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-master-password/{print $2;exit}' auto-deployment.ini`
			error "$mysql_pass" "Master Slave Relationship" "auto-deployment.ini" "mysql-master-password"
			local user=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-slave-username/{print $2;exit}' auto-deployment.ini`
			error "$user" "Master Slave Relationship" "auto-deployment.ini" "mysql-slave-username"
			local pass=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-slave-password/{print $2;exit}' auto-deployment.ini`
			error "$pass" "Master Slave Relationship" "auto-deployment.ini" "mysql-slave-password"
                        expect expect/mysql-relationship.exp ${master_username} ${master_host} ${master_password} ${mysql_user} ${mysql_pass} >/dev/null 2>&1
			expect expect/scp1.exp ${master_username} ${master_host} ${master_password} /usr/local/cmha/mysql-relationship.txt /tmp/ >/dev/null 2>&1
			scp_error "${ip} Relationship: /usr/local/cmha/mysql-relationship.txt"
                        while read LINE
                        do
                                if $(echo ${LINE} | grep -w "mysql-bin" >/dev/null); then
                                        slave_log_file=`echo ${LINE}`
                                elif $(echo ${LINE} | grep "^[0-9]*$" >/dev/null); then
                                        slave_log_pos=`echo ${LINE}`
                                fi
                        done < /tmp/mysql-relationship.txt
                        expect expect/mysql-change.exp ${slave_username} ${slave_host} ${slave_password} ${master_host} ${master_user} ${master_pass} ${slave_log_file} ${slave_log_pos} ${user} ${pass} ${master_port} >/dev/null 2>&1
		fi
	done
done
}

configure_keep_try(){
for i in $(seq 1 ${db_number})
do
        for j in master slave
        do
		 local ip=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&\$1~/mysql-'$j'-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $1}'`
                error "$ip" "Configure Keep Try" "auto-deployment.ini" "mysql-'$j'-ip-hostname"
                local password=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-'$j'-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $4}'`
                error "$password" "Configure Keep Try" "auto-deployment.ini" "mysql-'$j'-ip-hostname"
                local username=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-'$j'-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $3}'`	
		error "$username" "Configure Keep Try" "auto-deployment.ini" "mysql-'$j'-ip-hostname"
		local mysql_password=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-'$j'-password/{print $2;exit}' auto-deployment.ini`
                error "$mysql_password" "Install Mysql" "auto-deployment.ini" "mysql-'$j'-password"
                local mysql_username=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-'$j'-username/{print $2;exit}' auto-deployment.ini`
                error "$mysql_username" "Install Mysql" "auto-deployment.ini" "mysql-'$j'-username"
		expect expect/mysql-keep-try.exp ${username} ${ip} ${password}	${mysql_password} ${mysql_username} >/dev/null 2>&1
	done
done
	
}


create_cmha_check() {
for i in $(seq 1 ${db_number})
do
	local master_ip=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&\$1~/mysql-master-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $1}'`
	error "$master_ip" "Cmha Check" "auto-deployment.ini" "mysql-master-ip-hostname"
	local master_user=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&\$1~/mysql-master-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $3}'`
	error "$master_user" "Cmha Check" "auto-deployment.ini" "mysql-master-ip-hostname"
	local master_pass=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&\$1~/mysql-master-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $4}'`
	error "$master_pass" "Cmha Check" "auto-deployment.ini" "mysql-master-ip-hostname"
	local mysql_user=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-master-username/{print $2;exit}' auto-deployment.ini`
	error "$mysql_user" "Cmha Check" "auto-deployment.ini" "mysql-master-username"
	local mysql_pass=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-master-password/{print $2;exit}' auto-deployment.ini`
	error "$mysql_pass" "Cmha Check" "auto-deployment.ini" "mysql-master-password"
	expect expect/cmha-check.exp ${master_user} ${master_ip} ${master_pass} ${mysql_user} ${mysql_pass} >/dev/null 2>&1
done
}

install_haproxy_keepalived_ct(){
for i in $(seq 1 ${db_number})
do 
	for a in chap
	do
		for j in master slave
		do
			local ip=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&\$1~/'$a'-'$j'-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $1}'`
			error "$ip" "Install Haproxy Keepalived CT" "auto-deployment.ini" "$a-$j-ip-hostname"
			local hostname=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-'$j'-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $2}'`
			error "$hostname" "Install Haproxy Keepalived CT" "auto-deployment.ini" "$a-$j-ip-hostname"
			local password=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-'$j'-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $4}'`
			error "$password" "Install Haproxy Keepalived CT" "auto-deployment.ini" "$a-$j-ip-hostname"
			local virtual_ipaddress=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/chap-virtual-ipaddress/{print $2;exit}' auto-deployment.ini`
			error "$virtual_ipaddress" "Install Haproxy Keepalived CT" "auto-deployment.ini" "chap-virtual-ipaddress"
			local username=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-'$j'-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $3}'`
			error "$username" "Install Haproxy Keepalived CT" "auto-deployment.ini" "mysql-$j-ip-hostname"
			local servicename=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/servicename/{print $2;exit}' auto-deployment.ini`
			error "$servicename" "Install Haproxy Keepalived CT" "auto-deployment.ini" "servicename"
			local server_ip=`awk -F '=' '/\[cmha_server\]/{a=1}a==1&&$1~/server1-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $1}'`
			error "$server_ip" "Install Haproxy Keepalived CT" "auto-deployment.ini" "server1-ip-hostname"
			local port=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/mysql-master-port/{print $2;exit}' auto-deployment.ini`
			error "$port" "Install Haproxy Keepalived CT" "auto-deployment.ini" "mysql-master-port"
			local networkcard=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/chap-'$j'-network-card/{print $2;exit}' auto-deployment.ini`
			error "$networkcard" "Install Haproxy Keepalived CT" "auto-deployment.ini" "chap-$j-network-card"
			expect expect/syslog.exp ${username} ${ip} ${password} "/usr/local/cmha/log/consul-template.log" ${rsyslog_facility} >/dev/null 2>&1
			expect expect/expect.exp ${username} ${ip} ${password} mkdir /usr/local/cmha/consul-template.d/ >/dev/null 2>&1
			if [ $? -ne 0 ]; then
                                echo "${ip} mkdir /usr/local/cmha/consul-template.d/ failure"
                                exit 146
                        fi
			expect expect/expect.exp ${username} ${ip} ${password} mkdir /usr/local/cmha/scripts/ >/dev/null 2>&1
                	if [ $? -ne 0 ]; then
                        	echo "${ip} mkdir /usr/local/cmha/scripts/ failure"
                        	exit 146
                	fi
			expect expect/scp.exp ${username} ${ip} ${password} "package/haproxy-1.5.2-2.el6.x86_64.rpm" /usr/local/cmha/ >/dev/null 2>&1
			scp_error "${ip} Install haproxy: haproxy-1.5.2-2.el6.x86_64.rpm"
			expect expect/scp.exp ${username} ${ip} ${password} "chap-remove.sh" /usr/local/cmha/ >/dev/null 2>&1
			scp_error "${ip} Install haproxy: chap-remove.sh"
			expect expect/scp.exp ${username} ${ip} ${password} "consul.sh" /usr/local/cmha/ >/dev/null 2>&1
			scp_error "${ip} Install haproxy: consul.sh"
			expect expect/scp.exp ${username} ${ip} ${password} "consul-template.sh" /usr/local/cmha/ >/dev/null 2>&1
                        scp_error "${ip} Install haproxy: consul-template.sh"
			expect expect/scp.exp ${username} ${ip} ${password} "package/keepalived-1.2.13-5.el6_6.x86_64.rpm" /usr/local/cmha/ >/dev/null 2>&1
			scp_error "${ip} Install haproxy: keepalived-1.2.13-5.el6_6.x86_64.rpm"
			expect expect/scp.exp ${username} ${ip} ${password} "bin/consul-template" /usr/local/bin/ >/dev/null 2>&1
			scp_error "${ip} Install haproxy: consul-template"
			expect  expect/scp.exp ${username} ${ip} ${password} "keepalived.sh" /usr/local/cmha/scripts/ >/dev/null 2>&1
			scp_error "${ip} Install haproxy: keepalived.sh"
			while read line
			do
				if $(echo ${line} | grep -w "leader" >/dev/null); then
					sed -i "s/{{key \"service\/mysql\/leader\"}}/{{key \"service\/${servicename}\/leader\"}}/" conf/haproxy.ctmpl
					sed_error "$?" "Install Haproxy Keepalived CT" "conf/haproxy.ctmpl" "key"
				elif $(echo ${line} | grep -w "bind" >/dev/null); then
					sed -i "s/bind \*:3306/bind \*:${port}/" conf/haproxy.ctmpl
					sed_error "$?" "Install Haproxy Keepalived CT" "conf/haproxy.ctmpl" "bind"
				fi
			done < conf/haproxy.ctmpl
			expect expect/scp.exp ${username} ${ip} ${password} "conf/haproxy.ctmpl" /usr/local/cmha/consul-template.d/ >/dev/null 2>&1
			scp_error "Install haproxy:　conf/haproxy.ctmpl"
			echo "${ip} consul template file configure success"
			sed -i "s/{{key \"service\/${servicename}\/leader\"}}/{{key \"service\/mysql\/leader\"}}/" conf/haproxy.ctmpl
			sed_error "$?" "Install Haproxy Keepalived CT" "conf/haproxy.ctmpl" "key"
			sed -i "s/bind \*:${port}/bind \*:3306/" conf/haproxy.ctmpl
			sed_error "$?" "Install Hapeoxy Keepalived CT" "conf/haproxy.ctmpl" "bind"
			while read LINE	
			do
#				if $(echo ${LINE} | grep -w "router_id" >/dev/null); then
#					sed -i "s/${LINE}/router_id ${hostname}/" conf/keepalived.conf
#					sed_error "$?" "Install Haproxy Keepalived CT" "conf/keepalived.conf" "router_id"
				if $(echo ${LINE} | grep -w "script" >/dev/null); then
					sed -i "s/    script \"\/usr\/local\/cmha\/scripts\/keepalived.sh 192.168.2.1\"/    script \"\/usr\/local\/cmha\/scripts\/keepalived.sh ${server1_ip} ${server2_ip} ${server3_ip} ${servicename}\"/" conf/keepalived.conf
					sed_error "$?" "Install Haproxy Keepalived CT" "conf/keepalived.conf" "script"
				 elif [ $j == "master" ] && $(echo ${LINE} | grep -w "state" >/dev/null); then
                                        sed -i "s/${LINE}/state MASTER/" conf/keepalived.conf
                                        sed_error "$?" "Install Haproxy Keepalived CT" "conf/keepalived.conf" "state"
                                elif [ $j == "slave" ] && $(echo ${LINE} | grep -w "state" >/dev/null); then
                                        sed -i "s/${LINE}/state BACKUP/" conf/keepalived.conf
                                        sed_error "$?" "Install Haproxy Keepalived CT" "conf/keepalived.conf" "state"
				elif $(echo ${LINE} |grep -w "interface" >/dev/null); then
					sed -i "s/    interface eth0/    interface ${networkcard}/" conf/keepalived.conf
					sed_error "$?" "Install Haproxy Keepalived CT" "conf/keepalived.conf" "interface"
				elif $(echo ${LINE} | grep -w "virtual_router_id" >/dev/null); then
					sed -i "s/${LINE}/virtual_router_id $i/" conf/keepalived.conf
					sed_error "$?" "Install Haproxy Keepalived CT" "conf/keepalived.conf" "virtual_router_id"
				elif $(echo ${LINE} | grep -q '^[0-9".]\+$' >/dev/null); then
					sed -i "s/${LINE}/${virtual_ipaddress}/" conf/keepalived.conf
					sed_error "$?" "Install Haproxy Keepalived CT" "conf/keepalived.conf" "virtual_ipaddress"
				elif [ $j == "master" ] && $(echo ${LINE} | grep -w "priority" >/dev/null); then
					sed -i "s/${LINE}/priority 101/" conf/keepalived.conf
					sed_error "$?" "Install Haproxy Keepalived CT" "conf/keepalived.conf" "priority"
				elif [ $j == "slave" ] && $(echo ${LINE} | grep -w "priority" >/dev/null); then
					sed -i "s/${LINE}/priority 100/" conf/keepalived.conf
					sed_error "$?" "Install Haproxy Keepalived CT" "conf/keepalived.conf" "priority"
				fi 
			done < conf/keepalived.conf
			expect expect/haproxy.exp ${username} ${ip} ${password} >/dev/null 2>&1
			expect expect/scp.exp ${username} ${ip} ${password} "conf/keepalived.conf" /etc/keepalived/ >/dev/null 2>&1
			scp_error "${ip} Install haproxy: conf/keepalived.conf"
			sed -i "s/    script \"\/usr\/local\/cmha\/scripts\/keepalived.sh ${server1_ip} ${server2_ip} ${server3_ip} ${servicename}\"/    script \"\/usr\/local\/cmha\/scripts\/keepalived.sh 192.168.2.1\"/" conf/keepalived.conf
			sed_error "$?" "${ip} Install Haproxy Keepalived CT" "conf/keepalived.conf" "script"
			expect expect/consul-template.exp ${username} ${ip} ${password} ${enable_rsyslog} ${rsyslog_facility} >/dev/null 2>&1
			expect expect/expect.exp ${username} ${ip} ${password} "rpm -qa | grep haproxy" >/dev/null 2>&1
			if [ $? -ne 0 ]; then
				echo "${ip} haproxy install fail!"
				exit 123
			else
				echo "${ip} haproxy install success"
			fi
			expect expect/expect.exp ${username} ${ip} ${password} "rpm -qa | grep keepalived" >/dev/null 2>&1
			if [ $? -ne 0 ]; then
				echo "${ip} keepalived install fail!"
				exit 122
			else
				echo "${ip} keepalived install success"
			fi
			expect expect/pross.exp ${username} ${ip} ${password} "pgrep consul-template" >/dev/null 2>&1
        		if [ $? -ne 0 ];then
                		echo "${ip} consul template start fail!"
                		exit 126
			else
				echo "${ip} consul template start success"
        		fi
			expect expect/consul.exp ${username} ${ip} ${password} "service keepalived restart" >/dev/null 2>&1
			if [ $? -ne 0 ]; then
				echo "${ip} keepalived restart fail!"
				exit 124
			else
				echo "${ip} keepalived restart success"
			fi
		done
	done
done
}

configure_mha_handlers(){
for i in $(seq 1 ${db_number})
do
	for a in mysql
	do 
		for j in master slave
		do 
			local hostname=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-'$j'-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $2}'`
			error "$hostname" "Configure Mha Handlers" "auto-deployment.ini" "$a-$j-ip-hostname"
			if [ "$j" = "master" ]; then
                                other=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-slave-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $2}'`
#                               echo "aaaaa"$other_hostname
                                error "$other" "Configure Mha Handlers" "auto-deployment.ini" "$a-slave-ip-hostname"
                        else
                                other=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-master-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $2}'`
#                               echo "bbbb"$other_hostname
                                error "$other" "Configure Mha Handlers" "auto-deployment.ini" "$a-master-ip-hostname"
                        fi
			local ip=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&\$1~/'$a'-'$j'-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $1}'`
			error "$ip" "Configure Mha Handlers" "auto-deployment.ini" "$a-$j-ip-hostname"
			local hostname_password=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&\$1~/'$a'-'$j'-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $4}'`
			error "$hostname_password" "Configure Mha Handlers" "auto-deployment.ini" "$a-$j-ip-hostname"
			local port=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-'$j'-port/{print $2;exit}' auto-deployment.ini`
			error "$port" "Configure Mha Handlers" "auto-deployment.ini" "$a-$j-port"
			local username=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-'$j'-username/{print $2;exit}' auto-deployment.ini`
			error "$username" "Configure Mha Handlers" "auto-deployment.ini" "$a-$j-username"
			local password=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-'$j'-password/{print $2;exit}' auto-deployment.ini`
			error "$password" "Configure Mha Handlers" "auto-deployment.ini" "$a-$j-password"
			local check_username=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-master-cmha-check-username/{print $2;exit}' auto-deployment.ini`
	                error "$check_username" "Configure Bootstrap" "auto-deployment.ini" "$a-master-cmha-check-username"
        	        local check_password=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-master-cmha-check-password/{print $2;exit}' auto-deployment.ini`
                	error "$check_password" "Configure Bootstrap" "auto-deployment.ini" "$a-master-cmha-check-password"
			local servicename=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/servicename/{print $2;exit}' auto-deployment.ini`
			error "$servicename" "Configure Mha Handlers" "auto-deployment.ini" "servicename"
			local switch=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$j'-mha-switch/{print $2;exit}' auto-deployment.ini`
			error "$switch" "Configure Mha Handlers" "auto-deployment.ini" "$j-mha-switch"
			while read LINE
			do
				if $(echo ${LINE} | grep -w "hostname" >/dev/null); then
					sed -i "s/${LINE}/hostname = ${hostname}/" mha-handlers/conf/app.conf
					sed_error "$?" "Configure Mha Handlers" "mha-handlers/conf/app.conf" "hostname"
				elif $(echo ${LINE} | grep -w "otherhostname" >/dev/null); then
                                        sed -i "s/${LINE}/otherhostname = ${other}/" mha-handlers/conf/app.conf
                                        sed_error "$?" "Configure Mha Handlers" "mha-handlers/conf/app.conf" "otherhostname"
				elif $(echo ${LINE} | grep -w "ip" >/dev/null); then
					sed -i "s/${LINE}/ip = ${ip}/" mha-handlers/conf/app.conf
					sed_error "$?" "Configure Mha Handlers" "mha-handlers/conf/app.conf" "ip"
				elif $(echo ${LINE} | grep -w "port" >/dev/null); then
					sed -i "s/${LINE}/port = ${port}/" mha-handlers/conf/app.conf
					sed_error "$?" "Configure Mha Handlers" "mha-handlers/conf/app.conf" "port"
				elif $(echo ${LINE} | grep -w "username" >/dev/null); then
					sed -i "s/${LINE}/username = ${check_username}/" mha-handlers/conf/app.conf
					sed_error "$?" "Configure Mha Handlers" "mha-handlers/conf/app.conf" "check_username"
				elif $(echo ${LINE} | grep -w "password" >/dev/null); then
					sed -i "s/${LINE}/password = ${check_password}/" mha-handlers/conf/app.conf
					sed_error "$?" "Configure Mha Handlers" "mha-handlers/conf/app.conf" "check_password"
				elif $(echo ${LINE} | grep -w "datacenter" >/dev/null); then
					sed -i "s/${LINE}/datacenter = ${datacenter}/" mha-handlers/conf/app.conf
					sed_error "$?" "Configure Mha Handlers" "mha-handlers/conf/app.conf" "datacenter"
				elif $(echo ${LINE}| grep -w "service_ip" >/dev/null); then
					sed -i "s/${LINE}/service_ip = ${s_ip}/" mha-handlers/conf/app.conf
					sed_error "$?" "Configure Mha Handlers" "mha-handlers/conf/app.conf" "service_ip"
				elif $(echo ${LINE} | grep -w "servicename" >/dev/null); then
					sed -i "s/${LINE}/servicename = ${servicename}/" mha-handlers/conf/app.conf
					sed_error "$?" "Configure Mha Handlers" "mha-handlers/conf/app.conf" "servicename"
			#	elif $(echo ${LINE} | grep -w "switch" >/dev/null); then
			#		sed -i "s/${LINE}/switch = ${switch}/" mha-handlers/conf/app.conf
			#		sed_error "$?" "Configure Mha Handlers" "mha-handlers/conf/app.conf" "switch"
				fi	
			done < mha-handlers/conf/app.conf
			scp_expect ${username} ${ip} ${hostname_password} mha-handlers /usr/local/cmha/ >/dev/null 2>&1
		done
	done
done
}

configure_monitor_handlers(){
for i in $(seq 1 ${db_number})
do
	for a in mysql
	do
		for j in master slave
		do
			local ip=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-'$j'-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $1}'`
			error "$ip" "Configure Monitor Handlers" "auto-deployment.ini" "$a-$j-ip-hostname"
			local port=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-'$j'-port/{print $2;exit}' auto-deployment.ini`
			error "$port" "Configure Monitor Handlers" "auto-deployment.ini" "$a-$j-port"
			local username=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-'$j'-username/{print $2;exit}' auto-deployment.ini`
			error "$username" "Configure Monitor Handlers" "auto-deployment.ini" "$a-$j-username"
			if [ "$j" = "master" ]; then
				other_hostname=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-slave-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $2}'`
                        	error "$other_hostname" "Configure Monitor Handlers" "auto-deployment.ini" "$a-slave-ip-hostname"
				localhostname=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-'$j'-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $2}'`
				error "$localhostname" "Configure Monitor Handlers" "auto-deployment.ini" "$a-$j-ip-hostname" 
			else
				other_hostname=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-master-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $2}'`
                        	error "$other_hostname" "Configure Monitor Handlers" "auto-deployment.ini" "$a-master-ip-hostname"
				localhostname=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-'$j'-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $2}'`
				error "$localhostname" "Configure Monitor Handlers" "auto-deployment.ini" "$a-$j-ip-hostname"
			fi
			local host_password=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&\$1~/'$a'-'$j'-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $4}'`
			error "$host_password" "Configure Monitor Handlers" "auto-deployment.ini" "$a-$j-ip-hostname"
			local password=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-'$j'-password/{print $2;exit}' auto-deployment.ini`
			error "$password" "Configure Monitor Handlers" "auto-deployment.ini" "$a-$j-password"
			local check_username=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-master-cmha-check-username/{print $2;exit}' auto-deployment.ini`
	                error "$check_username" "Configure Bootstrap" "auto-deployment.ini" "$a-master-cmha-check-username"
        	        local check_password=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-master-cmha-check-password/{print $2;exit}' auto-deployment.ini`
                	error "$check_password" "Configure Bootstrap" "auto-deployment.ini" "$a-master-cmha-check-password"
			local servicename=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/servicename/{print $2;exit}' auto-deployment.ini`
			error "$servicename" "Configure Monitor Handlers" "auto-deployment.ini" "servicename"
			local switch=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$j'-monitor-switch/{print $2;exit}' auto-deployment.ini`
			error "$switch" "Configure Monitor Handlers" "auto-deployment.ini" "$j-monitor-switch"
			while read LINE
			do
				if $(echo ${LINE} | grep -w "ip" >/dev/null); then
					sed -i "s/${LINE}/ip = ${ip}/" monitor-handlers/conf/app.conf
					sed_error "$?" "Configure Monitor Handlers" "monitor-handlers/conf/app.conf" "ip"
				elif $(echo ${LINE} | grep -w "port" >/dev/null); then
					sed -i "s/${LINE}/port = ${port}/" monitor-handlers/conf/app.conf
					sed_error "$?" "Configure Monitor Handlers" "monitor-handlers/conf/app.conf" "port"
				elif $(echo ${LINE} | grep -w "username" >/dev/null); then
					sed -i "s/${LINE}/username = ${check_username}/" monitor-handlers/conf/app.conf
					sed_error "$?" "Configure Monitor Handlers" "monitor-handlers/conf/app.conf" "check_username"
				elif $(echo ${LINE} | grep -w "hostname" >/dev/null); then
                                        sed -i "s/${LINE}/hostname = ${localhostname}/" monitor-handlers/conf/app.conf
                                        sed_error "$?" "Configure Monitor Handlers" "monitor-handlers/conf/app.conf" "hostname"
				elif $(echo ${LINE} | grep -w "otherhostname" >/dev/null); then
                                        sed -i "s/${LINE}/otherhostname = ${other_hostname}/" monitor-handlers/conf/app.conf
                                        sed_error "$?" "Configure Monitor Handlers" "monitor-handlers/conf/app.conf" "otherhostname"
				elif $(echo ${LINE} | grep -w "password" >/dev/null); then
					sed -i "s/${LINE}/password = ${check_password}/" monitor-handlers/conf/app.conf
					sed_error "$?" "Configure Monitor Handlers" "monitor-handlers/conf/app.conf" "check_password"
				elif $(echo ${LINE} | grep -w "datacenter" >/dev/null); then
					sed -i "s/${LINE}/datacenter = ${datacenter}/" monitor-handlers/conf/app.conf
					sed_error "$?" "Configure Monitor Handlers" "monitor-handlers/conf/app.conf" "datacenter"
				elif $(echo ${LINE} | grep -w "service_ip" >/dev/null); then
					sed -i "s/${LINE}/service_ip = ${s_ip}/" monitor-handlers/conf/app.conf
					sed_error "$?" "Configure Monitor Handlers" "monitor-handlers/conf/app.conf" "service_ip"
				elif $(echo ${LINE} | grep -w "servicename" >/dev/null); then
					sed -i "s/${LINE}/servicename = ${servicename}/" monitor-handlers/conf/app.conf
					sed_error "$?" "Configure Monitor Handlers" "monitor-handlers/conf/app.conf" "servicename"
				elif $(echo ${LINE} | grep -w "switch_async" >/dev/null); then
					sed -i "s/${LINE}/switch_async = ${switch}/" monitor-handlers/conf/app.conf
					sed_error "$?" "Configure Monitor Handlers" "monitor-handlers/conf/app.conf" "switch"
				elif $(echo ${LINE} | grep -w "tag" >/dev/null) && [ $j == "master" ]; then
					sed -i "s/${LINE}/tag = slave/" monitor-handlers/conf/app.conf
					sed_error "$?" "Configure Monitor Handlers" "monitor-handlers/conf/app.conf" "tag"
				elif $(echo ${LINE} | grep -w "tag" >/dev/null) && [ $j == "slave" ]; then
					sed -i "s/${LINE}/tag = master/" monitor-handlers/conf/app.conf
					sed_error "$?" "Configure Monitor Handlers" "monitor-handlers/conf/app.conf" "tag"
				fi
			done < monitor-handlers/conf/app.conf
			scp_expect ${username} ${ip} ${host_password} monitor-handlers /usr/local/cmha/ >/dev/null 2>&1
		done 
	done
done
}

insert_kv_register(){
for i in $(seq 1 ${db_number})
do
	for a in mysql
	do
		for j in master slave
		do
			local servicename=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/servicename/{print $2;exit}' auto-deployment.ini`
			error "$servicename" "Insert Kv Register" "auto-deployment.ini" "servicename"
			local ip=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&\$1~/'$a'-'$j'-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $1}'`
			error "$ip" "Insert Kv Register" "auto-deployment.ini" "$a-$j-ip-hostname"
			local port=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-'$j'-port/{print $2;exit}' auto-deployment.ini`
			error "$port" "Insert Kv Register" "auto-deployment.ini" "$a-$j-port"
			local user=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-master-cmha-check-username/{print $2;exit}' auto-deployment.ini`
			#local user=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-'$j'-username/{print $2;exit}' auto-deployment.ini`
#			echo $user
			error "$user" "Insert Kv Register" "auto-deployment.ini" "'$a'-master-cmha-check-username"
			local password=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-master-cmha-check-password/{print $2;exit}' auto-deployment.ini`
			#local password=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-'$j'-password/{print $2;exit}' auto-deployment.ini`
#			echo $password
			error "$password" "Insert Kv Register" "auto-deployment.ini" "'$a'-master-cmha-check-password"
			if [ "$j" = "master" ]; then
				curl -X PUT http://${server1_ip}:8500/v1/kv/service/${servicename}/leader >/dev/null 2>&1
				if [ $? -ne 0 ]; then
					echo "${ip} Insert /service/${servicename}/leader fail!"
					exit 115
				else
					echo "${ip} Insert /service/${servicename}/leader success"
				fi
			fi
			curl -X PUT http://${ip}:8500/v1/agent/service/register -d "{\"Name\":\"${servicename}\",\"Tags\":[\"$j\"],\"Port\":${port},\"Check\":{\"Script\":\"/usr/local/cmha/scripts/mysqlcheck.sh ${ip} ${user} ${password}\",\"Interval\":\"3s\"}}" >/dev/null 2>&1
			if [ $? -ne 0 ]; then
				echo "Register ${servicename} service fail!"
				exit 116
			else
				echo "Register ${servicename} service success"
			fi
		done
	done
done
}

chap_register(){
for i in $(seq 1 ${db_number})
do
        for a in chap-master chap-slave
        do
                        local servicename=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/servicename/{print $2;exit}' auto-deployment.ini`
			error "$servicename" "Chap Register" "auto-deployment.ini" "servicename"
                        local ip=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&\$1~/'$a'-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $1}'`
			error "$ip" "Chap Register" "auto-deployment.ini" "$a-ip-hostname"
                        curl -X PUT http://${ip}:8500/v1/agent/service/register -d "{\"Name\":\"${servicename}\",\"Tags\":[\"$a\"],\"Check\":{\"Script\":\"/usr/local/cmha/scripts/keepalived.sh ${server1_ip} ${server2_ip} ${server3_ip} ${servicename}\",\"Interval\":\"10s\"}}" >/dev/null 2>&1
			if [ $? -ne 0 ]; then
				echo "${ip} Chap register ${servicename} service fail!"
				exit 118
			else
				echo "${ip} Chap register ${servicename} service success"
			fi
        done
done
}


configure_watch(){
for i in $(seq 1 ${db_number})
do
	for a in mysql
	do
		for j in master slave
		do
			local servicename=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/servicename/{print $2;exit}' auto-deployment.ini`
			error "$servicename" "Configure Watch" "auto-deployment.ini" "servicename"
			local username=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-'$j'-username/{print $2;exit}' auto-deployment.ini`
			error "$username" "Configure Watch" "auto-deployment.ini" "$a-$j-username"
                        local host_password=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&\$1~/'$a'-'$j'-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $4}'`
			error "$host_password" "Configure Watch" "auto-deployment.ini" "$a-$j-ip-hostname"
			local ip=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&\$1~/'$a'-'$j'-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $1}'`
			error "$ip" "Configure Watch" "auto-deployment.ini" "$a-$j-ip-hostname"
			while read LINE
			do
				if $(echo ${LINE} | grep -w "leader" >/dev/null); then
					sed -i "s/     \"key\": \"service\/mysql\/leader\",/     \"key\": \"service\/${servicename}\/leader\",/" conf/watch.json
					sed_error "$?" "Configure Watch" "conf/watch.json" "key"
				fi
			done < conf/watch.json
			while read line
			do
				if $(echo ${line} | grep -w "mysql" >/dev/null); then
					sed -i "s/     \"service\": \"mysql\",/     \"service\": \"${servicename}\",/" conf/watch-service.json
					sed_error "$?" "Configure Watch" "conf/watch-service.json" "service"
				elif $(echo ${line} | grep -w "tag" >/dev/null) && [ $j == "master" ]; then
					sed -i "s/     \"tag\": \"slave\",/     \"tag\": \"slave\",/" conf/watch-service.json
					sed_error "$?" "Configure Watch" "conf/watch-service.json" "tag"
				elif $(echo ${line} | grep -w "tag" >/dev/null) && [ $j == "slave" ]; then
					sed -i "s/     \"tag\": \"slave\",/     \"tag\": \"master\",/" conf/watch-service.json
					sed_error "$?" "Configure Watch" "conf/watch-service.json" "tag"
				fi
			done < conf/watch-service.json
			scp_expect ${username} ${ip} ${host_password} conf/watch.json /usr/local/cmha/consul.d/ >/dev/null 2>&1
			sed -i "s/     \"key\": \"service\/${servicename}\/leader\",/     \"key\": \"service\/mysql\/leader\",/" conf/watch.json
			sed_error "$?" "Configure Watch" "conf/watch.json" "key"
			scp_expect ${username} ${ip} ${host_password} conf/watch-service.json /usr/local/cmha/consul.d/ >/dev/null 2>&1
			sed -i "s/     \"service\": \"${servicename}\",/     \"service\": \"mysql\",/" conf/watch-service.json
			sed_error "$?" "Configure Watch" "conf/watch-service.json" "service"
			sed -i "s/     \"tag\": \"master\",/     \"tag\": \"slave\",/" conf/watch-service.json
			expect expect/consul.exp ${username} ${ip} ${host_password} "consul reload -rpc-addr=${ip}:8400" >/dev/null 2>&1
			if [ $? -ne 0 ]; then
				echo "${ip} consul reload fail!"
				exit 124
			else
				echo "${ip} consul reload  success"
			fi
		done
	done
done
}

configure_bootstrap(){
for i in $(seq 1 ${db_number})
do
	for a in mysql
	do
		local hostname=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-master-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $2}'`
		error "$hostname" "Configure Bootstrap" "auto-deployment.ini" "$a-master-ip-hostname"
		local other=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-slave-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $2}'`
#		echo "ccccc"$other
		error "$other" "Configure Bootstrap" "auto-deployment.ini" "$a-slave-ip-hostname"
		local ip=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&\$1~/'$a'-master-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $1}'`
		error "$ip"　"Configure Bootstrap" "auto-deployment.ini" "$a-master-ip-hostname"
		local port=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-master-port/{print $2;exit}' auto-deployment.ini`
		error "$port" "Configure Bootstrap" "auto-deployment.ini" "$a-master-port"
		local username=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-master-username/{print $2;exit}' auto-deployment.ini`
		error "$username" "Configure Bootstrap" "auto-deployment.ini" "$a-master-username"
		local password=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-master-password/{print $2;exit}' auto-deployment.ini`
		error "$password" "Configure Bootstrap" "auto-deployment.ini" "$a-master-password"
		local check_username=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-master-cmha-check-username/{print $2;exit}' auto-deployment.ini`
		error "$check_username" "Configure Bootstrap" "auto-deployment.ini" "$a-master-cmha-check-username"
		local check_password=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/'$a'-master-cmha-check-password/{print $2;exit}' auto-deployment.ini`
		error "$check_password" "Configure Bootstrap" "auto-deployment.ini" "$a-master-cmha-check-password"
		local servicename=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&$1~/servicename/{print $2;exit}' auto-deployment.ini`
		error "$servicename" "Configure Bootstrap" "auto-deployment.ini" "servicename"
		local host_password=`awk -F '=' '/\[cmha_agent'$i'\]/{a=1}a==1&&\$1~/'$a'-master-ip-hostname/{print $2;exit}' auto-deployment.ini |awk -F ',' '{print $4}'`
		error "$host_password" "Configure Bootstrap" "auto-deployment.ini" "$a-master-ip-hostname"
		while read LINE
		do
			if $(echo ${LINE} | grep -w "hostname" >/dev/null); then
				sed -i "s/${LINE}/hostname = ${hostname}/" bootstrap/conf/app.conf
				sed_error "$?" "Configure Bootstrap" "bootstrap/conf/app.conf" "hostname"
			elif $(echo ${LINE} | grep -w "otherhostname" >/dev/null); then
			#	echo ${LINE} "otherhostname = "${other}
				echo ${other}
				sed -i "s/${LINE}/otherhostname = ${other}/" bootstrap/conf/app.conf
                                #sed -i "s/${LINE}/otherhostname = aasss/" bootstrap/conf/app.conf
                                sed_error "$?" "Configure Bootstrap" "bootstrap/conf/app.conf" "otherhostname"
			elif $(echo ${LINE} | grep -w "ip" >/dev/null); then
				sed -i "s/${LINE}/ip = ${ip}/" bootstrap/conf/app.conf
				sed_error "$?" "Configure Bootstrap" "bootstrap/conf/app.conf" "ip"
			elif $(echo ${LINE} | grep -w "port" >/dev/null); then
				sed -i "s/${LINE}/port = ${port}/" bootstrap/conf/app.conf
				sed_error "$?" "Configure Bootstrap" "bootstrap/conf/app.conf" "port"
			elif $(echo ${LINE} | grep -w "username" >/dev/null); then
				sed -i "s/${LINE}/username = ${check_username}/" bootstrap/conf/app.conf
				sed_error "$?" "Configure Bootstrap" "bootstrap/conf/app.conf" "check_username"
			elif $(echo ${LINE} | grep -w "password" >/dev/null); then
				sed -i "s/${LINE}/password = ${check_password}/" bootstrap/conf/app.conf
				sed_error "$?" "Configure Bootstrap" "bootstrap/conf/app.conf" "check_password"
			elif $(echo ${LINE} | grep -w "datacenter" >/dev/null); then
				sed -i "s/${LINE}/datacenter = ${datacenter}/" bootstrap/conf/app.conf
				sed_error "$?" "Configure Bootstrap" "bootstrap/conf/app.conf" "datacenter"
			elif $(echo ${LINE} | grep -w "service_ip" >/dev/null); then
				sed -i "s/${LINE}/service_ip = ${s_ip}/g" bootstrap/conf/app.conf
				sed_error "$?" "Configure Bootstrap" "bootstrap/conf/app.conf" "service_ip"
			elif $(echo ${LINE} | grep -w "servicename" >/dev/null); then
				sed -i "s/${LINE}/servicename = ${servicename}/" bootstrap/conf/app.conf
				sed_error "$?" "Configure Bootstrap" "bootstrap/conf/app.conf" "servicename"
			fi
		done < bootstrap/conf/app.conf
		scp_expect ${username} ${ip} ${host_password} bootstrap /usr/local/cmha/ >/dev/null 2>&1
		expect expect/bootstrap.exp ${username} ${ip} ${host_password} >/dev/null 2>&1
	done
done
}
prog=auto-deployment
case $1 in 
	cs)
		echo "Configure consul server and start consul server...."
		read_consul_server_config
		sleep 5
		result=$(curl -X GET "http://${server1_ip}:8500/v1/status/leader")
		echo ${result} >/tmp/leader.txt
		leader=`grep ":" /tmp/leader.txt |sed 's/\"//g'`
		if $(echo ${leader} | grep -q '^[0-9.:]\+$' >/dev/null); then
			echo "leader exits,consul server success."
		else
			echo "no cluster leader,consul server fail."
		fi
	;;
	ca)
		echo "Configure consul agent and start consul agent...."
		read_consul_agent_config
	;;
	db)
		echo "Install mysql...."
		install_mysql
		echo "Configure master slave relationship...."
		configure_master_slave_relationship
		configure_keep_try
	;;
	chap)
		echo "Install chap...."
		create_cmha_check
		insert_kv_register
		install_haproxy_keepalived_ct
		chap_register
		configure_bootstrap
		configure_mha_handlers
		configure_monitor_handlers
		configure_watch
		echo "install chap success."
	;;
	*)
	echo $"Usage: $prog {cs|ca|db|chap}"
	exit 2
esac
