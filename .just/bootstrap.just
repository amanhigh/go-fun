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
_sync:
    just _title "GO" "Go Module Syncing"
    go work sync
