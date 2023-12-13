### Variables
.DEFAULT_GOAL := help
POD_LABEL=app.kubernetes.io/component=primary
POD_NAME=$$(kubectl get pod --selector=$(POD_LABEL) -o name)
K8_CMD=kubectl exec -it $(POD_NAME)

DBNAME=play
RESTORE_DB_NAME=$(DBNAME)_restore
CONFIG=/tmp/mysql.cnf

MYSQL=$(K8_CMD) -- mysql --defaults-extra-file=$(CONFIG)
MYSQLDB= $(MYSQL) $(DBNAME)

### Basic
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@printf "\033[1;33m Local Binary: -MYSQL=mysql \033[0m \n"
	@printf "\033[1;32m prepare/all Before Use \033[0m \n"

### Info
show-db: ## Show available databases
	@printf "\033[1;32m Available Databases \033[0m \n"
	@$(MYSQL) -e 'show databases;'

show-tables: ## Show available tables
	@printf "\033[1;32m Available Tables \033[0m \n"
	@$(MYSQLDB) -e 'show tables;'

show-process: ## Show MySQL process
	@printf "\033[1;32m MySQL Process \033[0m \n"
	@$(MYSQL) -e 'show processlist;'

### Create
create-db: ## Create database
	@printf "\033[1;32m Creating Database: $(DBNAME) \033[0m \n"
	@$(MYSQL) -e 'create database $(DBNAME);'

create-clone: ## Create cloned database
	@$(MAKE) --no-print-directory DBNAME=$(RESTORE_DB_NAME) create-db

seed: ## Seed database
	@printf "\033[1;32m Seeding Database: $(DBNAME) \033[0m \n"
	@$(K8_CMD) -- bash -c 'mysql --defaults-extra-file=$(CONFIG) $(DBNAME) < /tmp/seed.sql'

### Delete
delete-db: ## Delete database
	@printf "\033[1;32m Deleting Database: $(DBNAME) \033[0m \n"
	@-$(MYSQL) -e 'drop database $(DBNAME);'

delete-clone: ## Delete cloned database
	@$(MAKE) --no-print-directory DBNAME=$(RESTORE_DB_NAME) delete-db
 
### Housekeeping
prepare: ## Prepare Pod for Test Setup
	@printf "\033[1;32m Preparing Pod: $(POD_NAME) \033[0m \n"
	@kubectl cp seed.sql `basename $(POD_NAME)`:/tmp
	@kubectl cp mysql.cnf `basename $(POD_NAME)`:/tmp

backup: ## Backup Database
	@printf "\033[1;32m Backing up Database: $(DBNAME) \033[0m \n"
	@$(K8_CMD) -- bash -c 'mysqldump --defaults-extra-file=$(CONFIG) $(DBNAME) --routines > /tmp/backup.sql'

restore: ## Restore Database
	@printf "\033[1;32m Restoring Database: $(RESTORE_DB_NAME) \033[0m \n"
	@$(K8_CMD) -- bash -c 'mysql --defaults-extra-file=$(CONFIG) $(RESTORE_DB_NAME) < /tmp/backup.sql'

### Exec
exec: ## Connect to MySQL
	$(MYSQL)
ssh: ## SSH to MySQL
	$(K8_CMD) bash

### Workflows
info: show-db show-tables

infos: info show-process

setup: create-db seed ## Setup Database

clean: delete-db delete-clone ## Clean up database

clone: delete-clone backup create-clone restore ## Clone Database

reset: clean setup info ## Reset Database

all: prepare reset clone ## All Targets
