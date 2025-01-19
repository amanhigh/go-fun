include ../../common/tools/base.mk

### Variables
SERVICE_MAKE_OPTS=--no-print-directory -f ./services.mk
BOOTABLE = traefik\|webui
XTRA="webui" # Should not be empty

DOCKER_HOSTS=127.0.0.1 docker httpbin.docker dashy.docker resty.docker webui.docker pdf.docker paperless.docker app.docker \
mysqladmin.docker redisadmin.docker cron.docker prometheus.docker grafana.docker jaeger.docker kiali.docker \
ldapadmin.docker webssh.docker elastic.docker kibana.docker webssh2.docker sshwifty.docker nginx.docker portainer.docker \
consul.docker opa.docker sonar.docker metabase.docker clarity.docker locust.docker
NS:=$(shell kubectl get sa -o=jsonpath='{.items[0]..metadata.namespace}')
ANS_FILE:=/tmp/k8-svc.txt
ANS:=$(shell cat $(ANS_FILE))

### Cleanup
clear-kiali:
	@if [ "$(call istio_installed)" = "true" ]; then \
		printf $(_WARN) "Clearing" "Kiali"; \
		kubectl delete kiali --all --all-namespaces 2> /dev/null || true; \
	else \
		printf $(_INFO) "Skipping cleanup" "Kiali not installed"; \
	fi

clear-istio:
	@if [ "$(call istio_installed)" = "true" ]; then \
		printf $(_WARN) "Clearing" "Istio"; \
		helm delete -n istio-system $$(helm list -n istio-system -q) > $(OUT) 2>&1 || true; \
		kubectl delete namespace istio-system --grace-period=0 --force > $(OUT) 2>&1 || true; \
		kubectl label namespace default istio-injection- > $(OUT) 2>&1 || true; \
		kubectl wait --for=delete namespace/istio-system --timeout=60s > $(OUT) 2>&1 || true; \
		printf $(_INFO) "Completed" "Istio cleanup"; \
	else \
		printf $(_INFO) "Skipping cleanup" "Istio not installed"; \
	fi

clear-ns:
	printf $(_WARN) "Resetting Namespace" "$(NS)"
	kubectl delete --all all --namespace=$(NS) > $(OUT)

clear-helms:
	printf $(_WARN) "Deleting" "Permanent Helms"
	-helm delete $$(helm list --short) > $(OUT)

clear-helm:
	printf $(_WARN) "Deleting" "Deployed Helms"
	#Clear CRD's (Needed before Helm Deletion)
	#TASK: Add Mysql CRD's
		
	-helm delete $$(helm list --short | grep -v "$(BOOTABLE)") > $(OUT)

### Management
install:
	printf $(_TITLE) "Services" "Installing"
	$(MAKE) $(SERVICE_MAKE_OPTS) $(ANS)

update: ## Update Services via Helm
	printf $(_TITLE) "Services" "Updating"
	$(MAKE) $(SERVICE_MAKE_OPTS) CMD=upgrade $(ANS)

bootstrap:
	$(MAKE) $(SERVICE_MAKE_OPTS) $(subst \|, ,$(BOOTABLE)) $(XTRA)
	kubectl wait --for=condition=Ready pod -l app.kubernetes.io/name=traefik --timeout=2m; > $(OUT)

select: ## Select/Configure Service
	printf $(_TITLE) "Services" "Selecting"
	$(MAKE) $(SERVICE_MAKE_OPTS)  | awk '{print $$1}' | grep -v "help\|make" | gum choose --limit 5 > $(ANS_FILE)
	$(MAKE) --no-print-directory info setup

### Info
info-svc:
	printf $(_TITLE) "Services" "Info"
	echo $(ANS)

### Prepare
helm-add:
	printf $(_TITLE) "Helm" "Adding"
	helm repo add onechart https://chart.onechart.dev
	helm repo add go-fun https://amanhigh.github.io/go-fun
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
	helm repo add gabe565 https://charts.gabe565.com
	helm repo add apiclarity https://openclarity.github.io/apiclarity
	helm repo add deliveryhero https://charts.deliveryhero.io/
	helm repo add elastic https://helm.elastic.co
	helm repo add vector https://helm.vector.dev

helm-update:
	printf $(_TITLE) "Helm" "Updating"
	helm repo update > $(OUT)


hosts:
	if ! grep -q "docker" /etc/hosts; then \
		printf $(_TITLE) "Hosts" "Adding"; \
		echo "$(DOCKER_HOSTS)" | sudo tee -a /etc/hosts; \
	else \
		printf $(_TITLE) "Hosts" "Updating"; \
		sudo sed -i 's/.*docker/$(subst $(eval) ,\ ,$(DOCKER_HOSTS))/g' /etc/hosts; \
	fi

helm: helm-add helm-update

### Workflows
prepare: helm hosts ## One Time Setup
info: info-svc ## Service Info
infos: info ## Service Extended Info
setup: clear-helm install ## Clean and Resetup Configured Services
clean: clear-helm ## Stop Services
reset: clean clear-istio clear-kiali clear-ns bootstrap info setup ## Clean & Resetup All Services 
all: prepare reset update infos ## Run All Targets

# Function to check if Istio is installed
define istio_installed
$(shell kubectl get namespace istio-system >/dev/null 2>&1 && echo "true" || echo "false")
endef