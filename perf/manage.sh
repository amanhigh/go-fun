#!/usr/bin/env bash

function install() {
    echo -en '\033[1;34m **Installing Dependencies ** \033[0m \n'
    sudo apt-get install -y --force-yes build-essential libssl-dev;

    echo -en '\033[1;34m **Building Wrk ** \033[0m \n'
    cd ~/perf/wrkd;
    make &> /dev/null;

    echo -en '\033[1;34m **Linking Wrk ** \033[0m \n'
    ln -s ~/perf/wrkd/wrk ~/perf/wrk

    echo -en '\033[1;34m **Setting Permissions ** \033[0m \n'
    chmod 755 ~/perf/wrk

    echo -en "\033[1;32m Install Done \033[0m \n"
}

function upload() {
    echo -en "\033[1;34m Cleaning Perf \033[0m \n"
    ssh $ip "rm -rf ~/perf"

    echo -en "\033[1;34m Uploading Setup as $user to $ip \033[0m \n"
    rsync -avzh . $user@$ip:~/perf
    echo -en "\033[1;32m Upload Complete \033[0m \n"

    echo -en "\033[1;34m Unpacking \033[0m \n"
    ssh $ip "tar -jxf perf.tar.bz2 2> /dev/null"
    ssh $ip "cd ~/perf; tar -jxf wrkd.tar.bz2 2> /dev/null"

    ssh $ip "~/perf/manage.sh install"
}

function sync() {
    watch -n1 "rsync -avzh . $user@$ip:~/perf"
}

case $1 in
    install) install ;;
    upload)
    if [ "$#" -lt 2 ]; then
        echo "Usage upload <ip> <optional user>"
    fi

    ip=${2:-10.34.238.175}
    user=${3:-amanpreet.singh}
    upload $ip $user
    ;;
    sync)
    if [ "$#" -lt 2 ]; then
        echo "Usage sync <ip> <optional user>"
    fi

    ip=${2:-10.34.238.175}
    user=${3:-amanpreet.singh}
    sync
    ;;
    *) echo "Usage: build/run/install"; ;;
esac