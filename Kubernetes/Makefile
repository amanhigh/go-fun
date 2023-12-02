### Variables
.DEFAULT_GOAL := help
PORT=8091
CPU=3
MEMORY=3096
XTRA_BOOT=""

MINI_BKP_FILE=~/Downloads/mini-bkp.txt
MINI_CURRENT_BKP_FILE=/tmp/mini-bkp
#Current Directory where Make File is Present
MAKE_FILE_DIR=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
SERVICE_FILES=$(MAKE_FILE_DIR)/services/files

### Basic
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

stop: ## Stop Minikube
	printf "\033[1;32m Deleting Minikube Clusters \033[0m \n"
	minikube -p minikube delete;

### Lifecycle
#Additional Flags: --kubernetes-version v1.23.0
#BUG: Metric Server Not Working K9s Pulse.
start: ## Start Minikube
	minikube  -p minikube start --memory=$(MEMORY) --cpus=$(CPU) \
	--cache-images=true --mount-string="$(SERVICE_FILES):/etc/files" --mount --host-only-cidr='24.1.1.100/24'\
	--addons metrics-server;

wait:
	printf "\033[1;34m Waiting for Minikube to be Ready \033[0m \n";
	sleep 30
	kubectl wait --for=condition=Ready pod -l k8s-app=kube-dns -n kube-system --timeout=60s

info: ## Minikube Info
	printf "\033[1;32m Minikube Setup \033[0m \n";
	printf "\033[1;33m Run 'minikube tunnel' for Emulating ELB\033[0m \n";
	printf "\033[1;33m Dashboard: http://localhost:$(PORT)/api/v1/namespaces/kubernetes-dashboard/services/http:kubernetes-dashboard:/proxy/# \033[0m \n";
	printf "\033[1;33m Swagger: http://localhost:$(PORT)/swagger-ui \033[0m \n";
	printf "\033[1;33m K9S:  k9s --context minikube \033[0m \n";
	printf "\033[1;33m Context: `kubectl config current-context;`\033[0m \n";
	printf "\033[1;34m Mounted: $(SERVICE_FILES) \033[0m \n"

### Addons
dashboard: ## Minikube Dashboard
	printf "\033[1;32m Minikube Dashboard & Addons \033[0m \n";
	minikube -p minikube dashboard --port=$(PORT) --url=true > /dev/null 2>&1 &

istio:
	XTRA_BOOT="ISTIO"

bootstrap: ## Bootstrap basic Services
	cd $(MAKE_FILE_DIR)/services;\
	XTRA_BOOT=$(XTRA_BOOT) ./service.zsh -b;\

### Port Forwarding
port-fwd:
	printf "\033[1;34m Please enter password for Port 80 Forward \033[0m \n";
	sudo -v
	screen -S mini -dmL bash -c 'sudo kubectl port-forward deployment/traefik 80:8000'

port-ls:
	screen -ls

port-kill:
	screen -XS mini quit

port: port-kill port-fwd port-ls ## Port Forward for Minikube

proxy:
	kubectl proxy --port=$(PORT);

### Backup and Restore



### Workflow

reset: stop start wait bootstrap ## Reset Minikube


## k9s
# k9s --readonly , -n <namespace>, -l <loglevel>
# k9s :pu (pulse), :dp (deployments), :po (pods), :ns (namespace),:rb (Role Bindings) ,:a (aliases), :hpa (autoscaler)
# Switch Namespace :po <namespace> to see pods of that namespace
# Portforward: Select Pod, Shift+f (Create PF), f (Show PF)

## kompose: brew install kompose
# Convert:  kompose convert -f jira.yml
# Apply: kubectl apply $(ls jira*.yaml | awk ' { print " -f " $NF } ')

## Docker
# Docker Start: sudo service docker start
# Docker Image to Minikube: eval $(minikube docker-env); docker build -t fun-app .

## Kubectl
# Port Forward (Local Port 9090 to Container Port 8080) - kubectl port-forward `kubectl get pods -o name | grep fun-app | head  -1` 9090:8080
# Logs - kubectl logs `kubectl get pods -o name | grep fun-app | head  -1` --since=1m -f
# Login - kubectl -it exec `kubectl get pods -o name | grep fun-app | head  -1` bash
# Delete All - kubectl delete all --all

## Minkikube
# Tunnel (Emulate Load Balancer) - minikube tunnel
# List Emulated Services URL's - minikube service list

# https://minikube.sigs.k8s.io/docs/handbook/pushing/
# Load Image from Host - minikube image load amanfdk/controller
# Minikube Docker ENV Connect - eval $(minikube docker-env)

## Minikube Images (Cache: ~/.minikube/cache)
# Image List - minikube image ls
# Cache Save - minikube image save docker.io/bitnami/mysql:8.0.32-debian-11-r0 --daemon
# Cache Load - minikube image load docker.io/bitnami/mysql:8.0.32-debian-11-r0 --daemon