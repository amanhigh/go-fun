kubectl apply -f bookinfo.yaml
kubectl apply -f bookinfo-gateway.yaml

export INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].nodePort}')
export SECURE_INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="https")].nodePort}')
export INGRESS_HOST=$(minikube ip)
export GATEWAY_URL=$INGRESS_HOST:$INGRESS_PORT

echo -en "\033[1;32m http://localhost:8091/api/v1/namespaces/istio-system/services/istio-ingressgateway:80/proxy/productpage \033[0m \n"
echo -en "\033[1;32m http://$GATEWAY_URL/productpage (ELB)\033[0m \n"
