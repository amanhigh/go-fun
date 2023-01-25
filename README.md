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
Need to have running K8/Istio Setup. Refer [Kubernetes](#Kubernetes) below.

* Via Github
    * Setup <br/>
    `helm repo add go-fun https://amanhigh.github.io/go-fun` <br/>
    `helm install -n fun-app fun-app go-fun/fun-app` <br/>

    * Cleanup: <br/>
    `helm -n fun-app delete fun-app`

* Via Local
    * Setup <br/>
    `go-fun/components/fun-app/charts/setup.sh`
    
    * Clean <br/>
    `go-fun/components/fun-app/charts/clean.sh`


    
Open http://localhost:9000/metrics (Minikube: Run "minikube tunnel")


<br/> ![](common/images/fun-app/helm.gif)

### Development Container
 *After Helm Setup*, run Development Remote Container.
 It is configured to Auto Reload Code Changes

* Start:<br/>
    Run `devspace -n fun-app dev` <br/>
    Open http://localhost:8080/metrics

* Cleanup: `devspace -n fun-app purge` <br/>

* Override Vars:  `devspace list vars --var DB="mysql-primary",RATE_LIMIT=-1`

* Load Test (From Vegeta Container):  `echo 'GET http://app:8080/person/all' | vegeta attack | vegeta report`

* Log Analyzer: ``kubectl -n fun-app logs `kubectl get pods -n fun-app -o name | grep app | head  -1` --since=1m -f | goaccess --log-format='%^ %d - %t | %s | %~%D | %b | %~%h | %^ | %m %U' --date-format='%Y/%m/%d' --time-format '%H:%M:%S'``


<br/> ![](common/images/fun-app/devcode.gif)

### Log Analyzer
Monitor Logs via [GoAccess](https://github.com/allinurl/goaccess)

* Terminal Access: `go run main.go | goaccess --log-format='%^ %d - %t | %s | %~%D | %b | %~%h | %^ | %m %U' --date-format='%Y/%m/%d' --time-format '%H:%M:%S'`
* Web Access: Add Flags to Above Command `-o report.html --real-time-html`. Open report.html in Browser and it should auto refresh.

**Custom Log Monitoring**
- Identify [Date and Time Format](https://www.freebsd.org/cgi/man.cgi?query=strftime&sektion=3)). Eg. `date '+%Y/%m/%d - %H:%M:%S'` for 2023/01/23 - 14:38:2
- Identify [Log Format](https://goaccess.io/man#custom-log). Eg. 
- Define Date/Time Format, For Log Format start with initial fields and progress further for easy debug.
- Debug Mode: `-l debug.log`

**Useful Fields**
- Mandatory Fields: %d (Date), %h (Host), %r/%m %U (Request)
- Skip: Ignore (%^) , Skip Space (%~)
- DateTime: (%x/--datetime-format) OR Time (%t/--date-format) + Date (%d/--time-format)
- Host: IP (%h) OR Virtual Host (%v)
- Request: Full With Quotes (%r) or Method (%m), URL (%U), Query (%q), PROTOCOL (%H),
- Response: Status Code (%s), Size (%b)
- Latency: MicroSecond (%D), MilliSecond.MicroSecond (%T), MilliSecond With Decimal (%L)
- User Info: User-Agent (%u), Referer (%R)

TODO: Add GIF

## Kubernetes
To ease development and easy setup of dependencies we use Kubernetes. Also [K9S](https://github.com/derailed/k9s) provides easy interface to manage containers, see logs etc. [Helms](https://github.com/helm/helm) are used to setup various services which application can depend on.

### Minikube
To setup kubernetes there are multiple options available like minikube, kind, k89, k3s etc. In this project we are using [minikube](https://minikube.sigs.k8s.io/docs/).


Script and Multiselect can be used to enable Istio, Ingress Gateway etc.
* Setup - `./go-fun/Kubernetes/mini.sh`
* Teardown - `./go-fun/Kubernetes/clean.sh`

<br/> ![](common/images/fun-app/minikube.gif)

### Services
Package has multiple service which can be setup on top of Minikube. This helps in easy setup of complex dependencies like Mysql Cluster, Mongo, Prometheus, Sonar and many more ...

Service Script allows you multiple flags to set, create and teardown the setup.

Flags (Multiple flags can be passed together)
* Set (s) - Allows you set Service Reciepe.
* Install (i) - Installs Helms
* Delete (d) - Deletes & Clears all Helms
* Reset (r) - Clear all Resources in Current Namespace & Helms

Eg.
* Set & Install - `./go-fun/Kubernetes/services/services.sh -si`
* Destroy & Install - `./go-fun/Kubernetes/services/services.sh -di` </br>
(Needs Set to be already done)

<br/> ![](common/images/fun-app/k8-service.gif)


## TODO
- Message Queue
- Traces
- Swagger