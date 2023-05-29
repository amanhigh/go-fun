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
        echo "\033[1;33m Show Databases \033[0m \n";
        kubectl exec -it $(kubectl get pods -l app=mysql -o jsonpath='{.items[0].metadata.name}') -c mysql -- /bin/sh -c 'mysql -h 127.0.0.1 -u root -proot -e "show databases;"'
        echo "\033[1;33m Mysql Health Check from Sidecar \033[0m \n";
        kubectl exec -it $(kubectl get pods -l app=mysql -o jsonpath='{.items[0].metadata.name}') -c sidecar -- sh -c 'if pgrep mysqld >/dev/null 2>&1; then echo "MySQL process is running"; else echo "MySQL process is not running"; fi'
        echo "\033[1;33m Kill DB from Sidecar \033[0m \n";
        kubectl exec -it $(kubectl get pods -l app=mysql -o jsonpath='{.items[0].metadata.name}') -c sidecar -- sh -c 'pkill mysqld'
        ;;
    *)
        echo "Invalid option selected."
        ;;
esac
