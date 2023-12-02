### Help
# Silent: -s
# Paraller Jobs: -j2

### Variables
.DEFAULT_GOAL := help
BUILD_OPTS := CGO_ENABLED=1 GOOS=linux GOARCH=amd64
COMPONENT_DIR := ./components

FUN_IMAGE_TAG := amanfdk/fun-app

.PHONY: sync test

### Basic
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

sync: ## Sync Go Modules
	go work sync

### Testing
test-operator: ## Run operator tests
	make -C $(COMPONENT_DIR)/operator/ test

test-fun-cover: ## Run Fun Server with Coverage
	$(COMPONENT_DIR)/fun-app/it/cover.zsh run > /dev/null 2>&1 &

test-unit: ## Run unit tests
	ginkgo -r '--label-filter=!setup' -cover .

test-clean:
	$(COMPONENT_DIR)/fun-app/it/cover.zsh clean

test: test-operator test-fun-cover test-unit ## Run all tests

### Builds
build-fun: ## Build Fun App
	$(BUILD_OPTS) go build -o $(COMPONENT_DIR)/fun-app/fun $(COMPONENT_DIR)/fun-app/main.go

build-kohan:
	$(BUILD_OPTS) go build -o $(COMPONENT_DIR)/kohan/kohan $(COMPONENT_DIR)/kohan/main.go

build-clean:
	rm $(COMPONENT_DIR)/fun-app/fun
	rm $(COMPONENT_DIR)/kohan/kohan

build: build-fun build-kohan ## Build all Binaries

### Helm
helm-add: ## Add Helm Repos
	helm repo add onechart https://chart.onechart.dev
	helm repo add stakater https://stakater.github.io/stakater-charts
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo add istio https://istio-release.storage.googleapis.com/charts
	helm repo add kiali https://kiali.org/helm-charts
	helm repo add opa https://open-policy-agent.github.io/kube-mgmt/charts
	helm repo add hashicorp https://helm.releases.hashicorp.com
	helm repo add portainer https://portainer.github.io/k8s/
	helm repo add traefik https://traefik.github.io/charts
	helm repo add hashicorp https://helm.releases.hashicorp.com
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	helm repo add grafana https://grafana.github.io/helm-charts
	helm repo add jaegertracing https://jaegertracing.github.io/helm-charts
	helm repo add bitpoke https://helm-charts.bitpoke.io

helm-update: ## Update Helm Repos
	helm repo update

helm-build: ## Build Helm Charts
	helm dependency build $(COMPONENT_DIR)/fun-app/charts/

helm-package: helm-build ## Package Helm Charts
	helm package $(COMPONENT_DIR)/fun-app/charts/ -d $(COMPONENT_DIR)/fun-app/charts

### Local Setup
setup-hosts: ## Setup Hosts
	DOCKER_HOSTS="127.0.0.1 docker httpbin.docker dashy.docker resty.docker app.docker\
	mysqladmin.docker redisadmin.docker prometheus.docker grafana.docker jaeger.docker kiali.docker\
	ldapadmin.docker webssh.docker webssh2.docker sshwifty.docker nginx.docker portainer.docker\
	consul.docker opa.docker sonar.docker";\
	echo $$DOCKER_HOSTS | sudo tee -a /etc/hosts

setup-tools: ## Setup Tools	for Local Environment
	go install github.com/onsi/ginkgo/v2/ginkgo

#HACK: Add Make to Readme
setup: setup-hosts setup-tools ## Setup Local Environment

### Docker
docker-fun: build-fun
	docker build -t $(FUN_IMAGE_TAG) -f $(COMPONENT_DIR)/fun-app/Dockerfile $(COMPONENT_DIR)/fun-app

docker-build: docker-fun ## Build Docker Images

### Workflows
all: sync test build helm-package docker-build ## Run Complete Build Process

clean: test-clean build-clean ## Clean up Residue