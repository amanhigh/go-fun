### Variables
.DEFAULT_GOAL := help
OUT := /dev/null

POD_LABEL=app.kubernetes.io/component=primary
POD_NAME=$$(kubectl get pod --selector=$(POD_LABEL) -o name)
K8_CMD=kubectl exec -it $(POD_NAME)

DBNAME=play
RESTORE_DB_NAME=$(DBNAME)_restore
USER=postgres
PASS=root

URL=postgresql://$(USER):$(PASS)@localhost
PG=$(K8_CMD) -- psql $(URL)
PGDB=$(K8_CMD) -- psql $(URL)/$(DBNAME)

### Basic
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	printf "\033[1;33m Local Binary: -MYSQL=mysql \033[0m \n"
	printf "\033[1;32m FirstTime: prepare/all, OUT=/dev/stdout (Debug) \033[0m \n"

### Info
show-dblist: ## Show available databases
	printf "\033[1;32m Available Databases \033[0m \n"
	$(PG) -l;

show-db: ## Show Current Database
	printf "\033[1;32m Current Database \033[0m \n"
	$(PGDB) -c '\c';

show-schema: ## Show available schemas
	printf "\033[1;32m Available Schemas \033[0m \n"
	$(PGDB) -c '\dn'

show-tables: ## Show available tables
	printf "\033[1;32m Available Tables \033[0m \n"
	$(PGDB) -c '\dt'

show-index: ## Show available indexes
	printf "\033[1;32m Available Indexes \033[0m \n"
	$(PGDB) -c '\di'

show-sequence: ## Show available sequences
	printf "\033[1;32m Available Sequences \033[0m \n"
	$(PGDB) -c '\d'

### Create
create-db: ## Create database
	printf "\033[1;32m Creating Database: $(DBNAME) \033[0m \n"
	$(PG) -c 'create database $(DBNAME);'

create-clone: ## Create cloned database
	$(MAKE) --no-print-directory DBNAME=$(RESTORE_DB_NAME) create-db

seed: ## Seed database
	printf "\033[1;32m Seeding Database: $(DBNAME) \033[0m \n"
	$(K8_CMD) -- bash -c 'psql $(URL)/$(DBNAME) -f /tmp/seed.sql'

### Delete
delete-db: ## Delete database
	printf "\033[1;32m Deleting Database: $(DBNAME) \033[0m \n"
	-$(PG) -c 'drop database $(DBNAME);'

delete-backup:
	printf "\033[1;32m Deleting Backup Directory \033[0m \n"
	$(K8_CMD) -- bash -c 'rm -rf /tmp/backup'

delete-clone: delete-backup ## Delete cloned database
	$(MAKE) --no-print-directory DBNAME=$(RESTORE_DB_NAME) delete-db

### Housekeeping
# Passfile: https://tableplus.com/blog/2019/09/how-to-use-pgpass-in-postgresql.html
prepare: ## Prepare Pod for Test Setup
	printf "\033[1;32m Preparing Pod: $(POD_NAME) \033[0m \n"
	kubectl cp ../mysql/seed.sql `basename $(POD_NAME)`:/tmp/

backup: ## Backup Database
	printf "\033[1;32m Backing up Database: $(DBNAME) \033[0m \n"
	$(K8_CMD) -- bash -c 'pg_dump $(URL)/$(DBNAME) -v -Fd -j 2 -f /tmp/backup -Z4' > $(OUT)

# https://www.youtube.com/playlist?list=PL73JD1WuNceTCuWnSzj4iq2GWAnk6j2eL
restore: ## Restore Database
	printf "\033[1;32m Restoring Database: $(RESTORE_DB_NAME) \033[0m \n"
	$(K8_CMD) -- bash -c 'pg_restore -d $(URL)/$(RESTORE_DB_NAME) -v -Fd -j 2 /tmp/backup' > $(OUT)

### Exec
exec: ## Connect to MySQL
	$(PG)
ssh: ## SSH to MySQL
	$(K8_CMD) bash

### Workflows
info: show-db show-tables

infos: info show-dblist show-schema show-index show-sequence

setup: create-db seed ## Setup Database

clean: delete-db delete-clone ## Clean up database

clone: delete-clone backup create-clone restore ## Clone Database

reset: clean setup info ## Reset Database

all: prepare reset clone infos ## Run All Targets