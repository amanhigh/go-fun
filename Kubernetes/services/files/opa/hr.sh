echo -en "\033[1;34m Apply Demo Policy \033[0m \n"
curl -s -X PUT --data-binary @./example.rego opa-opa-kube-mgmt:8181/v1/policies/example > /dev/null
echo -en "\033[1;34m Policy: http://localhost:8181/v1/policies/example \033[0m \n"
echo -en "\033[1;34m Data: http://localhost:8181/v1/data/httpapi/authz \033[0m \n"
echo -en "\033[1;34m Data Metrics: http://localhost:8181/v1/data/httpapi/authz?metrics=true \033[0m \n"
sleep 2

echo -en "\033[1;32m Alice Checks Her Own Salary \033[0m \n"
curl --user alice:password localhost:5000/finance/salary/alice

echo -en "\033[1;32m Alice Manager (Bob) Checker Her Salary \033[0m \n"
curl --user bob:password localhost:5000/finance/salary/alice

echo -en "\033[1;33m Bob cannot check Charlie Salary (Not His Manager)\033[0m \n"
curl --user bob:password localhost:5000/finance/salary/charlie

echo -en "\033[1;33m HR Trying to access Alice Salary (Without Policy)\033[0m \n"
curl --user david:password localhost:5000/finance/salary/alice

echo -en "\033[1;34m Apply HR Policy \033[0m \n"
curl -s -X PUT --data-binary @./hr.rego opa-opa-kube-mgmt:8181/v1/policies/example-hr > /dev/null
sleep 2

echo -en "\033[1;32m HR Able to Check all Salaries \033[0m \n"
curl --user david:password localhost:5000/finance/salary/alice
curl --user david:password localhost:5000/finance/salary/bob
curl --user david:password localhost:5000/finance/salary/charlie
curl --user david:password localhost:5000/finance/salary/david

#-------------------------------------------------
echo -en "\033[1;32m OPA Eval (Http API): True \033[0m \n"
curl -X POST 'http://opa-opa-kube-mgmt:8181/v1/data/httpapi/authz/allow' \
--data '{
  "input": {
    "user": "alice",
    "path": ["finance", "salary", "alice"],
    "method": "GET"
  }
}'

echo -en "\033[1;32m \nIs Manager Check: True \033[0m \n"
curl -X POST 'http://opa-opa-kube-mgmt:8181/v1/data/httpapi/authz/is_manager' \
--data '{
  "input": {
    "user": "bob"
  }
}'

echo ""

#Cleanup
curl -s -X DELETE opa-opa-kube-mgmt:8181/v1/policies/example-hr > /dev/null