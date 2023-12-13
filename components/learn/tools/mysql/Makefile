### Variables
.DEFAULT_GOAL := help
POD_LABEL=app.kubernetes.io/component=primary
K8_CMD=kubectl exec -it $$(kubectl get pod --selector=$(POD_LABEL) -o name)
MYSQL=$(K8_CMD) -- mysql -uroot -proot

### Basic
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@printf "\033[1;33m Local Binary: -MYSQL=mysql \033[0m \n"

### Info
show-db: ## Show available databases
	$(MYSQL) -e 'show databases;'

exec: ## Connect to MySQL
	$(MYSQL)