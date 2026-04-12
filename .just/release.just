[group('info')]
[doc('Show release tags and Docker image information')]
info: _info-release _info-docker

[group('info')]
[doc('Show release info plus devspace information')]
infos: info space-info

[confirm("Release Helm charts?")]
[group('release')]
[doc('Package and commit the Helm chart release')]
release-helm:
    #!/usr/bin/env bash
    set -euo pipefail
    if [[ -z "${VER:-}" ]]; then
        echo "VER not set. Eg. 1.1.0"
        exit 1
    fi
    just _title "RELEASE" "Release Helm Charts: ${VER}"
    VERSION="${VER}" just _helm-package
    git add components/fun-app/charts/Chart.yaml
    git commit -m "Helm Released: ${VER}"
    just _info "RELEASE" "Release: https://github.com/amanhigh/go-fun/actions/workflows/release.yml"

[confirm("Release Go modules and tags?")]
[group('release')]
[doc('Release Go module tags for models, common, and the Fun App')]
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

[confirm("Delete published Go module tags?")]
[group('release')]
[doc('Delete previously published Go module tags')]
unrelease:
    #!/usr/bin/env bash
    set -euo pipefail
    if [[ -z "${VER:-}" ]]; then
        echo "VER not set. Eg. v1.1.0"
        exit 1
    fi
    just _warn "DELETE" "Release: ${VER}"
    just _info-release
    git tag -d models/${VER}
    git push --delete origin models/${VER}
    git tag -d common/${VER}
    git push --delete origin common/${VER}
    git tag -d ${VER}
    git push --delete origin ${VER}
    just _info-release

[group('info')]
_info-release:
    just _info "INFO" "Go Modules"
    git tag | grep "models" | tail -2
    git tag | grep "common" | tail -2
    git tag | grep "v" | grep -v "/" | tail -2

[group('info')]
_info-docker:
    just _info "INFO" "FunApp DockerHub: https://hub.docker.com/r/amanfdk/fun-app/tags"
    curl -s "https://hub.docker.com/v2/repositories/amanfdk/fun-app/tags/?page_size=25&page=1&name&ordering" | jq -r '.results[]|.name' | head -3
    just _info "INFO" "Docker Images: amanfdk/fun-app"
    docker images | grep fun-app

[group('release')]
_helm-package:
    #!/usr/bin/env bash
    set -euo pipefail
    if [[ -n "${VERSION:-}" || -n "${VER:-}" ]]; then
        VERSION="${VERSION:-${VER:-}}" make -C components/fun-app/charts package
    else
        make -C components/fun-app/charts package
    fi

[group('release')]
_release-models:
    #!/usr/bin/env bash
    set -euo pipefail
    just _title "RELEASE" "Release Models: ${VER}"
    git tag models/${VER}
    git tag | grep models | tail -2
    just _title "RELEASE" "Pushing Tags"
    git push --tags && just _title "RELEASE" "Models Released: ${VER}"

[group('release')]
_release-common:
    #!/usr/bin/env bash
    set -euo pipefail
    just _title "RELEASE" "Bump Models: ${VER}"
    pushd ./common
    go get -u github.com/amanhigh/go-fun/models@${VER}
    git add go.* && git commit -m "Bumping Models: ${VER}"
    popd
    just _title "RELEASE" "Release Common: ${VER}"
    git tag common/${VER}
    git tag | grep common | tail -2
    just _title "RELEASE" "Pushing Tags"
    git push --tags && just _title "RELEASE" "Common Released: ${VER}"

[group('release')]
_release-fun:
    #!/usr/bin/env bash
    set -euo pipefail
    just _title "RELEASE" "Bump Common: ${VER}"
    pushd ./components/fun-app
    go get -u github.com/amanhigh/go-fun/common@${VER}
    git add go.* && git commit -m "Bumping Common: ${VER}"
    popd
    just _title "RELEASE" "Release Fun: ${VER}"
    git tag ${VER}
    just _info-release
    just _title "RELEASE" "Pushing Tags"
    git push --tags && just _title "RELEASE" "Fun Released: ${VER}"
