### Variables
.DEFAULT_GOAL := help
CMD=install

#TODO: Add Locust
# Bootstrap: helm show values bitnami/postgresql > postgres.yml
# Debug: find . | entr -s "helm template elasticsearch bitnami/elasticsearch -f elasticsearch.yml > debug.txt;make setup"
# sudo kubefwd svc | awk '{ if($2 ~ /Port-Forward/) {print $0" URL: http://"$4"/"} else {print}}'

### Basic
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

### Bootstrap
traefik: ## Traefik
	-helm $(CMD) traefik traefik/traefik -f traefik.yml > /dev/null
	kubectl apply -f ./files/traefik/middleware.yml
	# kubectl apply -f ./files/traefik/ingress.yml
	printf $(_INFO) "Dashboard" "http://docker:9000/dashboard/#/"
	printf $(_INFO) "HealthCheck" "http://docker:9000/ping"
	printf $(_INFO) "Ingress" "http://docker:8000/mysqladmin"

dashy: ## Dashy
	-helm $(CMD) dashy onechart/onechart -f dashy.yml > /dev/null
	printf $(_INFO) "Dashy" "http://dashy.docker/"

### Http
httpbin: ## Httpbin
	-helm $(CMD) httpbin onechart/onechart -f httpbin.yml > /dev/null
	printf $(_INFO) "Swagger" "http://httpbin.docker"
	printf $(_INFO) "httpbin" "http://httpbin.docker/anything"
	printf $(_INFO) "curl http://httpbin:8810/headers"

loader: ## Stress Testing
	-helm $(CMD) vegeta onechart/onechart -f vegeta.yml > /dev/null
	printf $(_INFO) "Check Logs for Output"
	printf $(_INFO) "Vegeta" "echo 'GET http://nginx' | vegeta attack | vegeta report"

app: ## Fun Application
	-helm $(CMD) app onechart/onechart -f app.yml > /dev/null
	printf $(_INFO) "App Metrics" "http://app.docker/app/metrics"
	printf $(_INFO) "App All" "http://app.docker/app/person/all"
	printf $(_INFO) "Vegeta" "echo 'GET http://app:9001/person/all' | vegeta attack | vegeta report"

proxy: ## Proxy Servers
	printf $(_TITLE) "Nginx"
	-helm $(CMD) nginx bitnami/nginx -f nginx.yml > /dev/null
	printf $(_INFO) "Nginx" "http://nginx.docker/"
	printf $(_INFO) "Server Blocks" "Refer Server Blocks for More Help.."

	printf $(_TITLE) "Resty"
	-helm $(CMD) resty onechart/onechart -f resty.yml > /dev/null
	printf $(_INFO) "Resty" "http://resty.docker/"
	printf $(_INFO) "Resty Example" "http://resty.docker/example"
	printf $(_INFO) "Resty NDTV" "http://resty.docker/ndtv"

	printf $(_TITLE) "Squid"
	-helm $(CMD) squid onechart/onechart -f squid.yml > /dev/null
	printf $(_INFO) "ALLOW" "curl -x localhost:3128 www.google.com"
	printf $(_INFO) "DENY" "curl -x localhost:3128 www.fb.com"

	printf $(_TITLE) "TinyProxy"
	-helm $(CMD) tinyproxy stakater/application -f tinyproxy.yml > /dev/null
	printf $(_INFO) "curl localhost:8888 tinyproxy.stats"

cron: ## Cron Server
	-helm $(CMD) cron onechart/onechart -f rundeck.yml > /dev/null
	printf $(_INFO) "Health" "http://cron.docker/health"
	printf $(_INFO) "Cron" "http://cron.docker"
	printf $(_INFO) "Credentials" "Username/Password: admin/admin"

portainer: ## Portainer
	-helm $(CMD) portainer portainer/portainer -f portainer.yml
	printf $(_INFO) "Portainer" "http://portainer.docker/"

webui: ## Open Web UI
	-helm $(CMD) webui onechart/onechart -f webui.yml > /dev/null
	printf $(_INFO) "WebUI" "http://webui.docker/"

pdf: ## Open Stirling PDF
	-helm $(CMD) pdf onechart/onechart -f pdf.yml > /dev/null
	printf $(_INFO) "PDF" "http://pdf.docker/"

paperless: postgres redis ## Paperless NGX
	-helm $(CMD) paperless gabe565/paperless-ngx -f paperless.yml > /dev/null
	printf $(_INFO) "Paperless" "http://paperless.docker/"

### Security
opa: ## Open Policy Agent
	-helm $(CMD) opa opa/opa-kube-mgmt -f opa.yml > /dev/null
	-helm $(CMD) opa-demo onechart/onechart -f opa-demo.yml > /dev/null
	printf $(_INFO) "curl --user david:password http://opa.docker/finance/salary/david"
	printf $(_DETAIL) "Demo (opa-demo)" "/demo/hr.sh"
	printf $(_DETAIL) "Demo (opa-demo)" "/demo/authz.sh"
	printf $(_DETAIL) "Docker (Localhost)" "./demo/docker.sh"

vault: ## Hashicorp Vault
	-helm $(CMD) vault hashicorp/vault -f vault.yml > /dev/null
	printf $(_INFO) "vault status"
	printf $(_DETAIL) "/demo/vault.sh"

sonar: ## Sonar
	-helm $(CMD) sonar bitnami/sonarqube -f sonar.yml > /dev/null
	printf $(_INFO) "http://sonar.docker/"
	printf $(_DETAIL) "Login" "aman/aman (Need 5GB Mem)"

webshell: ## Web Shell
	-helm $(CMD) sshwifty onechart/onechart -f sshwifty.yml  > /dev/null
	-helm $(CMD) webssh onechart/onechart -f webssh.yml > /dev/null
	-helm $(CMD) webssh2 onechart/onechart -f webssh2.yml > /dev/null
	# FIXME: Fix Config for Web SSH
	printf $(_INFO) "Sshwifty" "http://sshwifty.docker/"
	printf $(_INFO) "Webssh" "http://webssh.docker/"
	printf $(_INFO) "Webssh" "http://webssh2.docker/"

### Databases
mysql-admin:
	-helm $(CMD) mysql-admin bitnami/phpmyadmin -f phpmyadmin.yml > /dev/null
	printf $(_INFO) "URL" "http://mysqladmin.docker/"

# BUG: Metabase doesn't start on h2, move pg with snapshot volume ?
metabase:
	-helm $(CMD) metabase onechart/onechart -f metabase.yml > /dev/null
	printf $(_INFO) "URL" "http://metabase.docker/"
	printf $(_INFO) "Login" "email/aman"
	printf $(_INFO) "DB" "sudo cp -r h2.db /tmp/mini/metabase"

mysql: metabase ## MySQL
	-helm $(CMD) mysql bitnami/mysql -f mysql.yml > /dev/null
	printf $(_INFO) "MySQL(3306) Login" "mysql-primary/mysql-secondary, root/root, aman/aman"

postgres: metabase ## PostgreSQL
	-helm $(CMD) postgres bitnami/postgresql -f postgres.yml > /dev/null
	printf $(_INFO) "Postgres(5432) Login" "pg-primary/pg-read, postgres/root, aman/aman"

mongo: ## Mongo
	-helm $(CMD) mongo bitnami/mongodb -f mongo.yml > /dev/null
	printf $(_INFO) "Cli" "mongosh -u root -p root --host localhost  < /etc/files/scripts/mongo.js"
	printf $(_INFO) "Endpoint" "mongo-mongodb:27017 (Standalone Mode Only)"

redis: ## Redis
	-helm $(CMD) redis bitnami/redis -f redis.yml > /dev/null
	-helm $(CMD) redis-admin onechart/onechart -f redis-admin.yml > /dev/null
	printf $(_INFO) "Cli" "redis-cli -c incr mycounter"
	printf $(_INFO) "Cli" "redis-cli -c set mypasswd lol"
	printf $(_INFO) "Cli" "redis-cli -c get mypasswd"
	printf $(_INFO) "Commander" "http://redisadmin.docker/"

ldap: ## LDAP Server
	-helm $(CMD) ldap onechart/onechart -f ldap.yml > /dev/null
	-helm $(CMD) ldap-admin onechart/onechart -f ldap-admin.yml > /dev/null
	printf $(_INFO) "CMD" "ldapsearch -H ldap://localhost:3891 -xLL -D 'cn=admin,dc=example,dc=com' -b 'dc=example,dc=com' -W '(cn=admin)'"
	printf $(_INFO) "UI" "http://ldapadmin.docker/"
	printf $(_INFO) "Admin Login" "Username:cn=admin,dc=example,dc=com Password: admin"

mysql-op: ## Mysql Operator
	-helm $(CMD) mysql-operator bitpoke/mysql-operator -f bitspoke.yml > /dev/null
	kubectl apply -f ./files/bitspoke/secret.yml
	kubectl apply -f ./files/bitspoke/cluster.yml
	-helm $(CMD) mysql-admin bitnami/phpmyadmin > /dev/null
	printf $(_INFO) "Mysql Info" "kubectl get mysql; kubectl describe mysql mysql-operator;"
	printf $(_INFO) "Mysql Clear" "kubectl delete mysql mysql-operator;"
	printf $(_INFO) "Login" "root/root, aman/aman [Host: mysql]"

consul: ## Consul
	-helm $(CMD) consul hashicorp/consul -f consul.yml > /dev/null
	printf $(_INFO) "URL" "http://consul.docker/"

etcd: ## Etcd
	-helm $(CMD) etcd bitnami/etcd -f etcd.yml > /dev/null
	printf $(_INFO) "Demo" "./demo/demo.sh"

zookeeper: ## Zookeeper
	-helm $(CMD) zookeeper bitnami/zookeeper -f zookeeper.yml > /dev/null
	printf $(_INFO) "Demo" "/demo/demo.sh"

### Telemetry
elk: ## ElasticSearch Kibana Logstash
	#FIXME: Logstash in ELK
	# helm $(CMD) logstash bitnami/logstash -f logstash.yml > /dev/null
	-helm $(CMD) elasticsearch bitnami/elasticsearch -f elasticsearch.yml > /dev/null
	-helm $(CMD) kibana bitnami/kibana -f kibana.yml > /dev/null
	printf $(_TITLE) "ELK needs CPU: 4, Memory: 10Gig"
	printf $(_INFO) "ElasticSearch" "http://elastic.docker/_cluster/health?pretty"
	printf $(_INFO) "ES Master" "http://docker:9200"
	printf $(_INFO) "Kibana" "http://kibana.docker"

monitor: ## Prometheus, Grafana and Jaeger
	-helm $(CMD) prometheus prometheus-community/prometheus -f prometheus.yml > /dev/null
	printf $(_INFO) "Prometheus Server" "http://prometheus.docker/"
	printf $(_INFO) "Prometheus Query" "http://prometheus.docker/api/v1/query"

	-helm $(CMD) grafana grafana/grafana -f grafana.yml > /dev/null
	printf $(_INFO) "Grafana Login" "http://grafana.docker/login (aman/aman)"
	printf $(_INFO) "Datasource" "http://grafana.docker/datasources/new"
	printf $(_INFO) "Add Datasource Prometheus" "http://prometheus-server"

	-helm $(CMD) jaeger jaegertracing/jaeger -f jaeger.yml > /dev/null
	printf $(_INFO) "Jaeger" "http://jaeger.docker/"

### Istio
istio: ## Istio Service Mesh
	-helm $(CMD) istio-base istio/base -n istio-system --create-namespace > /dev/null
	-helm $(CMD) istiod istio/istiod -n istio-system -f istio.yml > /dev/null
	printf $(_TITLE) "Istio" "Enabled Istio for Default Namespace"
	# kubectl label namespace default istio-injection-
	kubectl label namespace default istio-injection=enabled --overwrite

kiali: ## Kiali Dashboard
	-helm $(CMD) kiali-operator kiali/kiali-operator -f kiali.yml > /dev/null
	#Create Kiali CRD
	kubectl apply -f ./files/istio/kiali-crd.yml
	printf $(_INFO) "Kiali" "http://kiali.docker/kiali"

### Formatting
_INFO := "\033[33m[%s]\033[0m %s\n"  # Yellow text for "printf"
_DETAIL := "\033[34m[%s]\033[0m %s\n"  # Blue text for "printf"
_TITLE := "\033[32m[%s]\033[0m %s\n" # Green text for "printf"
_WARN := "\033[31m[%s]\033[0m %s\n" # Red text for "printf"