echo -en "\033[1;32m Configuring Istio: Primary \033[0m \n"
minikube profile minikube
istioctl manifest apply -f istio-primary.yaml
kubectl label namespace default istio-injection=enabled

#istioctl manifest apply --set profile=demo
#-- Dasboards
#istioctl dashboard kiali