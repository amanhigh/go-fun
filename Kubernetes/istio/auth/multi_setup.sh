echo -en "\033[1;32m Setting Up Foo(Minikube Cluster) with Proxy \033[0m \n"
kubectl --context minikube create ns foo
kubectl --context minikube label namespace foo istio-injection=enabled
kubectl --context minikube apply -f httpbin.yaml -n foo
kubectl --context minikube apply -f sleep.yaml -n foo

echo -en "\033[1;32m Setting Up Bar(Secondary Cluster) with Proxy \033[0m \n"
kubectl --context secondary create ns bar
kubectl --context secondary label namespace bar istio-injection=enabled
kubectl --context secondary apply -f httpbin.yaml -n bar
kubectl --context secondary apply -f sleep.yaml -n bar