### Variables
.DEFAULT_GOAL := help
OUT := /dev/null

### Basic
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	printf $(_TITLE) "FirstTime: prepare/all, OUT=/dev/stdout (Debug) "

### Formatting
_INFO := "\033[33m[%s]\033[0m %s\n"  # Yellow text for "printf"
_DETAIL := "\033[34m[%s]\033[0m %s\n"  # Blue text for "printf"
_TITLE := "\033[32m[%s]\033[0m %s\n" # Green text for "printf"
_WARN := "\033[31m[%s]\033[0m %s\n" # Red text for "printf"