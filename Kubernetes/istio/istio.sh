echo -en "\033[1;32m Configuring Istio: Primary \033[0m \n"
minikube profile minikube
istioctl manifest apply -f istio-primary.yaml
kubectl label namespace default istio-injection=enabled

echo -en "\033[1;32m Creating Certs \033[0m \n"
kubectl create secret generic cacerts -n istio-system \
  --from-file=./certs/ca-cert.pem --from-file=./certs/ca-key.pem \
  --from-file=./certs/root-cert.pem --from-file=./certs/cert-chain.pem

echo -en "\033[1;32m Create Kialia Admin Login \033[0m \n"
kubectl apply -f ./components/kiali.yaml

#istioctl manifest apply --set profile=demo
#-- Dasboards
#istioctl dashboard kiali
