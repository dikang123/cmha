#!/bin/sh
<<<<<<< HEAD
# version 1.1.5-Beta
=======
# version 1.1.5-Beta.1
>>>>>>> 126d33b0306a2de4f2f5445489f9e46636c7c67e
# consul-template        Start the consul-template daemon
#
# chkconfig: 345 99 10
# description: Starts the consul-template daemon
#
# processname: consul-template

# Source function library.
. /etc/rc.d/init.d/functions

RETVAL=0

# Default variables
CONSUL_BIN="/usr/local/bin/consul-template"
CONSUL_CONF="/usr/local/cmha/consul-template.d"
CONSUL_TEMPLATE_PID=`pidof consul-template`
# called.
start(){
if [ "$CONSUL_TEMPLATE_PID" != "" ];then
        status consul-template
        RETVAL=$?
	exit $RETVAL
else
	echo -n "Starting Consul-template daemon: "
        daemon --check consul-template "$CONSUL_BIN -config=$CONSUL_CONF/consul-template.conf > /dev/null 2>&1 &"
        echo
fi
}
stop(){
if [ "$CONSUL_TEMPLATE_PID" = "" ];then
        status consul-template
        RETVAL=$?
        exit $RETVAL
else    
        echo -n "Stoping Consul-template daemon: "
        killproc consul-template
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
        status consul-template
        RETVAL=$?
        ;;
  restart)
        $0 stop
        $0 start
        RETVAL=$?
        ;;
  *)
        echo "Usage: consul-template {start|stop|status|restart}"
        exit 1
esac

exit $REVAL
