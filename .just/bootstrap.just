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
    set -u

    just _title "DOCTOR" "Checking external prerequisites"

    if [[ -n "${WAYLAND_DISPLAY:-}" ]]; then
        just _info "DISPLAY" "Wayland session detected"
    elif [[ -n "${DISPLAY:-}" ]]; then
        just _info "DISPLAY" "X11 session detected"
    elif command -v xvfb-run >/dev/null 2>&1 && command -v xauth >/dev/null 2>&1; then
        just _info "DISPLAY" "xvfb-run and xauth available"
    else
        just _warn "DISPLAY" "Clipboard tests need a display or xvfb-run"
        just _detail "DISPLAY" "Debian/Ubuntu: sudo apt-get update && sudo apt-get install -y xvfb xauth"
    fi

    if command -v helm >/dev/null 2>&1; then
        just _info "HELM" "available for chart packaging and prepare-k8s"
    else
        just _warn "HELM" "Install Helm for chart packaging and prepare-k8s"
    fi

    if command -v kubectl >/dev/null 2>&1; then
        just _info "KUBECTL" "available for Kubernetes workflows"
    else
        just _warn "KUBECTL" "Kubernetes deployment workflows will not work"
    fi

    if command -v docker >/dev/null 2>&1; then
        just _info "DOCKER" "available for image builds"
    else
        just _warn "DOCKER" "Docker image build commands will not work"
    fi

[group('setup')]
_sync:
    just _title "GO" "Go Module Syncing"
    go work sync
