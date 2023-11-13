#!/bin/zsh
# helm repo add onechart https://chart.onechart.dev
# helm repo add stakater https://stakater.github.io/stakater-charts
# helm repo add bitnami https://charts.bitnami.com/bitnami

# helm repo update
# sudo kubefwd svc | awk '{ if($2 ~ /Port-Forward/) {print $0" URL: http://"$4"/"} else {print}}'

# Vars
CMD="install"
ANS_FILE=/tmp/k8-svc.txt

process()
{
    SVC_LIST=${1:-`cat $ANS_FILE`}

    for SVC in $( echo "$SVC_LIST" )
    do
        echo "\033[1;32m \n $SVC \033[0m \n"
        case $SVC in
        MYSQL)
            helm $CMD mysql bitnami/mysql -f mysql.yml > /dev/null
            helm $CMD mysql-admin bitnami/phpmyadmin -f phpmyadmin.yml > /dev/null
            echo "\033[1;33m http://mysqladmin.docker/ Login: mysql-primary, root/root \033[0m \n"
            ;;

        MONGO)
            helm $CMD mongo bitnami/mongodb -f mongo.yml > /dev/null
            echo "\033[1;33m mongosh -u root -p root --host localhost  < /etc/files/scripts/mongo.js \033[0m \n"
            echo "\033[1;33m Svc Endpoint: mongo-mongodb:27017 (Standalone Mode Only) \033[0m \n"
            ;;

        REDIS)
            helm $CMD redis bitnami/redis -f redis.yml > /dev/null
            helm $CMD redis-admin onechart/onechart -f redis-admin.yml

            echo "\033[1;33m redis-cli -c incr mycounter \033[0m \n"
            echo "\033[1;33m redis-cli -c set mypasswd lol \033[0m \n"
            echo "\033[1;33m redis-cli -c get mypasswd \033[0m \n"
            echo "\033[1;33m Commander: http://redisadmin.docker/ \033[0m \n"
            ;;
        ISTIO)
            #helm repo add istio https://istio-release.storage.googleapis.com/charts
            helm $CMD istio-base istio/base -n istio-system --create-namespace > /dev/null
            helm $CMD istiod istio/istiod -n istio-system -f istio.yml > /dev/null
            # helm install istio-ingress istio/gateway -n istio-system --wait

            echo "\033[1;32m Enabled Istio for Default Namespace \033[0m \n"
            # kubectl label namespace default istio-injection-
            kubectl label namespace default istio-injection=enabled --overwrite
            ;;
        KIALI)
            #helm repo add kiali https://kiali.org/helm-charts
            helm $CMD kiali-operator kiali/kiali-operator -f kiali.yml > /dev/null
            #Create Kiali CRD
            kubectl apply -f ./files/istio/kiali-crd.yml
            echo "\033[1;33m Kiali: http://kiali.docker/kiali \033[0m \n";
            ;;

        PROXY)
            echo "\033[1;32m Nginx \033[0m \n"
            helm $CMD nginx bitnami/nginx -f nginx.yml > /dev/null

            echo "\033[1;33m http://nginx.docker/ \033[0m \n"
            echo "\033[1;33m Refer Server Blocks for More Help.. \033[0m \n"
            
            echo "\033[1;32m Resty \033[0m \n"
            helm $CMD resty onechart/onechart -f resty.yml > /dev/null

            echo "\033[1;33m http://resty.docker/ \033[0m \n"
            echo "\033[1;33m http://resty.docker/example \033[0m \n"
            echo "\033[1;33m http://resty.docker/ndtv \033[0m \n"
            
            echo "\033[1;32m Squid \033[0m \n"
            helm $CMD squid onechart/onechart -f squid.yml > /dev/null
            echo "\033[1;33m ALLOW: curl -x localhost:3128 www.google.com \033[0m \n"
            echo "\033[1;33m DENY: curl -x localhost:3128 www.fb.com \033[0m \n"

            echo "\033[1;32m TinyProxy \033[0m \n"
            helm $CMD tinyproxy stakater/application -f tinyproxy.yml > /dev/null

            echo "\033[1;33m curl -x localhost:8888 tinyproxy.stats \033[0m \n"
            ;;

        LOADER)
            helm $CMD vegeta onechart/onechart -f vegeta.yml > /dev/null
            echo "\033[1;33m Check Logs for Output \033[0m \n"
            echo "\033[1;33m echo 'GET http://nginx' | vegeta attack | vegeta report \033[0m \n"
            ;;

        HTTPBIN)
            helm $CMD httpbin onechart/onechart -f httpbin.yml > /dev/null
            echo "\033[1;33m Swagger: http://httpbin.docker \033[0m \n"
            echo "\033[1;33m http://httpbin.docker/anything \033[0m \n"
            echo "\033[1;33m curl http://httpbin:8810/headers \033[0m \n"
            ;;

        CRON)
            # FIXME: Cron
            helm $CMD cron onechart/onechart -f cron.yml > /dev/null
            echo "\033[1;33m Check Logs for Output \033[0m \n"
            ;;

        DASHY)
            helm $CMD dashy onechart/onechart -f dashy.yml > /dev/null
            echo "\033[1;33m http://dashy.docker/ \033[0m \n"
            ;;
    
        APP)
            helm $CMD app onechart/onechart -f app.yml > /dev/null
            echo "\033[1;33m http://app.docker/app/metrics\n http://app.docker/app/person/all \033[0m \n"
            echo "\033[1;33m echo 'GET http://app:9001/person/all' | vegeta attack | vegeta report \033[0m \n"
            ;;

        OPA)
            # helm repo add opa https://open-policy-agent.github.io/kube-mgmt/charts
            helm $CMD opa opa/opa-kube-mgmt -f opa.yml > /dev/null
            helm $CMD opa-demo onechart/onechart -f opa-demo.yml > /dev/null
            
            echo "\033[1;33m curl --user david:password http://opa.docker/finance/salary/david \033[0m \n"

            echo "\033[1;33m Demo (opa-demo): /demo/hr.sh \033[0m \n"
            echo "\033[1;33m Demo (opa-demo): /demo/authz.sh \033[0m \n"
            echo "\033[1;33m Docker (Localhost): ./demo/docker.sh \033[0m \n"
            ;;

        VAULT)
            # helm repo add hashicorp https://helm.releases.hashicorp.com
            helm $CMD vault hashicorp/vault -f vault.yml > /dev/null
            echo "\033[1;33m vault status \033[0m \n"
            echo "\033[1;33m /demo/vault.sh \033[0m \n"
            ;;
        PORTAINER)
            # helm repo add portainer https://portainer.github.io/k8s/
            helm $CMD portainer portainer/portainer -f portainer.yml
            echo "\033[1;33m http://portainer.docker/ \033[0m \n"
            ;;
        TRAEFIK)
            # helm repo add traefik https://traefik.github.io/charts
            helm $CMD traefik traefik/traefik -f traefik.yml > /dev/null
            kubectl apply -f ./files/traefik/middleware.yml
            # kubectl apply -f ./files/traefik/ingress.yml

            echo "\033[1;33m Dashboard: http://docker:9000/dashboard/#/ \033[0m \n"
            echo "\033[1;33m HealthCheck: http://docker:9000/ping \033[0m \n"
            echo "\033[1;33m Ingress: http://docker:8000/mysqladmin \033[0m \n"
            echo "\033[1;33m PortForward(80): sudo kubectl port-forward deployment/traefik 80:8000 > /dev/null &\033[0m \n"
            ;;

        CONSUL)
            # helm repo add hashicorp https://helm.releases.hashicorp.com
            helm $CMD consul hashicorp/consul -f consul.yml > /dev/null
            echo "\033[1;33m http://consul.docker/ \033[0m \n"
            ;;

        ETCD)
            helm $CMD etcd bitnami/etcd -f etcd.yml > /dev/null
            echo "\033[1;33m ./demo/demo.sh \033[0m \n"
            ;;

        SONAR)
            helm $CMD sonar bitnami/sonarqube -f sonar.yml > /dev/null
            echo "\033[1;33m http://sonar.docker/ \033[0m \n"
            echo "\033[1;33m Login: aman/aman (Need 5GB Mem) \033[0m \n"
            ;;

        ZOOKEEPER)
            helm $CMD zookeeper bitnami/zookeeper -f zookeeper.yml > /dev/null
            echo "\033[1;33m /demo/demo.sh \033[0m \n"
            ;;
        
        ELK)
            #FIXME: Complete Compose Setup
            helm $CMD logstash bitnami/logstash -f logstash.yml > /dev/null
            helm $CMD elasticsearch bitnami/elasticsearch -f elasticsearch.yml > /dev/null
            # helm $CMD kibana bitnami/kibana -f kibana.yml > /dev/null
            echo "\033[1;33m /demo/demo.sh \033[0m \n"
            ;;

        MONITOR)
            #helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
            helm $CMD prometheus prometheus-community/prometheus -f prometheus.yml > /dev/null
            echo "\033[1;33m Prometheus Server: http://prometheus.docker/ \033[0m \n"
            echo "\033[1;33m Prometheus Query: http://prometheus.docker/api/v1/query \033[0m \n"

             # helm repo add grafana https://grafana.github.io/helm-charts
            helm $CMD grafana grafana/grafana -f grafana.yml > /dev/null
            echo "\033[1;33m http://grafana.docker/login (aman/aman) \033[0m \n"
            echo "\033[1;33m Datasource: http://grafana.docker/datasources/new \033[0m \n"
            echo "\033[1;33m Add Datasource Prometheus: http://prometheus-server \033[0m \n"

            #helm repo add jaegertracing https://jaegertracing.github.io/helm-charts
            helm $CMD jaeger jaegertracing/jaeger -f jaeger.yml > /dev/null
            echo "\033[1;33m http://jaeger.docker/ \033[0m \n"

            ;;
        
        LDAP)
            helm $CMD ldap onechart/onechart -f ldap.yml > /dev/null
            helm $CMD ldap-admin onechart/onechart -f ldap-admin.yml > /dev/null
            echo "\033[1;33m CMD: ldapsearch -H ldap://localhost:3891 -xLL -D 'cn=admin,dc=example,dc=com' -b 'dc=example,dc=com' -W '(cn=admin)' \033[0m \n"
            echo "\033[1;33m UI: http://ldapadmin.docker/ \033[0m \n"
            echo "\033[1;33m Admin Login: Username:cn=admin,dc=example,dc=com Password: admin \033[0m \n"
            ;;
        MYSQL-OP)
            #helm repo add bitpoke https://helm-charts.bitpoke.io
            helm $CMD mysql-operator bitpoke/mysql-operator -f bitspoke.yml > /dev/null
            kubectl apply -f ./files/bitspoke/secret.yml
            kubectl apply -f ./files/bitspoke/cluster.yml
            helm $CMD mysql-admin bitnami/phpmyadmin > /dev/null
            echo "\033[1;33m Mysql Info: kubectl get mysql; kubectl describe mysql mysql-operator; \033[0m \n"
            echo "\033[1;33m Mysql Clear: kubectl delete mysql mysql-operator; \033[0m \n"
            echo "\033[1;33m Login: root/root, aman/aman [Host: mysql] \033[0m \n"
            ;;
        WEBSHELL)
            helm $CMD sshwifty onechart/onechart -f sshwifty.yml  > /dev/null
            helm $CMD webssh onechart/onechart -f webssh.yml > /dev/null
            #FIXME: Fix Config for Web SSH
            helm $CMD webssh2 onechart/onechart -f webssh2.yml > /dev/null
            echo "\033[1;33m Sshwifty: http://sshwifty.docker/ \033[0m \n"
            echo "\033[1;33m Webssh: http://webssh.docker/ \033[0m \n"
            echo "\033[1;33m Webssh: http://webssh2.docker/ \033[0m \n"
            ;;

        *)
            #TODO: Add Locust
            echo "\033[1;34m Service Not Supported: $SVC \033[0m \n"
            ;;
        esac
    done
}

delete()
{
    echo "\033[1;32m Clearing all Helms \033[0m \n"
    #Clear CRD's (Needed before Helm Deletion)
    kubectl delete kiali --all --all-namespaces 2> /dev/null
    #HACK: Add Mysql CRD's
        
    #Exclude Permanent Helms
    helm delete $(helm list --short | grep -v "traefik\|dashy") 2> /dev/null
}


# Flags
while getopts 'dusrib' OPTION; do
  case "$OPTION" in
    r)
        NS=$(kubectl get sa -o=jsonpath='{.items[0]..metadata.namespace}')
        echo "\033[1;32m Restting Namespace: $NS \033[0m \n"
        #Delete Resources
        kubectl delete --all all --namespace=$NS
        #Process Normal Delete
        delete
        #Helm Clear Remaining
        helm delete $(helm list --short)
        #Istio Clear
        helm delete -n istio-system $(helm list --short -n istio-system)
        ;;
    d)
        # Process Delete
        delete
        ;;
    i)
        echo "\033[1;32m Helm: Install \033[0m \n"
        process
        ;;
    u)
        echo "\033[1;32m Helm: Upgrade \033[0m \n"
        CMD="upgrade"
        process
        ;;
    b)
        echo "\033[1;32m Bootstraping Base Services \033[0m \n"

        process "TRAEFIK DASHY $XTRA_BOOT"
        
        echo "\033[1;32m Attempting Traefik Portforward \033[0m \n";
        kubectl wait --for=condition=Ready pod -l app.kubernetes.io/name=traefik --timeout=2m
        # kubectl port-forward deployment/traefik 9000:9000 > /dev/null &
        # kubectl port-forward deployment/traefik 8000:8000 > /dev/null &
        ;;
    s)
        # Prompt
        answers=`gum choose MYSQL MONGO REDIS APP PROXY LOADER CRON HTTPBIN VAULT OPA CONSUL LDAP ETCD SONAR PORTAINER ZOOKEEPER ELK ISTIO KIALI MONITOR WEBSHELL MYSQL-OP --limit 5`
        echo $answers > $ANS_FILE    
        echo "\033[1;32m Service Set \033[0m \n"
        echo "\033[1;33m $answers \033[0m \n"
        ;;
    ?)
        echo "\033[1;32m script usage: $0 [-s] [-i] [-u] [-d] \033[0m \n"
        echo "\033[1;33m [-s] Set \033[0m \n"
        echo "\033[1;33m [-i] Install \033[0m \n"
        echo "\033[1;33m [-u] Upgrade \033[0m \n"
        echo "\033[1;33m [-d] Delete \033[0m \n"
        echo "\033[1;33m [-r] Reset \033[0m \n"
        exit 1
        ;;
  esac
done