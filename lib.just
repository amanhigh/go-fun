set quiet := true

# Shared tool binaries
tailwind := "tailwindcss"
tsc := "tsc"
esbuild := "esbuild"
templ := "go run github.com/a-h/templ/cmd/templ@latest"
goimports := "go run golang.org/x/tools/cmd/goimports@latest"
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

# Format Go code in specified directory (recursive).
[group('go')]
format dir=".":
    just _info "GO" "Format"
    cd {{dir}} && {{goimports}} -w .

# Generate templ files in specified directory (recursive).
[group('go')]
templ dir=".":
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

# Run server with Air hot reload in specified directory.
[group('run')]
_air_run dir port="8080":
    just _title "Air" "Starting server"
    just _detail "URL" "http://localhost:{{port}}"
    cd {{dir}} && go run github.com/air-verse/air@latest -c .air.toml

# Generate Tailwind CSS bundle.
[group('css')]
css dir=".":
    just _css_sources {{dir}}
    just _info "CSS" "Generating Tailwind bundle"
    cd {{dir}} && {{tailwind}} build --input ./assets/css/input.css --output ./assets/css/app.css

# Type-check TypeScript sources.
[group('js')]
tsc dir=".":
    #!/usr/bin/env bash
    if command -v {{tsc}} >/dev/null 2>&1; then
        just _info "JS" "Type-checking TypeScript"
        set -euo pipefail
        cd {{dir}}/assets/js && {{tsc}} --noEmit --strict --skipLibCheck --lib ES2021,DOM --target ES2020 --module ESNext --moduleResolution bundler *.ts
    else
        just _warn "JS" "TypeScript not found, skipping type-check"
    fi

# Bundle TypeScript into JavaScript.
[group('js')]
js dir=".":
    just tsc {{dir}}
    just _info "JS" "Bundling with esbuild"
    cd {{dir}} && {{esbuild}} assets/js/input.ts --bundle --outfile=assets/js/app.js --format=iife --target=es2020
