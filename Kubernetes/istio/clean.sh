minikube profile minikube
kubectl delete -f ./components/*
istioctl manifest generate -f istio-primary.yaml | kubectl delete -f -
minikube profile secondary
istioctl manifest generate -f istio-secondary.yaml | kubectl delete -f -