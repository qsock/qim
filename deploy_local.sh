#!/bin/bash

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

git checkout master
git pull --rebase origin master

make p

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

    \cp $DIR/sample.service $DIR/$1.service
    \sed 's/sample/'$1'/g' $DIR/$1.service -i
    sudo cp $DIR/$1.service /usr/lib/systemd/system/
    sudo systemctl daemon-reload
    \rm $DIR/$1.service
    \mkdir -p /srv/$1/bin/ /srv/$1/logs/ /srv/$1/conf/

    # 拷贝配置文件
    \cp $DIR/config/dev/$1.toml /srv/$1/conf/
    # rm 防止file busy
    \rm /srv/$1/bin/$1
    \cp $GOBIN/$1 /srv/$1/bin/
    sudo systemctl restart $1 && sudo systemctl is-active $1
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