### Help
# Silent: -s

### Variables
.DEFAULT_GOAL := help

.PHONY: sync test

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

sync: ## Sync Go Modules
	go work sync

### Testing
test-operator: ## Run operator tests
	make -C ./components/operator/ test

test-integration: ## Run integration tests
	sh -c './components/fun-app/it/cover.zsh run > /dev/null 2>&1 &'

test-unit: ## Run unit tests
	ginkgo -r '--label-filter=!setup' -cover .

test: test-operator test-unit ## Run all tests

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
	helm dependency build ./components/fun-app/charts/

helm-package: helm-build ## Package Helm Charts
	helm package ./components/fun-app/charts/ -d ./components/fun-app/charts

### Builds
build-fun:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o ./components/fun-app/fun ./components/fun-app/main.go

build-kohan:
	go build -o ./components/kohan/kohan ./components/kohan/main.go

# Docker
docker-fun:
	docker build -t amanfdk/fun-app -f ./components/fun-app/Dockerfile .

docker-build: docker-fun ## Build Docker Images