### Variables
.DEFAULT_GOAL := help
OUT := /dev/null

CMD := sudo gor
PORT := 8085
REQUEST_FILE := requests

_INFO := "\033[33m[%s]\033[0m %s\n"  # Yellow text for "printf"
_DETAIL := "\033[34m[%s]\033[0m %s\n"  # Blue text for "printf"
_TITLE := "\033[32m[%s]\033[0m %s\n" # Green text for "printf"
_WARN := "\033[31m[%s]\033[0m %s\n" # Red text for "printf"

### Basic
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@printf $(_TITLE) "FirstTime" "prepare/all, OUT=/dev/stdout (Debug)"

stdout: ## Recording Stdout
	printf $(_TITLE) "Listening on Port" "$(PORT)"
	$(CMD) --input-raw :$(PORT) --output-stdout

replay: ## Replay Port $(PORT) to 8080
	printf $(_TITLE) "Replaying Traffic from $(PORT)" "to 8080"
	$(CMD) --input-raw :$(PORT) --output-http="http://localhost:8080"

save: ## Save Traffic to File
	printf $(_TITLE) "Saving Traffic from $(PORT)" "to File"
	$(CMD) --input-raw :$(PORT) --output-file=$(REQUEST_FILE).gor

load: ## Load Traffic from File
	printf $(_TITLE) "Loading Traffic from File" "to 8080"
	$(CMD) --input-file $(REQUEST_FILE)_0.gor --output-http="http://localhost:8080" --stats --output-http-stats

info: ## Info
	printf $(_TITLE) "Info"
	printf $(_INFO) "Guide" "https://github.com/buger/goreplay/wiki/Getting-Started"

clean: ## Cleanup
	sudo rm $(REQUEST_FILE)*

### Workflows
prepare: ## Onetime Setup

### Useful Flags ##
## Rate Control 
# Limit: --output-http=”http://localhost:8001|10%" // Limits rate to 10% of incoming traffic
# SpeedUp: --output-http=”http://localhost:8001" // Replay faster than original traffic, effectively doubling the load on the system

## Request Control
# IncludeURL: --http-allow-url /api
# Param (Tainting Request): --http-set-param PERF_TEST=true