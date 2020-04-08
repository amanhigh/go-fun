istioctl manifest apply --set profile=demo
kubectl label namespace default istio-injection=enabled

istioctl dashboard kiali &