echo -en "\033[1;32m AuthZ: Server Mode \033[0m \n"
curl -s -X PUT --data-binary @./authz.rego http://opa-opa-kube-mgmt:8181/v1/policies/gofun/authz > /dev/null
curl -s -X PUT --data-binary @./authz.json http://opa-opa-kube-mgmt:8181/v1/data/gofun/authz > /dev/null

echo -en "\033[1;34m Policy: http://localhost:8181/v1/policies/gofun/authz \033[0m \n"
echo -en "\033[1;34m Data: http://localhost:8181/v1/data/gofun/authz \033[0m \n"

#FIXME: Allow Returning False
curl -X POST --data-binary @./input.json 'http://opa-opa-kube-mgmt:8181/v1/data/gofun/authz/allow'
echo ""