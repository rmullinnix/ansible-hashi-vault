#!/bin/bash
#
# vault        Manage the vault agent
#       
# chkconfig:   2345 95 95
# description: Vault is a secret store connected to consul
# processname: vault
# config: {{ vault_app_path }}/conf/vault.hcl
# pidfile: /var/run/vault.pid
 
### BEGIN INIT INFO
# Provides:       vault
# Required-Start: $local_fs $network
# Required-Stop:
# Should-Start:
# Should-Stop:
# Default-Start: 2 3 4 5
# Default-Stop:  0 1 6
# Short-Description: Manage the vault agent
# Description: Valut is a secret store connected to consul
### END INIT INFO
. /etc/rc.status
rc_reset
 
prog="vault"
user="{{ consul_user }}"
group="{{ consul_group }}"
exec="{{ vault_app_path }}/bin/$prog"
pidfile="{{ pid_path }}/$prog.pid"
lockfile="{{ lock_path }}/subsys/$prog"
logfile="{{ inf_log_path }}/$prog.log"
conffile="{{ vault_app_path }}/conf/vault.hcl"
 
export GOMAXPROCS=${GOMAXPROCS:-2}
 
start() {
    [ -x $exec ] || exit 5
    
    [ -f $conffile ] || exit 6
 
    umask 033
 
    touch $logfile $pidfile
    chown $user:$group $logfile $pidfile
 
    ulimit -v unlimited

    echo -n $"Starting $prog: $1"

    ## use the start daemone or nohup
    /sbin/start_daemon -f \
        -p $pidfile \
	-u $user \
           $exec server -config=$conffile -log-level={{ vault_log_level }} 2>&1 >> $logfile & echo $! > $pidfile

    ##	nohup $exec server -config=$conffile -log-level={{ vault_log_level }} 2>&1 >> $logfile & echo $! > $pidfile &
    
    RETVAL=$?
    echo
    
    [ $RETVAL -eq 0 ] && touch $lockfile
    
    return $RETVAL
}
 
stop() {
    echo -n $"Shutting down $prog: "
    ## graceful shutdown with SIGINT
    /sbin/start-stop-daemon -K -p $pidfile -u $user -x $exec -s 15
    RETVAL=$?
    echo
    [ $RETVAL -eq 0 ] && rm -f $lockfile
    return $RETVAL
}
 
restart() {
    stop
    sleep 4
    start
}
 
reload() {
    echo -n $"Reloading $prog: "
    /sbin/start-stop-daemon -K -p $pidfile -u $user -x $exec -s 1
    RETVAL=$?
    return $RETVAL
}
 
force_reload() {
    restart
}
 
rh_status() {
    checkproc -p $pidfile $exec
    rc_status -v
    $exec version
}
 
rh_status_q() {
    rh_status >/dev/null 2>&1
}
 
case "$1" in
    start)
        $1
	rc_status -v
        ;;
    stop)
        $1
	rc_status -v
        ;;
    restart)
        $1
	rc_status -v
        ;;
    reload)
        $1
	rc_status -v
        ;;
    force-reload)
        force_reload
        ;;
    status)
        rh_status
        ;;
    *)
        echo $"Usage: $0 {start|stop|status|restart|reload|force-reload}"
        exit 2
esac
 
rc_exit
