# Run the Fun App binary.
# TODO: Create Just File for Building and Running Fun App
[group('run')]
run: _build-fun
    just _title "RUN" "Running Fun App"
    out="${OUT:-/dev/stdout}"
    bin/fun > "$out"

[group('build')]
_build-fun:
    just _title "BUILD" "Building Fun App"
    mkdir -p bin
    CGO_ENABLED=0 GOARCH=amd64 go build -o bin/fun components/fun-app/main.go

[group('build')]
_build-fun-cover:
    just _title "BUILD" "Building Fun App with Coverage"
    mkdir -p bin
    CGO_ENABLED=0 GOARCH=amd64 go build -cover -o bin/fun components/fun-app/main.go

[group('build')]
_build-kohan:
    just _title "BUILD" "Building Kohan"
    mkdir -p bin
    CGO_ENABLED=1 GOARCH=amd64 go build -o bin/kohan components/kohan/main.go

[group('core')]
[doc('Build the Fun App and Kohan binaries')]
build: format lint _build-fun _build-kohan

[group('clean')]
_clean-build:
    just _warn "CLEAN" "Cleaning Build"
    -rm -rf bin
    -make -C components/operator/ clean

[group('clean')]
[doc('Clean build and test artifacts')]
clean: _clean-test _clean-build

[group('info')]
[doc('Show release tags and Docker image information')]
info: _info-release _info-docker

[group('info')]
[doc('Show release info plus devspace information')]
infos: info space-info




