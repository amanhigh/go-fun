set shell := ["bash", "-cu"]

import './lib.just'

default:
	just --list

[group('go')]
generate: _generate-mocks _generate-swagger _generate-templ

[group('go')]
load: _load

[group('go')]
profile: _profile

[group('go')]
info: _info-release _info-docker

[group('go')]
infos: info _space-info

[group('go')]
install: _install-kohan

[group('go')]
prepare: _setup-tools setup-k8 _install-deadcode

[group('go')]
setup: _setup

[group('go')]
reset: _reset

[group('go')]
all: _all

[group('go')]
release-docker: docker-build
    #!/usr/bin/env bash
    set -euo pipefail
    if [[ -z "${VER:-}" ]]; then
        echo "VER not set. Use VER=v1.0.5"
        exit 1
    fi
    just _title "RELEASE" "Release Docker Images: ${VER}"
    just _detail "RELEASE" "Docker Tag"
    docker tag amanfdk/fun-app:latest amanfdk/fun-app:${VER}
    just _detail "RELEASE" "Docker Push"
    docker push amanfdk/fun-app:latest
    docker push amanfdk/fun-app:${VER}

[group('go')]
release-helm:
    #!/usr/bin/env bash
    set -euo pipefail
    if [[ -z "${VER:-}" ]]; then
        echo "VER not set. Eg. 1.1.0"
        exit 1
    fi
    just _title "RELEASE" "Release Helm Charts: ${VER}"
    just _helm-package
    git add components/fun-app/charts/Chart.yaml
    git commit -m "Helm Released: ${VER}"
    just _info "RELEASE" "Release: https://github.com/amanhigh/go-fun/actions/workflows/release.yml"

[group('go')]
release: _info-release
    #!/usr/bin/env bash
    set -euo pipefail
    if [[ -z "${VER:-}" ]]; then
        echo "VER not set. Eg. v1.1.0"
        exit 1
    fi
    just _release-models
    just _release-common
    just _release-fun

[group('go')]
run: _build-fun
    just _title "RUN" "Running Fun App"
    out="${OUT:-/dev/stdout}"
    bin/fun > "$out"

[group('go')]
analyse:
    just _title "Analyse" "Fun App Logs"
    OUT=/dev/stdout just run 2>/dev/null | grep GIN | goaccess --log-format='%^ %d - %t | %s | %~%D | %b | %~%h | %^ | %m %U' --date-format='%Y/%m/%d' --time-format='%H:%M:%S'

[group('go')]
watch:
    #!/usr/bin/env bash
    set -euo pipefail
    cmd="${CMD:-ls}"
    just _title "Watch (entr)" "${cmd}"
    find . | entr -s "date +%M:%S; ${cmd}"

[group('go')]
space: space-purge
    just _title "Starting" "Devspace"
    devspace use namespace fun-app
    devspace dev

[group('go')]
space-purge:
    just _title "Purging" "Devspace"
    -devspace purge

[group('go')]
space-test:
    just _title "TEST" "Devspace Tests"
    devspace run ginkgo
    CMD="devspace run fun-test" just watch

[group('go')]
setup-k8:
    just _title "SETUP" "Kubernetes Setup"
    make -C ./Kubernetes/services helm hosts

[group('go')]
docker-build: _docker-fun
    just _info "Docker Hub" "https://hub.docker.com/r/amanfdk/fun-app/tags"

[group('go')]
pack:
    just _title "Pack" "Repository"
    repomix --style markdown .

[group('go')]
verify: _test-focus
    just _info "INFO" "just watch cmd='just verify'"

[group('go')]
unrelease:
    #!/usr/bin/env bash
    set -euo pipefail
    if [[ -z "${VER:-}" ]]; then
        echo "VER not set. Eg. v1.1.0"
        exit 1
    fi
    just _warn "DELETE" "Release: ${VER}"
    just _info-release
    if just _confirm; then
        git tag -d models/${VER}
        git push --delete origin models/${VER}
        git tag -d common/${VER}
        git push --delete origin common/${VER}
        git tag -d ${VER}
        git push --delete origin ${VER}
    fi
    just _info-release

[group('go')]
lint: _lint-fix _lint-ci

[group('go')]
_build-fun:
    just _title "BUILD" "Building Fun App"
    mkdir -p bin
    CGO_ENABLED=0 GOARCH=amd64 go build -o bin/fun components/fun-app/main.go

[group('go')]
_build-fun-cover:
    just _title "BUILD" "Building Fun App with Coverage"
    mkdir -p bin
    CGO_ENABLED=0 GOARCH=amd64 go build -cover -o bin/fun components/fun-app/main.go

[group('go')]
_build-kohan:
    just _title "BUILD" "Building Kohan"
    mkdir -p bin
    CGO_ENABLED=1 GOARCH=amd64 go build -o bin/kohan components/kohan/main.go

[group('go')]
build: format lint _build-fun _build-kohan

[group('go')]
_test-operator:
    just _title "TEST" "Running Operator Tests"
    mkdir -p /tmp/cover/operator
    make -C components/operator/ test GOCOVERDIR=/tmp/cover/operator

[group('go')]
_test-unit:
    just _title "TEST" "Running Unit Tests"
    mkdir -p /tmp/cover/unit
    {{ginkgo}} -r --label-filter=\!setup\ \&\&\ \!slow --skip-package=components/fun-app/it -cover . -- -test.gocoverdir=/tmp/cover/unit

[group('go')]
test-slow:
    just _title "TEST" "Running Slow Tests"
    mkdir -p /tmp/cover/slow
    {{ginkgo}} -r '--label-filter=slow' -cover . -- -test.gocoverdir=/tmp/cover/slow

[group('go')]
_run-fun-cover: _build-fun-cover
    just _title "COVER" "Running Fun App with Coverage"
    mkdir -p /tmp/cover/integration
    GOCOVERDIR=/tmp/cover/integration PORT=8085 bin/fun >/dev/null 2>&1 &

[group('go')]
_cover-report:
    #!/usr/bin/env bash
    set -euo pipefail
    if [ ! -f "/tmp/cover/coverage-combined.out" ]; then
        just _warn "COVER" "No coverage profile found. Run 'just cover' first."
        exit 1
    fi
    just _title "COVER" "Coverage Analysis Report"
    go tool cover -func=/tmp/cover/coverage-combined.out
    overall=$(go tool cover -func=/tmp/cover/coverage-combined.out 2>/dev/null | tail -1 | awk '{print $3}')
    echo "Overall Coverage: $overall"
    echo ""
    just _title "COVER" "Packages by Coverage (Lowest to Highest)"
    go tool cover -func=/tmp/cover/coverage-combined.out 2>/dev/null | \
        awk '/github.com\/amanhigh\/go-fun\// { \
            gsub(/github.com\/amanhigh\/go-fun\//, "", $1); \
            gsub(/\/[^\/]*\.go:.*/, "", $1); \
            gsub(/%/, "", $3); \
            if ($1 != prev_pkg) { \
                if (prev_pkg != "") print prev_pkg, int(total/count); \
                prev_pkg = $1; total = $3; count = 1; \
            } else { \
                total += $3; count++; \
            } \
        } \
        END { if (prev_pkg != "") print prev_pkg, int(total/count) }' | \
        sort -k2 -n | \
        awk 'BEGIN{critical=0; low=0; medium=0; good=0} \
        { \
            pct = int($2); \
            if (pct == 0) { \
                print "\033[31m🔴 " sprintf("%-40s %6s%%", $1, $2) "\033[0m ← NEEDS TESTS!"; \
                critical++; \
            } else if (pct < 25) { \
                print "\033[31m🟠 " sprintf("%-40s %6s%%", $1, $2) "\033[0m ← CRITICAL"; \
                low++; \
            } else if (pct < 50) { \
                print "\033[33m🟡 " sprintf("%-40s %6s%%", $1, $2) "\033[0m ← LOW"; \
                low++; \
            } else if (pct < 75) { \
                print "\033[34m🔵 " sprintf("%-40s %6s%%", $1, $2) "\033[0m ← MEDIUM"; \
                medium++; \
            } else { \
                print "\033[32m🟢 " sprintf("%-40s %6s%%", $1, $2) "\033[0m ← GOOD"; \
                good++; \
            } \
        } \
        END { \
            print ""; \
            print "\033[32m[Summary]\033[0m"; \
            print "🔴 Critical (0%): " critical " packages"; \
            print "🟠 Low (<50%): " low " packages"; \
            print "🔵 Medium (50-75%): " medium " packages"; \
            print "🟢 Good (≥75%): " good " packages"; \
        }'

[group('go')]
_cover-analyse: combine-coverage _cover-report
    just _title "COVER" "Analysing Coverage Reports"
    go tool cover -func=/tmp/cover/coverage-combined.out
    echo ""
    go tool cover -html=/tmp/cover/coverage-combined.out -o /tmp/coverage.html
    just _info "HTML Report: file:///tmp/coverage.html"
    just _info "Vscode" "go.apply.coverprofile /tmp/cover/coverage-combined.out"

[group('go')]
cover: _clean-test _test-unit _run-fun-cover _cover-analyse

[group('go')]
test: cover _test-operator

[group('go')]
cover-report: _cover-report

[group('go')]
combine-coverage:
    just _title "COVER" "Combining Binary Coverage Data"
    coverage_dirs=""; \
    for dir in /tmp/cover/unit /tmp/cover/integration /tmp/cover/operator /tmp/cover/slow; do \
        if [ -d "$dir" ] && [ -n "$(ls -A "$dir" 2>/dev/null)" ]; then \
            coverage_dirs="$coverage_dirs,$dir"; \
            just _info "COVER" "Found coverage data in $(basename "$dir")"; \
        fi; \
    done; \
    if [ -n "$coverage_dirs" ]; then \
        coverage_dirs=${coverage_dirs#,}; \
        just _info "COVER" "Merging coverage from: $coverage_dirs"; \
        go tool covdata textfmt -i=$coverage_dirs -o /tmp/cover/coverage-combined.out; \
        just _title "COVER" "Combined coverage created: /tmp/cover/coverage-combined.out"; \
    else \
        just _warn "COVER" "No coverage data found to combine"; \
    fi

[group('go')]
_clean-test:
    just _warn "CLEAN" "Cleaning Tests"
    rm -rf /tmp/cover

[group('go')]
_clean-build:
    just _warn "CLEAN" "Cleaning Build"
    -rm -rf bin
    -make -C components/operator/ clean

[group('go')]
clean: _clean-test _clean-build

[group('go')]
_sync:
    just _title "GO" "Go Module Syncing"
    go work sync

[group('go')]
_lint-ci:
    #!/usr/bin/env bash
    set -euo pipefail
    if [[ "${GITHUB_ACTIONS:-}" == "true" ]]; then
        just _warn "LINT" "Skipping in GitHub Actions Environment"
    else
        just _title "LINT" "Golang CLI"
        go work edit -json | jq -r '.Use[].DiskPath' | xargs -I{} {{golangci_lint}} run {}/...
    fi

[group('go')]
_lint-fix:
    just _title "LINT" "Auto-fixing Modernize Issues"
    go work edit -json | jq -r '.Use[].DiskPath' | xargs -I{} {{golangci_lint}} run --enable-only=modernize --fix {}/...

[group('go')]
_generate-mocks:
    just _title "Generate" "Mocks"
    {{mockery}}

[group('go')]
_install-deadcode:
    just _title "Installing" "DeadCode"
    go install golang.org/x/tools/cmd/deadcode@latest

[group('go')]
_generate-swagger:
    just _title "Generate" "Swagger"
    cd components/fun-app && {{swag}} i --parseDependency true
    just _info "Swagger" "http://localhost:8080/swagger/index.html"

[group('go')]
_generate-templ:
    just _title "Generate" "Template Files"
    just _templ common/ui
    just _templ components/learn
    just _templ components/kohan
    just format

[group('go')]
_info-release:
    just _info "INFO" "Go Modules"
    git tag | grep "models" | tail -2
    git tag | grep "common" | tail -2
    git tag | grep "v" | grep -v "/" | tail -2

[group('go')]
_info-docker:
    just _info "INFO" "FunApp DockerHub: https://hub.docker.com/r/amanfdk/fun-app/tags"
    curl -s "https://hub.docker.com/v2/repositories/amanfdk/fun-app/tags/?page_size=25&page=1&name&ordering" | jq -r '.results[]|.name' | head -3
    just _info "INFO" "Docker Images: amanfdk/fun-app"
    docker images | grep fun-app

[group('go')]
_install-kohan:
    just _title "Installing" "Kohan"
    CGO_ENABLED=1 GOARCH=amd64 go install ./components/kohan

[group('go')]
_helm-package:
    VERSION="${VERSION:-${VER:-}}" make -C components/fun-app/charts package

[group('go')]
_setup-tools:
    just _title "Setting up" "Tools"
    go install github.com/onsi/ginkgo/v2/ginkgo
    go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.2
    go install github.com/swaggo/swag/cmd/swag
    go install golang.org/x/tools/cmd/goimports@latest
    go install github.com/vektra/mockery/v3@v3.5.5
    go install github.com/a-h/templ/cmd/templ@latest
    go install github.com/air-verse/air@latest

[group('go')]
_docker-fun: _build-fun
    just _title "Docker" "Building FunApp"
    cp bin/fun components/fun-app/fun
    docker buildx build -t amanfdk/fun-app -f components/fun-app/Dockerfile components/fun-app 2>/dev/null
    rm components/fun-app/fun

[group('go')]
_docker-fun-run: _docker-fun
    just _title "Docker" "Running FunApp"
    docker run -it amanfdk/fun-app

[group('go')]
_docker-fun-exec:
    just _title "Docker" "Execing Into FunApp"
    docker run -it --entrypoint /bin/sh amanfdk/fun-app

[group('go')]
_docker-fun-clean:
    just _warn "Docker" "Deleting FunApp"
    -docker rmi -f `docker images amanfdk/fun-app -q`

[group('go')]
_space-info:
    just _title "INFO" "Devspace"
    devspace list vars --var DB="mysql-primary",RATE_LIMIT=-1
    just _detail "INFO" "http://localhost:8080/metrics"
    just _detail "INFO" "Login: devspace enter"

[group('go')]
_test-focus:
    just _title "TEST" "Running Focus Tests"
    {{ginkgo}} --focus "should create & get person" components/fun-app/it

[group('go')]
_load:
    just _title "Load" "Test Fun App"
    make -C components/fun-app/it all

[group('go')]
_profile:
    #!/usr/bin/env bash
    set -euo pipefail
    endpoint="${ENDPOINT:-http://localhost:8080}"
    just _title "PROFILE" "ENDPOINT=${endpoint} | http://app.docker/app"
    just _detail "PROFILE" "Profiling Heap"
    go tool pprof -http=:8001 "${endpoint}/debug/pprof/heap" 2>/dev/null &
    just _detail "PROFILE" "Profiling CPU"
    go tool pprof -http=:8000 --seconds=30 "${endpoint}/debug/pprof/profile" 2>/dev/null
    just _warn "PROFILE" "Killing Profilers"
    kill %1

[group('go')]
_confirm:
    #!/usr/bin/env bash
    set -euo pipefail
    if [[ -z "${CI:-}" ]]; then
        reply=""
        read -p "⚠ Are you sure? [y/n] > " -r
        if [[ ! ${REPLY} =~ ^[Yy]$ ]]; then
            just _warn "KO" "Stopping"
            exit 1
        fi
        just _title "OK" "Continuing"
    fi

[group('go')]
_release-models:
    #!/usr/bin/env bash
    set -euo pipefail
    just _title "RELEASE" "Release Models: ${VER}"
    if just _confirm; then
        git tag models/${VER}
        git tag | grep models | tail -2
    fi
    just _title "RELEASE" "Pushing Tags"
    if just _confirm; then
        git push --tags && just _title "RELEASE" "Models Released: ${VER}"
    fi

[group('go')]
_release-common:
    #!/usr/bin/env bash
    set -euo pipefail
    just _title "RELEASE" "Bump Models: ${VER}"
    if just _confirm; then
        pushd ./common
        go get -u github.com/amanhigh/go-fun/models@${VER}
        git add go.* && git commit -m "Bumping Models: ${VER}"
        popd
    fi
    just _title "RELEASE" "Release Common: ${VER}"
    if just _confirm; then
        git tag common/${VER}
        git tag | grep common | tail -2
    fi
    just _title "RELEASE" "Pushing Tags"
    if just _confirm; then
        git push --tags && just _title "RELEASE" "Common Released: ${VER}"
    fi

[group('go')]
_release-fun:
    #!/usr/bin/env bash
    set -euo pipefail
    just _title "RELEASE" "Bump Common: ${VER}"
    if just _confirm; then
        pushd ./components/fun-app
        go get -u github.com/amanhigh/go-fun/common@${VER}
        git add go.* && git commit -m "Bumping Common: ${VER}"
        popd
    fi
    just _title "RELEASE" "Release Fun: ${VER}"
    if just _confirm; then
        git tag ${VER}
        just _info-release
    fi
    just _title "RELEASE" "Pushing Tags"
    if just _confirm; then
        git push --tags && just _title "RELEASE" "Fun Released: ${VER}"
    fi

[group('go')]
_lint-dead:
    just _title "LINT" "DeadCode"
    go work edit -json | jq -r '.Use[].DiskPath' | sed 's|^\./||' | grep -vE "common|models|components/learn" | xargs -I{} {{deadcode}} github.com/amanhigh/go-fun/{}/...

[group('go')]
_setup: _sync test generate build _lint-dead _helm-package docker-build

[group('go')]
_reset: _setup info clean

[group('go')]
_all: prepare _docker-fun-clean install _reset infos test-slow
