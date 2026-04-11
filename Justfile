set shell := ["bash", "-cu"]

import './lib.just'

default:
	just --list

[group('go')]
_sync:
    just _title "Go Module Syncing"
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
_template:
    just _title "Generate" "Template Files"
    just _templ common/ui
    just _templ components/learn
    just _templ components/kohan
    just format
