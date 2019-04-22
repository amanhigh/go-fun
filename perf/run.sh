#!/bin/bash
echo "Usage: perf.sh <script> <host:port> <concurrency> <timeout (minutes)>"
script="$1.lua"
host=${2:-localhost\:8080}
con=${3:-10}
if [ $con == 1 ];then
       th=1;
    else
        th=2;
fi
min=${4:-1}
echo -en "\n\033[1;32m Running Pref(${script}) on $host for $min Minute with $con Concurrency ($th Threads) \033[0m \n"
wrk -t$th -c$con -d${min}m --latency -s ./lua/${script} http://$host
