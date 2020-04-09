export MAIN_CLUSTER_NAME=minikube
export REMOTE_CLUSTER_NAME=secondary

export MAIN_CLUSTER_NETWORK=minikube-network
export REMOTE_CLUSTER_NETWORK=secondary-network

echo -en "\033[1;32m Configuring Istio: Primary \033[0m \n"
minikube profile minikube
istioctl manifest apply --set profile=default
istioctl manifest apply -f istio-primary.yaml --set values.global.mtls.enabled=true
kubectl label namespace default istio-injection=enabled

echo -en "\033[1;32m Creating Certs \033[0m \n"
kubectl create secret generic cacerts -n istio-system \
  --from-file=./certs/ca-cert.pem --from-file=./certs/ca-key.pem \
  --from-file=./certs/root-cert.pem --from-file=./certs/cert-chain.pem

#istioctl manifest apply --set profile=demo
#-- Dasboards
#istioctl dashboard kiali
