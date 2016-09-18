#!/bin/bash
###version 1.1.5-Beta.6
basedir=/usr/local/cmha
deployment_file="auto-deployment.ini"
servicename=`awk -F= '/servicename/ {print $2}' $deployment_file`
rsyslog_facility=`awk -F= '/rsyslog-facility/ {print $2}' $deployment_file`
#################自动scp文件或者目录到目标机器###################
AUTO_SCP() {
    expect -c "set timeout -1;
                spawn scp -r -o StrictHostKeyChecking=no ${@:2};
       expect {
           *assword:* {send -- $1\r;
                        expect { 
                       \"password: \" {exit 1}
                       \"No such file or directory\" {exit 2}
                       \"Permission denied\"
                       {exit 3}
                       eof
                        }
        } 
                    eof         {exit 1;}
                }
                "
status=$?
	if [ $status -ne 0 ];then
		echo "$(date "+%Y-%m-%d %H:%M:%S") [ERROR]: scp ${@:2} fail"
		exit 1
	fi
    return $status
}>>/dev/null
#################自动登陆节点执行命令###################
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
	if [ $status -ne 0 ];then
		echo "$(date "+%Y-%m-%d %H:%M:%S") [ERROR]: host $2 exec command $3 fail"
		exit 1
	fi
	return $status
}
#################安装consul 服务###################
install_consul(){
		for node in `./my_print_defaults -e ${deployment_file} $1 $2|awk -F= '/ip-hostname/ {print $2}'`
		do 
			local ip_node=`echo $node|awk -F"," '{print $1}'`
			local user=`echo $node|awk -F"," '{print $3}'`
			local password=`echo $node|awk -F"," '{print $4}'`
			AUTO_LOGIN_EXEC $password $user@$ip_node "mkdir -p $basedir/log"
			AUTO_LOGIN_EXEC $password $user@$ip_node "mkdir -p $basedir/consul.d"
			AUTO_LOGIN_EXEC $password $user@$ip_node "mkdir -p $basedir/consul_data"
			AUTO_LOGIN_EXEC $password $user@$ip_node "chmod 755 $basedir"
			AUTO_SCP $password ../bin/consul $user@$ip_node:/usr/local/bin
			AUTO_SCP $password ../package/jq-linux64 $user@$ip_node:/usr/local/bin/jq
			AUTO_SCP $password ../scripts/consul.sh $user@$ip_node:$basedir
			AUTO_SCP $password ../scripts/consul.sh $user@$ip_node:/etc/init.d/consul
			AUTO_SCP $password ../conf/tmp_conf/config.json.$ip_node $user@$ip_node:$basedir/consul.d/config.json
			AUTO_SCP $password ../conf/tmp_conf/Statistics_register.json.$ip_node $user@$ip_node:$basedir/Statistics_register.json
			config_rsyslog consul
			AUTO_LOGIN_EXEC $password $user@$ip_node "cd $basedir;./consul.sh start"	
			if [ $status -eq 0 ];then
				echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: host $ip_node install consul successfully"
				sleep 10
				if [ "$1" = "cs_node" ];then
			AUTO_LOGIN_EXEC $password $user@$ip_node "cd $basedir;curl -X PUT -d @Statistics_register.json http://$ip_node:8500/v1/agent/service/register"
				else
				
			AUTO_LOGIN_EXEC $password $user@$ip_node "cd $basedir;curl -X PUT -d @Statistics_register.json http://127.0.0.1:8500/v1/agent/service/register"
				fi
			fi		
		done 		
	return $status
}
install_alerts(){
	local alerts_service=$2
	for node in `./my_print_defaults -e ${deployment_file} $1|awk -F= '/ip-hostname/ {print $2}'`
	do
	 	local ip_node=`echo $node|awk -F"," '{print $1}'`
                local user=`echo $node|awk -F"," '{print $3}'`
                local password=`echo $node|awk -F"," '{print $4}'`
		if [ "$alerts_service" = "CS" ];then
                AUTO_SCP $password ../ui $user@$ip_node:$basedir
                AUTO_SCP $password ../conf/tmp_conf/dict.js.$ip_node $user@$ip_node:$basedir/ui/js/dict.js
                AUTO_SCP $password ../alerts-handler $user@$ip_node:$basedir
                AUTO_SCP $password ../conf/tmp_conf/alerts-handler.conf.$ip_node $user@$ip_node:$basedir/alerts-handler/conf/app.conf
		AUTO_SCP $password ../conf/tmp_conf/watch-CS-alerts.json $user@$ip_node:$basedir/consul.d
		AUTO_LOGIN_EXEC $password $user@$ip_node "curl -X PUT -d "enable" http://$ip_node:8500/v1/kv/cmha/service/CS/alerts/alert_boot"
		else 
		AUTO_SCP $password ../conf/tmp_conf/watch-${servicename}-alerts.json $user@$ip_node:$basedir/consul.d
		AUTO_LOGIN_EXEC $password $user@$ip_node "curl -X PUT -d "enable" http://$ip_node:8500/v1/kv/cmha/service/$servicename/alerts/alert_boot"
		fi
		local result=$(curl -s -X GET "http://${ip_node}:8500/v1/status/leader"|awk -F'"' '{print $2}')
		if [ "$result" != "${ip_node}:8300" ];then
		AUTO_LOGIN_EXEC $password $user@$ip_node "consul reload -rpc-addr=$ip_node:8400"
		fi
			if [ $status -eq 0 ];then
				echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: host $ip_node install alerts plugin successfully"
			fi		
	done
	return $status
}
#################配置rsyslog服务###################
config_rsyslog(){
	if [ "$1" = consul ];then
	AUTO_LOGIN_EXEC $password $user@$ip_node "echo $rsyslog_facility.*                                                $basedir/log/consul.log >>/etc/rsyslog.conf"
	elif [ "$1" = chap ];then

	AUTO_LOGIN_EXEC $password $user@$ip_node "echo $rsyslog_facility.* 												  $basedir/log/consul-template.log >>/etc/rsyslog.conf"
	fi
	AUTO_LOGIN_EXEC $password $user@$ip_node "service rsyslog restart"
}
#################安装mysql节点配置复制###################
install_mysql(){
	local datadir=`awk -F= '/datadir/ {print $2}' ${deployment_file}`
	local mysqlversion=`awk -F= '/mysqlversion/ {print $2}' ${deployment_file}`
	local mysql_user=`awk -F= '/mysql-master-username/ {print $2}' ${deployment_file}`
	local mysql_password=`awk -F= '/mysql-master-password/ {print $2}' ${deployment_file}`
	local register_ip=`grep -w "ip" ../bootstrap/conf/app.conf|awk '{print $3}'`
	for node in `./my_print_defaults -e ${deployment_file} db_node|awk -F= '/ip-hostname/ {print $2}'`
	do
			local ip_node=`echo $node|awk -F"," '{print $1}'`
			local user=`echo $node|awk -F"," '{print $3}'`
			local password=`echo $node|awk -F"," '{print $4}'`
			AUTO_LOGIN_EXEC $password $user@$ip_node "rpm -qa|grep libaio"
			AUTO_LOGIN_EXEC $password $user@$ip_node "groupadd mysql;useradd -r -g mysql mysql"
			AUTO_SCP $password ../package/$mysqlversion $user@$ip_node:$basedir
			AUTO_LOGIN_EXEC $password $user@$ip_node "tar zxvf $basedir/$mysqlversion -C $basedir/"
			local mysqldir=`echo $mysqlversion|awk -F".tar.gz" '{print $1}'`
			AUTO_LOGIN_EXEC $password $user@$ip_node "mkdir -p $datadir;chown -R mysql.mysql $datadir"
			AUTO_LOGIN_EXEC $password $user@$ip_node "cd $basedir;ln -s $mysqldir mysql;chown -R mysql.mysql mysql;cd $basedir/mysql/;cp support-files/mysql.server /etc/init.d/mysql"
			AUTO_SCP $password ../conf/tmp_conf/my.cnf.$ip_node $user@$ip_node:/etc/my.cnf
			AUTO_SCP $password ../conf/tmp_conf/mysql.sql.$ip_node $user@$ip_node:/tmp
			#AUTO_SCP $password ../conf/tmp_conf/mysql_replication.sh.$ip_node $user@$ip_node:/tmp
		echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: host $ip_node Start to initialize the MySQL database."		
			AUTO_LOGIN_EXEC $password $user@$ip_node "perl $basedir/mysql/scripts/mysql_install_db --user=mysql --basedir=/usr/local/cmha/mysql --datadir=$datadir"
		echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: host $ip_node initialize the MySQL database successfully."		
			AUTO_LOGIN_EXEC $password $user@$ip_node "service mysql start"
			AUTO_LOGIN_EXEC $password $user@$ip_node "echo 'export PATH=$PATH:$basedir/mysql/bin' >>/etc/profile"
			AUTO_LOGIN_EXEC $password $user@$ip_node "mkdir $basedir/scripts"
			AUTO_LOGIN_EXEC $password $user@$ip_node "mysqladmin -u$mysql_user password '$mysql_password' -S /tmp/mysql.sock"
			AUTO_LOGIN_EXEC $password $user@$ip_node "mysql -u$mysql_user -p$mysql_password </tmp/mysql.sql.$ip_node"
			AUTO_LOGIN_EXEC $password $user@$ip_node "rm -rf /tmp/mysql.sql.$ip_node"
			AUTO_SCP $password ../conf/tmp_conf/mysql_replication.sh.$ip_node $user@$ip_node:$basedir/mysql_replication.sh
			if [ "$ip_node" = "$register_ip" ];then
				AUTO_SCP $password ../bootstrap $user@$ip_node:$basedir
			fi
			
			if [ $status -eq 0 ];then
                        echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: host $ip_node MySQL install configure successfully"
                        fi
		done
}	
#################配置MySQL复制###################
configure_replication(){
	for node in `./my_print_defaults -e ${deployment_file} db_node|awk -F= '/ip-hostname/ {print $2}'`
	do
			local ip_node=`echo $node|awk -F"," '{print $1}'`
			local user=`echo $node|awk -F"," '{print $3}'`
			local password=`echo $node|awk -F"," '{print $4}'`
			AUTO_LOGIN_EXEC $password $user@$ip_node "sh $basedir/mysql_replication.sh"
			if [ $status -eq 0 ];then
			echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: host $ip_node mysql replication configure successfully"
			fi
	done
}

#################安装chap节点###################
install_chap(){
		for node in `./my_print_defaults -e ${deployment_file} chap_node|awk -F= '/ip-hostname/ {print $2}'`
		do
			local ip_node=`echo $node|awk -F"," '{print $1}'`
			local user=`echo $node|awk -F"," '{print $3}'`
			local password=`echo $node|awk -F"," '{print $4}'`	
			AUTO_SCP $password ../package/haproxy-1.5.2-2.el6.x86_64.rpm $user@$ip_node:/tmp
			AUTO_SCP $password ../package/keepalived-1.2.13-5.el6_6.x86_64.rpm $user@$ip_node:/tmp
			AUTO_LOGIN_EXEC $password $user@$ip_node "rpm -ivh /tmp/keepalived-1.2.13-5.el6_6.x86_64.rpm"
			AUTO_LOGIN_EXEC $password $user@$ip_node "rpm -ivh /tmp/haproxy-1.5.2-2.el6.x86_64.rpm"
			AUTO_SCP $password ../conf/tmp_conf/chap_register.json.$ip_node $user@$ip_node:$basedir/chap_register.json
			AUTO_LOGIN_EXEC $password $user@$ip_node "mkdir $basedir/consul-template.d"
			AUTO_LOGIN_EXEC $password $user@$ip_node "mkdir $basedir/scripts"
			AUTO_SCP $password ../scripts/consul-template.sh $user@$ip_node:$basedir
			AUTO_SCP $password ../scripts/consul-template.sh $user@$ip_node:/etc/init.d/consul-template
			AUTO_SCP $password ../conf/haproxy.ctmpl $user@$ip_node:$basedir/consul-template.d
			AUTO_SCP $password ../conf/consul-template.conf $user@$ip_node:$basedir/consul-template.d
			AUTO_SCP $password ../bin/consul-template $user@$ip_node:/usr/local/bin
			AUTO_SCP $password ../conf/tmp_conf/keepalived.conf.$ip_node $user@$ip_node:/etc/keepalived/keepalived.conf
			AUTO_SCP $password ../conf/tmp_conf/haproxy.ctmpl $user@$ip_node:$basedir/consul-template.d
			AUTO_SCP $password ../scripts/keepalived.sh $user@$ip_node:$basedir/scripts
			AUTO_SCP $password ../scripts/keepalived_status.sh $user@$ip_node:$basedir/scripts
			AUTO_SCP $password ../cmha-cli $user@$ip_node:$basedir
			config_rsyslog chap
			AUTO_LOGIN_EXEC $password $user@$ip_node "$basedir/consul-template.sh start"
			AUTO_LOGIN_EXEC $password $user@$ip_node "/etc/init.d/keepalived start;ps -ef|grep keepalived"
		done
}
#################注册agent到cs###################
insert_kv_register(){

		local server_ip=`awk -F= '/server1-ip-hostname/ {print $2}' ${deployment_file}|awk -F"," '{print $1}'`
                 curl -X PUT http://${server_ip}:8500/v1/kv/cmha/service/${servicename}/db/leader >/dev/null 2>&1
                                if [ $? -eq 0 ]; then
                                  echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: Insert /service/${servicename}/leader to CS successfully" 
                                else
                                  echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: Insert /service/${servicename}/leader to CS fail"
                                   exit 1
                                fi
		for i in `./my_print_defaults -e ${deployment_file} db_node|awk -F"," '/ip-hostname/ {print $0}'`
		do
			local node=`echo $i|awk -F= '{print $2}'|awk -F"," '{print $1}'`
			local user=`echo $i|awk -F= '{print $2}'|awk -F"," '{print $3}'`
			local password=`echo $i|awk -F= '{print $2}'|awk -F"," '{print $4}'`
			AUTO_SCP $password ../mysqlcheck $user@$node:$basedir
			AUTO_SCP $password ../conf/tmp_conf/mysqlcheck.conf.$node $user@$node:$basedir/mysqlcheck/conf/app.conf
			AUTO_SCP $password ../conf/tmp_conf/db_register.json.$node $user@$node:$basedir/db_register.json
		AUTO_LOGIN_EXEC $password $user@$node "cd $basedir/;curl -X PUT -d @db_register.json http://127.0.0.1:8500/v1/agent/service/register"
			if [ $status -eq 0 ]; then
				echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: host $node Register MySQL to CS successfully"
			else
				echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: host $node Register MySQL to CS fail"
				exit 1
			fi

		done
			
}
configure_bootstrap(){
                local boot_ip=`awk -F= '/mysql-master-ip-hostname/ {print $2}' ${deployment_file} |awk -F"," '{print $1}'`
                local boot_user=`awk -F= '/mysql-master-ip-hostname/ {print $2}' ${deployment_file} |awk -F"," '{print $3}'`
                local boot_password=`awk -F= '/mysql-master-ip-hostname/ {print $2}' ${deployment_file} |awk -F"," '{print $4}'`
                       AUTO_LOGIN_EXEC $boot_password $boot_user@$boot_ip "cd $basedir/bootstrap;./bootstrap"
                         echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: host $boot_ip bootstrap successfully"
}

#################重启DB节点ca###################
configure_watch(){

		for node in `./my_print_defaults -e ${deployment_file} db_node|awk -F= '/ip-hostname/ {print $2}'`
		do
			local ip_node=`echo $node|awk -F"," '{print $1}'`
			local user=`echo $node|awk -F"," '{print $3}'`
			local password=`echo $node|awk -F"," '{print $4}'`
			AUTO_LOGIN_EXEC $password $user@$ip_node "cd $basedir;mkdir scripts;mkdir mha-handlers;mkdir monitor-handlers"
			AUTO_SCP $password ../mha-handlers $user@$ip_node:$basedir
			AUTO_SCP $password ../monitor-handlers $user@$ip_node:$basedir
			AUTO_SCP $password ../conf/tmp_conf/watch.json.$ip_node $user@$ip_node:$basedir/consul.d/watch.json
			AUTO_SCP $password ../conf/tmp_conf/watch-service.json.$ip_node $user@$ip_node:$basedir/consul.d/watch-service.json
			AUTO_SCP $password ../conf/tmp_conf/mha-handlers.conf.$ip_node $user@$ip_node:$basedir/mha-handlers/conf/app.conf
			AUTO_SCP $password ../conf/tmp_conf/monitor-handlers.conf.$ip_node $user@$ip_node:$basedir/monitor-handlers/conf/app.conf
			AUTO_LOGIN_EXEC $password $user@$ip_node "consul reload -rpc-addr=127.0.0.1:8400"
		done
}
#################注册chap到cs中###################
chap_register(){
		local vip=`./my_print_defaults -e ${deployment_file} chap_node|awk -F= '/chap-virtual-ipaddress/ {print $2}'`
	for i in `./my_print_defaults -e ${deployment_file} chap_node|awk -F= '/ip-hostname/ {print $2}'`
	do
			local node=`echo $i|awk -F"," '{print $1}'`
                        local user=`echo $i|awk -F"," '{print $3}'`
                        local password=`echo $i|awk -F"," '{print $4}'`
			AUTO_LOGIN_EXEC $password $user@$node "cd $basedir/;curl -X PUT -d @chap_register.json http://127.0.0.1:8500/v1/agent/service/register"
			AUTO_LOGIN_EXEC $password $user@$node "curl -X PUT -d "$vip" http://127.0.0.1:8500/v1/kv/cmha/service/${servicename}/chap/VIP"
		if [ $status -eq 0 ]; then
			echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: Register chap node $node to CS successfully"
			AUTO_SCP $password ../conf/tmp_conf/haproxy_check.json $user@$node:$basedir/consul.d
			AUTO_SCP $password ../scripts/Statistics_haproxy.sh $user@$node:$basedir/scripts
			AUTO_LOGIN_EXEC $password $user@$node "consul reload -rpc-addr=127.0.0.1:8400"
		else
			echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: Register chap node $node to CS fail"
			exit 1
		fi 
	done
}
#################根据传入参数安装指定的服务###################
case $1 in 
	cs)
		echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: Configure CS and start CS ...."
		install_consul cs_node
		install_alerts cs_node CS
		sleep 10
		for i in `./my_print_defaults -e ${deployment_file} cs_node|awk -F= '{print $2}'|awk -F"," '{print $1}'`
		do
		result=$(curl -s -X GET "http://${i}:8500/v1/status/leader")
		if $(echo $result|grep -q "$i" >/dev/null); then
			echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: $result is leader ,Leader election success."
			echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: Please execute (curl -s -X GET http://${i}:8500/v1/status/leader) view leader"
			echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: CS configuration complete"
		fi
		done
	;;
	ca)
		echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: Configure CA node and start CA...."
		install_consul chap_node db_node	
		if [ $status -eq 0 ];then
			echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: CA configuration complete"
		else
			echo "$(date "+%Y-%m-%d %H:%M:%S") [ERROR]: Configure CA node and start CA fail"	
			exit 1
		fi
	;;
	db)
 		echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: Installation and configuration MySQL..."
		install_mysql
		echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: install MySQL successfully"
		echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: Start the configuration of MySQL replication..."
		configure_replication
		echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: DB configuration complete"
		insert_kv_register
		sleep 10
		configure_bootstrap
		sleep 20
		configure_watch
	;;
	chap)
		install_chap
		echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: install chap successfully"	
		echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: Install chap nodes and register services to CS...."
		chap_register
		sleep 5
		#configure_bootstrap
		#sleep 20
		#configure_watch
		install_alerts cs_node $servicename
		echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: chap installation is complete"
		
	;;
	*)
	echo $"Usage: $0 {cs|ca|db|chap}"
	exit 2
esac

