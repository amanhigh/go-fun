PORT=8091
#Use minikube config set vm-driver virtualbox/docker
minikube -p minikube delete;
minikube  -p minikube start \
--memory=3096 --cpus=2 \
# TODO: Fix on WSL
#--host-only-cidr='24.1.1.100/24' \
--cache-images=true \
--extra-config="apiserver.enable-swagger-ui=true" \
--extra-config="apiserver.service-account-api-audiences=api" \
--extra-config="apiserver.service-account-issuer=api" \
--extra-config="apiserver.service-account-key-file=/var/lib/minikube/certs/sa.pub" \
--extra-config="apiserver.service-account-signing-key-file=/var/lib/minikube/certs/sa.key";

minikube -p minikube ssh 'sudo cat /var/lib/minikube/certs/sa.pub'
minikube -p minikube dashboard --url=true &

minikube addons enable metrics-server

echo -en "\033[1;33m Run 'minikube tunnel’ for Emulating ELB\033[0m \n"
echo -en "\033[1;32m Dashboard: http://localhost:$PORT/api/v1/namespaces/kubernetes-dashboard/services/http:kubernetes-dashboard:/proxy/# \033[0m \n"
echo -en "\033[1;32m Swagger: http://localhost:$PORT/swagger-ui \033[0m \n"
echo -en "\033[1;33m Context: `kubectl config current-context; k9s --context minikube`\033[0m \n"
kubectl proxy --port=$PORT


## k9s
# k9s --readonly , -n <namespace>, -l <loglevel>
# k9s :pu (pulse), :dp (deployments), :po (pods), :ns (namespace),:rb (Role Bindings) ,:a (aliases)
# Switch Namespace :po <namespace> to see pods of that namespace

# Helpful Commands
# Docker Start: sudo service docker start
# Docker Image to Minikube: eval $(minikube docker-env); docker build -t fun-app .
# Port Forward (Local Port 9090 to Container Port 8080) - kubectl port-forward `kubectl get pods -o name | grep fun-app | head  -1` 9090:8080
# Logs - kubectl logs `kubectl get pods -o name | grep fun-app | head  -1` -f
# Login - kubectl -it exec `kubectl get pods -o name | grep fun-app | head  -1` bash
# Delete All - kubectl delete all --all
# Tunnel (Emulate Load Balancer) - minikube tunnel
# List Emulated Services URL's - minikube service list