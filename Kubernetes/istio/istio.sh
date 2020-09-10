echo -en "\033[1;32m Configuring Istio: Primary \033[0m \n"
minikube profile minikube
istioctl install -f istio-primary.yaml
#istioctl install --set profile=demo
kubectl label namespace default istio-injection=enabled

echo -en "\033[1;32m Creating Certs \033[0m \n"
kubectl create secret generic cacerts -n istio-system \
  --from-file=./certs/ca-cert.pem --from-file=./certs/ca-key.pem \
  --from-file=./certs/root-cert.pem --from-file=./certs/cert-chain.pem

echo -en "\033[1;32m Create Kialia Admin Login \033[0m \n"
kubectl apply -f ./components/kiali.yaml
kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.7/samples/addons/kiali.yaml

#istioctl manifest apply --set profile=demo
#-- Dasboards
istioctl dashboard kiali

echo -en "\033[1;32m Gateway Filter Chain Check CMD \033[0m \n"
echo 'istioctl proxy-config listeners -n istio-system $(kubectl get pod -l app=istio-ingressgateway -n istio-system -o jsonpath={.items..metadata.name}) --port 80 -o json | jq ".[0].filterChains[0].filters"'