echo -en "\033[1;32m Configuring Istio: Primary \033[0m \n"
minikube profile minikube
# https://istio.io/v1.1/docs/setup/kubernetes/additional-setup/config-profiles/
istioctl install --set profile=default

echo -en "\033[1;32m Enabled Istio for Default Namespace \n"
kubectl label namespace default istio-injection=enabled

echo -en "\033[1;32m Setting Up Kiali \033[0m \n"
# sleep 10
# kubectl apply -f ./components/kiali.yaml
# helm install --namespace istio-system --set auth.strategy="anonymous" --repo https://kiali.org/helm-charts kiali-server kiali-server


#-- Dasboards
# istioctl dashboard kiali

# kubectl port-forward svc/kiali 20001:20001 -n istio-system
# echo -en "\033[1;32m Kiali: https://localhost:20001/ \033[0m \n"