helm delete -n fun-app fun-mysqladmin
helm delete -n fun-app fun-mysql
helm delete -n fun-app fun-redis
helm delete -n fun-app fun-app

kubectl delete ns fun-app