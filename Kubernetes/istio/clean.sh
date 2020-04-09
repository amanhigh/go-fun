minikube profile minikube
istioctl manifest generate -f istio-primary.yaml | kubectl delete -f -
minikube profile secondary
istioctl manifest generate -f istio-secondary.yaml | kubectl delete -f -