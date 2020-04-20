echo -en "\033[1;32m Cleaning Spire \033[0m \n"
kubectl -f ./spire-server.yaml delete
kubectl -f ./spire-agent.yaml delete

echo -en "\033[1;32m Cleaning Test Client \033[0m \n"
kubectl delete -f sleep.yaml -n foo
kubectl delete ns foo
