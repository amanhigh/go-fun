answers=`gum choose RESET MYSQL MONGO REDIS --limit 5`
for SVC in $answers
do
    echo -en "\033[1;32m \n $SVC \033[0m \n"
    case $SVC in
    RESET)
        echo -en "\033[1;32m Clearing all Helms \033[0m \n"
        helm delete $(helm list --short)
        ;;

    MYSQL)
        helm install mysql bitnami/mysql -f mysql.yml > /dev/null
        helm install mysql-admin bitnami/phpmyadmin > /dev/null
        echo -en "\033[1;33m Login: mysql-primary, root/root \033[0m \n"
        ;;

    MONGO)
        helm install mongo bitnami/mongodb -f mongo.yml > /dev/null
        echo -en "\033[1;33m mongosh -u root -p root --host localhost  < /etc/files/scripts/mongo.js \033[0m \n"
        ;;

    REDIS)
        #TODO: Move commander generic helm
        helm install redis bitnami/redis -f redis.yml > /dev/null
        dman run rn5 1 'redis-cli -c incr mycounter'
        echo -en "\033[1;33m redis-cli -c incr mycounter \033[0m \n"
        echo -en "\033[1;33m redis-cli -c set mypasswd lol \033[0m \n"
        echo -en "\033[1;33m redis-cli -c get mypasswd \033[0m \n"
        ;;

    *)
        echo -en "\033[1;34m Service Not Supported: $SVC \033[0m \n"
        ;;
    esac
done
