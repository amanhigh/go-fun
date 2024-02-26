### Help
# Tutorial: https://makefiletutorial.com/
# Silent: -s, Keepgoing -k, 
# Paraller Jobs: -j2
### Calls
# Override Vars: make test-it COVER_DIR=./test
# Call Target: $(MAKE) --no-print-directory XTRA=ISTIO bootstrap
# Make In Directory: make -C /path/to/dir
# Continue Step or error: Start with `-`. Eg. -rm test.txt
### Variables
# SHELL Var in Make: CUR_DIR := $(shell pwd) (Outside Target)
# Make Var in SHELL: $(eval RESTORE_DB_NAME := $(DBNAME)_restore)

### Variables
.DEFAULT_GOAL := help

BUILD_OPTS := CGO_ENABLED=0 GOARCH=amd64
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
	printf $(_TITLE) "FirstTime: prepare/all, OUT=/dev/stdout (Debug)"

sync:
	printf $(_TITLE) "Go Module Syncing"
	go work sync

### Testing
# https://golangci-lint.run/usage/quick-start/
# FIXME: Use Configuration - https://golangci-lint.run/usage/configuration/
lint: ## Lint the Code
	printf $(_TITLE) "Running Linting"
	go work edit -json | jq -r '.Use[].DiskPath'  | xargs -I{} golangci-lint run {}/...

test-operator:
	printf $(_TITLE) "Running Operator Tests"
	make -C $(COMPONENT_DIR)/operator/ test > $(OUT)

test-unit:
	printf $(_TITLE) "Running Unit Tests"
	ginkgo -r '--label-filter=!setup && !slow' -cover . > $(OUT)

test-slow: ## Run slow tests
	printf $(_TITLE) "Running Slow Tests"
	ginkgo -r '--label-filter=slow' -cover . > $(OUT)

cover-analyse:
	printf $(_TITLE) "Analysing Coverage Reports"
	# Generate Cover Profile
	go tool covdata textfmt -i=$(COVER_DIR) -o $(PROFILE_FILE)
	
	# Analyse Cover Profile
	go tool cover -func=$(PROFILE_FILE) > $(OUT)

	printf $(_TITLE) "Package Summary";
	# Analyse Report and Print Coverage
	go tool covdata percent -i=$(COVER_DIR);
	echo "";
	printf $(_INFO) "Vscode" "go.apply.coverprofile $(PROFILE_FILE)";

test-it: run-fun-cover test-unit cover-analyse

test-clean:
	printf $(_WARN) "Cleaning Tests"
	rm -rf $(COVER_DIR)

profile: ## Run Profiling
	printf $(_TITLE) "Running Profiling on Port 8080"
	printf $(_DETAIL) "Profiling Heap"
	go tool pprof -http=:8001 http://localhost:8080/debug/pprof/heap 2> $(OUT) &\
	printf $(_DETAIL) "Profiling CPU"
	go tool pprof -http=:8000 --seconds=30 http://localhost:8080/debug/pprof/profile 2> $(OUT);\
	printf $(_WARN) "Killing Profilers"
	kill %1;

### Builds
swag-fun:
	printf $(_TITLE) "Generating Swagger"
	cd $(FUN_DIR);\
	swag i --parseDependency true > $(OUT);\
	printf $(_INFO) "Swagger" "http://localhost:8080/swagger/index.html";

build-fun: swag-fun
	printf $(_TITLE) "Building Fun App"
	$(BUILD_OPTS) go build -o $(FUN_DIR)/fun $(FUN_DIR)/main.go

build-fun-cover:
	printf $(_TITLE) "Building Fun App with Coverage"
	$(BUILD_OPTS) go build -cover -o $(FUN_DIR)/fun $(FUN_DIR)/main.go

build-kohan:
	printf $(_TITLE) "Building Kohan"
	$(BUILD_OPTS) go build -o $(COMPONENT_DIR)/kohan/kohan $(COMPONENT_DIR)/kohan/main.go

build-clean:
	printf $(_WARN) "Cleaning Build"
	rm "$(FUN_DIR)/fun";
	rm "$(COMPONENT_DIR)/kohan/kohan";

### Helpers
confirm:
	@if [[ -z "$(CI)" ]]; then \
		REPLY="" ; \
		read -p "⚠ Are you sure? [y/n] > " -r ; \
		if [[ ! $$REPLY =~ ^[Yy]$$ ]]; then \
			printf $(_WARN) "KO" "Stopping" ; \
			exit 1 ; \
		else \
			printf $(_TITLE) "OK" "Continuing" ; \
			exit 0; \
		fi \
	fi

### Release
release-models:
	printf $(_TITLE) "Release Models: $(VER)"
	@if $(MAKE) --no-print-directory confirm ; then \
		git tag models/$(VER) ; \
		git tag | grep models | tail -2 ;
	fi

	printf $(_TITLE) "Pushing Tags";
	@if $(MAKE) --no-print-directory confirm ; then \
		git push --tags && printf $(_TITLE) "Models Released: $(VER)" ; \
	fi

release-common:
	printf $(_TITLE) "Bump Models: $(VER)";
	@if $(MAKE) --no-print-directory confirm ; then \
		pushd ./common ; \
		go get -u github.com/amanhigh/go-fun/models@$(VER); \
		git add go.* && git commit -m "Bumping Models: $(VER)"; \
		popd; \
	fi

	printf $(_TITLE) "Release Common: $(VER)";
	@if $(MAKE) --no-print-directory confirm ; then \
		git tag common/$(VER) ; \
		git tag | grep common | tail -2 ;
	fi

	printf $(_TITLE) "Pushing Tags";
	@if $(MAKE) --no-print-directory confirm ; then \
		git push --tags && printf $(_TITLE) "Common Released: $(VER)" ; \
	fi

release-fun:
	printf $(_TITLE) "Bump Common: $(VER)";
	@if $(MAKE) --no-print-directory confirm ; then \
		pushd ./components/fun-app ; \
		go get -u github.com/amanhigh/go-fun/common@$(VER); \
		git add go.* && git commit -m "Bumping Common: $(VER)"; \
		popd; \
	fi

	printf $(_TITLE) "Release Fun: $(VER)";
	@if $(MAKE) --no-print-directory confirm ; then \
		git tag $(VER) ; \
		$(MAKE) info-release ; \
	fi

	printf $(_TITLE) "Pushing Tags";
	@if $(MAKE) --no-print-directory confirm ; then \
		git push --tags && printf $(_TITLE) "Fun Released: $(VER)" ; \
	fi

release-docker: docker-build ## Release Docker Images
ifndef VER
	$(error VER not set. Eg. v1.1.0)
endif
	printf $(_TITLE) "Release Docker Images: $(VER)"
	printf $(_DETAIL) "Docker Tag"
	docker tag amanfdk/fun-app:latest amanfdk/fun-app:$(VER)

	printf $(_DETAIL) "Docker Push"
	docker push amanfdk/fun-app:latest
	docker push amanfdk/fun-app:$(VER)

unrelease: ## Revoke Release of Golang Packages
ifndef VER
	$(error VER not set. Eg. v1.1.0)
endif
	printf $(_WARN) "Deleting" "Release: $(VER)"
	@if $(MAKE) --no-print-directory confirm ; then \
		git tag -d models/$(VER) ; \
		git push --delete origin models/$(VER); \
		git tag -d common/$(VER) ; \
		git push --delete origin common/$(VER); \
		git tag -d $(VER) ; \
		git push --delete origin $(VER); \
	fi
	$(MAKE) --no-print-directory info-release

release: info-release ## Release Golang Packages
ifndef VER
	$(error VER not set. Eg. v1.1.0)
endif
	$(MAKE) --no-print-directory release-models;
	$(MAKE) --no-print-directory release-common;
	$(MAKE) --no-print-directory release-fun;

### Info
info-release:
	printf $(_INFO) "Go Modules"
	git tag | grep "models" | tail -2
	git tag | grep "common" | tail -2
	git tag | grep "v" | grep -v "/" | tail -2

info-docker:
	printf $(_INFO) "FunApp DockerHub: https://hub.docker.com/r/amanfdk/fun-app/tags"
	curl -s "https://hub.docker.com/v2/repositories/$(FUN_IMAGE_TAG)/tags/?page_size=25&page=1&name&ordering" | jq -r '.results[]|.name' | head -3
	printf $(_INFO) "Docker Images: $(FUN_IMAGE_TAG)"
	docker images | grep fun-app

### Runs
run: build-fun ## Run Fun App
	printf $(_TITLE) "Running: Fun App"
	$(FUN_DIR)/fun > $(OUT)

load: ## Load Test Fun App
	printf $(_TITLE) "Load Test: Fun App"
	make -C $(FUN_DIR)/it all > $(OUT)

# make watch CMD=ls
watch: ## Watch (entr): `make watch CMD=ls`
	printf $(_TITLE) "Watch (entr): $(CMD)"
	find . | entr -s "date +%M:%S; $(CMD)"

# Guide - https://dustinspecker.com/posts/go-combined-unit-integration-code-coverage/
run-fun-cover: build-fun-cover
	printf $(_TITLE) "Running Fun App with Coverage"
	mkdir -p $(COVER_DIR)
	GOCOVERDIR=$(COVER_DIR) PORT=8085 $(FUN_DIR)/fun > $(OUT) 2>&1 &

### Helm
helm-package:
	$(MAKE) -C $(FUN_DIR)/charts package

### Local Setup
setup-tools:
	printf $(_TITLE) "Setting up Tools"
	go install github.com/onsi/ginkgo/v2/ginkgo
	go install github.com/swaggo/swag/cmd/swag

setup-k8: ## Kubernetes Setup
	printf $(_TITLE) "Setting up Kubernetes"
	$(MAKE) -C ./Kubernetes/services helm hosts

### Docker
docker-fun: build-fun
	printf $(_TITLE) "Building FunApp Docker Image"
	docker buildx build -t $(FUN_IMAGE_TAG) -f $(FUN_DIR)/Dockerfile $(FUN_DIR) 2> $(OUT)

docker-fun-run: docker-fun
	printf $(_TITLE) "Running FunApp Docker Image"
	docker run -it amanfdk/fun-app

docker-fun-exec:
	printf $(_TITLE) "Execing Into FunApp Docker Image"
	docker run -it --entrypoint /bin/sh amanfdk/fun-app

docker-fun-clean:
	printf $(_WARN) "Deleting FunApp Docker Image"
	docker rmi -f `docker images $(FUN_IMAGE_TAG)  -q` > $(OUT)

### Devspace
space: space-purge ## Setup Devspace
	printf $(_TITLE) "Starting Devspace"
	devspace use namespace fun-app
	devspace dev

space-purge: ## Purge Devspace
	printf $(_TITLE) "Purging Devspace"
	-devspace purge > $(OUT)

space-info:
	printf $(_TITLE) "Info Devspace"
	devspace list vars --var DB="mysql-primary",RATE_LIMIT=-1
	printf $(_DETAIL) "http://localhost:8080/metrics"
	printf $(_DETAIL) "Login: devspace enter"

space-test: ## Gink Tests Devspace (Watch Mode)
	printf $(_TITLE) "Devspace Tests"
	devspace run ginkgo > $(OUT)
	$(MAKE) watch CMD="devspace run fun-test"

docker-build: docker-fun ## Build Docker Images

### Workflows
test: test-operator test-it ## Run all tests (Excludes test-slow)
build: build-fun build-kohan ## Build all Binaries

info: info-release info-docker ## Repo Information
infos: info space-info ## Repo Extended Information
prepare: setup-tools setup-k8 # One Time Setup

setup: sync test build helm-package docker-build # Build and Test
clean: test-clean build-clean ## Clean up Residue

reset: setup info clean ## Setup with Info and Clean
all: prepare docker-fun-clean reset infos test-slow ## Run All Targets
	printf $(_TITLE) "******* Complete BUILD Successful ********"

### Formatting
_INFO := "\033[33m[%s]\033[0m %s\n"  # Yellow text for "printf"
_TITLE := "\033[32m[%s]\033[0m %s\n" # Green text for "printf"
_WARN := "\033[31m[%s]\033[0m %s\n" # Red text for "printf"
_DETAIL := "\033[34m[%s]\033[0m %s\n"  # Blue text for "printf"