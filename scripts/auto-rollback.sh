#!/bin/bash
###version 1.1.5-Beta.6
basedir=/usr/local/cmha
deployment_file="auto-deployment.ini"
servicename=`awk -F= '/servicename/ {print $2}' $deployment_file`
#################auto login execute command###################
AUTO_LOGIN_EXEC(){
expect << EOF >>/dev/null
        set timeout -1 
        spawn ssh -o StrictHostKeyChecking=no $2
        expect { 
                "password"  { send "$1\n" } 
                         expect { 
                       \"password: \" {exit 1}
                       \"No such file or directory\" {exit 2}
                       \"Permission denied\"
                       {exit 3}
                       eof
                        }
                } 
        expect "]#" 
        send "$3 \n" 
        expect "]#" 
        send "exit\n" 
        expect eof
        catch wait result;
		exit [lindex \$result 3]  
EOF
status=$?
	return $status
}


#################rollback consul###################
rollback_consul(){
		for node in `./my_print_defaults -e ${deployment_file} $1 $2|awk -F= '/ip-hostname/ {print $2}'`
		do 
			local ip_node=`echo $node|awk -F"," '{print $1}'`
			local user=`echo $node|awk -F"," '{print $3}'`
			local password=`echo $node|awk -F"," '{print $4}'`
			AUTO_LOGIN_EXEC $password $user@$ip_node "killall consul"
			AUTO_LOGIN_EXEC $password $user@$ip_node "rm -rf $basedir;rm -rf /usr/local/bin/consul"
			AUTO_LOGIN_EXEC $password $user@$ip_node "rm -rf /etc/init.d/consul"
			AUTO_LOGIN_EXEC $password $user@$ip_node "rm -rf /usr/local/bin/jq"
			rollback_rsyslog consul	
			echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: host $ip_node rollback consul successfully"	
		done 		
}

#################rollback rsyslog###################
rollback_rsyslog(){
	if [ "$1" = consul ];then
	AUTO_LOGIN_EXEC $password $user@$ip_node "sed -i '/consul.log/'d /etc/rsyslog.conf"
	elif [ "$1" = chap ];then

	AUTO_LOGIN_EXEC $password $user@$ip_node "sed -i '/consul-template.log/'d /etc/rsyslog.conf"
	fi
	AUTO_LOGIN_EXEC $password $user@$ip_node "service rsyslog restart"
}
#################rollback mysql###################
rollback_mysql(){
	local datadir=`awk -F= '/datadir/ {print $2}' ${deployment_file}`
	local mysqlversion=`awk -F= '/mysqlversion/ {print $2}' ${deployment_file}`
	local register_ip=`grep -w "ip" ../bootstrap/conf/app.conf|awk '{print $3}'`
	for node in `./my_print_defaults -e ${deployment_file} db_node|awk -F= '/ip-hostname/ {print $2}'`
	do
			if [ "$datadir" = "" -o "$datadir" = "/" ];then
				echo "$(date "+%Y-%m-%d %H:%M:%S") [ERROR]: datadir Can not be empty"
				exit 1	
			fi
			local ip_node=`echo $node|awk -F"," '{print $1}'`
			local user=`echo $node|awk -F"," '{print $3}'`
			local password=`echo $node|awk -F"," '{print $4}'`
			AUTO_LOGIN_EXEC $password $user@$ip_node "pkill mysqld"
			AUTO_LOGIN_EXEC $password $user@$ip_node "rm -rf /etc/my.cnf"
			AUTO_LOGIN_EXEC $password $user@$ip_node "rm -rf /etc/init.d/mysql"
			AUTO_LOGIN_EXEC $password $user@$ip_node "rm -rf $datadir"
			AUTO_LOGIN_EXEC $password $user@$ip_node "rm -rf $basedir/mysql"
			AUTO_LOGIN_EXEC $password $user@$ip_node "sed -i '/\\\/usr\\\/local\\\/cmha\\\/mysql\\\/bin/'d /etc/profile"
			AUTO_LOGIN_EXEC $password $user@$ip_node "rm -rf $basedir/scripts"
			if [ "$ip_node" = "$register_ip" ];then
				AUTO_LOGIN_EXEC $password $user@$ip_node "rm -rf $basedir/bootstrap" 
			fi
			local mysqldir=`echo $mysqlversion|awk -F".tar.gz" '{print $1}'`
			AUTO_LOGIN_EXEC $password $user@$ip_node "rm -rf $basedir/$mysqldir"
			AUTO_LOGIN_EXEC $password $user@$ip_node "rm -rf $basedir/mha-handlers"
			AUTO_LOGIN_EXEC $password $user@$ip_node "rm -rf $basedir/consul.d/watch.json"
			AUTO_LOGIN_EXEC $password $user@$ip_node "rm -rf $basedir/consul.d/watch-service.json"
			AUTO_LOGIN_EXEC $password $user@$ip_node "rm -rf $basedir/monitor-handlers"
			AUTO_LOGIN_EXEC $password $user@$ip_node "rm -rf $basedir/$mysqlversion"
			AUTO_LOGIN_EXEC $password $user@$ip_node "rm -rf $basedir/mysql_replication.sh"
			AUTO_LOGIN_EXEC $password $user@$ip_node "rm -rf $basedir/mysqlcheck"
			AUTO_LOGIN_EXEC $password $user@$ip_node "rm -rf $basedir/db_register.json"
			AUTO_LOGIN_EXEC $password $user@$ip_node "userdel -f mysql"
		done
}	
#################rollback chap###################
rollback_chap(){
		for node in `./my_print_defaults -e ${deployment_file} chap_node|awk -F= '/ip-hostname/ {print $2}'`
		do
			local ip_node=`echo $node|awk -F"," '{print $1}'`
			local user=`echo $node|awk -F"," '{print $3}'`
			local password=`echo $node|awk -F"," '{print $4}'`	
			AUTO_LOGIN_EXEC $password $user@$ip_node "killall consul-template"
			AUTO_LOGIN_EXEC $password $user@$ip_node "killall keepalived" 
			AUTO_LOGIN_EXEC $password $user@$ip_node "killall haproxy"
			AUTO_LOGIN_EXEC $password $user@$ip_node "rm -rf $basedir/scripts"
			AUTO_LOGIN_EXEC $password $user@$ip_node "rpm -e keepalived"
			AUTO_LOGIN_EXEC $password $user@$ip_node "rpm -e haproxy"
			AUTO_LOGIN_EXEC $password $user@$ip_node "rm -rf $basedir/consul-template.d" 
			AUTO_LOGIN_EXEC $password $user@$ip_node "rm -rf $basedir/consul-template.sh"
			AUTO_LOGIN_EXEC $password $user@$ip_node "rm -rf /etc/init.d/consul-template"
			AUTO_LOGIN_EXEC $password $user@$ip_node "rm -rf /usr/local/bin/consul-template"
			AUTO_LOGIN_EXEC $password $user@$ip_node "rm -rf /etc/haproxy"
			AUTO_LOGIN_EXEC $password $user@$ip_node "rm -rf /etc/keepalived"
			AUTO_LOGIN_EXEC $password $user@$ip_node "rm -rf $basedir/chap_register.json"
			AUTO_LOGIN_EXEC $password $user@$ip_node "rm -rf $basedir/cmha-cli"
			rollback_rsyslog chap
		done
}


#################rollback kv register from cs###################
rollback_kv_register(){

		local server_ip=`awk -F= '/server1-ip-hostname/ {print $2}' ${deployment_file}|awk -F"," '{print $1}'`
                
		for i in `./my_print_defaults -e ${deployment_file} db_node|awk -F= '/ip-hostname/ {print $2}'`
		do
			local ip_node=`echo $i|awk -F"," '{print $1}'`
                        local user=`echo $i|awk -F"," '{print $3}'`
                        local password=`echo $i|awk -F"," '{print $4}'`
			#local node=`echo $i|awk -F= '{print $2}'|awk -F"," '{print $1}'`
			AUTO_LOGIN_EXEC $password $user@$ip_node "curl -X DELETE http://127.0.0.1:8500/v1/agent/service/deregister/${servicename}"
			AUTO_LOGIN_EXEC $password $user@$ip_node "curl -X DELETE http://127.0.0.1:8500/v1/agent/service/deregister/Statistics"
			echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: deregister db node $ip_node from CS successfully"

		done
		 curl -X DELETE http://${server_ip}:8500/v1/kv/cmha/service/${servicename}/db/leader >/dev/null 2>&1
         echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: rollback ${servicename} service leader from CS successfully"
			
}
#################rollback chap register from cs###################

rollback_chap_register(){

	for i in `./my_print_defaults -e ${deployment_file} chap_node|awk -F= '/ip-hostname/ {print $2}'`
	do
		local ip_node=`echo $i|awk -F"," '{print $1}'`
		local user=`echo $i|awk -F"," '{print $3}'`
		local password=`echo $i|awk -F"," '{print $4}'`
		#local node=`echo $i|awk -F"-" '{print $7}'`
		AUTO_LOGIN_EXEC $password $user@$ip_node "curl -X DELETE http://127.0.0.1:8500/v1/agent/service/deregister/${servicename}"
		AUTO_LOGIN_EXEC $password $user@$ip_node "curl -X DELETE http://127.0.0.1:8500/v1/agent/service/deregister/Statistics"
		echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: deregister chap node $ip_node from CS successfully"
	done
}
#################EXEC###################
read -t 60 -p "Do you want to remove $1 [Y/N]:"
      case $REPLY in  
            Y|y)   
		case $1 in
        	cs)
                	echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: start rollback CS ...."
                	rollback_consul cs_node
                	echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: rollback CS complete"
        	;;
        	ca)
                	echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: start rollback CA...."
                	rollback_chap_register
			rollback_kv_register
                	rollback_consul chap_node db_node
                	echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: rollback CA complete"
        	;;
        	db)
                	echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: start rollback MySQL..."
                	rollback_kv_register
                	rollback_mysql
                	echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: DB rollback complete"
        	;;
        	chap)
                	echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: start rollback chap...."
                	rollback_chap_register
                	rollback_chap
                	echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: chap rollback complete"      
        	;;
        	*)
        	echo $"Usage: $0 {chap|db|ca|cs}"
        	exit 2
		esac
                ;;  
       N|n)  
	exit 0
       ;;  
       *)  
          echo -e "input parameter error !!"  
          exit 1 
       ;; 
      esac   
