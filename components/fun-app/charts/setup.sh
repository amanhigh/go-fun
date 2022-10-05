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
helm install -n fun-app fun-app . --set rateLimit.perMin=150 > /dev/null
echo -en "\033[1;33m FunApp: http://localhost:8091/api/v1/namespaces/fun-app/services/fun-app:9000/proxy/metrics \033[0m \n"

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



