[group('run')]
[doc('Run the focused verification flow for the Fun App')]
verify: _test-focus
    just _info "INFO" "just watch cmd='just verify'"

[group('quality')]
_test-operator:
    just _title "TEST" "Running Operator Tests"
    mkdir -p /tmp/cover/operator
    just _run "make -C components/operator/ test GOCOVERDIR=/tmp/cover/operator"

[group('quality')]
_test-unit:
    just _title "TEST" "Running Unit Tests"
    mkdir -p /tmp/cover/unit
    just _run "{{ginkgo}} -r --label-filter=\!setup\ \&\&\ \!slow --skip-package=components/fun-app/it -cover . -- -test.gocoverdir=/tmp/cover/unit"

[group('quality')]
[doc('Run slow test suites')]
test-slow:
    just _title "TEST" "Running Slow Tests"
    mkdir -p /tmp/cover/slow
    just _run "{{ginkgo}} -r '--label-filter=slow' -cover . -- -test.gocoverdir=/tmp/cover/slow"

[group('quality')]
_run-fun-cover: _build-fun-cover
    just _title "COVER" "Running Fun App with Coverage"
    mkdir -p /tmp/cover/integration
    just _run "GOCOVERDIR=/tmp/cover/integration PORT=8085 bin/fun &"

[group('quality')]
_cover-report:
    #!/usr/bin/env bash
    set -euo pipefail
    if [ ! -f "/tmp/cover/coverage-combined.out" ]; then
        just _warn "COVER" "No coverage profile found. Run 'just cover' first."
        exit 1
    fi
    just _title "COVER" "Coverage Analysis Report"
    just _run-no-output "go tool cover -func=/tmp/cover/coverage-combined.out"
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

[group('quality')]
[doc('Run unit, integration, and coverage analysis workflows')]
cover: _clean-test _test-unit _run-fun-cover _combine-coverage
    just _title "COVER" "Analysing Coverage Reports"
    just _cover-report
    echo ""
    just _run-no-output "go tool cover -html=/tmp/cover/coverage-combined.out -o /tmp/coverage.html"
    just _info "COVER" "HTML Report: file:///tmp/coverage.html"
    just _info "Vscode" "go.apply.coverprofile /tmp/cover/coverage-combined.out"

[group('core')]
[doc('Run the main test suite including operator tests')]
test: cover _test-operator

[group('quality')]
_combine-coverage:
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

[group('clean')]
_clean-test:
    just _warn "CLEAN" "Cleaning Tests"
    rm -rf /tmp/cover

[group('quality')]
_test-focus:
    just _title "TEST" "Running Focus Tests"
    just _run "{{ginkgo}} --focus 'should create & get person' components/fun-app/it"
