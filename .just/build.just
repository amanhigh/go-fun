[group('build')]
_build-fun:
    just components/fun-app/build

[group('build')]
_build-fun-cover:
    just components/fun-app/build-cover

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
    -just _run "make -C components/operator/ clean"

[group('clean')]
[doc('Clean build and test artifacts')]
clean: _clean-test _clean-build

[group('info')]
[doc('Show release tags and Docker image information')]
info: _info-release _info-docker

[group('info')]
[doc('Show release info plus devspace information')]
infos: info space-info


