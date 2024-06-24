include ../../../common/tools/base.mk

### Variables
URL := http://localhost:8080/v1
DURATION := 5s

RESULT_FILE := result.bin
ID_FILE := id.txt

BATCH_SIZE := 10
TOTAL_RECORDS := 500

### Clean
clean-file:
	printf $(_TITLE) "Cleaning File"
	-rm $(ID_FILE)
	-rm $(RESULT_FILE)

### Load Test
load-read: ## Read Load Test
	printf $(_TITLE) "Get Load Test"
	URL=$(URL) jq -R -ncM --arg url "$$URL" 'inputs as $$id | {method: "GET", url: "\(env.URL)/person/\($$id)"}' < $(ID_FILE) \
	| vegeta attack -format=json -duration=$(DURATION) | vegeta report

load-list: ## List Load Test
	printf $(_TITLE) "List Load Test";
	URL=$(URL) BATCH_SIZE=$(BATCH_SIZE) TOTAL_RECORDS=$(TOTAL_RECORDS) \
	jq -ncM 'range(0; env.TOTAL_RECORDS | tonumber; env.BATCH_SIZE | tonumber) | {method: "GET", url: "\(env.URL)/person?offset=\(. | tostring)&limit=\(env.BATCH_SIZE)&sort_by=name&order=asc"}' \
	| vegeta attack -format=json -duration=$(DURATION) | vegeta report

load-write: ## Write Load Test
	printf $(_TITLE) "Write Load Test"; \
	URL=$(URL) jq -ncM --slurpfile a payload.json \
	'while(true; .+1) as $$i | $$a | $$a[0].name="\($$a[0].name)-\($$i)" | {method: "POST", header: {"Content-Type": ["application/json"]}, url: "\(env.URL)/person", body: .[0] | @base64}' \
	| vegeta attack -lazy -format=json -duration=$(DURATION) -output=$(RESULT_FILE)
	printf $(_DETAIL) "Extracting Id's"
	vegeta encode $(RESULT_FILE) | jq -r '.body | @base64d | fromjson | .id' > $(ID_FILE)
	printf $(_DETAIL) "Test Result"
	vegeta report $(RESULT_FILE)

load: load-write load-list load-read

### Info
info-load:
	printf $(_TITLE) "Load Info"
	printf $(_INFO) "URL" "$(URL)"
	printf $(_INFO) "DURATION" "$(DURATION)"
	printf $(_INFO) "BATCH_SIZE" "$(BATCH_SIZE)"
	printf $(_INFO) "TOTAL_RECORDS" "$(TOTAL_RECORDS)"

### Workflows
info: info-load ## Info
infos: info ## Extended Info
prepare: ## Onetime Setup
setup: info load ## Load Test (All API's)
clean: clean-file ## Clean
reset: setup clean ## Reset
all:prepare reset ## Run All Targets