[group('ops')]
[doc('Run the Fun App integration load tests')]
load:
    just _title "Load" "Test Fun App"
    just _run "make -C components/fun-app/it all"

[group('ops')]
[doc('Open heap and CPU pprof sessions for the running app')]
profile:
    #!/usr/bin/env bash
    set -euo pipefail
    endpoint="${ENDPOINT:-http://localhost:8080}"
    just _title "PROFILE" "ENDPOINT=${endpoint} | http://app.docker/app"
    just _detail "PROFILE" "Profiling Heap"
    just _run-no-output "go tool pprof -http=:8001 '${endpoint}/debug/pprof/heap' &"
    just _detail "PROFILE" "Profiling CPU"
    just _run-no-output "go tool pprof -http=:8000 --seconds=30 '${endpoint}/debug/pprof/profile'"
    just _warn "PROFILE" "Killing Profilers"
    kill %1

[group('ops')]
[doc('Analyse Fun App logs with GoAccess')]
analyse: _build-fun
    just _title "Analyse" "Fun App Logs"
    bin/fun 2>/dev/null | grep GIN | goaccess --log-format='%^ %d - %t | %s | %~%D | %b | %~%h | %^ | %m %U' --date-format='%Y/%m/%d' --time-format='%H:%M:%S'

# TODO: Replace with air tool enstead of entr ?
[group('run')]
[doc('Watch files and rerun a command with entr')]
watch:
    #!/usr/bin/env bash
    set -euo pipefail
    cmd="${CMD:-ls}"
    just _title "Watch (entr)" "${cmd}"
    find . | entr -s "date +%M:%S; ${cmd}"
