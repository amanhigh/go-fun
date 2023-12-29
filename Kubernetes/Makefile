### Variables
.DEFAULT_GOAL := help
OUT := /dev/null

PORT=8091
CPU=3
MEMORY=4096
XTRA=""

MINI_BKP_FILE=~/Downloads/mini-bkp.txt
MINI_CURRENT_BKP_FILE=/tmp/mini-bkp

#Current Directory where Make File is Present
MAKE_FILE_DIR=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
SERVICE_FILES=$(MAKE_FILE_DIR)/services/files

### Basic
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	printf $(_INFO) "FirstTime: prepare/all, OUT=/dev/stdout (Debug)"

### Lifecycle
#Additional Flags: --kubernetes-version v1.23.0
prepare:
	printf $(_TITLE) "Preparing" "Minikube"
	mkdir -p /tmp/mini
	chmod 777 /tmp/mini
	$(MAKE) -C ./services prepare > $(OUT)
	
start: ## Start Minikube
	printf $(_TITLE) "Starting" "Minikube: $(CPU) CPUs, $(MEMORY) MB Memory"
	minikube  -p minikube start --memory=$(MEMORY) --cpus=$(CPU) --addons=metrics-server \
	--mount-string="$(SERVICE_FILES):/etc/files" --mount --cache-images=true \
	--host-only-cidr='24.1.1.100/24' > $(OUT);

	screen -S mini-mount -dmL bash -c 'minikube mount /tmp/mini:/tmp/mini' -L ""

wait:
	printf $(_INFO) "Waiting" "Minikube to be Ready"
	sleep 30
	kubectl wait --for=condition=Ready pod -l k8s-app=kube-dns -n kube-system --timeout=60s > $(OUT)

stop: ## Stop Minikube
	printf $(_WARN) "Deleting" "Minikube Clusters"
	minikube -p minikube delete > $(OUT)

## Info
info-mini: ## Minikube Info
	printf $(_TITLE) "Minikube Info" ""
	printf $(_INFO) "Context" "`kubectl config current-context;`"
	printf $(_DETAIL) "Mounted" "$(SERVICE_FILES)"
	printf $(_INFO) "Dashboard" "http://docker:9000/dashboard/#/"

info-extra: ## Minikube Extra Info
	printf $(_TITLE) "Minikube Extra Info" ""
	printf $(_INFO) "Emulating ELB" "Run 'minikube tunnel'"
	printf $(_INFO) "Dashboard" "http://localhost:$(PORT)/api/v1/namespaces/kubernetes-dashboard/services/http:kubernetes-dashboard:/proxy/#"
	printf $(_INFO) "Swagger" "http://localhost:$(PORT)/swagger-ui"
	printf $(_INFO) "K9S" "k9s --context minikube"

### Addons
dashboard: ## Minikube Dashboard
	printf $(_TITLE) "Minikube Dashboard & Addons" ""
	minikube -p minikube dashboard --port=$(PORT) --url=true > /dev/null 2>&1 &

bootstrap: ## Bootstrap basic Services
	printf $(_TITLE) "Bootstraping Services" ""
	$(MAKE) -C ./services bootstrap > $(OUT)

.PHONY: istio
istio: ## Reset Minikube with Istio
	$(MAKE) XTRA=istio reset

### Port Forwarding
port-fwd:
	printf $(_INFO) "Please enter password for Port 80 Forward" ""
	sudo -v
	screen -S mini -dmL bash -c 'sudo kubectl port-forward deployment/traefik 80:8000'
	screen -S traefik -dmL bash -c 'kubectl port-forward deployment/traefik 9000:9000'

port-ls:
	printf $(_TITLE) "Listing Screens" ""
	screen -ls

port-kill:
	printf $(_WARN) "Killing Port Forward" ""
	screen -XS mini quit
	screen -XS traefik quit

proxy:
	kubectl proxy --port=$(PORT);

port: port-kill port-fwd port-ls ## Port Forward for Minikube

### Backup and Restore
backup:
	# FIXME: Backup Flow
	test -f $(MINI_BKP_FILE) || touch $(MINI_BKP_FILE)
	minikube image ls | grep -v none | tee $(MINI_CURRENT_BKP_FILE)
	printf $(_INFO) "MinkubeImage Count:" "`wc -l $(MINI_CURRENT_BKP_FILE)`"

	# Append Image list to Master List
	cp $(MINI_BKP_FILE) /tmp/mini-bkp-old
	sort $(MINI_CURRENT_BKP_FILE) /tmp/mini-bkp-old | uniq | tee $(MINI_BKP_FILE)
	printf $(_INFO) "MasterList Count:" "`wc -l $(MINI_BKP_FILE)`"

	for IMG in $$(cat $(MINI_CURRENT_BKP_FILE)); do\
		CACHE_PATH="$${IMG%/*}";\
		IMAGE_ID="$${IMG##*/}";\
		IMAGE_CACHE_PATH="$${HOME}/.minikube/cache/images/$$(uname -m)/$${CACHE_PATH}";\
		IMAGE_CACHE_FILE="$${IMAGE_CACHE_PATH}/$$(echo $IMAGE_ID | sed 's/[:]/_/g')";\
		echo $${IMAGE_CACHE_FILE};\

		# if [ -f ${IMAGE_CACHE_FILE} ]; then
		# 	printf "\033[1;34m Skipping IMAGE: ${IMG} -> ${IMAGE_CACHE_FILE}\033[0m \n"
		# else
		#	printf "\033[1;33m Caching IMAGE: ${IMG} \033[0m \n"
		# 	mkdir -p ${IMAGE_CACHE_PATH}
		# 	minikube image save --daemon ${IMG}
		# fi
	done
restore:
	for IMG in $$(cat $(MINI_BKP_FILE)); do \
		printf $(_INFO) "Loading Image:" "$${IMG}"; \
		minikube image load --daemon $${IMG}; \
	done

### Workflows
info: info-mini ## Minikube Info
infos: info info-extra ## Minikube Extended Info
setup: start wait bootstrap port-fwd ## Setup Minikube
clean: stop ## Stop Minikube
# HACK:#B Fixe Readme
reset: clean setup info ## Rebuild Minikube
all: prepare istio dashboard infos port ## Run All Targets

### Formatting
_INFO := "\033[33m[%s]\033[0m %s\n"  # Yellow text for "printf"
_DETAIL := "\033[34m[%s]\033[0m %s\n"  # Blue text for "printf"
_TITLE := "\033[32m[%s]\033[0m %s\n" # Green text for "printf"
_WARN := "\033[31m[%s]\033[0m %s\n" # Red text for "printf"

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