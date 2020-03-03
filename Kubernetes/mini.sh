minikube delete;
minikube start \
--extra-config="apiserver.enable-swagger-ui=true" \
--extra-config="apiserver.service-account-api-audiences=api" \
--extra-config="apiserver.service-account-issuer=api" \
--extra-config="apiserver.service-account-key-file=/var/lib/minikube/certs/sa.pub" \
--extra-config="apiserver.service-account-signing-key-file=/var/lib/minikube/certs/sa.key";

minikube dashboard &
kubectl proxy --port=8091 &
echo -en "\033[1;32m Dashboard: http://localhost:8091/api/v1/namespaces/kubernetes-dashboard/services/http:kubernetes-dashboard:/proxy/# \033[0m \n"
echo -en "\033[1;32m Swagger: http://localhost:8091/swagger-ui \033[0m \n"


# Helpful Commands
# Docker Imageto Minikube: eval $(minikube docker-env); docker build -t fun-app .
# Port Forward (Local Port 9090 to Container Port 8080) - kubectl port-forward `kubectl get pods -o name | grep fun-app | head  -1` 9090:8080
# Logs - kubectl logs `kubectl get pods -o name | grep fun-app | head  -1` -f
# Login - kubectl -it exec `kubectl get pods -o name | grep fun-app | head  -1` bash
#
#