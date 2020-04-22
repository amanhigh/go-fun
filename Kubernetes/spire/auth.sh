JWT_PATH=/tmp
echo -en "\033[1;32m Generating Tokens (x509 SVID) \033[0m \n"
kubectl exec $(kubectl get pod -l app=sleep -n foo -o jsonpath={.items..metadata.name}) \
  -c sleep -n foo -- /opt/spire/spire-agent api fetch -socketPath /run/spire/sockets/agent.sock

echo -en "\033[1;32m Generating Tokens (JWT SVID) \033[0m \n"
kubectl exec $(kubectl get pod -l app=sleep -n foo -o jsonpath={.items..metadata.name}) \
  -c sleep -n foo -- sh -c "/opt/spire/spire-agent api fetch jwt -audience spiffe://example.org/ns/foo/sa/my-httpbin -socketPath /run/spire/sockets/agent.sock > $JWT_PATH/jwt;
cat $JWT_PATH/jwt;"

echo -en "\033[1;32m Decoding JWT Token \033[0m \n"
JWT_TOKEN=$(kubectl exec $(kubectl get pod -l app=sleep -n foo -o jsonpath={.items..metadata.name}) \
  -c sleep -n foo -- head -2 $JWT_PATH/jwt | tail -1)
jq -R 'split(".") | .[1] | @base64d | fromjson' <<< "$JWT_TOKEN"

echo -en "\033[1;33m Validating Token (Invalid Audience) \033[0m \n"
kubectl exec $(kubectl get pod -l app=sleep -n foo -o jsonpath={.items..metadata.name}) \
  -c sleep -n foo -- /opt/spire/spire-agent api validate jwt -audience my-httpbin -socketPath /run/spire/sockets/agent.sock -svid $JWT_TOKEN

echo -en "\033[1;32m Validating Token (Valid Audience) \033[0m \n"
kubectl exec $(kubectl get pod -l app=sleep -n foo -o jsonpath={.items..metadata.name}) \
  -c sleep -n foo -- /opt/spire/spire-agent api validate jwt -audience spiffe://example.org/ns/foo/sa/my-httpbin -socketPath /run/spire/sockets/agent.sock -svid $JWT_TOKEN
