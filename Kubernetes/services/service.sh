# Vars
CMD="install"

# Prompt
answers=`gum choose MYSQL MONGO REDIS PROXY CRON --limit 5`

# Flags
while getopts 'du' OPTION; do
  case "$OPTION" in
    d)
        echo -en "\033[1;32m Clearing all Helms \033[0m \n"
        helm delete $(helm list --short)
        ;;
    u)
        echo -en "\033[1;32m Switching to Upgrade \033[0m \n"
        CMD="upgrade"
        ;;
    ?)
      echo "script usage: $0 [-u] [-r]" >&2
      exit 1
      ;;
  esac
done

for SVC in $answers
do
    echo -en "\033[1;32m \n $SVC \033[0m \n"
    case $SVC in
    MYSQL)
        helm $CMD mysql bitnami/mysql -f mysql.yml > /dev/null
        helm $CMD mysql-admin bitnami/phpmyadmin > /dev/null
        echo -en "\033[1;33m Login: mysql-primary, root/root \033[0m \n"
        ;;

    MONGO)
        helm $CMD mongo bitnami/mongodb -f mongo.yml > /dev/null
        echo -en "\033[1;33m mongosh -u root -p root --host localhost  < /etc/files/scripts/mongo.js \033[0m \n"
        ;;

    REDIS)
        helm $CMD redis bitnami/redis -f redis.yml > /dev/null
        helm $CMD redis-admin onechart/onechart -f redis-admin.yml

        echo -en "\033[1;33m redis-cli -c incr mycounter \033[0m \n"
        echo -en "\033[1;33m redis-cli -c set mypasswd lol \033[0m \n"
        echo -en "\033[1;33m redis-cli -c get mypasswd \033[0m \n"
        echo -en "\033[1;33m Commander: http://localhost:8081/ \033[0m \n"
        ;;

    PROXY)
        echo -en "\033[1;32m Resty \033[0m \n"
        helm $CMD resty onechart/onechart -f resty.yml > /dev/null

        echo -en "\033[1;33m Commander: http://localhost:8090/ \033[0m \n"
        echo -en "\033[1;33m Commander: http://localhost:8090/example \033[0m \n"
        echo -en "\033[1;33m Commander: http://localhost:8090/ndtv \033[0m \n"

        ;;
    CRON)
        # TODO: Fix Cron
        helm $CMD cron onechart/onechart -f cron.yml > /dev/null
        echo -en "\033[1;33m Check Logs for Output \033[0m \n"
        ;;

    *)
        echo -en "\033[1;34m Service Not Supported: $SVC \033[0m \n"
        ;;
    esac
done
