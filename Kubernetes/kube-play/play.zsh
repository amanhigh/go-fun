if [ -z "$1" ]; then
  echo "Usage: $0 [FileName] <Flag:d>"
  exit 1
fi

REPLY=$1

KUBE_CMD="kubectl "
APPLY_CMD="$KUBE_CMD apply -f"
DELETE_CMD="$KUBE_CMD delete -f"
SELECTED_CMD=$APPLY_CMD
CREATE=0

if [ "$2" = "-d" ]; then
  echo "\033[1;31m DELETING: $REPLY \033[0m \n";
  SELECTED_CMD=$DELETE_CMD
fi

case $REPLY in
    basic.yml)
        echo "\033[1;32m Basic Service with Deployment \033[0m \n";
        $SELECTED_CMD $REPLY
        echo "\033[1;33m Show Pods \033[0m \n";
        kubectl get pods -l app=nginx
        echo "\033[1;33m Show Services \033[0m \n";
        kubectl get service -l app=nginx
        ;;
    mysql.yml)
        echo "\033[1;32m Deploying Mysql with Busybox (Sidecar) \033[0m \n";
        echo "\033[1;34m Mysql Client Installed in Sidecar Using Post LifeCycle \033[0m \n";
        $SELECTED_CMD $REPLY
        echo "\033[1;33m Show Services \033[0m \n";
        kubectl get service -l app=mysql
        #TODO Run in Sidecar
        echo "\033[1;33m Mysql Health Check from Sidecar \033[0m \n";
        kubectl exec -it $(kubectl get pods -l app=mysql -o jsonpath='{.items[0].metadata.name}') -c sidecar -- sh -c 'if pgrep mysqld >/dev/null 2>&1; then echo "MySQL process is running"; else echo "MySQL process is not running"; fi'
        echo "\033[1;33m Show Databases \033[0m \n";
        kubectl exec -it $(kubectl get pods -l app=mysql -o jsonpath='{.items[0].metadata.name}') -c mysql -- /bin/sh -c 'mysql -h 127.0.0.1 -u root -proot -e "show databases;"'
        echo "\033[1;33m Master Status \033[0m \n";
        kubectl exec -it $(kubectl get pods -l app=mysql -o jsonpath='{.items[0].metadata.name}') -c mysql -- /bin/sh -c 'mysql -h 127.0.0.1 -u root -proot -e "SHOW MASTER STATUS\G;"'
        echo "\033[1;33m Kill DB from Sidecar \033[0m \n";
        kubectl exec -it $(kubectl get pods -l app=mysql -o jsonpath='{.items[0].metadata.name}') -c sidecar -- sh -c 'pkill mysqld'
        ;;
    slave.yml)
        echo "\033[1;32m Deploying Mysql Slave pointing to mysql-service \033[0m \n";
        $SELECTED_CMD $REPLY
        echo "\033[1;33m Show Services \033[0m \n";
        kubectl get service -l app=mysql-slave
        echo "\033[1;33m Show Databases \033[0m \n";
        kubectl exec -it $(kubectl get pods -l app=mysql-slave -o jsonpath='{.items[0].metadata.name}') -c mysql -- /bin/sh -c 'mysql -h 127.0.0.1 -u root -proot -e "show databases;"'
        echo "\033[1;33m Slave Status (Seconds_Behind_Master,Last_IO_Error) \033[0m \n";
        kubectl exec -it $(kubectl get pods -l app=mysql-slave -o jsonpath='{.items[0].metadata.name}') -c mysql -- /bin/sh -c 'mysql -h 127.0.0.1 -u root -proot -e "SHOW SLAVE STATUS\G;"'
        ;;
    stateful.yml)
        echo "\033[1;32m Deploying Nginx with Stateful Set \033[0m \n";
        $SELECTED_CMD $REPLY
        echo "\033[1;33m Show Stateful Sets \033[0m \n";
        kubectl get statefulset -l app=nginx
        echo "\033[1;33m Show Headless Service (No Cluster IP) \033[0m \n";
        kubectl get service -l app=nginx
        echo "\033[1;33m Show Pods \033[0m \n";
        kubectl get pods -l app=nginx
        echo "\033[1;33m Volume Claims \033[0m \n";
        kubectl get pvc -l app=nginx
        echo "\033[1;33m Host Names \033[0m \n";
        for i in 0 1 2; do kubectl exec -c nginx "nginx-statefulset-$i" -- sh -c 'hostname'; done
        echo "\033[1;33m Deploying Sidecar \033[0m \n";
        #Debug Sidecar attach: kubectl run -i --tty --image busybox:1.28 debug --restart=Never --rm
        # kubectl exec -it nginx-statefulset-0 -c sidecar -- sh -c 'nslookup nginx-statefulset-0'
        echo "\033[1;33m Delete Pod \033[0m \n";
        gum confirm "Delete Pods?" --timeout=5s --default="No" && kubectl delete pod -l app=nginx && sleep 20

        echo "\033[1;33m Host Names (Same Post Delete) \033[0m \n";
        for i in 0 1 2; do kubectl exec -c nginx "nginx-statefulset-$i" -- sh -c 'hostname'; done
        
        echo "\033[1;33m Scale from 3 to 5 \033[0m \n";
        gum confirm "Scale Up?" --timeout=5s --default="No" && kubectl scale sts nginx-statefulset --replicas=5 && sleep 5
        
        echo "\033[1;33m Show Pods (Post Scaling Up) \033[0m \n";
        kubectl get pods -l app=nginx

        echo "\033[1;33m Scale from 5 to 2 \033[0m \n";
        gum confirm "Scale Down?" --timeout=5s --default="No" && kubectl patch sts nginx-statefulset -p '{"spec":{"replicas":2}}' && sleep 15
        
        echo "\033[1;33m Show Pods (Post Scaling Down) \033[0m \n";
        kubectl get pods -l app=nginx
        ;;
    *)
        echo "Invalid option selected."
        ;;
esac
