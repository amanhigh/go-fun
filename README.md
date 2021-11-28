[![Actions Status](https://github.com/amanhigh/go-fun/workflows/Build/badge.svg)](https://github.com/amanhigh/go-fun/actions)
[![codecov](https://codecov.io/gh/amanhigh/go-fun/branch/master/graph/badge.svg)](https://codecov.io/gh/amanhigh/go-fun)


# Go Fun
Experiments and Fun Go Lang and Frameworks. Also includes tools like docker, k8, istio and perf (wrk)

## Build
Use goreleaser for build test. Install if not already installed

goreleaser build --snapshot --rm-dist
goreleaser release --snapshot --skip-publish --rm-dist

## Testing
Testing is handled via Ginkgo. To run all unit tests excluding any integration tests.

`ginkgo -r '--label-filter=!it' .`

## FunApp
Sample Funapp which is rest based app with various tools and tests required as sample.

By default it runs without any dependencies with in memory database which can be configured via ENV Variables.

Go
`go run ./components/fun-app/`

Docker (Post Build step)
`docker run amanfdk/fun-app`
