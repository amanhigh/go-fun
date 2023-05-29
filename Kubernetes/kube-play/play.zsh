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
  echo "delete"
  SELECTED_CMD=$DELETE_CMD
fi

echo $SELECTED_CMD
echo $REPLY
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
        mysql
        ;;
    *)
        echo "Invalid option selected."
        ;;
esac
