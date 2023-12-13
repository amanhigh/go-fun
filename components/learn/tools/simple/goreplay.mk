### Variables
.DEFAULT_GOAL := help
OUT := /dev/null

CMD := sudo gor
PORT := 8085

### Basic
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@printf "\033[1;32m FirstTime: prepare/all, OUT=/dev/stdout (Debug) \033[0m \n"

stdout: ## Recording Stdout
	printf "\033[1;32m Listening on Port: $(PORT) \033[0m \n"
	$(CMD) --input-raw :$(PORT) --output-stdout

replay: ## Replay Port $(PORT) to 8080
	printf "\033[1;32m Replaying Traffic from $(PORT) to 8080 \033[0m \n"
	$(CMD) --input-raw :$(PORT) --output-http="http://localhost:8080"

save: ## Save Traffic to File
	printf "\033[1;32m Saving Traffic from $(PORT) to File \033[0m \n"
	$(CMD) --input-raw :$(PORT) --output-file=requests.gor

load: ## Load Traffic from File
	printf "\033[1;32m Loading Traffic from File to 8080 \033[0m \n"
	$(CMD) --input-file requests_0.gor --output-http="http://localhost:8080" --stats --output-http-stats

info: ## Info
	printf "\033[1;32m Info \033[0m \n"
	printf "\033[1;33m https://github.com/buger/goreplay/wiki/Getting-Started \033[0m \n"
### Workflows
prepare: ## Onetime Setup
setup: ## Setup
clean: ## Clean

reset: clean setup info ## Reset
all:prepare reset ## Run All Targets

### Useful Flags ##
## Rate Control 
# Limit: --output-http=”http://localhost:8001|10%" // Limits rate to 10% of incoming traffic
# SpeedUp: --output-http=”http://localhost:8001" // Replay faster than original traffic, effectively doubling the load on the system

## Request Control
# IncludeURL: --http-allow-url /api
# Param (Tainting Request): --http-set-param PERF_TEST=true