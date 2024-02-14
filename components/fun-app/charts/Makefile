### Variables
.DEFAULT_GOAL := help
OUT := /dev/null

### Basic
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	printf $(_TITLE) "FirstTime: prepare/all, OUT=/dev/stdout (Debug)"

### Setup
setup-vegeta:
	printf $(_TITLE) "Running Vegeta Attack"
	kubectl run vegeta --image="peterevans/vegeta" -- sh -c "sleep 10000"

setup-fun:
	printf $(_TITLE) "Installing Fun App Helm Chart"
	helm install fun-app . --set rateLimit.perMin=150

#### Helm
helm-add:
	helm repo add bitnami https://charts.bitnami.com/bitnami


helm-build:
	printf $(_TITLE) "Updating Helm Dependencies"
	helm dependency update
	helm dependency build

package: helm-build ## Package Fun Helm Chart
	printf $(_TITLE) "Packaging Fun App Helm Chart"
	helm package . -d .

#### Info
# BUG: Fix Fun App via Tunnel
info-fun:
	printf $(_INFO) "Fun App Info"
	printf $(_DETAIL) "minikube tunnel"
	printf $(_DETAIL) "http://localhost:9090/metrics"

info-fun-notes:
	printf $(_INFO) "Fun App Notes"
	helm get notes fun-app

info-vegeta:
	printf $(_INFO) "Vegeta Info"
	printf $(_DETAIL) "Login Vegeta Pod"
	printf $(_DETAIL) "HELM: echo 'GET http://fun-app:9090/person/all' | vegeta attack | vegeta report"
	printf $(_DETAIL) "DEVSPACE: echo 'GET http://app:8080/person/all' | vegeta attack | vegeta report"

info-istio:
	printf $(_INFO) "Istio Info"
	printf $(_DETAIL) "Grafana Import: /fun-app/it/dashboard.json"
	printf $(_DETAIL) "Prometheus: : http://localhost:9090/graph?g0.expr=rate(fun_app_person_count%5B5m%5D)&g0.tab=0&g0.stacked=0&g0.show_exemplars=0&g0.range_input=5m&g1.expr=fun_app_person_create_time_bucket&g1.tab=0&g1.stacked=1&g1.show_exemplars=1&g1.range_input=1h&g2.expr=rate(fun_app_person_create_time_count%5B5m%5D)&g2.tab=0&g2.stacked=0&g2.show_exemplars=0&g2.range_input=1h"

### Clean
clean-fun:
	printf $(_TITLE) "Deleting Fun App Helm Chart"
	-helm delete fun-app

clean-vegeta:
	printf $(_TITLE) "Deleting Vegeta Attack"
	-kubectl delete pod vegeta

### Workflows
info: info-fun info-vegeta ## Info on Fun Helm Installation
infos: info info-istio info-fun-notes ## Extended Info
prepare: helm-add ## Onetime Setup
setup: setup-fun setup-vegeta ## Setup Fun via Helm
clean: clean-fun ## Clean Fun Helm
reset: clean setup info ## Reset Fun Helm
all: prepare reset package ## Run All Targets

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
