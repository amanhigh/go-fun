# helm repo add onechart https://chart.onechart.dev
# helm repo add stakater https://stakater.github.io/stakater-charts

# helm repo update

# Vars
CMD="install"

# Prompt
answers=`gum choose MYSQL MONGO REDIS PROXY LOADER CRON HTTPBIN VAULT OPA CONSUL --limit 5`

# Flags
while getopts 'du' OPTION; do
  case "$OPTION" in
    d)
        echo -en "\033[1;32m Clearing all Helms \033[0m \n"
        helm delete $(helm list --short)
        ;;
    u)
        echo -en "\033[1;32m Switching to Upgrade \033[0m \n"
        CMD="upgrade"
        ;;
    ?)
      echo "script usage: $0 [-u] [-r]" >&2
      exit 1
      ;;
  esac
done

for SVC in $answers
do
    echo -en "\033[1;32m \n $SVC \033[0m \n"
    case $SVC in
    MYSQL)
        helm $CMD mysql bitnami/mysql -f mysql.yml > /dev/null
        helm $CMD mysql-admin bitnami/phpmyadmin > /dev/null
        echo -en "\033[1;33m Login: mysql-primary, root/root \033[0m \n"
        ;;

    MONGO)
        helm $CMD mongo bitnami/mongodb -f mongo.yml > /dev/null
        echo -en "\033[1;33m mongosh -u root -p root --host localhost  < /etc/files/scripts/mongo.js \033[0m \n"
        echo -en "\033[1;33m Svc Endpoint: mongo-mongodb:27017 (Standalone Mode Only) \033[0m \n"
        ;;

    REDIS)
        helm $CMD redis bitnami/redis -f redis.yml > /dev/null
        helm $CMD redis-admin onechart/onechart -f redis-admin.yml

        echo -en "\033[1;33m redis-cli -c incr mycounter \033[0m \n"
        echo -en "\033[1;33m redis-cli -c set mypasswd lol \033[0m \n"
        echo -en "\033[1;33m redis-cli -c get mypasswd \033[0m \n"
        echo -en "\033[1;33m Commander: http://localhost:8081/ \033[0m \n"
        ;;

    PROXY)
        echo -en "\033[1;32m Nginx \033[0m \n"
        helm $CMD nginx bitnami/nginx -f nginx.yml > /dev/null

        echo -en "\033[1;33m http://localhost:8081/ \033[0m \n"
        echo -en "\033[1;33m Refer Server Blocks for More Help.. \033[0m \n"
        
        echo -en "\033[1;32m Resty \033[0m \n"
        helm $CMD resty onechart/onechart -f resty.yml > /dev/null

        echo -en "\033[1;33m http://localhost:8090/ \033[0m \n"
        echo -en "\033[1;33m http://localhost:8090/example \033[0m \n"
        echo -en "\033[1;33m http://localhost:8090/ndtv \033[0m \n"
        
        echo -en "\033[1;32m TinyProxy \033[0m \n"
        helm $CMD tinyproxy stakater/application -f tinyproxy.yml > /dev/null

        echo -en "\033[1;33m curl -x localhost:8888 tinyproxy.stats \033[0m \n"
        ;;

    LOADER)
        helm $CMD vegeta onechart/onechart -f vegeta.yml > /dev/null
        echo -en "\033[1;33m Check Logs for Output \033[0m \n"
        echo -en "\033[1;33m echo 'GET http://nginx' | vegeta attack | vegeta report \033[0m \n"
        ;;

    HTTPBIN)
        helm $CMD httpbin onechart/onechart -f httpbin.yml > /dev/null
        echo -en "\033[1;33m Swagger: http://localhost:8810 \033[0m \n"
        echo -en "\033[1;33m http://localhost:8810/anything \033[0m \n"
        ;;

    CRON)
        # TODO: Fix Cron
        helm $CMD cron onechart/onechart -f cron.yml > /dev/null
        echo -en "\033[1;33m Check Logs for Output \033[0m \n"
        ;;

    OPA)
        # helm repo add opa https://open-policy-agent.github.io/kube-mgmt/charts
        helm $CMD opa opa/opa-kube-mgmt -f opa.yml > /dev/null
        helm $CMD opa-demo onechart/onechart -f opa-demo.yml > /dev/null
        
        echo -en "\033[1;33m http://localhost:8181/ \033[0m \n"

        echo -en "\033[1;33m Demo: ./demo/hr.sh \033[0m \n"
        echo -en "\033[1;33m Demo: ./demo/authz.sh \033[0m \n"
        echo -en "\033[1;33m Localhost: ./demo/docker.sh \033[0m \n"
        ;;

    VAULT)
        # helm repo add hashicorp https://helm.releases.hashicorp.com
        helm $CMD vault hashicorp/vault -f vault.yml > /dev/null
        echo -en "\033[1;33m vault status \033[0m \n"
        echo -en "\033[1;33m /demo/vault.sh \033[0m \n"
        ;;

    CONSUL)
        # helm repo add hashicorp https://helm.releases.hashicorp.com
        helm $CMD consul hashicorp/consul -f consul.yml > /dev/null
        echo -en "\033[1;33m http://localhost:8500/ \033[0m \n"
        ;;

    *)
        echo -en "\033[1;34m Service Not Supported: $SVC \033[0m \n"
        ;;
    esac
done
