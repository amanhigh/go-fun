### Variables
.DEFAULT_GOAL := help
OUT := /dev/null

### Basic
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	printf $(_TITLE) "FirstTime: prepare/all, OUT=/dev/stdout (Debug)"

#### Helm
helm-build:
	helm dependency update
	helm dependency build

helm-fun:
	printf $(_TITLE) "Installing Fun App Helm Chart"
	helm install fun-app . --set rateLimit.perMin=150

#### Info
info-fun:
	printf $(_INFO) "Helm Fun App Info"

### Clean
clean-fun:
	-helm delete fun-app

### Workflows
info: info-fun ## Info on Fun Helm Installation
infos: info ## Extended Info
prepare: ## Onetime Setup
setup: helm-build helm-fun ## Setup Fun via Helm
clean: clean-fun ## Clean Fun Helm
reset: clean setup info ## Reset Fun Helm
all:prepare reset ## Run All Targets

### Formatting
_INFO := "\033[33m[%s]\033[0m %s\n"  # Yellow text for "printf"
_DETAIL := "\033[34m[%s]\033[0m %s\n"  # Blue text for "printf"
_TITLE := "\033[32m[%s]\033[0m %s\n" # Green text for "printf"
_WARN := "\033[31m[%s]\033[0m %s\n" # Red text for "printf"

### Helpful Commands
# helm init fun-app - Bootstrap Charts
# helm template . - Preview Charts with Values
# helm lint . - Check Errors
# helm show values <Chart Name> - Show configurable values

# helm install -n <Namespace> <Chart Name> . [--set <key>=<value>]
# helm upgrade -n <Namespace> <Chart Name> . [--set <key>=<value>]

# helm status -n <Namespace> <Chart Name>
# helm history -n <Namespace> <Chart Name>
# helm rollback -n <Namespace> <Chart Name> [Revision]
# helm delete -n <Namespace> <Chart Name>


# helm list -n <Namespace>

# helm repo list
# helm dependency list
# helm dependency update
# helm dependency build
