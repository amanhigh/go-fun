#!/usr/bin/env bash
baseCmd="docker-compose \
`cat /tmp/docker-manage-v1` action"

finalCmd="$baseCmd"

function start()
{
     echo -en "\033[1;34m Starting Docker Setup \033[0m \n"
     eval "${finalCmd//action/up -d}"
}

function restart()
{
    if [ "$#" -lt 1 ]; then
    	echo "Usage: $0 restart <service name> "
    	return -1;
	fi
    svcName=$1

    echo -en "\033[1;34m Restarting Docker Service $svcName \033[0m \n"
    eval "${finalCmd//action/up -d --force-recreate --no-deps} $svcName"
}

function stop()
{
     eval "${baseCmd//action/stop}"
}

function kill()
{
     eval "${baseCmd//action/rm -svf}"
}

function reset()
{
    #Stop and Clean Containers
    docker-clean stop

    #Fire up fresh Setup
    start

    echo -en "\033[1;34m Running Seed Script \033[0m \n"
}

function ps()
{
    eval "watch -n1 ${baseCmd//action/ps}"
}

function logs()
{
    eval "${baseCmd//action/logs} ${1}"
}

function login()
{
    if [ "$#" -lt 2 ]; then
    	echo "Usage: $0 login <service name> <container number> Eg. $0 login hcs 1"
    	return -1;
	fi
	docker exec -it compose_$1_$2 /bin/bash
}

function build-nocache()
{
    if [ "$#" -lt 1 ]; then
    	echo "Usage $0 build <Image>"
    	return -1;
	fi

	imageName=$1

    docker build ./$imageName -t $imageName:latest --no-cache
}

function set()
{
    if [ "$#" -lt 1 ]; then
    	echo "Usage $0 set <Path>"
    	return -1;
	fi

	echo -en "\033[1;34m Setting '$1' to /tmp/docker-manage-v1 \033[0m \n"
    echo $1 > /tmp/docker-manage-v1
}

case "$1" in
start)
    start
    ;;
restart)
    restart $2
    ;;
stop)
    stop
    ;;
kill)
    kill
    ;;
reset)
    reset
    ;;
ps)
    ps
    ;;
logs)
    logs "$2"
    ;;
set)
    set "$2"
    ;;
login)
    login $2 $3
    ;;
build)
    build-nocache $2
    restart $2
    ;;
*)  echo "
Valid docker-manage Options:

* start/restart <srvc>/stop/kill
* reset/ps/login <srvc> <no.>
* build <Image Name>
* set <file| -f ./compose/grafana.yml>"
    ;;
esac
