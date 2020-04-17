export ISTIOD_REMOTE_EP=$(kubectl --context minikube -n istio-system get svc istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
echo -en "\033[1;33m "ISTIOD_REMOTE_EP is ${ISTIOD_REMOTE_EP}"\033[0m \n"

echo -en "\033[1;32m Configuring Istio: Secondary \033[0m \n"
minikube profile secondary
istioctl manifest apply -f istio-secondary.yaml

#-- Dasboards
#istioctl dashboard kiali
