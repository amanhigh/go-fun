
answers=`gum choose RESET MYSQL MONGO --limit 5`
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

    *)
        echo -en "\033[1;34m Mongo \033[0m \n"
        ;;
    esac
done
