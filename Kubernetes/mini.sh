minikube delete;
minikube start \
--extra-config="apiserver.enable-swagger-ui=true" \
--extra-config="apiserver.service-account-api-audiences=api" \
--extra-config="apiserver.service-account-issuer=api" \
--extra-config="apiserver.service-account-key-file=/var/lib/minikube/certs/sa.pub" \
--extra-config="apiserver.service-account-signing-key-file=/var/lib/minikube/certs/sa.key";

minikube dashboard &
kubectl proxy --port=8091 &
echo -en "\033[1;32m Swagger: http://localhost:8091/swagger-ui \033[0m \n"
