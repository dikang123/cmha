#!/bin/bash
####version 0.6.0
deployment_file="auto-deployment.ini"
if [ ! -f ${deployment_file} ]; then
	echo "Error: auto deployment Profiles not exist,please configuration"	
	exit 100
fi
bootstrap_expect=`grep -w "bootstrap_expect" auto-deployment.ini | awk -F '=' '{print $2}'` 
if [ "${bootstrap_expect}" = "" ]; then
	echo "Error: auto-deployment.ini: bootstrap_expect is empty or categories error!"
	exit 138
fi
if [ ! -d "conf/" ]; then
	mkdir conf/
fi
config_server(){
	if [ ${bootstrap_expect} -eq 3 ]; then
		echo "{" >conf/config.json
		echo "  \"bootstrap_expect\": ${bootstrap_expect}," >>conf/config.json
		echo "  \"server\": false," >>conf/config.json
		echo "  \"datacenter\": \"dc1\"," >>conf/config.json	
		echo "  \"data_dir\": \"/tmp/consul\"," >>conf/config.json
		echo "  \"node_name\": \"rhel64\"," >>conf/config.json
		echo "  \"disable_update_check\": true," >>conf/config.json
		echo "  \"log_level\": \"INFO\"," >>conf/config.json
		echo "  \"enable_syslog\": true," >>conf/config.json
		echo "  \"syslog_facility\": \"LOCAL5\"," >>conf/config.json
		echo "  \"addresses\": {" >>conf/config.json
		echo "    \"http\": \"192.168.2.96\"," >>conf/config.json
		echo "    \"rpc\": \"192.168.2.96\"" >>conf/config.json
		echo "  }," >>conf/config.json
		echo "  \"start_join\": [" >>conf/config.json
		echo "    \"192.168.2.96\"," >>conf/config.json
		echo "    \"192.168.2.97\"," >>conf/config.json
		echo "    \"192.168.2.98\"" >>conf/config.json
		echo " ]" >>conf/config.json
		echo "}" >>conf/config.json
	elif [ ${bootstrap_expect} -eq 5 ]; then
		echo "{" >conf/config.json
		echo "  \"bootstrap_expect\": ${bootstrap_expect}," >>conf/config.json
   		echo "  \"server\": false," >>conf/config.json
	        echo "  \"datacenter\": \"dc1\"," >>conf/config.json
        	echo "  \"data_dir\": \"/tmp/consul\"," >>conf/config.json
	        echo "  \"node_name\": \"rhel64\"," >>conf/config.json
        	echo "  \"disable_update_check\": true," >>conf/config.json
	        echo "  \"log_level\": \"INFO\"," >>conf/config.json
		echo "  \"enable_syslog\": true," >>conf/config.json
                echo "  \"syslog_facility\": \"LOCAL5\"," >>conf/config.json
        	echo "  \"addresses\": {" >>conf/config.json
	        echo "    \"http\": \"192.168.2.96\"," >>conf/config.json
        	echo "    \"rpc\": \"192.168.2.96\"" >>conf/config.json
	        echo "  }," >>conf/config.json
        	echo "  \"start_join\": [" >>conf/config.json
	        echo "    \"192.168.2.96\"," >>conf/config.json
        	echo "    \"192.168.2.97\"," >>conf/config.json
	        echo "    \"192.168.2.98\"," >>conf/config.json
		echo "    \"192.168.2.99\"," >>conf/config.json
		echo "    \"192.168.2.100\""  >>conf/config.json
        	echo " ]" >>conf/config.json
	        echo "}" >>conf/config.json
	else
		echo "Error: auto-deployment.ini: bootstrap_expect value write error,3 or 5!"
		exit 139
	fi
	echo "cmha template file configuration success!"
}

watch(){
	echo "{" >conf/watch.json
	echo "  \"watches\": [" >>conf/watch.json
	echo "   {" >>conf/watch.json
	echo "     \"type\": \"key\"," >>conf/watch.json
	echo "     \"key\": \"service/mysql/leader\"," >>conf/watch.json
	echo "     \"handler\": \"/usr/local/cmha/mha-handlers/mha-handlers\"" >>conf/watch.json
	echo "   }" >>conf/watch.json
	echo "  ]" >>conf/watch.json
	echo "}" >>conf/watch.json
	echo "cmha watch kv template file configuration success!"
}

watch_service(){
	echo "{" >conf/watch-service.json
	echo "  \"watches\": [" >>conf/watch-service.json
	echo "   {" >>conf/watch-service.json
	echo "     \"type\": \"service\"," >>conf/watch-service.json
	echo "     \"service\": \"mysql\"," >>conf/watch-service.json
	echo "     \"tag\": \"slave\"," >>conf/watch-service.json
	echo "     \"handler\": \"/usr/local/cmha/monitor-handlers/monitor-handlers\"" >>conf/watch-service.json
	echo "   }" >>conf/watch-service.json
	echo "  ]" >>conf/watch-service.json
	echo "}" >>conf/watch-service.json
	echo "cmha watch service template file configuration success!"
}

haproxy(){
	echo "global" >conf/haproxy.ctmpl
	echo "	pidfile /var/run/haproxy.pid" >>conf/haproxy.ctmpl
	echo "	daemon" >>conf/haproxy.ctmpl
	echo "	user nobody" >>conf/haproxy.ctmpl
	echo "	group nobody" >>conf/haproxy.ctmpl
	echo "        maxconn 4096" >>conf/haproxy.ctmpl
	echo "	spread-checks 3" >>conf/haproxy.ctmpl
	echo "	quiet" >>conf/haproxy.ctmpl
	echo "defaults" >>conf/haproxy.ctmpl
	echo "	mode tcp" >>conf/haproxy.ctmpl
	echo "	option dontlognull" >>conf/haproxy.ctmpl
	echo "	option tcp-smart-accept" >>conf/haproxy.ctmpl
	echo "	option tcp-smart-connect" >>conf/haproxy.ctmpl
	echo "	retries 3" >>conf/haproxy.ctmpl
	echo "	option redispatch" >>conf/haproxy.ctmpl
	echo "        maxconn 4096" >>conf/haproxy.ctmpl
	echo "	timeout check 3500ms" >>conf/haproxy.ctmpl
	echo "	timeout queue 3500ms" >>conf/haproxy.ctmpl
	echo "	timeout connect 3500ms" >>conf/haproxy.ctmpl
	echo "	timeout client 10000ms" >>conf/haproxy.ctmpl
	echo "	timeout server 10000ms" >>conf/haproxy.ctmpl
	echo "	userlist STATSUSERS" >>conf/haproxy.ctmpl
	echo "	group admin users admin" >>conf/haproxy.ctmpl
	echo "	user admin insecure-password 1111" >>conf/haproxy.ctmpl
	echo "	user stats insecure-password 1111" >>conf/haproxy.ctmpl
	echo "listen admin_page 0.0.0.0:9600" >>conf/haproxy.ctmpl
	echo "	mode http" >>conf/haproxy.ctmpl
	echo "	stats enable" >>conf/haproxy.ctmpl
	echo "	stats refresh 60s" >>conf/haproxy.ctmpl
	echo "	stats uri /" >>conf/haproxy.ctmpl
	echo "	acl AuthOkay_ReadOnly http_auth(STATSUSERS)" >>conf/haproxy.ctmpl
	echo "	acl AuthOkay_Admin http_auth_group(STATSUSERS) admin" >>conf/haproxy.ctmpl
	echo "	stats http-request auth realm admin_page unless AuthOkay_ReadOnly" >>conf/haproxy.ctmpl
	echo "	stats admin if AuthOkay_Admin" >>conf/haproxy.ctmpl
	echo "listen innosql" >>conf/haproxy.ctmpl
	echo "	bind *:3306" >>conf/haproxy.ctmpl
	echo "	timeout client 60000ms" >>conf/haproxy.ctmpl
	echo "	timeout server 60000ms" >>conf/haproxy.ctmpl
	echo "	balance roundrobin" >>conf/haproxy.ctmpl
	echo "	{{key \"service/mysql/leader\"}}" >>conf/haproxy.ctmpl
	echo "cmha haproxy template file configuration success!"
}

keepalived(){
	echo "vrrp_script chk_haproxy_ct {" >conf/keepalived.conf
	echo "    script \"/usr/local/cmha/scripts/keepalived.sh 192.168.2.1\"" >>conf/keepalived.conf
	echo "    interval 2" >>conf/keepalived.conf
	echo "    weight 2" >>conf/keepalived.conf
	echo "}" >>conf/keepalived.conf
	echo "vrrp_instance haproxy {" >>conf/keepalived.conf
	echo "    state MASTER" >>conf/keepalived.conf
	echo "    interface eth0" >>conf/keepalived.conf
	echo "    virtual_router_id 1" >>conf/keepalived.conf
	echo "    priority 100" >>conf/keepalived.conf
	echo "    advert_int 1" >>conf/keepalived.conf
	echo "    authentication {" >>conf/keepalived.conf
	echo "        auth_type PASS" >>conf/keepalived.conf
	echo "        auth_pass 1111" >>conf/keepalived.conf
	echo "    }" >>conf/keepalived.conf
	echo "    virtual_ipaddress {" >>conf/keepalived.conf
	echo "	192.168.2.110" >>conf/keepalived.conf
	echo "    }" >>conf/keepalived.conf
	echo "    track_script {" >>conf/keepalived.conf
	echo "	chk_haproxy_ct" >>conf/keepalived.conf
	echo "    }" >>conf/keepalived.conf
	echo "}" >>conf/keepalived.conf
	echo "cmha keepalived template file configuration success!"
}
mysql(){
	echo "[mysqld]" >conf/my.cnf
	echo "user = mysql" >>conf/my.cnf 
	echo "basedir = /usr/local/cmha/mysql" >>conf/my.cnf
	echo "datadir = /var/lib/mysql" >>conf/my.cnf 
	echo "port            = 3306" >>conf/my.cnf
	echo "server-id       = 1" >>conf/my.cnf
	echo "core-file" >>conf/my.cnf 
	echo "log_bin = mysql-bin.log" >>conf/my.cnf 
	echo "binlog_format = ROW" >>conf/my.cnf 
	echo "auto_increment_increment = 2" >>conf/my.cnf 
	echo "auto_increment_offset = 1" >>conf/my.cnf 
	echo "pid_file = mysql.pid" >>conf/my.cnf 
	echo "socket          = /tmp/mysql.sock" >>conf/my.cnf
	echo "log_output = FILE" >>conf/my.cnf 
	echo "character_set_server = utf8" >>conf/my.cnf 
	echo "slow_query_log_file = mysql-slow.log" >>conf/my.cnf 
	echo "log_error = mysql-err.log" >>conf/my.cnf 
	echo "query_cache_type = 0" >>conf/my.cnf 
	echo "long_query_time = 3" >>conf/my.cnf 
	echo "max_connections = 1024" >>conf/my.cnf 
	echo "max_connect_errors = 1024" >>conf/my.cnf 
	echo "local_infile = 0" >>conf/my.cnf 
	echo "general_log = OFF" >>conf/my.cnf 
	echo "slow_query_log = ON" >>conf/my.cnf 
	echo "relay-log = mysqld-relay-bin" >>conf/my.cnf 
	echo "expire_logs_days = 15" >>conf/my.cnf 
	echo "innodb_io_capacity = 500 # SSD 2000 ~ 20000" >>conf/my.cnf
	echo "innodb_flush_method = O_DIRECT" >>conf/my.cnf 
	echo "innodb_file_format = Barracuda" >>conf/my.cnf 
	echo "innodb_file_format_max = Barracuda" >>conf/my.cnf 
	echo "innodb_log_file_size = 1G" >>conf/my.cnf 
	echo "innodb_file_per_table = ON" >>conf/my.cnf 
	echo "innodb_lock_wait_timeout = 5" >>conf/my.cnf 
	echo "innodb_buffer_pool_size = 512M" >>conf/my.cnf 
	echo "innodb_print_all_deadlocks = ON" >>conf/my.cnf 
	echo "innodb_additional_mem_pool_size = 32M" >>conf/my.cnf 
	echo "innodb_data_file_path = ibdata1:512M:autoextend" >>conf/my.cnf 
	echo "innodb_autoextend_increment = 64" >>conf/my.cnf 
	echo "innodb_thread_concurrency = 0" >>conf/my.cnf 
	echo "innodb_old_blocks_time = 1000" >>conf/my.cnf 
	echo "innodb_buffer_pool_instances = 8" >>conf/my.cnf 
	echo "thread_cache_size = 200" >>conf/my.cnf
	echo "innodb_lru_scan_depth = 512" >>conf/my.cnf 
	echo "innodb_flush_neighbors = 1" >>conf/my.cnf 
	echo "innodb_checksum_algorithm = crc32" >>conf/my.cnf 
	echo "table_definition_cache = 400" >>conf/my.cnf 
	echo "innodb_buffer_pool_dump_at_shutdown = ON" >>conf/my.cnf 
	echo "innodb_buffer_pool_load_at_startup = ON" >>conf/my.cnf 
	echo "innodb_read_io_threads = 4" >>conf/my.cnf 
	echo "innodb_adaptive_flushing = ON" >>conf/my.cnf 
	echo "innodb_log_buffer_size = 8388608" >>conf/my.cnf 
	echo "innodb_purge_threads = 4" >>conf/my.cnf 
	echo "performance_schema = ON" >>conf/my.cnf 
	echo "innodb_write_io_threads = 4" >>conf/my.cnf  
	echo "skip-name-resolve = ON" >>conf/my.cnf
	echo "skip_external_locking = ON" >>conf/my.cnf 
	echo "max_allowed_packet = 1M" >>conf/my.cnf
	echo "table_open_cache = 4" >>conf/my.cnf
	echo "log-bin=mysql-bin" >>conf/my.cnf 
	echo "binlog_format=ROW" >>conf/my.cnf
	echo "innodb_flush_log_at_trx_commit = 1" >>conf/my.cnf
	echo "sync_binlog = 1" >>conf/my.cnf
	echo "loose-rpl_semi_sync_slave_enabled = 1" >>conf/my.cnf
	echo "loose-rpl_semi_sync_master_enabled = 1" >>conf/my.cnf
	echo "loose-rpl_semi_sync_master_commit_after_ack = 1" >>conf/my.cnf
	echo "loose-rpl_semi_sync_master_keepsyncrepl=0" >>conf/my.cnf
	echo "loose-rpl_semi_sync_master_timeout = 10" >>conf/my.cnf
	echo "loose-rpl_semi_sync_master_trysyncrepl=0" >>conf/my.cnf
	echo "loose-rpl_reverse_recover_enabled = ON" >>conf/my.cnf 
	echo "loose-rpl_reverse_recover_mate_host=192.168.2.96" >>conf/my.cnf
	echo "loose-rpl_reverse_recover_mate_port = 3306" >>conf/my.cnf
	echo "loose-rpl_reverse_recover_mate_user = repl" >>conf/my.cnf
	echo "loose-rpl_reverse_recover_mate_passwd = 111111" >>conf/my.cnf 
	echo "relay_log_info_repository = TABLE # slave SQL thread crash safe" >>conf/my.cnf 
	echo "master_info_repository = FILE" >>conf/my.cnf 
	echo "relay_log_recovery = ON" >>conf/my.cnf
	echo "max_allowed_packet = 16M" >>conf/my.cnf
	echo "[mysql]" >>conf/my.cnf
	echo "socket = /tmp/mysql.sock" >>conf/my.cnf 
	echo "cmha db template file configuration success!"
}

bootstrap_conf(){
	echo "appname = bootstrap" >bootstrap/conf/app.conf
	echo "runmode = dev" >>bootstrap/conf/app.conf
	echo "loglevel = 7" >>bootstrap/conf/app.conf
	echo "hostname = cmha-node1" >>bootstrap/conf/app.conf
	echo "otherhostname = cmha-node2" >>bootstrap/conf/app.conf
	echo "ip = 192.168.2.95" >>bootstrap/conf/app.conf
	echo "port = 3306" >>bootstrap/conf/app.conf
	echo "username = check" >>bootstrap/conf/app.conf
	echo "password = check" >>bootstrap/conf/app.conf
	echo "datacenter = dc1" >>bootstrap/conf/app.conf
	echo "token =" >>bootstrap/conf/app.conf
	echo "service_ip = 192.168.2.96;192.168.2.97;192.168.2.98" >>bootstrap/conf/app.conf
	echo "service_port = 8500" >>bootstrap/conf/app.conf
	echo "servicename = mysql" >>bootstrap/conf/app.conf
	echo "format = hap"  >>bootstrap/conf/app.conf
	echo "cmha bootstrap template file configuration success!"
}


mha_handlers_conf(){
	echo "appname = mha-handlers" >mha-handlers/conf/app.conf
	echo "runmode = dev" >>mha-handlers/conf/app.conf
	echo "loglevel = 7" >>mha-handlers/conf/app.conf
	echo "hostname = cmha-node1" >>mha-handlers/conf/app.conf
	echo "ip = 192.168.2.95" >>mha-handlers/conf/app.conf
	echo "port = 3306" >>mha-handlers/conf/app.conf
	echo "username = check" >>mha-handlers/conf/app.conf
	echo "password = check" >>mha-handlers/conf/app.conf
	echo "datacenter = dc1" >>mha-handlers/conf/app.conf
	echo "token =" >>mha-handlers/conf/app.conf
	echo "service_ip = 192.168.2.96;192.168.2.97;192.168.2.98" >>mha-handlers/conf/app.conf
	echo "service_port = 8500" >>mha-handlers/conf/app.conf
	echo "servicename = mysql" >>mha-handlers/conf/app.conf
	echo "format = hap" >>mha-handlers/conf/app.conf
	echo "switch = on" >>mha-handlers/conf/app.conf
	echo "cmha mha-handlers template file configuration success!"
}

monitor_handlers_conf(){
        echo "appname = monitor-handlers" >monitor-handlers/conf/app.conf
        echo "runmode = dev" >>monitor-handlers/conf/app.conf
        echo "loglevel = 7" >>monitor-handlers/conf/app.conf
        echo "ip = 192.168.2.95" >>monitor-handlers/conf/app.conf
        echo "port = 3306" >>monitor-handlers/conf/app.conf
        echo "username = check" >>monitor-handlers/conf/app.conf
        echo "password = check" >>monitor-handlers/conf/app.conf
        echo "datacenter = dc1" >>monitor-handlers/conf/app.conf
	echo "otherhostname=cmha-node2" >>monitor-handlers/conf/app.conf
        echo "token =" >>monitor-handlers/conf/app.conf
	echo "service_ip = 192.168.2.96;192.168.2.97;192.168.2.98" >>monitor-handlers/conf/app.conf
        echo "service_port = 8500" >>monitor-handlers/conf/app.conf
        echo "servicename = mysql" >>monitor-handlers/conf/app.conf
        echo "tag = slave" >>monitor-handlers/conf/app.conf
        echo "switch_async = on" >>monitor-handlers/conf/app.conf
	echo "cmha monitor-handlers template file configuration success!"
}

config_server
watch
watch_service
haproxy
keepalived
mysql
bootstrap_conf
mha_handlers_conf
monitor_handlers_conf
mysqlversion=`grep -w "mysqlversion" auto-deployment.ini |awk -F '=' '{print $2}'`
if [ ! -f "conf/config.json" ]; then
	echo "Error: consul template file not exits"
	exit 101
else
	echo "consul template file exits"
fi
if [ ! -f "conf/watch.json" ]; then
	echo "Error: consul watch kv template file not exits"
	exit 102
else
	echo "consul watch kv template file exits"
fi
if [ ! -f "conf/watch-service.json" ]; then
	echo "Error: consul watch service template file not exits"
	exit 103
else
	echo "consul watch service template file exits"
fi
if [ ! -f "conf/haproxy.ctmpl" ]; then
	echo "Error: haproxy template file not exits"
	exit 104
else
	echo "haproxy template file exits"
fi
if [ ! -f "conf/keepalived.conf" ]; then
	echo "Error: keepalived template file not exits"
	exit 105
else
	echo "keepalived template file exits"
fi
if [ ! -f "conf/my.cnf" ]; then
	echo "Error: mysql template file not exits"
	exit 106
else
	echo "mysql template file exits"
fi
if [ ! -x "bootstrap/bootstrap" ]; then
	echo "Error: bootstrap binary not exits"
	exit 107
else
	echo "bootstrap binary file exits"
fi
if [ ! -f "bootstrap/conf/app.conf" ]; then
	echo "Error: bootstrap config file not exits"
	exit 108
else
	echo "bootstrap config file exits"
fi
if [ ! -d "bootstrap/logs/" ]; then
	echo "Error: bootstrap logs dir not exits"
	exit 109
else
	echo "bootstrap logs dir exits"
fi
if [ ! -x "expect/bootstrap.exp" ]; then
	echo "Error: bootstrap.exp not exits"
	exit 110
else
	echo "bootstrap.exp exits"	
fi
if [ ! -x "expect/cmha-check.exp" ]; then
	echo "Error: cmha-check.exp not exits"
	exit 110
else
	echo "cmha-check.exp exits"
fi
if [ ! -x "bin/consul" ]; then
	echo "Error: consul binary not exits"
	exit 112
else
	echo "consul binary exits"
fi
if [ ! -x "expect/consul.exp" ]; then
	echo "Error: consul.exp not exits"
	exit 113
else
	echo "consul.exp exits"
fi
if [ ! -x "bin/consul-template" ]; then
	echo "Error: consul-template binary not exits"
	exit 114
else
	echo "consul-template binary exits"
fi
if [ ! -x "expect/consul-template.exp" ]; then
	echo "Error: consul-template.exp not exits"
	exit 115
else
	echo "consul-template.exp exits"
fi
if [ ! -x "expect/expect.exp" ]; then
	echo "Error: expect.exp not exits"
	exit 116
else
	echo "expect.exp exits"
fi
if [ ! -x "consul.sh" ]; then
	echo "Error: consul.sh not exits"
	exit 119
else
	echo "consul.sh exits"
fi
if [ ! -x "expect/haproxy.exp" ]; then
	echo "Error: haproxy.exp not exits"
	exit 118
else
	echo "haproxy.exp exits"
fi
if [ ! -x "expect/scp.exp" ]; then
	echo "Error: scp.exp not exits"
	exit 118
else
	echo "scp.exp exits"
fi
if [ ! -x "expect/scp1.exp" ]; then
	echo "Error: scp1.exp not exits"
	exit 118
else
	echo "scp1.exp exits"
fi
if [ ! -x "expect/pross.exp" ];then
	echo "Error: pross.exp not exits"
	exit 118
else
	echo "expect.exp exits"
fi
if [ ! -x "keepalived.sh" ]; then
	echo "Error: keepalived.sh not exits"
	exit 120
else
	echo "keepalived.sh exits"
fi
if [ ! -x "mha-handlers/mha-handlers" ]; then
	echo "Error: mha-handlers binary not exits"
	exit 121
else
	echo "mha-handlers binary exits"
fi
if [ ! -f "mha-handlers/conf/app.conf" ]; then
	echo "Error: mha-handlers config file not exits"
	exit 122
else
	echo "mha-handlers config file exits"
fi
if [ ! -d "mha-handlers/logs/" ]; then
	echo "Error: mha-handlers logs dir not exits"
	exit 123
else
	echo "mha-handlers logs dir exits"
fi
if [ ! -x "monitor-handlers/monitor-handlers" ]; then
	echo "Error: monitor-handlers binary not exits"
	exit 124
else
	echo "monitor-handlers binary exits"
fi
if [ ! -f "monitor-handlers/conf/app.conf" ]; then
	echo "Error: monitor-handlers config file not exits"
	exit 125
else
	echo "monitor-handlers config file exits"
fi
if [ ! -d "monitor-handlers/logs" ]; then
	echo "Error: monitor-handlers logs dir not exits"
	exit 126
else
	echo "monitor-handlers logs dir exits"
fi
if [ ! -x "expect/mysql-change.exp" ]; then
	echo "Error: mysql-change.exp not exits"
	exit 128
else
	echo "mysql-change.exp exits"
fi
if [ ! -x "mysqlcheck.sh" ]; then
	echo "Error: mysqlcheck.sh not exits"
	exit 129 
else
	echo "mysqlcheck.sh exits"
fi
if [ ! -x "expect/mysql.exp" ]; then
	echo "Error: mysql.exp not exits"
	exit 130
else
	echo "mysql.exp exits"
fi
if [ ! -x "expect/mysql-install-plugin.exp" ]; then
	echo "Error: mysql-install-plugin.exp not exits"
	exit 131
else
	echo "mysql-install-plugin.exp exits"
fi
if [ ! -x "expect/mysql-relationship.exp" ]; then
	echo "Error: mysql-relationship.exp not exits"
	exit 132
else
	echo "mysql-relationship.exp exits"
fi
if [ ! -f "package/haproxy-1.5.2-2.el6.x86_64.rpm" ]; then
	echo "Error: haproxy-1.5.2-2.el6.x86_64.rpm not exits"
	exit 135
else
	echo "haproxy-1.5.2-2.el6.x86_64.rpm exits"
fi
if [ ! -f "package/keepalived-1.2.13-5.el6_6.x86_64.rpm" ]; then
	echo "Error: keepalived-1.2.13-5.el6_6.x86_64.rpm not exits"
	exit 136
else
	echo "keepalived-1.2.13-5.el6_6.x86_64.rpm exits"
fi
if [ ! -f "package/${mysqlversion}" ]; then
	echo "Error: ${mysqlversion} not exits"
	exit 137
else
	echo "${mysqlversion} exits"
	echo "OK"
fi
