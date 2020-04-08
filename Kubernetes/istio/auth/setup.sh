#kubectl delete ns foo bar legacy

echo -en "\033[1;32m Setting Up Foo, Bar with Proxy \033[0m \n"
kubectl create ns foo
kubectl apply -f <(istioctl kube-inject -f httpbin.yaml) -n foo
kubectl apply -f <(istioctl kube-inject -f sleep.yaml) -n foo
kubectl create ns bar
kubectl apply -f <(istioctl kube-inject -f httpbin.yaml) -n bar
kubectl apply -f <(istioctl kube-inject -f sleep.yaml) -n bar
kubectl create ns legacy
kubectl apply -f httpbin.yaml -n legacy
kubectl apply -f sleep.yaml -n legacy