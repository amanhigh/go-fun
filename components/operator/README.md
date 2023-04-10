# Memcached Operatior

## Description
This is a Trial K8 Controller Project to learn about it. Tutorial Followed can be found [here](https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/). For Operatior Golang has been used which gives most flexibility.

This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/) 
which provides a reconcile function responsible for synchronizing resources untile the desired state is reached on the cluster. More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## Setup
[Cheatsheet](https://sdk.operatorframework.io/docs/overview/cheat-sheet/)\
[Layout](https://sdk.operatorframework.io/docs/overview/project-layout/)

Steps followed
* **Init** -  Domain is used for group CRD and Repo is used for Golang Module Management (go.mod generation). Add Generated Module to Go Work File.\
`operator-sdk init --domain aman.com --repo github.com/amanhigh/go-fun/components/operator`\
`go work use ./components/operator/`

* **Controller** - Generate [Controller](https://book.kubebuilder.io/cronjob-tutorial/controller-overview.html) and Types. Type/[API](https://book.kubebuilder.io/cronjob-tutorial/new-api.html) will be *MemCached* with Version *v1alpha1* available under group *cache.aman.com*. Generated Go Files are under ./api, ./controllers and  ./config has yaml files.\
`operator-sdk create api --group cache --version v1alpha1 --kind Memcached --resource --controller`

* **Image Plugin** - Helps to control Docker File. Command Guides Creation of Docker file with image, command, user specifications.  This Generates [Controller](https://github.com/operator-framework/operator-sdk/blob/latest/testdata/go/v3/memcached-operator/controllers/memcached_controller.go), Type Specs and its Test. \
`operator-sdk create api --group cache --version v1alpha1 --kind Memcached --plugins="deploy-image/v1-alpha" --image=memcached:1.4.36-alpine --image-container-command="memcached,-m=64,modern,-v" --run-as-user="1001"`

* **Models** - Made Model Modification (ContainerPort addition) and do [autogeneration](https://book.kubebuilder.io/cronjob-tutorial/other-api-files.html). Updated Test and [Spec](config/samples/cache_v1alpha1_memcached.yaml) to include *containerPort: 8443*\
`make generate`

* **Manifests** - Generate Manifests CRD [(cache.aman.com_memcacheds.yaml)](config/crd/bases/cache.aman.com_memcacheds.yaml), RBAC [(role.yaml)](config/rbac/role.yaml).If you are editing the API definitions, generate the manifests such as CRs or CRDs using.\
 `make manifests`

* **Docker** - Updated Docker Base Image or any other Changes in [DockerFile](./Dockerfile). Fix Image Name `IMG ?= amanfdk/operator:latest` in [MakeFile](./Makefile).\
\
`make docker-build docker-push`\
`docker images`

**NOTE:** Run `make --help` for more information on all potential `make` targets

## Deployment
Project can be run in following ways

### Testing
* Without Cluster - `make test` (*This will interally run generate and manifest targets as well.*)
* With Cluster - `export USE_EXISTING_CLUSTER=true && ginkgo .` <br/>
[Requires Running Minikube Cluster with Context Set. Test should be run from Test Suit Directory] <br/>
This Includes extra tests which can't run on envtest as it simulates limited K8 Functions.



### Outside Cluster
Run `make install run` to run without cluster.

### On cluster
You’ll need a Kubernetes cluster which can be local (kind/minikube) or remote. 

**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

#### --Setup--
1. Deploy Operator: Operator is deployed in *operator-system* Namespace.\
`make deploy`

2. Install Instances (Current Namespace) of Custom Resources:\
`kubectl apply -f config/samples/cache_v1alpha1_memcached.yaml`

#### --Cleanup--
1. Remove Cluster:  `kubectl delete -f config/samples/cache_v1alpha1_memcached.yaml`
2. Remove Operator and CRD's: `make undeploy`

## Resources
### Naming Convetions
* API Group: <group>.<domain> eg. cache.aman.com
* Resource: GroupVersionKind [GVK](https://book.kubebuilder.io/cronjob-tutorial/gvks.html). Eg.
```
    apiVersion: cache.aman.com/v1alpha1
    kind: Memcached
```
* Spec: Golang Counterpart of Resource. Eg. `v1aplha1/Memcached`
    * Fields: Fields in Spec should follow CamelCase. Can also be ommited when empty. Eg. ``json:"containerPort,omitempty"``

### Markers
Markers are Golang Comments/Tags which hint Kubebuilder Generator.
* Kind: Type is a Kind `//+kubebuilder:object:root=true`
* Group: Tags Go Package to hold & Generaet Kind Objects `+kubebuilder:object:generate=true`
* [RBAC](https://book.kubebuilder.io/reference/markers/rbac.html): Generates [ClusterRole](config/rbac/role.yaml) for Controller & its Perms. `
    * CRD: `//+kubebuilder:rbac:groups=cache.aman.com,resources=memcacheds,verbs=get;list;watch;create;update;patch;delete`. Permission can also be on subresource like `memcacheds/status, memcacheds/finalizers`
    * Default Resources: `//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch`
