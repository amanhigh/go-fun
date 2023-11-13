#!/bin/zsh
echo "\033[1;32m Bash Fun \033[0m \n"

sleep 5 &
echo "Bring Sleep to Foreground and Wait"
wait

# Script Directory
SCRIPT_DIR=`dirname $0`
echo "\033[1;32m Script is Located in $SCRIPT_DIR \033[0m"