#!/bin/zsh
echo "\033[1;32m Bash Fun \033[0m \n"

# Script Directory
SCRIPT_DIR=`dirname $0`
echo "\033[1;32m Script is Located in $SCRIPT_DIR \033[0m"

# Date Formatting
formatted_date=$(date +"%Y-%m-%d %H:%M:%S")
echo "Formatted date: $formatted_date"

# Backgroud Process
sleep 5 &
echo "Bring Sleep to Foreground and Wait"
wait