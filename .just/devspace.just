[doc('Start the devspace development environment')]
[group('devspace')]
space: space-purge
    just _title "Starting" "Devspace"
    devspace use namespace fun-app
    devspace dev

[doc('Purge the current devspace environment')]
[group('devspace')]
space-purge:
    just _title "Purging" "Devspace"
    -devspace purge

[doc('Run the devspace Ginkgo test workflow')]
[group('devspace')]
space-test:
    just _title "TEST" "Devspace Tests"
    devspace run ginkgo
    CMD="devspace run fun-test" just watch

[group('devspace')]
[doc('Display information about the devspace environment')]
space-info:
    just _title "INFO" "Devspace"
    devspace list vars --var DB="mysql-primary",RATE_LIMIT=-1
    just _detail "INFO" "http://localhost:8080/metrics"
    just _detail "INFO" "Login: devspace enter"
