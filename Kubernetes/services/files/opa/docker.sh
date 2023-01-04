# Docker OPA Eval (Run in current Directory)
#------------------------------------------------
echo ""
echo -en "\033[1;32m AuthZ: OPA Eval \033[0m \n"
docker run -it openpolicyagent/opa eval '1*2+3'

echo -en "\033[1;32m Server Policy Evaluation  \033[0m \n"
docker run -it -v `pwd`:/inputs openpolicyagent/opa eval -i /inputs/input.json -d /inputs/authz.rego -d /inputs/authz.json 'data.gofun.authz' -f pretty --profile

echo -en "\033[1;32m Policy Testing \033[0m \n"
docker run -it -v `pwd`:/inputs openpolicyagent/opa test /inputs -v