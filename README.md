[![Actions Status](https://github.com/amanhigh/go-fun/workflows/Build/badge.svg)](https://github.com/amanhigh/go-fun/actions)
[![codecov](https://codecov.io/gh/amanhigh/go-fun/branch/master/graph/badge.svg)](https://codecov.io/gh/amanhigh/go-fun)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/amanhigh/go-fun)](https://github.com/amanhigh/go-fun/releases)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/amanhigh/go-fun)
[![Go Report Card](https://goreportcard.com/badge/github.com/amanhigh/go-fun)](https://goreportcard.com/report/github.com/amanhigh/go-fun)


# Go Fun
Experiments & Fun  with Go Lang and its Frameworks. Also includes tools like docker, k8, istio, observability, and perf.

## Build
Use goreleaser for build test. [Install](https://goreleaser.com/install/) if not already installed

goreleaser build --snapshot --rm-dist
goreleaser release --snapshot --skip-publish --rm-dist

## Testing
Testing is handled via [Ginkgo](https://github.com/onsi/ginkgo). To run all unit tests excluding ones require external setup.

`ginkgo -r '--label-filter=!setup' .`

## Load Test
* brew install gum vegeta
* Run `cd ./components/fun-app/it/;./load.sh`

## FunApp
Sample Funapp which is rest based app with various tools and tests required as sample.

By default it runs without any dependencies with in memory [sqlite3](https://github.com/mattn/go-sqlite3) database which can be configured via ENV Variables.


------
### Direct Run
`go run ./components/fun-app/` 

<br/> ![](common/images/fun-app/go-run.gif)

### Vscode Run
* Checkout Code
* Run FunApp Test Configuration
* Configure [ENV](components/fun-app/.env) if required

<br/> ![](common/images/fun-app/vscode-run.gif)

### Docker Run
`docker run amanfdk/fun-app`
<br/>
[Docker Hub](https://hub.docker.com/r/amanfdk/fun-app)


### K8/Istio Run
- Setup: <br/>
`helm repo add go-fun https://amanhigh.github.io/go-fun` <br/>
`helm install -n fun-app fun-app go-fun/fun-app` <br/>
Open http://localhost:9000/metrics (Minikube: Run "minikube tunnel")

    
- Cleanup: <br/>
 `helm -n fun-app delete fun-app`

<br/> ![](common/images/fun-app/helm.gif)

### Development Container
 *After Helm Setup*, run Development Remote Container.
 It is configured to Auto Reload Code Changes

* Start:<br/>
    Run `devspace -n fun-app dev` <br/>
    Open http://localhost:8080/metrics

* Cleanup: `devspace -n fun-app purge`

<br/> ![](common/images/fun-app/devcode.gif)



## TODO
- Message Queue
- Traces
- Swagger