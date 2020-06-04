echo -en "\033[1;34m Apply Demo Policy \033[0m \n"
curl -s -X PUT --data-binary @../opa/example.rego docker:8181/v1/policies/example > /dev/null
echo -en "\033[1;34m Policy: http://docker:8181/v1/policies/example \033[0m \n"
echo -en "\033[1;34m Data: http://docker:8181/v1/data/httpapi/authz \033[0m \n"
echo -en "\033[1;34m Data Metrics: http://docker:8181/v1/data/httpapi/authz?metrics=true \033[0m \n"
sleep 2

echo -en "\033[1;32m Alice Checks Her Own Salary \033[0m \n"
curl --user alice:password docker:5000/finance/salary/alice

echo -en "\033[1;32m Alice Manager (Bob) Checker Her Salary \033[0m \n"
curl --user bob:password docker:5000/finance/salary/alice

echo -en "\033[1;33m Bob cannot check Charlie Salary (Not His Manager)\033[0m \n"
curl --user bob:password docker:5000/finance/salary/charlie

echo -en "\033[1;33m HR Trying to access Alice Salary (Without Policy)\033[0m \n"
curl --user david:password docker:5000/finance/salary/alice

echo -en "\033[1;34m Apply HR Policy \033[0m \n"
curl -s -X PUT --data-binary @../opa/hr.rego docker:8181/v1/policies/example-hr > /dev/null
sleep 2

echo -en "\033[1;32m HR Able to Check all Salaries \033[0m \n"
curl --user david:password docker:5000/finance/salary/alice
curl --user david:password docker:5000/finance/salary/bob
curl --user david:password docker:5000/finance/salary/charlie
curl --user david:password docker:5000/finance/salary/david

#------------------------------------------------
echo -en "\033[1;32m OPA EVAL (Command Line) \033[0m \n"
docker run -it openpolicyagent/opa eval '1*2+3'

echo -en "\033[1;32m Server Policy Evaluation  \033[0m \n"
docker run -it -v $PWD/../opa/:/inputs openpolicyagent/opa eval -i /inputs/input.json -d /inputs/authz.rego -d /inputs/authz.json 'data.gofun.authz' -f pretty --profile

echo -en "\033[1;32m Policy Testing \033[0m \n"
docker run -it -v $PWD/../opa/:/inputs openpolicyagent/opa test /inputs -v

#-------------------------------------------------
echo -en "\033[1;32m OPA Eval (Http API) \033[0m \n"
curl -X POST 'http://docker:8181/v1/data/httpapi/authz/allow' \
--data-raw '{
  "input": {
    "user": "alice",
    "path": ["finance", "salary", "alice"],
    "method": "GET"
  }
}'

echo -en "\033[1;32m \nIs Manager Check \033[0m \n"
curl -X POST 'http://docker:8181/v1/data/httpapi/authz/is_manager' \
--data-raw '{
  "input": {
    "user": "bob"
  }
}'

echo ""

#Cleanup
curl -s -X DELETE docker:8181/v1/policies/example-hr > /dev/null
