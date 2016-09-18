#!/bin/sh
# version 1.1.5-Beta
# consul        Start the consul daemon
#
# chkconfig: 345 99 10
# description: Starts the Consul daemon
#
# processname: consul

# Source function library.
. /etc/rc.d/init.d/functions

RETVAL=0

# Default variables
CONSUL_BIN="/usr/local/bin/consul"
CONSUL_CONF="/usr/local/cmha/consul.d"
CONSUL_PID=`pidof consul`
CONSUL_RPC=`grep -w "rpc" $CONSUL_CONF/config.json|awk '{print $2}' |sed 's/\"//g'`
# called.
start(){
if [ "$CONSUL_PID" != "" ];then
        status consul
        RETVAL=$?
	exit $RETVAL
else
	echo -n "Starting Consul daemon: "
        daemon --check consul "$CONSUL_BIN agent -config-dir=$CONSUL_CONF > /dev/null 2>&1 &"
        echo
fi
}
stop(){
if [ "$CONSUL_PID" = "" ];then
        status consul
        RETVAL=$?
        exit $RETVAL
else    
        echo -n "Stoping Consul daemon: "
        killproc consul
        echo    
fi
}
reload(){
if [ "$CONSUL_PID" = "" ];then
	status consul
	RETVAL=$?
	exit $RETVAL
else
	echo -n "Reload Consul daemon: "
        daemon --check consul "$CONSUL_BIN reload -rpc-addr=${CONSUL_RPC}:8400"
        RETVAL=$?
        echo
fi
}
case "$1" in
  start)
	start
        ;;
  stop)
	stop
        ;;
  status)
        status consul
        RETVAL=$?
        ;;
  restart)
        $0 stop
        $0 start
        RETVAL=$?
        ;;
 reload )
	reload
        ;;
  *)
        echo "Usage: consul {start|stop|status|restart|reload}"
        exit 1
esac

exit $REVAL
