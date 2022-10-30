# helm repo add onechart https://chart.onechart.dev
# helm repo add stakater https://stakater.github.io/stakater-charts

# helm repo update
# sudo kubefwd svc

# Vars
CMD="install"
ANS_FILE=/tmp/k8-svc.txt

function process()
{

    for SVC in `cat $ANS_FILE`
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
            
            echo -en "\033[1;32m Squid \033[0m \n"
            helm $CMD squid onechart/onechart -f squid.yml > /dev/null
            echo -en "\033[1;33m ALLOW: curl -x localhost:3128 www.google.com \033[0m \n"
            echo -en "\033[1;33m DENY: curl -x localhost:3128 www.fb.com \033[0m \n"

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
    
        APP)
            helm $CMD app onechart/onechart -f app.yml > /dev/null
            echo -en "\033[1;33m http://localhost:7080/metrics\n http://localhost:7080/person/all \033[0m \n"
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

        ETCD)
            helm $CMD etcd bitnami/etcd -f etcd.yml > /dev/null
            echo -en "\033[1;33m ./demo/demo.sh \033[0m \n"
            ;;

        SONAR)
            helm $CMD sonar bitnami/sonarqube -f sonar.yml > /dev/null
            echo -en "\033[1;33m http://localhost:9000/ \033[0m \n"
            echo -en "\033[1;33m Login: aman/aman (Need 5GB Mem) \033[0m \n"
            ;;

        ZOOKEEPER)
            helm $CMD zookeeper bitnami/zookeeper -f zookeeper.yml > /dev/null
            echo -en "\033[1;33m /demo/demo.sh \033[0m \n"
            ;;
        
        ELK)
            #TODO: Complete Compose Setup
            helm $CMD logstash bitnami/logstash -f logstash.yml > /dev/null
            helm $CMD elasticsearch bitnami/elasticsearch -f elasticsearch.yml > /dev/null
            # helm $CMD kibana bitnami/kibana -f kibana.yml > /dev/null
            echo -en "\033[1;33m /demo/demo.sh \033[0m \n"
            ;;

        MONITOR)
            #helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
            helm $CMD prometheus prometheus-community/prometheus -f prometheus.yml > /dev/null
            echo -en "\033[1;33m Prometheus Server: http://localhost:9090/ \033[0m \n"

             # helm repo add grafana https://grafana.github.io/helm-charts
            helm $CMD grafana grafana/grafana -f grafana.yml > /dev/null
            echo -en "\033[1;33m http://localhost:3000/login (aman/aman) \033[0m \n"
            echo -en "\033[1;33m Datasource: http://localhost:3000/datasources/new \033[0m \n"
            echo -en "\033[1;33m Add Datasource Prometheus: http://prometheus-server \033[0m \n"
            ;;
        
        LDAP)
            helm $CMD ldap onechart/onechart -f ldap.yml > /dev/null
            helm $CMD ldap-admin onechart/onechart -f ldap-admin.yml > /dev/null
            echo -en "\033[1;33m UI: http://localhost:8030/ \033[0m \n"
            echo -en "\033[1;33m CMD: ldapsearch -H ldap://localhost:3891 -xLL -D 'cn=admin,dc=example,dc=com' -b 'dc=example,dc=com' -W '(cn=admin)' \033[0m \n"
            echo -en "\033[1;33m Admin Login: Username:cn=admin,dc=example,dc=com Password: admin \033[0m \n"
            ;;
        WEBSHELL)
            helm $CMD sshwifty onechart/onechart -f sshwifty.yml  > /dev/null
            helm $CMD webssh onechart/onechart -f webssh.yml > /dev/null
            #TODO: Fix Config
            helm $CMD webssh2 onechart/onechart -f webssh2.yml > /dev/null
            echo -en "\033[1;33m Sshwifty: http://localhost:8080/ \033[0m \n"
            echo -en "\033[1;33m Webssh: http://localhost:8182/ \033[0m \n"
            echo -en "\033[1;33m Webssh: http://localhost:2222/ \033[0m \n"
            ;;

        *)
            echo -en "\033[1;34m Service Not Supported: $SVC \033[0m \n"
            ;;
        esac
    done
}


# Flags
while getopts 'dusi' OPTION; do
  case "$OPTION" in
    d)
        echo -en "\033[1;32m Clearing all Helms \033[0m \n"
        helm delete $(helm list --short)
        ;;
    i)
        echo -en "\033[1;32m Helm: Install \033[0m \n"
        process
        ;;
    u)
        echo -en "\033[1;32m Helm: Upgrade \033[0m \n"
        process
        CMD="upgrade"
        ;;
    s)
        # Prompt
        answers=`gum choose MYSQL MONGO REDIS APP PROXY LOADER CRON HTTPBIN VAULT OPA CONSUL LDAP ETCD SONAR ZOOKEEPER ELK MONITOR WEBSHELL --limit 5`
        echo $answers > $ANS_FILE    
        echo -en "\033[1;32m Service Set \033[0m \n"
        echo -en "\033[1;33m $answers \033[0m \n"
        ;;
    ?)
        echo -en "\033[1;32m script usage: $0 [-s] [-i] [-u] [-d] \033[0m \n"
        echo -en "\033[1;33m [-s] Set \033[0m \n"
        echo -en "\033[1;33m [-i] Install \033[0m \n"
        echo -en "\033[1;33m [-u] Upgrade \033[0m \n"
        echo -en "\033[1;33m [-d] delete \033[0m \n"
        exit 1
        ;;
  esac
done