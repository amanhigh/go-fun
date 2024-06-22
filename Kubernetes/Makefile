### Variables
.DEFAULT_GOAL := help
OUT := /dev/null

PORT=8091
CPU=4
MEMORY=4096
XTRA=""

MINI_BKP_DIR=~/Downloads/minibkp

#Current Directory where Make File is Present
MAKE_FILE_DIR=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
SERVICE_FILES=$(MAKE_FILE_DIR)/services/files

### Basic
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	printf $(_INFO) "FirstTime: prepare/all, OUT=/dev/stdout (Debug)"
	printf $(_INFO) "Resources: CPU=4 MEMORY=10000"

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
info-mini:
	printf $(_TITLE) "Minikube Info" ""
	printf $(_INFO) "Context" "`kubectl config current-context;`"
	printf $(_DETAIL) "Mounted" "$(SERVICE_FILES)"
	printf $(_INFO) "Traefik" "http://docker:9000/dashboard/#/"

info-extra:
	printf $(_TITLE) "Minikube Extra Info" ""
	printf $(_INFO) "Emulating ELB" "Run 'minikube tunnel'"
	printf $(_INFO) "Dashboard" "http://localhost:$(PORT)/api/v1/namespaces/kubernetes-dashboard/services/http:kubernetes-dashboard:/proxy/#"
	printf $(_INFO) "Swagger" "http://localhost:$(PORT)/swagger-ui"
	printf $(_INFO) "K9S" "k9s --context minikube"

### Addons
dashboard: ## Minikube Dashboard
	printf $(_TITLE) "Minikube Dashboard & Addons" ""
	minikube -p minikube dashboard --port=$(PORT) --url=true > /dev/null 2>&1 &

bootstrap:
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
	-screen -XS mini quit
	-screen -XS traefik quit

proxy:
	kubectl proxy --port=$(PORT);

port: port-kill port-fwd port-ls ## Port Forward for Minikube

### Backup and Restore
backup:
	printf $(_TITLE) "Backup Minikube Images" "$(MINI_BKP_DIR)"
	
	# Create the backup directory if it doesn't exist
	mkdir -p $(MINI_BKP_DIR)
	minikube image ls | grep -v none > $(MINI_BKP_DIR)/current_image_list.txt

	printf $(_INFO) "Minikube Saving Images" "`wc -l $(MINrI_BKP_DIR)/current_image_list.txt | awk '{print $$1}'`"
	# Save images as tar files in the backup directory only if they do not already exist
	for IMG in $$(cat $(MINI_BKP_DIR)/current_image_list.txt); do \
		FILENAME=$$(echo $$IMG | sed 's/\//+/g').tar; \
		if [ ! -f $(MINI_BKP_DIR)/$$FILENAME ]; then \
			minikube image save $$IMG $(MINI_BKP_DIR)/$$FILENAME; \
			printf $(_DETAIL) "Saved:  $$IMG" "$$FILENAME"; \
		else \
			printf $(_WARN) "Skipped: $$IMG", "$$FILENAME"; \
		fi \
	done

restore:
	printf $(_TITLE) "Restore Minikube Images" "$(MINI_BKP_DIR)"

	# Generate a list of currently loaded images in Minikube
	minikube image ls | grep -v none > $(MINI_BKP_DIR)/loaded_image_list.txt

	# Load images from tar files in the backup directory only if they are not already loaded
	for FILE in $(MINI_BKP_DIR)/*.tar; do \
		IMG=$$(basename $$FILE .tar | sed 's/+/\//g'); \
		if grep -q -e "^$$IMG$$" $(MINI_BKP_DIR)/loaded_image_list.txt; then \
			printf $(_WARN) "Skipped: $$IMG" "$$FILE"; \
		else \
			minikube image load $$FILE; \
			printf $(_DETAIL) "Loaded: $$IMG" "$$FILE"; \
		fi \
	done


### Workflows
info: info-mini ## Minikube Info
infos: info info-extra ## Minikube Extended Info
setup: start wait bootstrap port-fwd ## Start, Bootstrap, Port Forward Minikube
clean: stop ## Stop Minikube
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