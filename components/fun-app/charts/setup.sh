kubectl create ns fun-app
echo -en "\033[1;32m Istio Enabled \033[0m \n"
kubectl label namespace fun-app istio-injection=enabled

echo -en "\033[1;32m Setup Mysql \033[0m \n"
# helm install --set auth.rootPassword=root --set auth.database=compute --set auth.username=aman --set auth.password=aman -n fun-app fun-mysql bitnami/mysql > /dev/null
echo -en "\033[1;33m K9S Shell to master, 'mysql -u root -p compute' (Password: root) \033[0m \n"
echo -en "\033[1;33m MysqlAdmin: http://localhost:8091/api/v1/namespaces/fun-app/services/fun-app-phpmyadmin:80/proxy/index.php?server=fun-app-mysql \033[0m \n"

echo -en "\033[1;32m Setup Redis \033[0m \n"
echo -en "\033[1;33m K9S Shell to master, 'redis-cli' OR 'redis-cli -h fun-app-redis-master-0'/fun-app-redis-replicas \033[0m \n"

echo -en "\033[1;32m Setup FunApp \033[0m \n"
# helm repo add bitnami https://charts.bitnami.com/bitnami
helm dependency build
helm install -n fun-app fun-app . --set rateLimit.perMin=150 > /dev/null
echo -en "\033[1;33m FunApp (Proxy): http://localhost:8091/api/v1/namespaces/fun-app/services/fun-app:9000/proxy/metrics \033[0m \n"
echo -en "\033[1;33m FunApp (minikube tunnel): http://localhost:9000/metrics \033[0m \n"

echo -en "\033[1;32m Vegeta Attack (Login Host) \033[0m \n"
kubectl run vegeta -n fun-app --image="peterevans/vegeta" -- sh -c "sleep 10000"
echo -en "\033[1;33m HELM: echo 'GET http://fun-app:9000/person/all' | vegeta attack | vegeta report \033[0m \n"
echo -en "\033[1;33m DEVSPACE: echo 'GET http://app:8080/person/all' | vegeta attack | vegeta report \033[0m \n"

echo -en "\033[1;32m Metrics (Istio Only) \033[0m \n"
echo -en "\033[1;33m Prometheus: http://localhost:9090/graph?g0.expr=rate(fun_app_person_count%5B5m%5D)&g0.tab=0&g0.stacked=0&g0.show_exemplars=0&g0.range_input=5m&g1.expr=fun_app_person_create_time_bucket&g1.tab=0&g1.stacked=1&g1.show_exemplars=1&g1.range_input=1h&g2.expr=rate(fun_app_person_create_time_count%5B5m%5D)&g2.tab=0&g2.stacked=0&g2.show_exemplars=0&g2.range_input=1h \033[0m \n"
echo -en "\033[1;33m Grafana Import: /fun-app/it/dashboard.json \033[0m \n"

### Helpful Commands
# helm init fun-app - Bootstrap Charts
# helm template . - Preview Charts with Values
# helm lint . - Check Errors
# helm show values <Chart Name> - Show configurable values

# helm install -n <Namespace> <Chart Name> . [--set <key>=<value>]
# helm upgrade -n <Namespace> <Chart Name> . [--set <key>=<value>]

# helm status -n <Namespace> <Chart Name>
# helm history -n <Namespace> <Chart Name>
# helm rollback -n <Namespace> <Chart Name> [Revision]
# helm delete -n <Namespace> <Chart Name>


# helm list -n <Namespace>

# helm repo list
# helm dependency list
# helm dependency update
# helm dependency build



