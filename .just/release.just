[confirm("Release Helm charts?")]
[group('release')]
[doc('Package and commit the Helm chart release')]
release-helm version:
    #!/usr/bin/env bash
    set -euo pipefail
    version={{quote(version)}}
    just _title "RELEASE" "Release Helm Charts: ${version}"
    just _helm-package "${version}"
    git add components/fun-app/charts/Chart.yaml
    git commit -m "Helm Released: ${version}"
    just _info "RELEASE" "Release: https://github.com/amanhigh/go-fun/actions/workflows/release.yml"

[confirm("Release Go modules and tags?")]
[group('release')]
[doc('Release Go module tags for models, common, and the Fun App')]
release version: _info-release
    #!/usr/bin/env bash
    set -euo pipefail
    version={{quote(version)}}
    just _release-models "${version}"
    just _release-common "${version}"
    just _release-fun "${version}"

[confirm("Delete published Go module tags?")]
[group('release')]
[doc('Delete previously published Go module tags')]
unrelease version:
    #!/usr/bin/env bash
    set -euo pipefail
    version={{quote(version)}}
    just _warn "DELETE" "Release: ${version}"
    just _info-release
    git tag -d "models/${version}"
    git push --delete origin "models/${version}"
    git tag -d "common/${version}"
    git push --delete origin "common/${version}"
    git tag -d "${version}"
    git push --delete origin "${version}"
    just _info-release

[group('release')]
_helm-package version='':
    #!/usr/bin/env bash
    set -euo pipefail
    version={{quote(version)}}
    if [[ -n "${version}" ]]; then
        just --set version "${version}" components/fun-app/charts/package
    else
        just components/fun-app/charts/package
    fi

[group('release')]
_release-models version:
    #!/usr/bin/env bash
    set -euo pipefail
    version={{quote(version)}}
    just _title "RELEASE" "Release Models: ${version}"
    git tag "models/${version}"
    git tag | grep models | tail -2
    just _title "RELEASE" "Pushing Tags"
    git push --tags && just _title "RELEASE" "Models Released: ${version}"

[group('release')]
_release-common version:
    #!/usr/bin/env bash
    set -euo pipefail
    version={{quote(version)}}
    just _title "RELEASE" "Bump Models: ${version}"
    pushd ./common
    go get -u "github.com/amanhigh/go-fun/models@${version}"
    git add go.* && git commit -m "Bumping Models: ${version}"
    popd
    just _title "RELEASE" "Release Common: ${version}"
    git tag "common/${version}"
    git tag | grep common | tail -2
    just _title "RELEASE" "Pushing Tags"
    git push --tags && just _title "RELEASE" "Common Released: ${version}"

[group('info')]
_info-release:
    just _info "INFO" "Go Modules"
    git tag | grep "models" | tail -2
    git tag | grep "common" | tail -2
    git tag | grep "v" | grep -v "/" | tail -2

[group('release')]
_release-fun version:
    #!/usr/bin/env bash
    set -euo pipefail
    version={{quote(version)}}
    just _title "RELEASE" "Bump Common: ${version}"
    pushd ./components/fun-app
    go get -u "github.com/amanhigh/go-fun/common@${version}"
    git add go.* && git commit -m "Bumping Common: ${version}"
    popd
    just _title "RELEASE" "Release Fun: ${version}"
    git tag "${version}"
    just _info-release
    just _title "RELEASE" "Pushing Tags"
    git push --tags && just _title "RELEASE" "Fun Released: ${version}"
