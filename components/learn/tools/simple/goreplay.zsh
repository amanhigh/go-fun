#!/bin/zsh
# https://github.com/buger/goreplay/wiki/Getting-Started

#Usage
if [ $# -eq 0 ]
    then
    echo "Usage: $0 stdout|replay|save|load"
    exit 1 
fi

action=$1

# Switch Action
case $action in
    "stdout")
        echo "\033[1;33m Listening on Port: 8085 \033[0m"
        sudo gor -input-raw :8085 --output-stdout
        ;;
    "replay")
        echo "\033[1;33m Replaying Traffic from 8085 to 8080 \033[0m"
        sudo gor --input-raw :8085 --output-http="http://localhost:8080"
        ;;
    "save")
        echo "\033[1;33m Saving Traffic from 8085 to File \033[0m"
        sudo gor --input-raw :8085 --output-file=requests.gor
        ;;
    "load")
        echo "\033[1;33m Saving Traffic from 8085 to File \033[0m"
        sudo gor --input-file=requests.gor --output-http="http://localhost:8080"
        ;;
    *)
        echo "\033[1;31m Invalid Action \033[0m"
        ;;
esac

