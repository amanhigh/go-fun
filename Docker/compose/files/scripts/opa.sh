echo -en "\033[1;34m Apply Demo Policy \033[0m \n"
curl -s -X PUT --data-binary @../opa/example.rego docker:8181/v1/policies/example > /dev/null
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

curl -s -X DELETE docker:8181/v1/policies/example-hr > /dev/null

#------------------------------------------------
echo -en "\033[1;32m OPA EVAL (Command Line) \033[0m \n"
docker run -it openpolicyagent/opa:0.11.0 eval '1*2+3'

echo -en "\033[1;32m Server Policy Evaluation  \033[0m \n"
docker run -it -v $PWD/../opa/:/inputs openpolicyagent/opa:0.11.0 eval -i /inputs/input.json -d /inputs/server.rego 'data.example.violation[x]'
