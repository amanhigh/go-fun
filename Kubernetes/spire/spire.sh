echo -en "\033[1;32m Setting Up Spire Server \033[0m \n"
kubectl -f ./spire-server.yaml apply
echo -en "\033[1;32m Setting Up Spire Agent \033[0m \n"
kubectl -f ./spire-agent.yaml apply

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
    -spiffeID spiffe://example.org/ns/default/sa/default \
    -parentID spiffe://example.org/ns/spire/sa/spire-agent \
    -selector k8s:ns:default \
    -selector k8s:sa:default

echo -en "\033[1;32m Registered Entries \033[0m \n"
kubectl exec -n spire spire-server-0 -- \
    /opt/spire/bin/spire-server entry show