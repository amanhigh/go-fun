kubectl create ns fun-app
kubectl label namespace fun-app istio-injection=enabled

echo -en "\033[1;32m Setup Mysql \033[0m \n"
#TODO: Master Slave Setup
helm install --set auth.rootPassword=root --set auth.database=compute --set auth.username=aman --set auth.password=aman --wait -n fun-app fun-mysql bitnami/mysql
helm install --wait -n fun-app fun-mysqladmin bitnami/phpmyadmin


echo -en "\033[1;32m Setup FunApp \033[0m \n"
kubectl apply -n fun-app -f .

echo -en "\033[1;32m MysqlAdmin: http://localhost:8091/api/v1/namespaces/fun-app/services/fun-mysqladmin-phpmyadmin:80/proxy/index.php?server=fun-mysql \033[0m \n"
echo -en "\033[1;32m FunApp: http://localhost:8091/api/v1/namespaces/fun-app/services/fun-app:9000/proxy/metrics \033[0m \n"