istioctl manifest apply --set profile=demo
kubectl label namespace default istio-injection=enabled

sleep 10
istioctl dashboard kiali &