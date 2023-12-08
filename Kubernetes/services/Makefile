### Variables
.DEFAULT_GOAL := help
ANS_FILE=/tmp/k8-svc.txt

### Basic
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


### Cleanup
clear-kiali:
	@printf "\033[1;32m Clearing Kiali \033[0m \n"
	kubectl delete kiali --all --all-namespaces 2> /dev/null

clear-istio:
	@printf "\033[1;32m Clearing Istio \033[0m \n"
	helm delete -n istio-system $(helm list --short -n istio-system)

NS:=$(shell kubectl get sa -o=jsonpath='{.items[0]..metadata.namespace}')
clear-ns:
	@printf "\033[1;32m Restting Namespace: $(NS) \033[0m \n"
	kubectl delete --all all --namespace=$(NS)

clear-helm:
	@printf "\033[1;32m Deleting all Helms \033[0m \n"
	helm delete $(helm list --short)

### Management
#HACK: Clear Istio
reset: delete clear-helmclear-ns ## Reset Kubernetes

delete: ## Delete deployed Helms
	@printf "\033[1;32m Deleting all Helms \033[0m \n"
	#Clear CRD's (Needed before Helm Deletion)
	#HACK: Clear Kiali if istio enabled
	#HACK: Add Mysql CRD's
		
	helm delete $(helm list --short | grep -v "traefik\|dashy") 2> /dev/null

install: ## Install Services
	@$(MAKE) --no-print-directory -f ./services.mk $$(cat $(ANS_FILE))

update: ## Update Services
	@$(MAKE) --no-print-directory -f ./services.mk CMD=upgrade $$(cat $(ANS_FILE))

select: ## Select Service
	@$(MAKE) --no-print-directory -f ./services.mk  | awk '{print $$1}' | grep -v "help\|make" | gum choose --limit 5 > $(ANS_FILE)
	@$(MAKE) info --no-print-directory

info: ## Info
	@printf "\033[1;32m Selected Services \033[0m \n"
	@cat $(ANS_FILE)