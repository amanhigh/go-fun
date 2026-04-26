[group('core')]
[doc('Run lint auto-fixes and full lint checks')]
lint: _lint-fix _lint-ci

[group('quality')]
_lint-ci:
    #!/usr/bin/env bash
    set -euo pipefail
    if [[ "${GITHUB_ACTIONS:-}" == "true" ]]; then
        just _warn "LINT" "Skipping in GitHub Actions Environment"
    else
        just _title "LINT" "Golang CLI"
        go work edit -json | jq -r '.Use[].DiskPath' | xargs -I{} {{golangci_lint}} run --show-stats=false {}/...
    fi

[group('quality')]
_lint-fix:
    just _title "LINT" "Auto-fixing Modernize Issues"
    go work edit -json | jq -r '.Use[].DiskPath' | xargs -I{} {{golangci_lint}} run --show-stats=false --enable-only=modernize --fix {}/...

[group('quality')]
[doc('Detect dead code in selected modules')]
lint-dead:
    just _title "LINT" "DeadCode"
    go work edit -json | jq -r '.Use[].DiskPath' | sed 's|^\./||' | grep -vE "common|models|components/learn" | xargs -I{} {{deadcode}} github.com/amanhigh/go-fun/{}/...
