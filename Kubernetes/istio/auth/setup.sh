kubectl create ns foo
kubectl apply -f <(istioctl kube-inject -f httpbin.yaml) -n foo
kubectl apply -f <(istioctl kube-inject -f sleep.yaml) -n foo
kubectl create ns bar
kubectl apply -f <(istioctl kube-inject -f httpbin.yaml) -n bar
kubectl apply -f <(istioctl kube-inject -f sleep.yaml) -n bar
kubectl create ns legacy
kubectl apply -f httpbin.yaml -n legacy
kubectl apply -f sleep.yaml -n legacy