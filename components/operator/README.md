# Memcached Operator
- ## Description
    This is a Trial K8 Controller Project to learn about it. Tutorial Followed can be found [here](https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/). For Operator Golang has been used which gives most flexibility.  
    
    This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)  
    
    It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/) which provides a reconcile function responsible for synchronizing resources untile the desired state is reached on the cluster.   
    
    More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)  
- ## Deployment
  This section details various ways to Run Operator and Test with and Without Cluster.  
	- ### With Cluster
	  You’ll need a Kubernetes cluster which can be local (kind/minikube) or remote.  
		- #### Setup
			- Deploy Cert Manager required for Webhooks.
			  `make deploy-cert`  
			- Deploy Operator: Operator is deployed in *operator-system* Namespace.
			  `make deploy`  
			- Install CRD: Create Custom Resources in Current Namespace.
			  `kubectl apply -f config/samples/cache_v1beta1_memcached.yaml`  
		- #### Cleanup: Remove Operator and CRD's:
		  `kubectl delete -f config/samples/cache_v1beta1_memcached.yaml`  
		  `make undeploy undeploy-cert`  
		- #### Integration Testing
			- Envtest has some [Limitation](https://book.kubebuilder.io/reference/envtest.html#namespace-usage-limitation) which are not there when Test on Cluster
			- Run in *Test Suit Directory*
			  `export USE_EXISTING_CLUSTER=true && ginkgo .`  
	- ### Without Cluster
		- Unit Test: `make test` (*This will internally run generate and manifest targets as well.*)
		- #### Outside Cluster
			- Place Certificates: `cp tls.crt tls.key /tmp/k8s-webhook-server/serving-certs`
			- Runs Cluster in Local: `MEMCACHED_IMAGE=memcached:1.4.36-alpine make install run`
			- Install CRD: Create Custom Resources in Current Namespace.
- ## Setup
    [Cheatsheet](https://sdk.operatorframework.io/docs/overview/cheat-sheet/)  -  [Layout](https://sdk.operatorframework.io/docs/overview/project-layout/)  
    - ### Init
        - Group: Combine related Resources Together. Format of <group>.<domain> eg. `cache.aman.com` for domain: `aman.com` is used for group CRD
        - Repo is used for Golang Module Management (go.mod generation).
        - Generate Go Module and add to Multi-Module Project (go.work)
            ```
            operator-sdk init --domain aman.com --repo github.com/amanhigh/go-fun/components/operator
            go work use ./components/operator/
            ```
    - ### Controller
        - Generate [Controller](https://book.kubebuilder.io/cronjob-tutorial/controller-overview.html),Types and [API](https://book.kubebuilder.io/cronjob-tutorial/new-api.html)
        - Type will be `MemCached` with Version `v1alpha1` available under group `cache.aman.com`.
        - Generated go Files are under ./api, ./controllers and  ./config has yaml files.
            
            `operator-sdk create api --group cache --version v1alpha1 --kind Memcached --resource --controller`  
    - #### Image Plugin
        - Helps to control Docker File. Command Guides Creation of Docker file with image, command, user specifications.
        - This Generates [Controller](https://github.com/operator-framework/operator-sdk/blob/latest/testdata/go/v3/memcached-operator/controllers/memcached_controller.go), Type Specs and its Test.
            `operator-sdk create api --group cache --version v1alpha1 --kind Memcached --plugins="deploy-image/v1-alpha" --image=memcached:1.4.36-alpine --image-container-command="memcached,-m=64,modern,-v" --run-as-user="1001"`  
    - #### Models
        - Made Model Modification (ContainerPort addition) and do [autogeneration](https://book.kubebuilder.io/cronjob-tutorial/other-api-files.html).
        - Updated Test and [Spec](config/samples/cache_v1alpha1_memcached.yaml) to include `containerPort: 8443`
            `make generate`  
    - #### Manifests
        - Generate Manifests CRD [(cache.aman.com_memcacheds.yaml)](config/crd/bases/cache.aman.com_memcacheds.yaml), RBAC [(role.yaml)](config/rbac/role.yaml).
        - If you are editing the API definitions, generate the manifests such as CRs or CRDs using.
            `make manifests`  
    - #### Docker
        - Updated Docker Base Image and any Modifications in [DockerFile](./Dockerfile).
        - Fix Image Name `IMG ?= amanfdk/operator:latest` in [MakeFile](./Makefile).
        - Push: `make docker-build docker-push` | Verify: `docker images`
        - Minkube Reload:  Added New Make Target to Build and Reload.
            `minikube-push` (Ensure **imagePullPolicy:IfNotPresent** is Set)  
    - #### WebHook
        - This step Generates Default ( Mutating ) and Validating [Webhooks](https://sdk.operatorframework.io/docs/building-operators/golang/webhook/).
        - #### Generate
            - Default and Validation Hooks
                `operator-sdk create webhook --group cache --version v1alpha1 --kind Memcached --defaulting --programmatic-validation`  
            - Adds `.SetupWebhookWithManager` in main.go, [Webhooks](./api/v1alpha1/memcached_webhook.go) are implemented and [registered](./config/default/webhookcainjection_patch.yaml).
            - It Generates [CertConfig](./config/certmanager/certificate.yaml), to generate two self signed certificates.
            - Webhook [Server](./config/default/manager_webhook_patch.yaml) will listen on part `9443` behind Webhook [Service](./config/webhook/service.yaml) `443`
        - Uncomment WebHook & CertManager (Not Required in [OLM](https://github.com/operator-framework/operator-sdk/issues/6257)) in [Default](./config/default/kustomization.yaml), [CRD](./config/crd/kustomization.yaml), [Manifest](./config/manifests/kustomization.yaml) [Kuztomize](https://book.kubebuilder.io/cronjob-tutorial/running-webhook.html#deploy-webhooks) Files
        - Deploy Cert [Manager](https://cert-manager.io/docs/installation/)
            `kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.11.0/cert-manager.yaml`  
        - Rebuild: Push: `make docker-build docker-push` | Verify: `docker images`
        - Setup with Cluster
            - Verify: `kubectl get validatingwebhookconfigurations`
            - Apply CRD: Verify Webhook Logs in `Operator` Pod
  - #### CRD Update
    - Introduce new [Version](https://vincenthou.medium.com/how-to-create-conversion-webhook-for-my-operator-with-operator-sdk-36f5ee0170de#aec0) of CRD `v1beta1` which will make `sidecarImage` configurable.
    - New Version`v1beta1`
      - `operator-sdk create api --group cache --version v1beta1 --kind Memcached`
      - Generates new Sample Spec, New Package [v1beta1](api/v1beta1) with Empty Types
      - Implemented Empty Spec and Type.
      - Marked v1beta1 Storage Version `//+kubebuilder:storageversion`.
      - Regenerate: `make generate manifests`
      - Updated Controller, Test to newer Version.
    - Webhook
      - Regenerate Webhook including Conversion: `operator-sdk create webhook --group cache --version v1beta1 --kind Memcached --defaulting --programmatic-validation --conversion`
      - Implemented New [Webhooks](api/v1beta1/memcached_webhook.go) and did Regenerate.
      - Deleted old Webhook for v1alpha1 as Conversion Hooks runs before Others.
    - Conversion
      - Added [memcache_version.go](api/v1beta1/memcached_conversion.go) to v1beta1 to make it Hub Version.
      - Added [memcache_version.go](api/v1alpha1/memcached_conversion.go) to v1alpha1 to Implement Conversion methods.
      - Deploy Older Version
        ```
        kubectl apply -f config/samples/cache_v1alpha1_memcached.yaml
        kubectl delete -f config/samples/cache_v1alpha1_memcached.yaml
        ```
- ### Rough Notes
    Temporary Notes Section
