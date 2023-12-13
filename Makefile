### Help
# Tutorial: https://makefiletutorial.com/
# Silent: -s
# Paraller Jobs: -j2
# Override Vars: make test-it COVER_DIR=./test
# Call Target: $(MAKE) --no-print-directory XTRA=ISTIO bootstrap
# Store Var: CUR_DIR := $(shell pwd) (Outside Target)
# Dynamic Var: $(eval RESTORE_DB_NAME := $(DBNAME)_restore)
# Continue Step or error: Start with `-`. Eg. -rm test.txt
# Make In Directory: make -C /path/to/dir

### Variables
.DEFAULT_GOAL := help
BUILD_OPTS := CGO_ENABLED=1 GOOS=linux GOARCH=amd64
COMPONENT_DIR := ./components
FUN_DIR := $(COMPONENT_DIR)/fun-app

COVER_DIR:= /tmp/cover
PROFILE_FILE:= $(COVER_DIR)/profile.out

FUN_IMAGE_TAG := amanfdk/fun-app

.PHONY: sync test

### Basic
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

sync: ## Sync Go Modules
	go work sync

### Testing
# https://golangci-lint.run/usage/quick-start/
# FIXME: Use Configuration - https://golangci-lint.run/usage/configuration/
lint: ## Lint the Code
	go work edit -json | jq -r '.Use[].DiskPath'  | xargs -I{} golangci-lint run {}/...

test-operator: ## Run operator tests
	make -C $(COMPONENT_DIR)/operator/ test

test-unit: ## Run unit tests
	ginkgo -r '--label-filter=!setup' -cover .

cover-analyse: ## Analyse Integration Coverage Reports
	@echo "Generating FunServer Cover Profile"
	# Generate Cover Profile
	go tool covdata textfmt -i=$(COVER_DIR) -o $(PROFILE_FILE)
	
	# Analyse Cover Profile
	go tool cover -func=$(PROFILE_FILE)

	@printf "\033[1;32m Package Summary \033[0m"
	# Analyse Report and Print Coverage
	go tool covdata percent -i=$(COVER_DIR)

	printf "\033[1;32m\n\n ******* Vscode: go.apply.coverprofile $(PROFILE_FILE) ******** \033[0m"


test-it: run-fun-cover test-unit cover-analyse ## Integration test coverage analyse

test-clean:
	@echo "Cleaning Coverage Reports"
	rm -rf $(COVER_DIR)

profile:
	go tool pprof -http=:8001 http://localhost:8080/debug/pprof/heap &\
	go tool pprof -http=:8000 --seconds=30 http://localhost:8080/debug/pprof/profile;\
	kill %1;

test: test-operator test-it ## Run all tests

### Builds
swag-fun: ## Swagger Generate: Fun App (Init/Update)
	cd $(FUN_DIR);\
	swag i --parseDependency true;\
	echo http://localhost:8080/swagger/index.html;

build-fun: swag-fun ## Build Fun App
	$(BUILD_OPTS) go build -o $(FUN_DIR)/fun $(FUN_DIR)/main.go

build-fun-cover: ## Build Fun App with Coverage
	$(BUILD_OPTS) go build -cover -o $(FUN_DIR)/fun $(FUN_DIR)/main.go

build-kohan:
	$(BUILD_OPTS) go build -o $(COMPONENT_DIR)/kohan/kohan $(COMPONENT_DIR)/kohan/main.go

build-clean:
	rm "$(FUN_DIR)/fun";
	rm "$(COMPONENT_DIR)/kohan/kohan";

build: build-fun build-kohan ## Build all Binaries

### Runs
run-fun: build-fun ## Run Fun App
	$(FUN_DIR)/fun

# Guide - https://dustinspecker.com/posts/go-combined-unit-integration-code-coverage/
run-fun-cover: build-fun-cover ## Run Fun App with Coverage
	mkdir -p $(COVER_DIR)
	GOCOVERDIR=$(COVER_DIR) PORT=8085 $(FUN_DIR)/fun > $(FUN_DIR)/funcover.log &

### Helm
helm-build: ## Build Helm Charts
	helm dependency build $(FUN_DIR)/charts/

helm-package: helm-build ## Package Helm Charts
	helm package $(FUN_DIR)/charts/ -d $(FUN_DIR)/charts

### Local Setup
setup-tools: ## Setup Tools	for Local Environment
	go install github.com/onsi/ginkgo/v2/ginkgo
	go install github.com/swaggo/swag/cmd/swag

setup-k8: ## Kubernetes Setup
	$(MAKE) -C ./Kubernetes/services helm hosts

#HACK: Add Make to Readme
setup: setup-tools setup-k8 ## Setup Local Environment

### Docker
docker-fun: build-fun
	docker build -t $(FUN_IMAGE_TAG) -f $(FUN_DIR)/Dockerfile $(FUN_DIR)

docker-build: docker-fun ## Build Docker Images

### Workflows
all: sync test build helm-package docker-build ## Run Complete Build Process
	@printf "\033[1;32m\n\n ******* Complete BUILD Successful ********\n \033[0m"

clean: test-clean build-clean ## Clean up Residue