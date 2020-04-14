echo -en "\033[1;32m Deleting Minikube Clusters \033[0m \n"
minikube -p secondary minikube
minikube -p secondary delete