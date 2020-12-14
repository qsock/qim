#!/bin/bash

# 确定user
#username=`whoami`
#if [ $username != "centos" ]; then
#    echo "You must run this script as centos."
#    exit 1
#fi;


# TODO fix for product
exit(1)

DIR="."

declare -A SERVERS

SERVERS[api_gateway]=api
SERVERS[file_server]=file_server
SERVERS[id_server]=id_server
SERVERS[msg_server]=msg_server
SERVERS[passport_server]=passport_server
SERVERS[user_server]=user_server
SERVERS[ws_server]=ws_server

cd $DIR
if git status --porcelain | read; then
    echo "Uncommitted or untracked files found"
    exit 1
fi

git checkout dev
git pull --rebase origin dev

#make dep
make p

# 集群部署就填写多个
DEFAULT_HOSTS="in-prod-common-goserver-2"
DEFAULT_SLEEP_TIME=2

declare -A SLEEP_TIME
#SLEEP_TIME[msg_server]=5

# 申明部署的机器
declare -A MACHINES
#MACHINES[push_server]="in-prod-common-push-1"


# Deploy function
deploy() {
    if [ -z "${SERVERS[$1]}" ]; then
        echo "No service: $1"
        exit 1
    fi;
    make ${SERVERS[$1]}
    if [ $? -ne 0 ]; then
        echo "Make ${SERVERS[$1]} failed "
        exit 1
    fi

    if [ ! -z "${MACHINES[$1]}" ]; then
        read -r -a HOSTS <<< ${MACHINES[$1]}
    else
        read -r -a HOSTS <<< $DEFAULT_HOSTS
    fi

    #开始部署
    for host in "${HOSTS[@]}"; do
        # 备份原有的bin文件,方便做恢复回滚
        ssh worker@$host "mkdir -p /srv/$1/bin/;mkdir -p /srv/$1/logs/;mkdir -p /srv/$1/conf/;/usr/bin/cp /srv/$1/bin/$1 /srv/$1/bin/$1.backup"
        # rsync 配置文件
        rsync --rsync-path="sudo rsync" -avCL $DIR/config/prod/$1.toml worker@$host:/srv/$1/conf/
        # rsync bin文件
        rsync --rsync-path="sudo rsync" -avCL $GOBIN/$1 worker@$host:/srv/$1/bin/

        \rm $DIR/$1.service

        \sed 's/sample/'$1'/g' $DIR/$1.service -i
        rsync --rsync-path="sudo rsync" -avCL $DIR/$1.service worker@$host:/usr/lib/systemd/system/
        ssh worker@$host "sudo systemctl daemon-reload && sudo systemctl restart $1 && sudo systemctl status $1"
        if [ $host == ${HOSTS[-1]} ]; then # 最后一台机器不 sleep
            break
        fi
        if [ ! -z "${SLEEP_TIME[$1]}" ]; then
            sleep ${SLEEP_TIME[$1]}
        else
            sleep $DEFAULT_SLEEP_TIME
        fi
    done;
}

if [ ! -z "$1" ]; then
    for i in "$@"; do
        deploy $i
    done
else # Deploy all servers
    for i in "${!SERVERS[@]}"; do
        deploy $i
    done
fi
