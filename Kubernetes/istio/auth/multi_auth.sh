minikube profile primary
#Auth Headers
echo -en "\033[1;32m Auth Headers: SPIFFE (primary.foo.sleep to primary.foo.httpbin) \033[0m \n"
kubectl exec $(kubectl get pod -l app=sleep -n foo -o jsonpath={.items..metadata.name}) -c sleep -n foo -- curl http://httpbin.foo:8000/headers -s #| grep X-Forwarded-Client-Cert

minikube profile secondary
#Auth Headers
echo -en "\033[1;32m Auth Headers: SPIFFE (secondary.bar.sleep to secondary.bar.httpbin) \033[0m \n"
kubectl exec $(kubectl get pod -l app=sleep -n bar -o jsonpath={.items..metadata.name}) -c sleep -n bar -- curl http://httpbin.bar:8000/headers -s #| grep X-Forwarded-Client-Cert