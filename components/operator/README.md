# Memcached Operatior

## Description
This is a Trial K8 Controller Project to learn about it. Tutorial Followed can be found [here](https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/). For Operatior Golang has been used which gives most flexibility.

## Setup

Steps followed
* **Init** -  Domain is used for group CRD and Repo is used for Golang Module Management (go.mod generation). Add Generated Module to Go Work File.\
`operator-sdk init --domain aman.com --repo github.com/amanhigh/go-fun/components/operator`\
`go work use ./components/operator/`

* **Controller** - Generate Controller and Types. Type/API will be *MemCached* with Version *v1alpha1* available under group *cache.aman.com*. Generated Go Files are under ./api, ./controllers and  ./config has yaml files.\
`operator-sdk create api --group cache --version v1alpha1 --kind Memcached --resource --controller`

* **Image Plugin** - Helps to control Docker File. Command Guides Creation of Docker file with image, command, user specifications.  This Generates [Controller](https://github.com/operator-framework/operator-sdk/blob/latest/testdata/go/v3/memcached-operator/controllers/memcached_controller.go), Type Specs and its Test. \
`operator-sdk create api --group cache --version v1alpha1 --kind Memcached --plugins="deploy-image/v1-alpha" --image=memcached:1.4.36-alpine --image-container-command="memcached,-m=64,modern,-v" --run-as-user="1001"`

* **Models** - Made Model Modification (ContainerPort addition) and do autogeneration. Updated [Spec](config/samples/cache_v1alpha1_memcached.yaml) to include *containerPort: 8443*\
`make generate`

* **Manifests** - Generate Manifests CRD (cache.aman.com_memcacheds.yaml), RBAC(role.yml).\
 `make manifests`

* **Docker** - Updated Docker Base Image & [File](./Dockerfile). Fixed Image Name (amanfdk/operator) in Make File.\
`make docker-build docker-push`\
`docker images`

## Deployment
Project can be run in following ways

### Outside Cluster
Run `make install run` to run without cluster.

### On cluster
You’ll need a Kubernetes cluster which can be local (kind/minikube) or remote. 

**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

#### --Setup--
1. Deploy Operator: Operator is deployed in *operator-system* Namespace.\
`make deploy`

2. Install Instances (Current Namespace) of Custom Resources:\
`kubectl apply -f config/samples/`

#### --Cleanup--
1. To delete the CRDs from the cluster:  `make uninstall`
2. UnDeploy the controller to the cluster: `make undeploy`

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/) 
which provides a reconcile function responsible for synchronizing resources untile the desired state is reached on the cluster 

### Test It Out
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

