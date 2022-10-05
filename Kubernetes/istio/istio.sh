echo -en "\033[1;32m Configuring Istio: Primary \033[0m \n"
minikube profile minikube
# https://istio.io/v1.1/docs/setup/kubernetes/additional-setup/config-profiles/
istioctl install -y --set profile=default > /dev/null

echo -en "\033[1;32m Enabled Istio for Default Namespace \033[0m \n"
kubectl label namespace default istio-injection=enabled

echo -en "\033[1;32m Installing Addons \033[0m \n"
kubectl apply -f $ISTIO_HOME/samples/addons > /dev/null


echo -en "\033[1;32m Istio Setup \033[0m \n"
echo -en "\033[1;33m istioctl dashboard kiali \n"
echo -en "\033[1;33m Kiali: http://localhost:8091/api/v1/namespaces/istio-system/services/kiali:20001/proxy/ \033[0m \n"
echo -en "\033[1;33m Grafana: http://localhost:8091/api/v1/namespaces/istio-system/services/grafana:3000/proxy/ \033[0m \n"
echo -en "\033[1;33m Jaeger: http://localhost:8091/api/v1/namespaces/istio-system/services/tracing:80/proxy/ \033[0m \n"

# Helm setup doesn't auto configure with istio etc hence not using for now.
# echo -en "\033[1;32m Setting Up Kiali \033[0m \n"
# helm install --set cr.create=true --set cr.namespace=istio-system \
# --namespace kiali-operator --create-namespace kiali-operator kiali/kiali-operator
# 

# kubectl port-forward svc/kiali 20001:20001 -n istio-system
# echo -en "\033[1;32m Kiali: https://localhost:20001/ \033[0m \n"

# HELM
# Check Deployment: helm list -n kiali-operator
# Check Configurations: helm show values kiali/kiali-operator
