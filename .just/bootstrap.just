[group('setup')]
[doc('Install the Kohan CLI locally')]
install:
    just _title "Installing" "Kohan"
    CGO_ENABLED=1 GOARCH=amd64 go install ./components/kohan

[group('setup')]
_setup-gotools:
    just _title "SETUP" "Installing ginkgo"
    GOBIN={{gobin}} go install {{ginkgo_install}}
    just _title "SETUP" "Installing golangci-lint"
    GOBIN={{gobin}} go install {{golangci_lint_install}}
    just _title "SETUP" "Installing swag"
    GOBIN={{gobin}} go install {{swag_install}}
    just _title "SETUP" "Installing goimports"
    GOBIN={{gobin}} go install {{goimports_install}}
    just _title "SETUP" "Installing mockery"
    GOBIN={{gobin}} go install {{mockery_install}}
    just _title "SETUP" "Installing templ"
    GOBIN={{gobin}} go install {{templ_install}}
    just _title "SETUP" "Installing air"
    GOBIN={{gobin}} go install {{air_install}}
    just _title "SETUP" "Installing deadcode"
    GOBIN={{gobin}} go install {{deadcode_install}}

[group('setup')]
_setup-k8:
    just _title "SETUP" "Kubernetes Setup"
    make -C ./Kubernetes/services helm hosts

[group('setup')]
_doctor:
    #!/usr/bin/env bash
    set -uo pipefail

    errors=0
    warnings=0

    ok() {
        printf '✓ %-24s %s\n' "$1" "$2"
    }

    error() {
        printf '✗ %-24s %s\n' "$1" "$2"
        errors=$((errors + 1))
    }

    warn() {
        printf '○ %-24s %s\n' "$1" "$2"
        warnings=$((warnings + 1))
    }

    command_status() {
        command -v "$1" >/dev/null 2>&1
    }

    expected_go=$(sed -n 's/^go[[:space:]]\+//p' go.work | head -n 1)
    if ! command_status go; then
        error "Go" "Required to build, test, and install repository tools"
    else
        actual_go=$(go version | awk '{print $3}' | sed 's/^go//')
        if [[ -n "$expected_go" && "$actual_go" != "$expected_go" ]]; then
            error "Go version" "Found $actual_go; repository requires $expected_go"
        else
            ok "Go" "$actual_go"
        fi
    fi

    for command_name in jq make; do
        if command_status "$command_name"; then
            ok "$command_name" "available"
        else
            error "$command_name" "Required by repository recipes"
        fi
    done

    if command_status cc || command_status gcc; then
        ok "C compiler" "available for the CGO-enabled Kohan build"
    else
        error "C compiler" "Required to build the Kohan binary"
    fi

    for tool in "{{ginkgo}}" "{{golangci_lint}}" "{{swag}}" "{{goimports}}" "{{mockery}}" "{{templ}}" "{{air}}" "{{deadcode}}"; do
        if [[ -x "$tool" ]]; then
            ok "$(basename "$tool")" "available"
        else
            error "$(basename "$tool")" "Install with just prepare"
        fi
    done

    if [[ -n "${WAYLAND_DISPLAY:-}" ]]; then
        ok "Desktop display" "Wayland session detected"
    elif [[ -n "${DISPLAY:-}" ]]; then
        ok "Desktop display" "X11 session detected"
    elif command_status xvfb-run && command_status xauth; then
        ok "Headless display" "xvfb-run and xauth available"
    else
        error "Headless display" "just test cannot run clipboard tests without a display"
        printf '  Install on Debian/Ubuntu: sudo apt-get update && sudo apt-get install -y xvfb xauth\n'
    fi

    if command_status helm; then
        ok "Helm" "available"
    else
        warn "Helm" "prepare-k8s and Helm packaging will not work"
    fi

    if command_status kubectl; then
        ok "kubectl" "available"
    else
        warn "kubectl" "Kubernetes deployment commands will not work"
    fi

    if command_status docker; then
        ok "Docker" "available"
    else
        warn "Docker" "Docker image build commands will not work"
    fi

    if [[ -r /etc/hosts ]] && grep -Eq '(^|[[:space:]])[^[:space:]]+\.docker([[:space:]]|$)' /etc/hosts; then
        ok "Kubernetes hosts" ".docker entries found in /etc/hosts"
    else
        warn "Kubernetes hosts" "local .docker ingress URLs will not resolve; run just prepare-k8s"
    fi

    printf '\nEnvironment check: %d error(s), %d warning(s)\n' "$errors" "$warnings"
    if (( errors > 0 )); then
        exit 1
    fi

[group('setup')]
_sync:
    just _title "GO" "Go Module Syncing"
    go work sync
