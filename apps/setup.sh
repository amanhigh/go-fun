#!/usr/bin/env bash
set -e

APPS_PATH=$(dirname $(cd .;pwd -P))/apps
echo -en "\033[1;32m Script is Idempotent can be run multiple times
 Please Run With Sudo as it requires creating soflinks in /etc/
 \033[0m \n"
echo -en "\033[1;34m Using APPS Source Path as $APPS_PATH \033[0m \n"

#Fun App
echo -en "\033[1;32m Fun App \033[0m \n"
FUNAPP_CONFIG=/etc/fun-app
rm -rf ${FUNAPP_CONFIG}; mkdir -p ${FUNAPP_CONFIG}

#sudo ln -s ${APPS_PATH}/components/fun-app/config.yml ${FUNAPP_CONFIG}/config.yml

mkdir -p /var/log/fun-app
chmod 777 /var/log/fun-app


