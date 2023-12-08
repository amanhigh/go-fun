### Variables
.DEFAULT_GOAL := help
CMD=install
ANS_FILE=/tmp/k8-svc.txt

#TODO: Add Locust
# Bootstrap: helm show values bitnami/postgresql > postgres.yml
# Debug: find . | entr -s "helm template elasticsearch bitnami/elasticsearch -f elasticsearch.yml > debug.txt;./service.zsh -di"
# sudo kubefwd svc | awk '{ if($2 ~ /Port-Forward/) {print $0" URL: http://"$4"/"} else {print}}'

### Basic
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

### Bootstrap
traefik: ## Traefik
	-helm $(CMD) traefik traefik/traefik -f traefik.yml > /dev/null
	kubectl apply -f ./files/traefik/middleware.yml
	@# kubectl apply -f ./files/traefik/ingress.yml
	@printf "\033[1;33m Dashboard: http://docker:9000/dashboard/#/ \033[0m \n"
	@printf "\033[1;33m HealthCheck: http://docker:9000/ping \033[0m \n"
	@printf "\033[1;33m Ingress: http://docker:8000/mysqladmin \033[0m \n"
	@printf "\033[1;33m PortForward(80): sudo kubectl port-forward deployment/traefik 80:8000 > /dev/null &\033[0m \n"

dashy: ## Dashy
	-helm $(CMD) dashy onechart/onechart -f dashy.yml > /dev/null
	@printf "\033[1;33m http://dashy.docker/ \033[0m \n"

### Http
httpbin: ## Httpbin
	-helm $(CMD) httpbin onechart/onechart -f httpbin.yml > /dev/null
	@printf "\033[1;33m Swagger: http://httpbin.docker \033[0m \n"
	@printf "\033[1;33m http://httpbin.docker/anything \033[0m \n"
	@printf "\033[1;33m curl http://httpbin:8810/headers \033[0m \n"

loader: ## Stress Testing
	-helm $(CMD) vegeta onechart/onechart -f vegeta.yml > /dev/null
	@printf "\033[1;33m Check Logs for Output \033[0m \n"
	@printf "\033[1;33m echo 'GET http://nginx' | vegeta attack | vegeta report \033[0m \n"

app: ## Fun Application
	-helm $(CMD) app onechart/onechart -f app.yml > /dev/null
	@printf "\033[1;33m http://app.docker/app/metrics\n http://app.docker/app/person/all \033[0m \n"
	@printf "\033[1;33m echo 'GET http://app:9001/person/all' | vegeta attack | vegeta report \033[0m \n"
    
proxy: ## Proxy Servers
	@printf "\033[1;32m Nginx \033[0m \n"
	-helm $(CMD) nginx bitnami/nginx -f nginx.yml > /dev/null
	@printf "\033[1;33m http://nginx.docker/ \033[0m \n"
	@printf "\033[1;33m Refer Server Blocks for More Help.. \033[0m \n"

	@printf "\033[1;32m Resty \033[0m \n"
	-helm $(CMD) resty onechart/onechart -f resty.yml > /dev/null
	@printf "\033[1;33m http://resty.docker/ \033[0m \n"
	@printf "\033[1;33m http://resty.docker/example \033[0m \n"
	@printf "\033[1;33m http://resty.docker/ndtv \033[0m \n"

	@printf "\033[1;32m Squid \033[0m \n"
	-helm $(CMD) squid onechart/onechart -f squid.yml > /dev/null
	@printf "\033[1;33m ALLOW: curl -x localhost:3128 www.google.com \033[0m \n"
	@printf "\033[1;33m DENY: curl -x localhost:3128 www.fb.com \033[0m \n"

	@printf "\033[1;32m TinyProxy \033[0m \n"
	-helm $(CMD) tinyproxy stakater/application -f tinyproxy.yml > /dev/null
	@printf "\033[1;33m curl -x localhost:8888 tinyproxy.stats \033[0m \n"

cron: ## Cron Server
	-helm $(CMD) cron onechart/onechart -f rundeck.yml > /dev/null
	@printf "\033[1;33m http://cron.docker/health \033[0m \n"
	@printf "\033[1;33m http://cron.docker \033[0m \n"
	@printf "\033[1;33m Username/Password: admin/admin \033[0m \n"

portainer: ## Portainer
	-helm $(CMD) portainer portainer/portainer -f portainer.yml
	@printf "\033[1;33m http://portainer.docker/ \033[0m \n"

### Security
opa: ## Open Policy Agent
	-helm $(CMD) opa opa/opa-kube-mgmt -f opa.yml > /dev/null
	-helm $(CMD) opa-demo onechart/onechart -f opa-demo.yml > /dev/null
	@printf "\033[1;33m curl --user david:password http://opa.docker/finance/salary/david \033[0m \n"
	@printf "\033[1;33m Demo (opa-demo): /demo/hr.sh \033[0m \n"
	@printf "\033[1;33m Demo (opa-demo): /demo/authz.sh \033[0m \n"
	@printf "\033[1;33m Docker (Localhost): ./demo/docker.sh \033[0m \n"

vault: ## Hashicorp Vault
	-helm $(CMD) vault hashicorp/vault -f vault.yml > /dev/null
	@printf "\033[1;33m vault status \033[0m \n"
	@printf "\033[1;33m /demo/vault.sh \033[0m \n"

sonar: ## Sonar
	-helm $(CMD) sonar bitnami/sonarqube -f sonar.yml > /dev/null
	@printf "\033[1;33m http://sonar.docker/ \033[0m \n"
	@printf "\033[1;33m Login: aman/aman (Need 5GB Mem) \033[0m \n"

webshell: ## Web Shell
	-helm $(CMD) sshwifty onechart/onechart -f sshwifty.yml  > /dev/null
	-helm $(CMD) webssh onechart/onechart -f webssh.yml > /dev/null
	#FIXME: Fix Config for Web SSH
	-helm $(CMD) webssh2 onechart/onechart -f webssh2.yml > /dev/null
	@printf "\033[1;33m Sshwifty: http://sshwifty.docker/ \033[0m \n"
	@printf "\033[1;33m Webssh: http://webssh.docker/ \033[0m \n"
	@printf "\033[1;33m Webssh: http://webssh2.docker/ \033[0m \n"

### Databases
mysql-admin:
	-helm $(CMD) mysql-admin bitnami/phpmyadmin -f phpmyadmin.yml > /dev/null
	@printf "\033[1;33m http://mysqladmin.docker/\033[0m \n"

mysql: mysql-admin ## MySQL
	-helm $(CMD) mysql bitnami/mysql -f mysql.yml > /dev/null
	@printf "\033[1;33m MySQL(3306) Login: mysql-primary, root/root \033[0m \n"

postgres: mysql-admin ## PostgreSQL
	-helm $(CMD) postgres bitnami/postgresql -f postgres.yml > /dev/null
	@printf "\033[1;33m Postgres(5432) Login: postgres-primary, postgres/root \033[0m \n"

mongo: ## Mongo
	-helm $(CMD) mongo bitnami/mongodb -f mongo.yml > /dev/null
	@printf "\033[1;33m mongosh -u root -p root --host localhost  < /etc/files/scripts/mongo.js \033[0m \n"
	@printf "\033[1;33m Svc Endpoint: mongo-mongodb:27017 (Standalone Mode Only) \033[0m \n"

redis: ## Redis
	-helm $(CMD) redis bitnami/redis -f redis.yml > /dev/null
	-helm $(CMD) redis-admin onechart/onechart -f redis-admin.yml > /dev/null
	@printf "\033[1;33m redis-cli -c incr mycounter \033[0m \n"
	@printf "\033[1;33m redis-cli -c set mypasswd lol \033[0m \n"
	@printf "\033[1;33m redis-cli -c get mypasswd \033[0m \n"
	@printf "\033[1;33m Commander: http://redisadmin.docker/ \033[0m \n"

ldap: ## LDAP Server
	-helm $(CMD) ldap onechart/onechart -f ldap.yml > /dev/null
	-helm $(CMD) ldap-admin onechart/onechart -f ldap-admin.yml > /dev/null
	@printf "\033[1;33m CMD: ldapsearch -H ldap://localhost:3891 -xLL -D 'cn=admin,dc=example,dc=com' -b 'dc=example,dc=com' -W '(cn=admin)' \033[0m \n"
	@printf "\033[1;33m UI: http://ldapadmin.docker/ \033[0m \n"
	@printf "\033[1;33m Admin Login: Username:cn=admin,dc=example,dc=com Password: admin \033[0m \n"

mysql-op: ## Mysql Operator
	-helm $(CMD) mysql-operator bitpoke/mysql-operator -f bitspoke.yml > /dev/null
	kubectl apply -f ./files/bitspoke/secret.yml
	kubectl apply -f ./files/bitspoke/cluster.yml
	-helm $(CMD) mysql-admin bitnami/phpmyadmin > /dev/null
	@printf "\033[1;33m Mysql Info: kubectl get mysql; kubectl describe mysql mysql-operator; \033[0m \n"
	@printf "\033[1;33m Mysql Clear: kubectl delete mysql mysql-operator; \033[0m \n"
	@printf "\033[1;33m Login: root/root, aman/aman [Host: mysql] \033[0m \n"

consul: ## Consul
	-helm $(CMD) consul hashicorp/consul -f consul.yml > /dev/null
	@printf "\033[1;33m http://consul.docker/ \033[0m \n"

etcd: ## Etcd
	-helm $(CMD) etcd bitnami/etcd -f etcd.yml > /dev/null
	@printf "\033[1;33m ./demo/demo.sh \033[0m \n"

zookeeper: ## Zookeeper
	-helm $(CMD) zookeeper bitnami/zookeeper -f zookeeper.yml > /dev/null
	@printf "\033[1;33m /demo/demo.sh \033[0m \n"

### Telemetry
elk: ## ElasticSearch Kibana Logstash
	# helm $(CMD) logstash bitnami/logstash -f logstash.yml > /dev/null
	-helm $(CMD) elasticsearch bitnami/elasticsearch -f elasticsearch.yml > /dev/null
	-helm $(CMD) kibana bitnami/kibana -f kibana.yml > /dev/null
	@printf "\033[1;33m ElasticSearch: http://elastic.docker/_cluster/health?pretty \033[0m \n"
	@printf "\033[1;33m Kibana: http://kibana.docker \033[0m \n"

monitor: ## Prometheus, Grafana and Jaeger
	-helm $(CMD) prometheus prometheus-community/prometheus -f prometheus.yml > /dev/null
	@printf "\033[1;33m Prometheus Server: http://prometheus.docker/ \033[0m \n"
	@printf "\033[1;33m Prometheus Query: http://prometheus.docker/api/v1/query \033[0m \n"

	-helm $(CMD) grafana grafana/grafana -f grafana.yml > /dev/null
	@printf "\033[1;33m http://grafana.docker/login (aman/aman) \033[0m \n"
	@printf "\033[1;33m Datasource: http://grafana.docker/datasources/new \033[0m \n"
	@printf "\033[1;33m Add Datasource Prometheus: http://prometheus-server \033[0m \n"

	-helm $(CMD) jaeger jaegertracing/jaeger -f jaeger.yml > /dev/null
	@printf "\033[1;33m http://jaeger.docker/ \033[0m \n"

### Istio
istio: ## Istio Service Mesh
	-helm $(CMD) istio-base istio/base -n istio-system --create-namespace > /dev/null
	-helm $(CMD) istiod istio/istiod -n istio-system -f istio.yml > /dev/null
	@printf "\033[1;32m Enabled Istio for Default Namespace \033[0m \n"
	# kubectl label namespace default istio-injection-
	kubectl label namespace default istio-injection=enabled --overwrite

kiali: ## Kiali Dashboard
	-helm $(CMD) kiali-operator kiali/kiali-operator -f kiali.yml > /dev/null
	#Create Kiali CRD
	kubectl apply -f ./files/istio/kiali-crd.yml
	@printf "\033[1;33m Kiali: http://kiali.docker/kiali \033[0m \n"
