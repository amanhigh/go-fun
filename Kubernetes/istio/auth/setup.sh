#kubectl delete ns foo bar legacy
minikube profile minikube

echo -en "\033[1;32m Setting Up Foo, Bar with Proxy \033[0m \n"
kubectl create ns foo
kubectl label namespace foo istio-injection=enabled
kubectl apply -f httpbin.yaml -n foo
kubectl apply -f sleep.yaml -n foo

kubectl create ns bar
kubectl label namespace bar istio-injection=enabled
kubectl apply -f httpbin.yaml -n bar
kubectl apply -f sleep.yaml -n bar

echo -en "\033[1;32m Setting Up Legacy Without Proxy \033[0m \n"
kubectl create ns legacy
kubectl apply -f httpbin.yaml -n legacy
kubectl apply -f sleep.yaml -n legacy


echo -en "\033[1;32m Setting Up Policies \033[0m \n"
kubectl apply -f foo-get-pod.yml