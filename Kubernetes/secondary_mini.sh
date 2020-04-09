PORT=8092
minikube -p secondary delete
minikube -p secondary start \
  --memory=2048 --cpus=4 \
  --extra-config="apiserver.service-account-api-audiences=api" \
  --extra-config="apiserver.service-account-issuer=api" \
  --extra-config="apiserver.service-account-key-file=/var/lib/minikube/certs/sa.pub" \
  --extra-config="apiserver.service-account-signing-key-file=/var/lib/minikube/certs/sa.key"

minikube -p secondary ssh 'sudo cat /var/lib/minikube/certs/sa.pub'
minikube -p secondary dashboard --url=true &

echo -en "\033[1;32m Dashboard: http://localhost:$PORT/api/v1/namespaces/kubernetes-dashboard/services/http:kubernetes-dashboard:/proxy/# \033[0m \n"
echo -en "\033[1;32m Swagger: http://localhost:$PORT/swagger-ui \033[0m \n"
echo -en "\033[1;33m Context: `kubectl config current-context`\033[0m \n"
kubectl proxy --port=$PORT
