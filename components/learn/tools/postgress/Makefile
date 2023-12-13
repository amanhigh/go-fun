### Variables
.DEFAULT_GOAL := help
POD_LABEL=app.kubernetes.io/component=primary
POD_NAME=$$(kubectl get pod --selector=$(POD_LABEL) -o name)
K8_CMD=kubectl exec -it $(POD_NAME)

DBNAME=play
RESTORE_DB_NAME=$(DBNAME)_restore
USER=postgres
PASS=root

URL=postgresql://$(USER):$(PASS)@localhost
PG=$(K8_CMD) -- psql $(URL)

### Basic
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	printf "\033[1;33m Local Binary: -MYSQL=mysql \033[0m \n"
	printf "\033[1;32m prepare/all Before Use \033[0m \n"

### Info
show-db: ## Show available databases
	printf "\033[1;32m Available Databases \033[0m \n"
	$(PG) -l;

show-tables: ## Show available tables
	printf "\033[1;32m Available Tables \033[0m \n"
	$(PG) -c '\c $(DBNAME) \dt'

### Create
create-db: ## Create database
	printf "\033[1;32m Creating Database: $(DBNAME) \033[0m \n"
	$(PG) -c 'create database $(DBNAME);'

create-clone: ## Create cloned database
	$(MAKE) --no-print-directory DBNAME=$(RESTORE_DB_NAME) create-db

seed: ## Seed database
	printf "\033[1;32m Seeding Database: $(DBNAME) \033[0m \n"
	$(K8_CMD) -- bash -c 'psql $(URL) -f /tmp/seed.sql'

### Delete
delete-db: ## Delete database
	printf "\033[1;32m Deleting Database: $(DBNAME) \033[0m \n"
	-$(PG) -c 'drop database $(DBNAME);'

delete-clone: ## Delete cloned database
	$(MAKE) --no-print-directory DBNAME=$(RESTORE_DB_NAME) delete-db
 
### Housekeeping
# Passfile: https://tableplus.com/blog/2019/09/how-to-use-pgpass-in-postgresql.html
prepare: ## Prepare Pod for Test Setup
	printf "\033[1;32m Preparing Pod: $(POD_NAME) \033[0m \n"
	kubectl cp ../mysql/seed.sql `basename $(POD_NAME)`:/tmp/

### Exec
exec: ## Connect to MySQL
	$(PG)
ssh: ## SSH to MySQL
	$(K8_CMD) bash

### Workflows
info: show-db show-tables

setup: create-db seed ## Setup Database

clean: delete-db## Clean up database

reset: clean setup info ## Reset Database

all: prepare reset ## All Targets