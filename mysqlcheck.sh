#!/bin/bash
####version 1.1.2
name=$1
dbname=`tr '[A-Z]' '[a-z]' <<<"$name"`
#basedir="/bsgchina/$dbname/MYBASE"
datadir="/tmp"
bindir="/usr/local/cmha/mysql/bin"

MYSQL=$bindir/mysql
MYSQL_USER=$2
MYSQL_PASSWORD=$3
CHECK_TIME=3
MYSQL_OK=1
HOST=$1
DATABASE=cmha_check
TIMEOUT=5

function check_mysql_health () {


${MYSQL} -u${MYSQL_USER} -p${MYSQL_PASSWORD} -h${HOST} --connect-timeout=${TIMEOUT} -D${DATABASE} <<EOF
set sql_log_bin=0;
set innodb_lock_wait_timeout=5;
begin;
select cmha_name from cmha_check;
delete from cmha_check;
insert into cmha_check(id,cmha_name,create_time) values(1,'cmha_check_insert',now());
update cmha_check set cmha_name='cmha_check_update' where id=1;
commit;
EOF
#$MYSQL  -u $MYSQL_USER -p$MYSQL_PASSWORD -h ${HOST} -e "show status;" > /dev/null 2>&1

if [ $? -eq 0 ]; then 
	MYSQL_OK=0;
else
	MYSQL_OK=1;
fi

	return $MYSQL_OK;
}

while [ $CHECK_TIME -ne 0 ]
do
	let "CHECK_TIME -= 1";
	check_mysql_health;
	if [ $MYSQL_OK = 0 ]; then
	#	CHECK_TIME=0;
		isyes=`${MYSQL} -u${MYSQL_USER} -p${MYSQL_PASSWORD} -h${HOST} -e 'show slave status\G' | grep Slave_IO_Running | awk -F: '{print $2}' |  sed 's/ //g'`
		if [ ${isyes} == "Yes" ]; then
			exit 0
		else
			exit 1
		fi
	fi

	if [ $MYSQL_OK -eq 1 ] && [ $CHECK_TIME -eq 0 ]; then
		exit 2;
	fi
	sleep 1;
done
