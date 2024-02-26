[![Actions Status](https://github.com/amanhigh/go-fun/actions/workflows/build.yml/badge.svg)](https://github.com/amanhigh/go-fun/actions?workflow=build)
[![codecov](https://codecov.io/gh/amanhigh/go-fun/branch/master/graph/badge.svg)](https://codecov.io/gh/amanhigh/go-fun)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/amanhigh/go-fun)](https://github.com/amanhigh/go-fun/releases)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/amanhigh/go-fun)
[![GitHub issues](https://img.shields.io/github/issues/amanhigh/go-fun)](https://github.com/amanhigh/go-fun/issues)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/amanhigh/go-fun)](https://github.com/amanhigh/go-fun/pulls)
[![Go Report Card](https://goreportcard.com/badge/github.com/amanhigh/go-fun)](https://goreportcard.com/report/github.com/amanhigh/go-fun)

# Go Fun
This repository follows the philosophy of Learning by Doing. It includes plays, experiments with Golang and its frameworks to learn. Later Kubernetes including Docker, Istio (Service Mesh), Performance testing is also included.  
- ## FunApp
	FunApp is a Sample Rest App which tries to use various golang Frameworks commonly required. It tries to follow good practices and standards. It runs without any dependencies with in memory [sqlite3](https://github.com/mattn/go-sqlite3) database by default.  
	- ### Setup
		- #### Onetime
			- We will use [Make](https://www.make.com/en) for Project Management.
			- Installs required dependencies and tools, `make prepare`
			- To check what all is available run `make` to display Help.
		- Run `make reset` will do Complete Build, Test, Coverage and Show Info.
		- Use `make info` (or `infos`) for displaying Info.
	- ### Play
		Easy ways to Play with **FunApp** without Dev Setup.  
		- Start
			- Golang: `make run`
				![Go Run](common/images/fun-app/go-run.gif)  
			- Docker [Image](https://hub.docker.com/r/amanfdk/fun-app): `docker run amanfdk/fun-app`
		- Testing
			- Unit and Integration Testing is done via [Ginkgo](https://github.com/onsi/ginkgo).
			- Run Tests: `make test test-slow` (Excludes that require separate setup.)
		- Performance Test
			[Vegeta](https://github.com/tsenart/vegeta) is the tool of choice here. [Gum](https://github.com/charmbracelet/gum) helps in prompts.  
			- Installation: `brew install gum vegeta`
			- Run: `make -C ./components/fun-app/it setup`
		- Linting
			- `make lint` will do Code Linting using [golangci-lint](https://github.com/golangci/golangci-lint)
	- ### Dev Setup
		This section guides on setup for Development Setup for this Repository. We are using Vscode and Kubernetes  for this example.  
		- ### Vscode
			- Use Command Palatte and run application using `> Select and Start Debugging` and then Select `FunApp`
			- Change [ENV](components/fun-app/.env) to override default Configuration.
			- ![Vscode Run](common/images/fun-app/vscode-run.gif)
		- ### Development Container
			- Development help you live debug an application in K8 Cluster.
			- It is configured to Auto Reload Code Changes.
			- This Sets Up Dev Container in `fun-app` Namespace.
			- Try:
				- Run `make space`, Open http://localhost:8080/metrics
				- Tests: `make space-test`
				- Cleanup: `make space-purge`
				- Check Environment Vars: `make infos`
				- Override Vars:  `devspace list vars --var DB="mysql-primary",RATE_LIMIT=10`
			- ![Devcode](common/images/fun-app/devcode.gif)
		- ### Kubernetes Cluster
			- Below Section requires running Kubernetes Cluster and Helm CLI. Charts can be used from local or from Github.
			- Via Github
				- Setup: `helm install -n fun-app fun-app go-fun/fun-app` (Onetime Setup Needed)
				- Cleanup: `helm -n fun-app delete fun-app`
			- Via Local
				- Deploys FunApp and Vegeta Container (for Load Test).
				- Setup: `make -C ./components/fun-app/charts setup` (or `reset`)
				- Clean: `make -C ./components/fun-app/charts clean` (or `info`)
			- Access
				- Open: http://localhost:9090/metrics  (Tunnel required for forwarding:  `minikube tunnel`)
				- Load Test (From Vegeta Container):  `echo 'GET http://fun-app:9090/person/all' | vegeta attack | vegeta report`
				- Log Analyzer : `make -C ./components/fun-app/charts analyse`
			- ![Helm](common/images/fun-app/helm.gif)
- ## Tools
	- ## Kubernetes
		To ease development and easy setup of dependencies we use Kubernetes. Also [K9S](https://github.com/derailed/k9s) provides easy interface to manage containers, see logs etc. [Helms](https://github.com/helm/helm) are used to setup various services which application can depend on.  
		- ### Minikube
			- To setup kubernetes there are multiple options available like minikube, kind, k89, k3s etc. In this project we are using [minikube](https://minikube.sigs.k8s.io/docs/).
			- This will also setup [traifik](https://github.com/traefik/traefik) ingress for easy access.
				- DNS Mapping via `/etc/hosts` is done as part of Onetime Prepration.
				- User needs to give sudo for port 80 forward.
			- Setup - `make -C ./Kubernetes setup` (or `reset` for clean and setup)
			- Teardown - `make -C ./Kubernetes clean`
			- Useful Targets: Istio Setup `istio`,  Info (`info` or `infos`), Portforward `port`, K8 Dashboard `dashboard`.
			- ![Minikube](common/images/fun-app/minikube.gif)
		- ### Services
			- Package has multiple service which can be setup on top of Minikube. This helps in easy setup of complex dependencies like Mysql Cluster, Load Testing, Mongo, Prometheus, Sonar and many more ...
			- Service Script allows you to configure a set of Services. This can then be installed, upgraded or deleted.
			- Select: `make -C ./Kubernetes/services select` will show a gum menu (Multiselect) to choose one or more services.
			- Setup: `make -C ./Kubernetes/services setup` (or `reset`) to deploy selected services.
			- Update: `make -C ./Kubernetes/services update` to update values or helm chart changes.
			- Stop: `make -C ./Kubernetes/services clean` to remove deployed helms. (Configured Services only)
			- Info: `make -C ./Kubernetes/services info` (or `infos`)
			- ![K8 Service](common/images/fun-app/k8-service.gif)
	- ### Log Analyzer
		- Monitor Logs via [GoAccess](https://github.com/allinurl/goaccess)
			- Terminal Access: `go run main.go | goaccess --log-format='%^ %d - %t | %s | %~%D | %b | %~%h | %^ | %m %U' --date-format='%Y/%m/%d' --time-format '%H:%M:%S'`
			- Web Access
				- Add Flags to Above Command `-o report.html --real-time-html`.
				- Open report.html in Browser and it should auto refresh.
		- **Useful Fields**
			- Mandatory Fields: %d (Date), %h (Host), %r/%m %U (Request)
			- Skip: Ignore (%^) , Skip Space (%~)
			- DateTime: (%x/--datetime-format) OR Time (%t/--date-format) + Date (%d/--time-format)
			- Host: IP (%h) OR Virtual Host (%v)
			- Request: Full With Quotes (%r) or Method (%m), URL (%U), Query (%q), PROTOCOL (%H),
			- Response: Status Code (%s), Size (%b)
			- Latency: MicroSecond (%D), MilliSecond.MicroSecond (%T), MilliSecond With Decimal (%L)
			- User Info: User-Agent (%u), Referrer (%R)
				
				TODO: Add GIF  
		- **Custom Log Monitoring**
			- Custom Logs require configuring various flags.
			- Identify [Date and Time Format](https://www.freebsd.org/cgi/man.cgi?query=strftime&sektion=3).
				- Configured via flag `--date-format` and `--time-format`
				- Verify Format, bash run: `date '+%Y/%m/%d - %H:%M:%S'` for output `2023/01/23 - 14:38:2`
			- Identify [Log Format](https://goaccess.io/man#custom-log)
				- Configured via flag `--log-format`
				- Start with initial fields and progress further for easy debug.
			- Debug Mode: `-l debug.log`
	- `TODO` Tree [Extension](https://marketplace.visualstudio.com/items?itemName=Gruntfuggly.todo-tree)
		- `TODO` Long Planned Enhancement in Code
		- `FIXME` Medium Sized fix or Unhandled Case.
		- `HACK` Remove Hacky way to Do Things.
		- `BUG` Small Bug which is observed but need Fix.
		- `XXX` Someday Tasks which may or may not be Done.
		- `#A` `#B` `#C` for Priortizations. Eg. `BUG: #A Some Message`
- ## Golang
  This section covers common practices and tools for Golang.  
	- ### Module Management
	  This is multi module project. Each module has its own go mod file. Modules can be managed using [semver](https://semver.org/) tags. Eg. v1.0.0  
		- Sync Modules using `make sync`. This is automatically done before builds.
		- Mod
			- New Module run `go mod init github.com/amanhigh/go-fun/components/fun-app` and to work using `go work use ./components/fun-app`
			- Link Module to new Release using `go mod tidy` or `go get -u github.com/amanhigh/go-fun/models`
			- Recursive Depdency Update `find . -name "go.mod" -execdir sh -c 'go get -u && go mod tidy';` (Run it in ProjectBase Dir)
		- Tags
			- See existing tags. `git tag | grep common`
			- Tag New Release. `git tag common/v1.0.0` followed by `git push --tags`
			- Remove Release involves deleting Tag with `git push --delete origin common/v1.0.0`
	- ### Release Management
	  Release management includes  build and release of Artifacts like binaries, dockers etc.  
		- Build Only - `make build docker-build` (add `clean` to remove residue)
		- Release - `make release release-docker VER=v1.0.3`
		- Delete Release - `make unrelease VER=v1.0.3` (Not Recommended)
