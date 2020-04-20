echo -en "\033[1;32m Generating Tokens (x509) \033[0m \n"
kubectl exec $(kubectl get pod -l app=sleep -n foo -o jsonpath={.items..metadata.name}) \
-c sleep -n foo -- /opt/spire/spire-agent api fetch -socketPath /run/spire/sockets/agent.sock

echo -en "\033[1;32m Generating Tokens (JWT) \033[0m \n"
kubectl exec $(kubectl get pod -l app=sleep -n foo -o jsonpath={.items..metadata.name}) \
-c sleep -n foo -- /opt/spire/spire-agent api fetch jwt -audience spiffe://example.org/ns/foo/sa/my-httpbin -socketPath /run/spire/sockets/agent.sock
