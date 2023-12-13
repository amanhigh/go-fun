### Help
# Tutorial: https://makefiletutorial.com/
# Silent: -s, Keepgoing -k, 
# Paraller Jobs: -j2
# Override Vars: make test-it COVER_DIR=./test
# Call Target: $(MAKE) --no-print-directory XTRA=ISTIO bootstrap
# Store Var: CUR_DIR := $(shell pwd) (Outside Target)
# Dynamic Var: $(eval RESTORE_DB_NAME := $(DBNAME)_restore)
# Continue Step or error: Start with `-`. Eg. -rm test.txt
# Make In Directory: make -C /path/to/dir

### Variables
.DEFAULT_GOAL := help
BUILD_OPTS := CGO_ENABLED=1 GOARCH=amd64
COMPONENT_DIR := ./components
FUN_DIR := $(COMPONENT_DIR)/fun-app

COVER_DIR:= /tmp/cover
PROFILE_FILE:= $(COVER_DIR)/profile.out

FUN_IMAGE_TAG := amanfdk/fun-app
OUT := /dev/null

.PHONY: sync test

### Basic
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@printf "\033[1;32m FirstTime: prepare/all, OUT=/dev/stdout (Debug) \033[0m \n"

sync: ## Sync Go Modules
	go work sync

### Testing
# https://golangci-lint.run/usage/quick-start/
# FIXME: Use Configuration - https://golangci-lint.run/usage/configuration/
lint: ## Lint the Code
	printf "\033[1;32m Running Linting \n\033[0m"
	go work edit -json | jq -r '.Use[].DiskPath'  | xargs -I{} golangci-lint run {}/...

test-operator: ## Run operator tests
	printf "\033[1;32m Running Operator Tests \n\033[0m"
	make -C $(COMPONENT_DIR)/operator/ test > $(OUT)

test-unit: ## Run unit tests
	printf "\033[1;32m Running Unit Tests \n\033[0m"
	ginkgo -r '--label-filter=!setup' -cover . > $(OUT)

cover-analyse: ## Analyse Integration Coverage Reports
	printf "\033[1;32m Analysing Coverage Reports \n\033[0m"
	# Generate Cover Profile
	go tool covdata textfmt -i=$(COVER_DIR) -o $(PROFILE_FILE)
	
	# Analyse Cover Profile
	go tool cover -func=$(PROFILE_FILE) > $(OUT)

	printf "\033[1;32m Package Summary \n\033[0m"
	# Analyse Report and Print Coverage
	go tool covdata percent -i=$(COVER_DIR)

	printf "\033[1;32m\n\n ******* Vscode: go.apply.coverprofile $(PROFILE_FILE) ******** \033[0m"


test-it: run-fun-cover test-unit cover-analyse ## Integration test coverage analyse

test-clean:
	printf "\033[1;31m Cleaning Tests \n\033[0m"
	rm -rf $(COVER_DIR)

profile: ## Run Profiling
	printf "\033[1;32m Running Profiling \n\033[0m"
	go tool pprof -http=:8001 http://localhost:8080/debug/pprof/heap &\
	go tool pprof -http=:8000 --seconds=30 http://localhost:8080/debug/pprof/profile;\
	kill %1;

### Builds
swag-fun: ## Swagger Generate: Fun App (Init/Update)
	printf "\033[1;32m Generating Swagger \n\033[0m"
	cd $(FUN_DIR);\
	swag i --parseDependency true > $(OUT);\
	printf "\033[1;33m http://localhost:8080/swagger/index.html \n\033[0m";

build-fun: swag-fun ## Build Fun App
	printf "\033[1;32m Building Fun App \n\033[0m"
	$(BUILD_OPTS) go build -o $(FUN_DIR)/fun $(FUN_DIR)/main.go

build-fun-cover: ## Build Fun App with Coverage
	printf "\033[1;32m Building Fun App with Coverage \n\033[0m"
	$(BUILD_OPTS) go build -cover -o $(FUN_DIR)/fun $(FUN_DIR)/main.go

build-kohan:
	printf "\033[1;32m Building Kohan \n\033[0m"
	$(BUILD_OPTS) go build -o $(COMPONENT_DIR)/kohan/kohan $(COMPONENT_DIR)/kohan/main.go

build-clean:
	printf "\033[1;31m Cleaning Build \n\033[0m"
	rm "$(FUN_DIR)/fun";
	rm "$(COMPONENT_DIR)/kohan/kohan";

### Runs
run-fun: build-fun ## Run Fun App
	printf "\033[1;32m Running Fun App \n\033[0m"
	@$(FUN_DIR)/fun > $(OUT)

# Guide - https://dustinspecker.com/posts/go-combined-unit-integration-code-coverage/
run-fun-cover: build-fun-cover ## Run Fun App with Coverage
	printf "\033[1;32m Running Fun App with Coverage \n\033[0m"
	mkdir -p $(COVER_DIR)
	GOCOVERDIR=$(COVER_DIR) PORT=8085 $(FUN_DIR)/fun > $(OUT) 2>&1 &

### Helm
helm-build: ## Build Helm Charts
	printf "\033[1;32m Building Helm Charts \n\033[0m"
	helm dependency build $(FUN_DIR)/charts/ > $(OUT);

helm-package: helm-build ## Package Helm Charts
	printf "\033[1;32m Packaging Helm Charts \n\033[0m"
	helm package $(FUN_DIR)/charts/ -d $(FUN_DIR)/charts

### Local Setup
setup-tools: ## Setup Tools	for Local Environment
	printf "\033[1;32m Setting up Tools \n\033[0m"
	go install github.com/onsi/ginkgo/v2/ginkgo
	go install github.com/swaggo/swag/cmd/swag

setup-k8: ## Kubernetes Setup
	printf "\033[1;32m Setting up Kubernetes \n\033[0m"
	$(MAKE) -C ./Kubernetes/services helm hosts

### Docker
docker-fun: build-fun
	printf "\033[1;32m Building FunApp Docker Image \033[0m"
	docker build -t $(FUN_IMAGE_TAG) -f $(FUN_DIR)/Dockerfile $(FUN_DIR) > $(OUT)

docker-build: docker-fun ## Build Docker Images

### Workflows
test: test-operator test-it ## Run all tests
build: build-fun build-kohan ## Build all Binaries

#HACK: Add Make to Readme
info:
prepare: setup-tools setup-k8 # Setup Tools

setup: sync test build helm-package docker-build # Build and Test
clean: test-clean build-clean ## Clean up Residue

reset: setup info ## Build and Show Info
all: prepare reset clean ## Run All Targets
	printf "\033[1;32m\n\n ******* Complete BUILD Successful ********\n \033[0m"
