minikube profile minikube
# Remove Istio Objects
istioctl x uninstall --purge
# Delete Istio Namespaces
kubectl delete ns istio-system