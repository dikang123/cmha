#!/bin/bash
####version 1.1.5-Beta.6
TOOL="deployment-check.sh"
LANG=C
L_ALL=C
deployment_file="auto-deployment.ini"
conf_dir=`pwd`/../conf
tmp_conf=$conf_dir/tmp_conf


#################主机连通性检查###################
check_host_connect(){
for i in `awk -F"," '/ip-hostname/ {print $1}' ${deployment_file} |awk -F= '{print $2}'`
do 
	ping -c 1 $i >>/dev/null
	if [ $? -eq 0 ];then
		echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: host $i ping check OK"
	else 
		echo "$(date "+%Y-%m-%d %H:%M:%S") [ERROR]: host $i connection unreasonable, host $i check fail"
		exit 1
	fi
done	
}
#################主机登陆连通性检查###################
login_check(){
	for h in `awk -F= '/ip-hostname/ {print $2}' ${deployment_file}`
	do
		local host=`echo $h|awk -F"," '{print $1}'`
		local user=`echo $h|awk -F"," '{print $3}'`
		local password=`echo $h|awk -F"," '{print $4}'`
		auto_login $password $user@$host >>/dev/null
		if [ $status -eq 0 ]; then
                echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: host $host login check OK!"
        elif [ $status -eq 1 ];then
                echo "$(date "+%Y-%m-%d %H:%M:%S") [ERROR]: host $host user name or password error, login check h fail!"
        		exit 1
        fi
	done
}
###########################自动登陆#############################
auto_login(){
expect << EOF 
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
        send "exit\n" 
        expect eof  
EOF
status=$?
    return $status
}
#################生成consul配置文件###################
consul_conf_generate(){
		local server="true"
		local log_level=`grep -w "log_level" ${deployment_file}|awk -F= '{print $2}'`
		local bootstrap_expect=`grep -w "bootstrap_expect" ${deployment_file}|awk -F= '{print $2}'`
		local enable_rsyslog=`grep -w "enable-rsyslog" ${deployment_file}|awk -F= '{print $2}'`
		local rsyslog_facility=`grep -w "rsyslog-facility" ${deployment_file}|awk -F= '{print $2}'`
		local cs_ip=`./my_print_defaults -e ${deployment_file} cs_node |awk -F= '{print $2}'|awk -F"," '{print $1";"}'|xargs|sed 's/[ ]//g'|sed 's/;$//'`
	for line in `awk -F= '/ip-hostname/ {print $2}' ${deployment_file}|awk -F"," '{print $1","$2}'`
	do
		local ip_add=`echo $line|awk -F"," '{print $1}'`
		local host_name=`echo $line|awk -F"," '{print $2}'`
		cp $conf_dir/config.json $tmp_conf/config.json.$ip_add
		sed -i "s/\$ip_address/$ip_add/g" $tmp_conf/config.json.$ip_add
		sed -i "s/\$mha/$datacenter/" $tmp_conf/config.json.$ip_add
		sed -i "s/\$log_level/$log_level/" $tmp_conf/config.json.$ip_add
		sed -i "s/\$enable_rsyslog/$enable_rsyslog/" $tmp_conf/config.json.$ip_add
		sed -i "s/\$rsyslog_facility/$rsyslog_facility/" $tmp_conf/config.json.$ip_add
		sed -i "s/\$bootstrap_expect/$bootstrap_expect/" $tmp_conf/config.json.$ip_add
		sed -i "s/\$hostname/$host_name/" $tmp_conf/config.json.$ip_add
	done

		#local host_name=`grep 'mysql-master-ip-hostname' ${deployment_file}|awk -F= '{print $2}'|awk -F"," '{print $2}'`
		#local otherhostname=`grep 'mysql-slave-ip-hostname' ${deployment_file}|awk -F= '{print $2}'|awk -F"," '{print $2}'`
		cp $conf_dir/watch-servicename-alerts.json $tmp_conf/watch-$servicename-alerts.json
		cp $conf_dir/watch-CS-alerts.json $tmp_conf/watch-CS-alerts.json
		sed -i "s/\$servicename/$servicename/g" $tmp_conf/watch-$servicename-alerts.json
		count=$bootstrap_expect
	for i in `./my_print_defaults -e ${deployment_file} cs_node|awk -F= '{print $2}'|awk -F"," '{print $1}'`
		do
			cp $conf_dir/dict.js $tmp_conf/dict.js.$i
			sed -i "s/\$CS_IP/$i/" $tmp_conf/dict.js.$i
			sed -i "s/false/$server/" $tmp_conf/config.json.$i
		        sed -i "s/\$http_addres/$i/g" $tmp_conf/config.json.$i
			cp $conf_dir/alerts-handler.conf $tmp_conf/alerts-handler.conf.$i
			sed -i "s/\$mha/$datacenter/" $tmp_conf/alerts-handler.conf.$i
			sed -i "s/\$ip_address/$i/g" $tmp_conf/alerts-handler.conf.$i
			sed -i "s/\$cs_ip/$cs_ip/g" $tmp_conf/alerts-handler.conf.$i
			cp $conf_dir/Statistics_register.json $tmp_conf/Statistics_register.json.$i
			sed -i "s/\$servicename/CS/g" $tmp_conf/Statistics_register.json.$i
			sed -i "s/\$ip_addr/$i/g" $tmp_conf/Statistics_register.json.$i
			#sed -i "s/\$hostname/$host_name/g" $tmp_conf/alerts-handler.conf.$i
			#sed -i "s/\$otherhostname/$otherhostname/g" $tmp_conf/alerts-handler.conf.$i
			for f in `ls $tmp_conf/config.json.*`
                        do
                        sed -i "s/cs_node$count/$i/" $f
                        done
			let count=$count-1
	done

	for i in `./my_print_defaults -e ${deployment_file} chap_node db_node|awk -F"," '/ip-hostname/ {print $1}'|awk -F= '{print $2}'`
	do
		sed -i "/bootstrap_expect/d" $tmp_conf/config.json.$i
		sed -i "/ui_dir/d" $tmp_conf/config.json.$i
		sed -i "s/\$http_addres/127.0.0.1/g" $tmp_conf/config.json.$i
		cp $conf_dir/Statistics_register.json $tmp_conf/Statistics_register.json.$i
		sed -i "s/\$servicename/$servicename/g" $tmp_conf/Statistics_register.json.$i
		sed -i "s/\$ip_addr/$i/g" $tmp_conf/Statistics_register.json.$i
		
	done
}
#################生成服务watch文件###################
watch_json_conf_generate(){

	for line in `./my_print_defaults -e ${deployment_file} db_node|awk -F"," '/ip-hostname/ {print $1}'|sed 's/=/\-/g'|awk -F"-" '{print $4"_"$7}'`
	do
		local ip_add=`echo $line|awk -F"_" '{print $2}'`
		local tag=`echo $line|awk -F"_" '{print $1}'`
		cp $conf_dir/watch.json $tmp_conf/watch.json.$ip_add
		cp $conf_dir/watch-service.json $tmp_conf/watch-service.json.$ip_add
		sed -i "s/\$servicename/$servicename/" $tmp_conf/watch.json.$ip_add
		sed -i "s/\$servicename/$servicename/" $tmp_conf/watch-service.json.$ip_add
		if [ "$tag" = "master" ];then
				tag="slave"
		elif [ "$tag" = "slave" ];then
				tag="master"
		fi
		sed -i "s/\$tag/$tag/" $tmp_conf/watch-service.json.$ip_add
	done
}
#################生成keepalived配置文件###################
keepalived_conf_generate(){

		local vip=`awk -F= '/virtual-ipaddress/ {print $2}' ${deployment_file} `
		local master_netcard=`awk -F= '/chap-master-network-card/ {print $2}' ${deployment_file}`
		local slave_netcard=`awk -F= '/chap-slave-network-card/ {print $2}' ${deployment_file}`
		local router_id=`awk -F"." '/chap-virtual-ipaddress/ {print $4}' ${deployment_file}`
		#local cs_ip=`./my_print_defaults -e ${deployment_file} cs_node|awk -F= '/ip-hostname/ {print $2}'|awk -F"," '{print $1}'|xargs`
		local pass=${servicename}@${router_id}
		local exec_script="$servicename"
	for line in `./my_print_defaults -e ${deployment_file} chap_node|awk -F"," '/ip-hostname/ {print $1}'|sed 's/=/\-/g'`
	do
		local ip_add=`echo $line|awk -F"-" '{print $7}'`
		local state=`echo $line|awk -F"-" '{print $4}'`
		cp $conf_dir/chap_register.json $tmp_conf/chap_register.json.$ip_add
		sed -i "s/\$servicename/$servicename/g" $tmp_conf/chap_register.json.$ip_add
		sed -i "s/\$ipaddr/$ip_add/" $tmp_conf/chap_register.json.$ip_add
		if [ "$state" == "master" ];then
		      local state="MASTER"
		      local priority=101
		      local netcard=$master_netcard
		      sed -i "s/\$tag/chap-master/" $tmp_conf/chap_register.json.$ip_add
		elif [ "$state" == "slave" ];then
		      local state="BACKUP"
		      local priority=100
		      local netcard=$slave_netcard
		      sed -i "s/\$tag/chap-slave/" $tmp_conf/chap_register.json.$ip_add
		fi
		cp $conf_dir/keepalived.conf $tmp_conf/keepalived.conf.$ip_add
		sed -i "s/\$vip/$vip/" $tmp_conf/keepalived.conf.$ip_add
		sed -i "s/\$state/$state/" $tmp_conf/keepalived.conf.$ip_add
		sed -i "s/\$priority/$priority/" $tmp_conf/keepalived.conf.$ip_add
		sed -i "s/\$netcard/$netcard/" $tmp_conf/keepalived.conf.$ip_add
		sed -i "s/\$router_id/$router_id/" $tmp_conf/keepalived.conf.$ip_add
		sed -i "s/\$scripts/$exec_script/" $tmp_conf/keepalived.conf.$ip_add
		sed -i "s/\$PASS/$pass/" $tmp_conf/keepalived.conf.$ip_add
		sed -i "s/\$servicename/$servicename/" $tmp_conf/keepalived.conf.$ip_add
	done	
		cp $conf_dir/haproxy.ctmpl $tmp_conf/haproxy.ctmpl
		sed -i "s/\$servicename/$servicename/" $tmp_conf/haproxy.ctmpl
		cp $conf_dir/haproxy_check.json $tmp_conf/haproxy_check.json
		sed -i "s/\$servicename/$servicename/" $tmp_conf/haproxy_check.json
}
#################生成my.cnf和monitor、mha handler配置文件###################
handler_and_my_conf_generate(){

		local master_ip=`grep 'mysql-master-ip-hostname' ${deployment_file}|awk -F= '{print $2}'|awk -F"," '{print $1}'`
		local slave_ip=`grep 'mysql-slave-ip-hostname' ${deployment_file}|awk -F= '{print $2}'|awk -F"," '{print $1}'`
		local port=`grep "mysql-master-port" ${deployment_file} |awk -F= '{print $2}'`
		local username=`grep 'mysql-master-cmha-check-username' ${deployment_file}|awk -F= '{print $2}'`
		local password=`grep 'mysql-master-cmha-check-password' ${deployment_file}|awk -F= '{print $2}'`
		#local service_ip=`./my_print_defaults -e ${deployment_file} cs_node |awk -F= '{print $2}'|awk -F"," '{print $1";"}'|xargs|sed 's/[ ]//g'|sed 's/;$//'`
		local format=`grep 'master-format' ${deployment_file}|awk -F= '{print $2}'`
		local mysql_user=`awk -F= '/master-usernam/ {print $2}' ${deployment_file}`
		local mysql_password=`awk -F= '/master-password/ {print $2}' ${deployment_file}`
	for t in monitor-handlers mha-handlers mysqlcheck
	do
		for i in $master_ip $slave_ip
		do
			cp $conf_dir/$t.conf ../conf/tmp_conf/$t.conf.$i
			sed -i "s/\$username/$username/" ../conf/tmp_conf/$t.conf.$i
			sed -i "s/\$password/$password/" ../conf/tmp_conf/$t.conf.$i
			sed -i "s/\$port/$port/" ../conf/tmp_conf/$t.conf.$i
			sed -i "s/\$ipaddress/$i/" ../conf/tmp_conf/$t.conf.$i
			sed -i "s/\$datacenter/$datacenter/" ../conf/tmp_conf/$t.conf.$i
			if [ $i = $master_ip ];then
				host_name=`grep 'mysql-master-ip-hostname' ${deployment_file}|awk -F= '{print $2}'|awk -F"," '{print $2}'`
				otherhostname=`grep 'mysql-slave-ip-hostname' ${deployment_file}|awk -F= '{print $2}'|awk -F"," '{print $2}'`
				sed -i "s/\$hostname/$host_name/" ../conf/tmp_conf/$t.conf.$i
				sed -i "s/\$otherhostname/$otherhostname/" ../conf/tmp_conf/$t.conf.$i
				if [ "$t" == "monitor-handlers" ];then
					sed -i "s/\$tag/slave/" ../conf/tmp_conf/$t.conf.$i
				else
					sed -i "s/\$format/$format/" ../conf/tmp_conf/$t.conf.$i
				fi
			elif [ $i = $slave_ip ];then
				host_name=`grep 'mysql-slave-ip-hostname' ${deployment_file}|awk -F= '{print $2}'|awk -F"," '{print $2}'`
				otherhostname=`grep 'mysql-master-ip-hostname' ${deployment_file}|awk -F= '{print $2}'|awk -F"," '{print $2}'`
				sed -i "s/\$hostname/$host_name/" ../conf/tmp_conf/$t.conf.$i
				sed -i "s/\$otherhostname/$otherhostname/" ../conf/tmp_conf/$t.conf.$i
				if [ "$t" == "monitor-handlers" ];then
					sed -i "s/\$tag/master/" ../conf/tmp_conf/$t.conf.$i
				else
					sed -i "s/\$format/$format/" ../conf/tmp_conf/$t.conf.$i
				fi				
			fi
			sed -i "s/\$servicename/$servicename/" ../conf/tmp_conf/$t.conf.$i
			#sed -i "s/\$service_ip/$service_ip/" ../$t/conf/app.conf.$i
		done
	done
	cp ../conf/tmp_conf/mha-handlers.conf.$master_ip ../bootstrap/conf/app.conf
	sed -i "s/mha-handlers/bootstrap/" ../bootstrap/conf/app.conf 
######################################my.cnf##############################
	local datadir=`grep "datadir" ${deployment_file} |awk -F= '{print $2}'`
	local ha_partner_password=`grep "slave-ha_partner_password" ${deployment_file} |awk -F= '{print $2}'`
	local ha_partner_port=`grep "slave-ha_partner_port" ${deployment_file} |awk -F= '{print $2}'`
	local ha_partner_user=`grep "slave-ha_partner_user" ${deployment_file} |awk -F= '{print $2}'`
	local check_user=`grep "cmha-check-username" ${deployment_file} |awk -F= '{print $2}'`
	local check_password=`grep "cmha-check-password" ${deployment_file} |awk -F= '{print $2}'`
	cp $conf_dir/my.cnf $tmp_conf/my.cnf.$master_ip
	cp $conf_dir/my.cnf $tmp_conf/my.cnf.$slave_ip
	sed -i "s/\$ha_partner_host/$slave_ip/" $tmp_conf/my.cnf.$master_ip
	sed -i "s/\$ha_partner_host/$master_ip/" $tmp_conf/my.cnf.$slave_ip
	sed -i "s/\$server-id/1/" $tmp_conf/my.cnf.$master_ip
	sed -i "s/\$server-id/2/" $tmp_conf/my.cnf.$slave_ip
	sed -i "s/\$auto_increment_offset/1/" $tmp_conf/my.cnf.$master_ip
	sed -i "s/\$auto_increment_offset/2/" $tmp_conf/my.cnf.$slave_ip
	for i in $master_ip $slave_ip
	do
		sed -i "s#\$datadir#$datadir#" $tmp_conf/my.cnf.$i
		sed -i "s/\$ha_partner_user/$ha_partner_user/" $tmp_conf/my.cnf.$i
		sed -i "s/\$ha_partner_password/$ha_partner_password/" $tmp_conf/my.cnf.$i
		sed -i "s/\$ha_partner_port/$ha_partner_port/" $tmp_conf/my.cnf.$i
		cp $conf_dir/mysql_replication.sh $conf_dir/tmp_conf/mysql_replication.sh.$i
		cp $conf_dir/mysql.sql $conf_dir/tmp_conf/mysql.sql.$i
		cp $conf_dir/db_register.json $conf_dir/tmp_conf/db_register.json.$i
		sed -i "s/\$check_user/$check_user/g" $conf_dir/tmp_conf/mysql_replication.sh.$i
		sed -i "s/\$check_password/$check_password/g" $conf_dir/tmp_conf/mysql_replication.sh.$i
		sed -i "s/\$repl_port/$ha_partner_port/g" $conf_dir/tmp_conf/mysql_replication.sh.$i
		sed -i "s/\$check_user/$check_user/g" $conf_dir/tmp_conf/db_register.json.$i
		sed -i "s/\$check_password/$check_password/g" $conf_dir/tmp_conf/db_register.json.$i
		sed -i "s/\$port/$ha_partner_port/g" $conf_dir/tmp_conf/db_register.json.$i
		sed -i "s/\$servicename/$servicename/g" $conf_dir/tmp_conf/db_register.json.$i
		sed -i "s/\$ipaddr/$i/g" $conf_dir/tmp_conf/db_register.json.$i
		sed -i "s/\$cmha_check_username/$username/" $conf_dir/tmp_conf/mysql.sql.$i
		sed -i "s/\$cmha_check_password/$password/" $conf_dir/tmp_conf/mysql.sql.$i
		sed -i "s/\$mysql_ha_partner_user/$ha_partner_user/" $conf_dir/tmp_conf/mysql.sql.$i 
		sed -i "s/\$mysql_ha_partner_password/$ha_partner_password/" $conf_dir/tmp_conf/mysql.sql.$i 
		sed -i "s/\$mysql_password/$mysql_password/" $conf_dir/tmp_conf/mysql.sql.$i
		sed -i "s/\$repl_user/$ha_partner_user/" $tmp_conf/mysql_replication.sh.$i
                sed -i "s/\$repl_password/$ha_partner_password/" $tmp_conf/mysql_replication.sh.$i
		if [ "$i" = "$master_ip" ];then
			sed -i "s/\$peer_host/$slave_ip/g" $conf_dir/tmp_conf/mysql_replication.sh.$i
			sed -i "s/\$mysql_ha_partner_host/$slave_ip/" $conf_dir/tmp_conf/mysql.sql.$i	
			sed -i "s/\$tag/master/" $conf_dir/tmp_conf/db_register.json.$i	
		elif [ "$i" = "$slave_ip" ];then
			sed -i "s/\$peer_host/$master_ip/g" $conf_dir/tmp_conf/mysql_replication.sh.$i
			sed -i "s/\$mysql_ha_partner_host/$master_ip/" $conf_dir/tmp_conf/mysql.sql.$i
			sed -i "s/\$tag/slave/" $conf_dir/tmp_conf/db_register.json.$i	
		fi
			sed -i "s/\$local_host/$i/g" $conf_dir/tmp_conf/mysql_replication.sh.$i
	done
}


#################检查模板文件和目录是否存在###################
check_template_file(){
	local template_file=(config.json keepalived.conf monitor-handlers.conf watch-service.json
haproxy.ctmpl mha-handlers.conf my.cnf watch.json)
	local cmha_file=(bootstrap/bootstrap monitor-handlers/monitor-handlers mha-handlers/mha-handlers alerts-handler/alerts-handler bin/consul bin/consul-template scripts/keepalived.sh scripts/keepalived_status.sh scripts/consul.sh scripts/consul-template.sh mysqlcheck/mysqlcheck package/keepalived-1.2.13-5.el6_6.x86_64.rpm package/mysql-5.6.27-v1-enterprise-commercial-linux-x86_64.tar.gz package/haproxy-1.5.2-2.el6.x86_64.rpm cmha-cli/cmha-cli)
	for f in ${template_file[*]}
	do
		if [ ! -f $conf_dir/$f ];then
			echo "$(date "+%Y-%m-%d %H:%M:%S") [ERROR]: template file $f not exits"
			exit 1
		else
			echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: template file $f check OK"
		fi
	done

	for f in ${cmha_file[*]}
	do
		if [ -e ../$f ];then
			echo "$(date "+%Y-%m-%d %H:%M:%S") [INFO]: cmha file $f check OK"
		else
			echo "$(date "+%Y-%m-%d %H:%M:%S") [ERROR]: cmha file $f not exits"
			exit 1
		fi
	done
}

#################检查模板文件和目录是否存在###################
MAIN(){
  	if [[ ! -x $(which expect) ]]; then
            echo "$(date "+%Y-%m-%d %H:%M:%S") [ERROR]: 'expect' command does not exist. please install expect"
            exit 2
    	fi

    if [ ! -f ${deployment_file} ]; then
		echo "$(date "+%Y-%m-%d %H:%M:%S") [ERROR]: auto deployment auto-deployment.ini not exist,please configuration"	
		exit 3
	fi

bootstrap_expect=`awk -F= '$1=="bootstrap_expect" {print $2}' ${deployment_file}` 
servicename=`grep servicename ${deployment_file} |awk -F= '{print $2}'`
datacenter=`grep -w "datacenter" ${deployment_file}|awk -F= '{print $2}'`
mysqlversion=`grep -w "mysqlversion" ${deployment_file} |awk -F '=' '{print $2}'`
	
	if [ -z ${bootstrap_expect} ];then
            echo "$(date "+%Y-%m-%d %H:%M:%S") [ERROR]: Passing parameters can not be empty of bootstrap_expect!!!"
            exit 2
    fi

    if [ -z  ${servicename} ];then
            echo "$(date "+%Y-%m-%d %H:%M:%S") [ERROR]: Passing parameters can not be empty of servicename!!!"
            exit 2
    fi
	
	if [ -z ${datacenter} ];then
            echo "$(date "+%Y-%m-%d %H:%M:%S") [ERROR]: Passing parameters can not be empty of datacenter!!!" 
            exit 2
    fi

    if [ "../package/$mysqlversion" != "$(ls ../package/$mysqlversion)" ];then
    		echo "$(date "+%Y-%m-%d %H:%M:%S") [ERROR]: There is no MySQL version of the installation package $mysqlversion!!!" 
    		exit 2
    fi


    if [ ! -d $tmp_conf ];then
		mkdir $tmp_conf
	fi
	check_host_connect
	login_check
	check_template_file
	consul_conf_generate
	watch_json_conf_generate
	keepalived_conf_generate
	handler_and_my_conf_generate
}

#----------------------------EXEC --------------------------------------------------------
if [ "${0##*/}" = "$TOOL" ] || [ "${0##*/}" = "bash" -a "${_:-""}" = "$0" ]; then
{
        MAIN "${@:-""}"
}
fi
		




