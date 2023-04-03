PORT=8091
MINI_BKP_FILE=~/Downloads/mini-bkp.txt
MINI_CURRENT_BKP_FILE=/tmp/mini-bkp
input=$1
answers=${input:-`gum choose MINIKUBE ISTIO BACKUP RESTORE CLEAN --limit 5`}

for SVC in $answers
do
    echo "\033[1;32m \n $SVC \033[0m \n"
    case $SVC in
    ISTIO)
        ./istio/istio.sh;
        ;;
    CLEAN)
        echo "\033[1;32m Deleting Minikube Clusters \033[0m \n"
        minikube -p minikube delete;
        minikube -p secondary delete
        exit 0
        ;;

    BACKUP)
        # TODO: Handle None Tags
        [[ ! -f $MINI_BKP_FILE ]] && touch $MINI_BKP_FILE
        minikube image ls | grep -v registry.k8s.io | grep -v none | tee $MINI_CURRENT_BKP_FILE
        echo "\033[1;33m Backup Count: `wc -l $MINI_CURRENT_BKP_FILE` \033[0m \n"

        
        # Append Image list to Master List
        cp $MINI_BKP_FILE /tmp/mini-bkp-old
        sort $MINI_CURRENT_BKP_FILE /tmp/mini-bkp-old | uniq | tee $MINI_BKP_FILE
        echo "\033[1;33m Image Count: `wc -l $MINI_BKP_FILE` \033[0m \n"

        for IMG in `cat $MINI_CURRENT_BKP_FILE`
        do 
            echo "\033[1;33m Caching Image: $IMG \033[0m \n"
            CPATH="${IMG%:*}";
            mkdir -p ~/.minikube/cache/images/$(uname -m)/$CPATH;
            minikube image save --daemon $IMG
        done
        exit 0
        ;;

    RESTORE)
        for IMG in `cat $MINI_BKP_FILE`
        do 
            echo "\033[1;33m Loading Image: $IMG \033[0m \n"
            minikube image load --daemon $IMG
        done
        exit 0
        ;;

    MINIKUBE)
        #Use minikube config set vm-driver virtualbox/docker
        minikube -p minikube delete;

        echo "\033[1;32m Creating Minikube Cluster \033[0m \n"
        FILE_PATH=`readlink -f ./services/files`
        
        #Additional Flags: --kubernetes-version v1.23.0
        minikube  -p minikube start --memory=3096 --cpus=3 --cache-images=true --mount-string="$FILE_PATH:/etc/files" --mount --host-only-cidr='24.1.1.100/24';
        
        # TODO: Fix on WSL
        # --extra-config="apiserver.enable-swagger-ui=true";
        # --extra-config="apiserver.service-account-api-audiences=api" \
        # --extra-config="apiserver.service-account-issuer=api" \
        # --extra-config="apiserver.service-account-key-file=/var/lib/minikube/certs/sa.pub" \
        # --extra-config="apiserver.service-account-signing-key-file=/var/lib/minikube/certs/sa.key";
        # minikube -p minikube ssh 'sudo cat /var/lib/minikube/certs/sa.pub'

        # echo "\033[1;32m Minikube Dashboard & Addons \033[0m \n";
        # minikube -p minikube dashboard --url=true > /dev/null &
        # minikube addons enable metrics-server;

        ;;

    *)
        echo "\033[1;34m Addon Not Supported: $SVC \033[0m \n"
        ;;
    esac
done

echo "\033[1;32m Minikube Setup \033[0m \n";
echo "\033[1;33m Run 'minikube tunnel' for Emulating ELB\033[0m \n";
echo "\033[1;33m Dashboard: http://localhost:$PORT/api/v1/namespaces/kubernetes-dashboard/services/http:kubernetes-dashboard:/proxy/# \033[0m \n";
echo "\033[1;33m Swagger: http://localhost:$PORT/swagger-ui \033[0m \n";
echo "\033[1;33m K9S:  k9s --context minikube \033[0m \n";
echo "\033[1;33m Context: `kubectl config current-context;`\033[0m \n";
# kubectl proxy --port=$PORT;

echo "\033[1;34m Waiting for Minikube to be Ready \033[0m \n";
sleep 20
kubectl wait --for=condition=Ready pod -l k8s-app=kube-dns -n kube-system

cd ./services
./service.zsh -b
cd -

echo "\033[1;34m Please enter password for Port 80 Forward \033[0m \n";
sudo kubectl port-forward deployment/traefik 80:8000

## k9s
# k9s --readonly , -n <namespace>, -l <loglevel>
# k9s :pu (pulse), :dp (deployments), :po (pods), :ns (namespace),:rb (Role Bindings) ,:a (aliases), :hpa (autoscaler)
# Switch Namespace :po <namespace> to see pods of that namespace
# Portforward: Select Pod, Shift+f (Create PF), f (Show PF)

## kompose: brew install kompose
# Convert:  kompose convert -f jira.yml
# Apply: kubectl apply $(ls jira*.yaml | awk ' { print " -f " $NF } ')

## Docker
# Docker Start: sudo service docker start
# Docker Image to Minikube: eval $(minikube docker-env); docker build -t fun-app .

## Kubectl
# Port Forward (Local Port 9090 to Container Port 8080) - kubectl port-forward `kubectl get pods -o name | grep fun-app | head  -1` 9090:8080
# Logs - kubectl logs `kubectl get pods -o name | grep fun-app | head  -1` --since=1m -f
# Login - kubectl -it exec `kubectl get pods -o name | grep fun-app | head  -1` bash
# Delete All - kubectl delete all --all

## Minkikube
# Tunnel (Emulate Load Balancer) - minikube tunnel
# List Emulated Services URL's - minikube service list

# https://minikube.sigs.k8s.io/docs/handbook/pushing/
# Load Image from Host - minikube image load amanfdk/controller
# Minikube Docker ENV Connect - eval $(minikube docker-env)

## Minikube Images (Cache: ~/.minikube/cache)
# Image List - minikube image ls
# Cache Save - minikube image save docker.io/bitnami/mysql:8.0.32-debian-11-r0 --daemon
# Cache Load - minikube image load docker.io/bitnami/mysql:8.0.32-debian-11-r0 --daemon