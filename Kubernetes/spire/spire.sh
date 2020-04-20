echo -en "\033[1;32m Setting Up Spire Server \033[0m \n"
kubectl -f ./spire-server.yaml apply
echo -en "\033[1;32m Setting Up Spire Agent \033[0m \n"
kubectl -f ./spire-agent.yaml apply

sleep 10
echo -en "\033[1;32m Mapping Agent Spiffe Id \033[0m \n"
kubectl exec -n spire spire-server-0 -- \
    /opt/spire/bin/spire-server entry create \
    -spiffeID spiffe://example.org/ns/spire/sa/spire-agent \
    -selector k8s_sat:cluster:demo-cluster \
    -selector k8s_sat:agent_ns:spire \
    -selector k8s_sat:agent_sa:spire-agent \
    -node

echo -en "\033[1;32m Mapping Workload Spiffe Id \033[0m \n"
kubectl exec -n spire spire-server-0 -- \
    /opt/spire/bin/spire-server entry create \
    -spiffeID spiffe://example.org/ns/foo/sa/my-sleep \
    -parentID spiffe://example.org/ns/spire/sa/spire-agent \
    -selector k8s:ns:foo \
    -selector k8s:sa:sleep

echo -en "\033[1;32m Registered Entries \033[0m \n"
kubectl exec -n spire spire-server-0 -- \
    /opt/spire/bin/spire-server entry show

echo -en "\033[1;32m Setting Up Test Client \033[0m \n"
kubectl create ns foo
kubectl apply -f sleep.yaml -n foo