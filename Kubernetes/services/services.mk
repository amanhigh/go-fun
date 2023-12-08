### Variables
.DEFAULT_GOAL := help
CMD=install
ANS_FILE=/tmp/k8-svc.txt

### Basic
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

### Services
httpbin: ## Httpbin
	helm $(CMD) httpbin onechart/onechart -f httpbin.yml > /dev/null
	@printf "\033[1;33m Swagger: http://httpbin.docker \033[0m \n"
	@printf "\033[1;33m http://httpbin.docker/anything \033[0m \n"
	@printf "\033[1;33m curl http://httpbin:8810/headers \033[0m \n"
cron: ## Cron
	helm $(CMD) cron onechart/onechart -f rundeck.yml > /dev/null
	@printf "\033[1;33m http://cron.docker \033[0m \n"
	@printf "\033[1;33m http://cron.docker/health \033[0m \n"