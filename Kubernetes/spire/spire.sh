echo -en "\033[1;32m Setting Up Spire Server \033[0m \n"
kubectl -f ./spire-server.yaml apply
echo -en "\033[1;32m Setting Up Spire Agent \033[0m \n"
kubectl -f ./spire-agent.yaml apply
