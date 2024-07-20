include ../../common/tools/base.mk

### Variables
SERVICE_MAKE_OPTS=--no-print-directory -f ./services.mk
XTRA=""
BOOTABLE = traefik\|webui

DOCKER_HOSTS=127.0.0.1 docker httpbin.docker dashy.docker resty.docker webui.docker pdf.docker paperless.docker app.docker \
mysqladmin.docker redisadmin.docker cron.docker prometheus.docker grafana.docker jaeger.docker kiali.docker \
ldapadmin.docker webssh.docker elastic.docker kibana.docker webssh2.docker sshwifty.docker nginx.docker portainer.docker \
consul.docker opa.docker sonar.docker metabase.docker clarity.docker
NS:=$(shell kubectl get sa -o=jsonpath='{.items[0]..metadata.namespace}')
ANS_FILE:=/tmp/k8-svc.txt
ANS:=$(shell cat $(ANS_FILE))

### Cleanup
clear-kiali:
	printf "\033[1;31m Clearing Kiali \033[0m \n"
	kubectl delete kiali --all --all-namespaces 2> /dev/null

clear-istio:
	printf "\033[1;31m Clearing Istio \033[0m \n"
	# Remove Istio components
	-helm delete -n istio-system $(shell helm list -n istio-system -q) > $(OUT)
	# Remove the istio-system namespace
	-kubectl delete namespace istio-system --grace-period=0 --force > $(OUT)
	# Remove Istio injection label from default namespace
	-kubectl label namespace default istio-injection- > $(OUT)
	# Wait for resources to be fully removed
	kubectl wait --for=delete namespace/istio-system --timeout=60s > $(OUT) || true
	printf "\033[1;32m Istio cleanup completed \033[0m \n"

clear-ns:
	printf "\033[1;31m Resetting Namespace: $(NS) \033[0m \n"
	kubectl delete --all all --namespace=$(NS) > $(OUT)

clear-helms:
	printf "\033[1;31m Deleting Permanent Helms \033[0m \n"
	-helm delete $$(helm list --short) > $(OUT)

clear-helm:
	printf "\033[1;31m Deleting Deployed Helms \033[0m \n"
	#Clear CRD's (Needed before Helm Deletion)
	#HACK: #B Clear Kiali if istio enabled
	#HACK: Add Mysql CRD's
		
	-helm delete $$(helm list --short | grep -v "$(BOOTABLE)") > $(OUT)

### Management
install:
	printf "\033[1;32m Installing Services \033[0m \n"
	$(MAKE) $(SERVICE_MAKE_OPTS) $(ANS)

update: ## Update Services via Helm
	printf "\033[1;32m Updating Services \033[0m \n"
	$(MAKE) $(SERVICE_MAKE_OPTS) CMD=upgrade $(ANS)

bootstrap:
	$(MAKE) $(SERVICE_MAKE_OPTS) $(subst \|, ,$(BOOTABLE)) $(XTRA)
	kubectl wait --for=condition=Ready pod -l app.kubernetes.io/name=traefik --timeout=2m; > $(OUT)

select: ## Select/Configure Service
	printf "\033[1;32m Select Services \033[0m \n"
	$(MAKE) $(SERVICE_MAKE_OPTS)  | awk '{print $$1}' | grep -v "help\|make" | gum choose --limit 5 > $(ANS_FILE)
	$(MAKE) --no-print-directory info setup

### Info
info-svc:
	printf "\033[1;32m Selected Services \033[0m \n"
	echo $(ANS)

### Prepare
helm-add:
	printf "\033[1;32m Adding Helm Repos \033[0m \n"
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

helm-update:
	printf "\033[1;32m Updating Helm Repos \033[0m \n"
	helm repo update > $(OUT)


hosts:
	if ! grep -q "docker" /etc/hosts; then \
		printf "\033[1;32m Adding Docker Hosts to /etc/hosts \033[0m \n"; \
		echo "$(DOCKER_HOSTS)" | sudo tee -a /etc/hosts; \
	else \
		printf "\033[1;33m Update Docker Hosts in /etc/hosts \033[0m \n"; \
		sudo sed -i 's/.*docker/$(subst $(eval) ,\ ,$(DOCKER_HOSTS))/g' /etc/hosts; \
	fi

helm: helm-add helm-update

### Workflows
prepare: helm hosts ## One Time Setup
info: info-svc ## Service Info
infos: info ## Service Extended Info
setup: clear-helm install ## Clean and Resetup Configured Services
#HACK: #C Clear Istio
clean: clear-helm ## Stop Services
reset: clean clear-helms clear-ns bootstrap info setup ## Clean & Resetup All Services 
all: prepare reset update infos ## Run All Targets