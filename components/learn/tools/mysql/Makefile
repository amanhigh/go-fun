### Variables
.DEFAULT_GOAL := help
POD_LABEL=app.kubernetes.io/component=primary
POD_NAME=$$(kubectl get pod --selector=$(POD_LABEL) -o name)
K8_CMD=kubectl exec -it $(POD_NAME)
DBNAME=play
USER=root
PASSWORD=root
MYSQL=$(K8_CMD) -- mysql -u $(USER) -p$(PASSWORD)
MYSQLDB= $(MYSQL) $(DBNAME)

### Basic
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@printf "\033[1;33m Local Binary: -MYSQL=mysql \033[0m \n"

### Info
show-db: ## Show available databases
	$(MYSQL) -e 'show databases;'

show-tables: ## Show available tables
	$(MYSQLDB) -e 'show tables;'

info: show-db show-tables

### Create
create-db: ## Create database
	$(MYSQL) -e 'create database $(DBNAME);'

seed: ## Seed database
	kubectl cp seed.sql mysql-primary-0:/tmp
	$(K8_CMD) -- bash -c 'mysql -u $(USER) -p$(PASSWORD) $(DBNAME) < /tmp/seed.sql'

setup: create-db seed ## Setup Database

### Delete
delete-db: ## Delete database
	$(MYSQL) -e 'drop database $(DBNAME);'
 
clean: delete-db ## Clean up database


### Workflows
reset: clean setup info ## Reset Database

exec: ## Connect to MySQL
	$(MYSQL)