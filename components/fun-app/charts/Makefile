include ../../../common/tools/base.mk

### Variables

### Setup
setup-vegeta:
	printf $(_TITLE) "Running Vegeta Attack"
	kubectl run vegeta --image="peterevans/vegeta" -- sh -c "sleep 10000" > $(OUT)

setup-fun:
	printf $(_TITLE) "Installing Fun App Helm Chart"
	helm install fun-app . --set rateLimit.perMin=150 > $(OUT)

setup-analyse:
	printf $(_TITLE) "Analyse Log: GoAccess"
	kubectl -n fun-app logs `kubectl get pods -n fun-app -o name | grep app | head  -1` --since=1m -f | goaccess --log-format='%^ %d - %t | %s | %~%D | %b | %~%h | %^ | %m %U' --date-format='%Y/%m/%d' --time-format '%H:%M:%S'

#### Helm
add:
	helm repo add bitnami https://charts.bitnami.com/bitnami > $(OUT)

update:
	printf $(_TITLE) "Updating Helm Dependencies"
	helm dependency update > $(OUT)
	helm dependency build > $(OUT)

# BUG: #C Introduce Helm Version.
package: update ## Package Helm Chart
	printf $(_TITLE) "Packaging Fun App Helm Chart"
	helm package . -d . > $(OUT)

#### Info
info-fun:
	printf $(_INFO) "Fun App Info"
	printf $(_DETAIL) "minikube tunnel"
	printf $(_DETAIL) "kubectl get svc app"
	printf $(_DETAIL) "http://<external_ip>:9001/metrics"
	printf $(_DETAIL) "http://app.docker/app/metrics"

info-fun-notes:
	printf $(_INFO) "Fun App Notes"
	helm get notes fun-app

info-vegeta:
	printf $(_INFO) "Vegeta Info"
	printf $(_DETAIL) "Login Vegeta Pod"
	printf $(_DETAIL) "HELM" "echo 'GET http://fun-app:9090/v1/person/all' | vegeta attack | vegeta report"
	printf $(_DETAIL) "DEVSPACE" "echo 'GET http://app:8080/v1/person/all' | vegeta attack | vegeta report"
	printf $(_DETAIL) "LOAD(MAKE)" "Portforward 8080: http://localhost:8080/metrics and fun-app/it Make"
	printf $(_DETAIL) "PROMETHEUS" "http://prometheus.docker/graph?g0.expr=funapp_create_person&g0.tab=0&g0.display_mode=stacked&g0.show_exemplars=0&g0.range_input=15m"

info-istio:
	printf $(_INFO) "Istio Info"
	printf $(_DETAIL) "Grafana Import: /fun-app/it/dashboard.json"
	printf $(_DETAIL) "Prometheus: : http://localhost:9090/graph?g0.expr=rate(fun_app_person_count%5B5m%5D)&g0.tab=0&g0.stacked=0&g0.show_exemplars=0&g0.range_input=5m&g1.expr=fun_app_person_create_time_bucket&g1.tab=0&g1.stacked=1&g1.show_exemplars=1&g1.range_input=1h&g2.expr=rate(fun_app_person_create_time_count%5B5m%5D)&g2.tab=0&g2.stacked=0&g2.show_exemplars=0&g2.range_input=1h"

### Clean
clean-fun:
	printf $(_WARN) "Deleting Fun App Helm Chart"
	-helm delete fun-app > $(OUT)

clean-vegeta:
	printf $(_WARN) "Deleting Vegeta Attack"
	-kubectl delete pod vegeta > $(OUT)

clean-package:
	printf $(_WARN) "Deleting Helm Package(s)"
	-rm *.tgz > $(OUT)

### Workflows
info: info-fun info-vegeta ## Helm Installation Info
infos: info info-istio info-fun-notes ## Helm Extended Info
prepare: add ## Onetime Setup
setup: setup-fun setup-vegeta ## Install Helm (Fun & Vegeta)
clean: clean-fun clean-package clean-vegeta ## Clean Helm
analyse: setup-analyse ## Analyse Log
reset: clean setup info ## Clean & Setup Helm
all: prepare reset package clean ## Run All Targets

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
