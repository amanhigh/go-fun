#Login Command -  kubectl exec -it $(kubectl get pod -l app=sleep -n foo -o jsonpath={.items..metadata.name}) -c sleep -n foo -- sh

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

echo -en "\033[1;32m Foo httpbin Requires JWT \033[0m \n"
kubectl apply -f - <<EOF
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: require-jwt
  namespace: foo
spec:
  selector:
    matchLabels:
      app: httpbin
  action: ALLOW
  rules:
  - from:
    - source:
       requestPrincipals: ["testing@secure.istio.io/testing@secure.istio.io"]
EOF
sleep 6

echo -en "\033[1;31m No JWT\033[0m \n"
kubectl exec $(kubectl get pod -l app=sleep -n foo -o jsonpath={.items..metadata.name}) -c sleep -n foo -- curl "http://httpbin.foo:8000/headers" -s -o /dev/null -w "%{http_code}\n"

TOKEN=$(curl https://raw.githubusercontent.com/istio/istio/release-1.5/security/tools/jwt/samples/demo.jwt -s)

echo -en "\033[1;32m Valid JWT Token \033[0m \n"
echo "Token: $(echo $TOKEN | cut -d '.' -f2 - | base64 -D -)"
kubectl exec $(kubectl get pod -l app=sleep -n foo -o jsonpath={.items..metadata.name}) -c sleep -n foo -- curl "http://httpbin.foo:8000/headers" -s -o /dev/null -H "Authorization: Bearer $TOKEN" -w "%{http_code}\n"

kubectl -n foo delete AuthorizationPolicy require-jwt


echo -en "\033[1;32m JWT Group Claim: Foo(httpbin) \033[0m \n"
kubectl apply -f - <<EOF
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: require-jwt-group
  namespace: foo
spec:
  selector:
    matchLabels:
      app: httpbin
  action: ALLOW
  rules:
  - from:
    - source:
       requestPrincipals: ["testing@secure.istio.io/testing@secure.istio.io"]
    when:
    - key: request.auth.claims[groups]
      values: ["group1"]
EOF

echo -en "\033[1;31m No JWT\033[0m \n"
kubectl exec $(kubectl get pod -l app=sleep -n foo -o jsonpath={.items..metadata.name}) -c sleep -n foo -- curl "http://httpbin.foo:8000/headers" -s -o /dev/null -w "%{http_code}\n"

TOKEN=$(curl https://raw.githubusercontent.com/istio/istio/release-1.5/security/tools/jwt/samples/groups-scope.jwt -s)

echo -en "\033[1;32m Valid JWT Group Token \033[0m \n"
echo "Token: $(echo $TOKEN | cut -d '.' -f2 - | base64 -D -)"
kubectl exec $(kubectl get pod -l app=sleep -n foo -o jsonpath={.items..metadata.name}) -c sleep -n foo -- curl "http://httpbin.foo:8000/headers" -s -o /dev/null -H "Authorization: Bearer $TOKEN" -w "%{http_code}\n"

kubectl -n foo delete AuthorizationPolicy require-jwt-group

echo -en "\033[1;32m Accessing K8 Service \033[0m \n"

echo -en "\033[1;33m Without K8 Token\033[0m \n"
kubectl exec $(kubectl get pod -l app=sleep -n foo -o jsonpath={.items..metadata.name}) -c sleep -n foo -- sh -c 'curl -sSk https://$KUBERNETES_SERVICE_HOST:$KUBERNETES_PORT_443_TCP_PORT/api/v1/namespaces/foo/pods/$HOSTNAME -o /dev/null -w "%{http_code}\n"'


echo -en "\033[1;33m With K8 Token\033[0m \n"
kubectl exec $(kubectl get pod -l app=sleep -n foo -o jsonpath={.items..metadata.name}) -c sleep -n foo -- sh -c 'curl -sSk https://$KUBERNETES_SERVICE_HOST:$KUBERNETES_PORT_443_TCP_PORT/api/v1/namespaces/foo/pods/$HOSTNAME -H "Authorization: Bearer $(cat /var/run/secrets/kubernetes.io/serviceaccount/token)" -o /dev/null -w "%{http_code}\n"'

