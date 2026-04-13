set quiet := true

# Shared tool binaries
tailwind := "tailwindcss"
tsc := "tsc"
esbuild := "esbuild"
gobin := `go env GOPATH 2>/dev/null` + "/bin"
templ_install := "github.com/a-h/templ/cmd/templ@latest"
templ := gobin + "/templ"
goimports_install := "golang.org/x/tools/cmd/goimports@latest"
goimports := gobin + "/goimports"
golangci_lint_install := "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.2"
golangci_lint := gobin + "/golangci-lint"
ginkgo_install := "github.com/onsi/ginkgo/v2/ginkgo@v2.13.2"
ginkgo := gobin + "/ginkgo"
deadcode_install := "golang.org/x/tools/cmd/deadcode@latest"
deadcode := gobin + "/deadcode"
swag_install := "github.com/swaggo/swag/cmd/swag@latest"
swag := gobin + "/swag"
mockery_install := "github.com/vektra/mockery/v3@v3.5.5"
mockery := gobin + "/mockery"
air_install := "github.com/air-verse/air@latest"
air := gobin + "/air"
log_file := "/tmp/go-fun-build.log"
root := `git rev-parse --show-toplevel`

# Print helpers (private)
_info task info:
    echo '{{ITALIC}}{{YELLOW}}[{{task}}]{{NORMAL}} {{info}} 🔔'

_detail task info:
    echo '{{CYAN}}[{{task}}]{{NORMAL}} {{info}}'

_title task info:
    echo '{{BOLD}}{{GREEN}}[{{task}}]{{NORMAL}} {{info}} 🎯'

_warn task info:
    echo '{{YELLOW}}[{{task}}]{{NORMAL}} {{info}} ⚠️'

_run command:
    if [ "${GF_CONCISE:-0}" != "1" ]; then \
        eval '{{command}}'; \
    else \
        { printf '\n=== %s ===\n' "$(date '+%Y-%m-%d %H:%M:%S')"; eval '{{command}}'; } >>'{{log_file}}' 2>&1 || { \
            just _warn "RUN" "Command failed. Last 40 log lines:"; \
            just _detail "RUN" "Detailed log: {{log_file}}"; \
            tail -n 40 '{{log_file}}' >&2 || true; \
            exit 1; \
        }; \
    fi

_run-no-output command:
    eval '{{command}}' 2>/dev/null

_format dir=".":
    just _info "GO" "Format"
    cd {{dir}} && {{goimports}} -w .

[group('build')]
_templ dir=".":
    just _info "Templ" "Generating templ files"
    cd {{dir}} && {{templ}} generate

# Generate CSS source list from templ files in specified directory.
_css_sources dir=".":
    #!/usr/bin/env bash
    set -euo pipefail
    cd {{dir}}
    TEMPLUI_PATH="$(go list -m -f '{{{{.Dir}}' github.com/templui/templui)"
    (
        find . -name "*.templ" -type f | sed 's|^\./||' | while read -r file; do
            dir=$(dirname "$file")
            echo "@source \"./${dir}/**/*.templ\";"
        done | sort -u
        echo "@source \"{{root}}/common/ui/**/*.templ\";"
        echo "@source \"${TEMPLUI_PATH}/components/**/*.templ\";"
    ) > ./assets/css/sources.generated.css

[group('build')]
_css dir=".":
    just _css_sources {{dir}}
    just _info "CSS" "Generating Tailwind bundle"
    cd {{dir}} && {{tailwind}} build --input ./assets/css/input.css --output ./assets/css/app.css

[group('build')]
_js dir=".":
    #!/usr/bin/env bash
    set -euo pipefail
    js_dir="{{dir}}/assets/js"
    bundle_entry="./assets/js/input.ts"
    bundle_output="./assets/js/app.js"

    if command -v {{tsc}} >/dev/null 2>&1; then
        just _info "JS" "Type-checking TypeScript"
        (cd "$js_dir" && {{tsc}} --noEmit --strict --skipLibCheck --lib ES2021,DOM --target ES2020 --module ESNext --moduleResolution bundler *.ts)
    else
        just _warn "JS" "TypeScript not found, skipping type-check"
    fi

    just _info "JS" "Bundling with esbuild"
    (cd {{dir}} && {{esbuild}} "$bundle_entry" --bundle --outfile="$bundle_output" --format=iife --target=es2020)

[group('run')]
_air-run dir port="8080":
    just _title "Air" "Starting server"
    just _detail "URL" "http://localhost:{{port}}"
    cd {{dir}} && {{air}} -c .air.toml
