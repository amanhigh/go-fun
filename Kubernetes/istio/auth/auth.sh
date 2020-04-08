#Run All Routes
echo -en "\033[1;32m All Routes \033[0m \n"
for from in "foo" "bar" "legacy"; do for to in "foo" "bar" "legacy"; do kubectl exec $(kubectl get pod -l app=sleep -n ${from} -o jsonpath={.items..metadata.name}) -c sleep -n ${from} -- curl "http://httpbin.${to}:8000/ip" -s -o /dev/null -w "sleep.${from} to httpbin.${to}: %{http_code}\n"; done; done

#Auth Headers
echo -en "\033[1;32m Auth Headers: SPIFFE (foo.sleep to foo.httpbin) \033[0m \n"
kubectl exec $(kubectl get pod -l app=sleep -n foo -o jsonpath={.items..metadata.name}) -c sleep -n foo -- curl http://httpbin.foo:8000/headers -s | grep X-Forwarded-Client-Cert

#Enable Global MTLS
echo -en "\033[1;32m Enabling MTLS Strict \033[0m \n"
kubectl apply -f - <<EOF
apiVersion: "security.istio.io/v1beta1"
kind: "PeerAuthentication"
metadata:
  name: "default"
  namespace: "istio-system"
spec:
  mtls:
    mode: STRICT
EOF

echo -en "\033[1;32m All Routes (STRICT) \033[0m \n"
for from in "foo" "bar" "legacy"; do for to in "foo" "bar" "legacy"; do kubectl exec $(kubectl get pod -l app=sleep -n ${from} -o jsonpath={.items..metadata.name}) -c sleep -n ${from} -- curl "http://httpbin.${to}:8000/ip" -s -o /dev/null -w "sleep.${from} to httpbin.${to}: %{http_code}\n"; done; done
kubectl -n istio-system delete PeerAuthentication default

#Enable Principle Auth
echo -en "\033[1;32m Auth Foo allows only Foo \033[0m \n"
kubectl apply -f - <<EOF
apiVersion: "security.istio.io/v1beta1"
kind: "AuthorizationPolicy"
metadata:
  name: "foo-only"
  namespace: foo
spec:
  rules:
  - from:
    - source:
        principals: ["cluster.local/ns/foo/sa/sleep"]
    to:
    - operation:
        methods: ["GET"]
EOF
sleep 2

echo -en "\033[1;32m All Routes (Foo Only) \033[0m \n"
for from in "foo" "bar" "legacy"; do for to in "foo" "bar" "legacy"; do kubectl exec $(kubectl get pod -l app=sleep -n ${from} -o jsonpath={.items..metadata.name}) -c sleep -n ${from} -- curl "http://httpbin.${to}:8000/ip" -s -o /dev/null -w "sleep.${from} to httpbin.${to}: %{http_code}\n"; done; done
kubectl -n foo delete AuthorizationPolicy foo-only

