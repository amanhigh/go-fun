
helm install --create-namespace --replace --set auth.rootPassword=root --set auth.username=aman --set auth.password=aman --wait -n fun-app fun-mysql bitnami/mysql

kubectl label namespace fun-app istio-injection=enabled