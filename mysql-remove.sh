#!/bin/bash
####version 1.1.4
#rm -fr /tmp/consul/
basedir=`awk -F '=' '/\[cmha_agent1\]/{a=1}a==1&&$1~/mysql-master-basedir/{print $2;exit}' auto-deployment.ini`
datadir=`awk -F '=' '/\[cmha_agent1\]/{a=1}a==1&&$1~/mysql-master-datadir/{print $2;exit}' auto-deployment.ini`
pkill consul
rm -fr /etc/consul.d/
rm -fr /usr/local/bin/consul
rm -fr /usr/local/cmha/
rm -fr /etc/init.d/mysql
rm -fr /etc/my.cnf
rm -fr ${basedir}
rm -fr ${datadir}
#rm -fr /var/lib/mysql/
kill -9 $(netstat -tlnp|grep mysqld|awk '{print $7}'|awk -F '/' '{print $1}')
userdel -f mysql
grep  "/usr/local/cmha/mysql/bin" /etc/profile 1>/dev/null
if [ $? -eq 0 ]; then
        sed -i '$d' /etc/profile
fi
