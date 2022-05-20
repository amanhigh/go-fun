helm delete -n fun-app fun-mysql
helm delete -n fun-app fun-redis
kubectl delete -n fun-app -f .


kubectl delete ns fun-app