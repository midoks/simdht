#!/bin/bash
# chkconfig: 2345 55 25
# description: simdht Service

### BEGIN INIT INFO
# Provides:          bt
# Required-Start:    $all
# Required-Stop:     $all
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: starts simdht
# Description:       starts the simdht
### END INIT INFO

PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin

# Source function library.
if [ -f /etc/init.d/functions ];then
  . /etc/init.d/functions
fi

if [ -f /etc/rc.d/init.d/functions ];then
  . /etc/rc.d/init.d/functions
fi

app_path={APP_PATH}
SERVICENAME="simdht"

im_start(){
    isStart=`ps -ef|grep 'simdht service' |grep -v grep|awk '{print $2}'`
    if [ "$isStart" == '' ];then
        echo -e "Starting simdht... \c"
        cd $app_path && daemon "${app_path}/simdht service &"
        isStart=""
        while [[ "$isStart" == "" ]];
        do
            echo -e ".\c"
            sleep 0.5
            isStart=$(lsof -n -P -i:25|grep LISTEN|grep -v grep|awk '{print $2}'|xargs)
            let n+=1
            if [ $n -gt 15 ];then
                break;
            fi
        done
        echo -e "\033[32mdone\033[0m"
    else
        echo "Starting ${SERVICENAME}(pid $(echo $isStart)) already running"
    fi
}

im_stop(){
	pids=`ps -ef|grep 'simdht service' |grep -v grep|awk '{print $2}'`
    arr=($pids)
    echo -e "Stopping ${SERVICENAME}... \c"
    for p in ${arr[@]}
    do
            kill -9 $p
    done
    echo -e "\033[32mdone\033[0m"
}

im_status(){
    isStart=`ps -ef|grep 'simdht service' |grep -v grep|awk '{print $2}'`
    if [ "$isStart" == '' ];then
      echo -e "${SERVICENAME} not running"
    else
      echo -e "${SERVICENAME}(pid $(echo $isStart)) already running"
    fi
}

case "$1" in
    'start') im_start;;
    'stop') im_stop;;
    'status') im_status;;
    'reload') 
        im_stop
        im_start;;
    'restart') 
        im_stop
        im_start;;
esac