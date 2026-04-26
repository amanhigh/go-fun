[doc('Release the Fun App Docker images to Docker Hub')]
[group('release')]
release-docker: docker-build
    #!/usr/bin/env bash
    set -euo pipefail
    if [[ -z "${VER:-}" ]]; then
        echo "VER not set. Use VER=v1.0.5"
        exit 1
    fi
    just _title "RELEASE" "Release Docker Images: ${VER}"
    just _detail "RELEASE" "Docker Tag"
    docker tag amanfdk/fun-app:latest amanfdk/fun-app:${VER}
    just _detail "RELEASE" "Docker Push"
    docker push amanfdk/fun-app:latest
    docker push amanfdk/fun-app:${VER}

[doc('Build the Fun App Docker image locally')]
[group('docker')]
docker-build: _docker-fun
    just _info "Docker Hub" "https://hub.docker.com/r/amanfdk/fun-app/tags"

[group('docker')]
_docker-fun: _build-fun
    just _title "Docker" "Building FunApp"
    cp bin/fun components/fun-app/fun
    just _run "docker buildx build -t amanfdk/fun-app -f components/fun-app/Dockerfile components/fun-app"
    rm components/fun-app/fun

[doc('Run the Fun App Docker image interactively')]
[group('docker')]
docker-run: _docker-fun
    just _title "Docker" "Running FunApp"
    docker run -it amanfdk/fun-app

[doc('Open a shell inside the Fun App Docker image')]
[group('docker')]
docker-exec:
    just _title "Docker" "Execing Into FunApp"
    docker run -it --entrypoint /bin/sh amanfdk/fun-app

[doc('Remove local Fun App Docker images')]
[group('clean')]
docker-clean:
    just _warn "Docker" "Deleting FunApp"
    -just _run "docker rmi -f $(docker images amanfdk/fun-app -q)"

[group('info')]
_info-docker:
    just _info "INFO" "FunApp DockerHub: https://hub.docker.com/r/amanfdk/fun-app/tags"
    curl -s "https://hub.docker.com/v2/repositories/amanfdk/fun-app/tags/?page_size=25&page=1&name&ordering" | jq -r '.results[]|.name' | head -3
    just _info "INFO" "Docker Images: amanfdk/fun-app"
    docker images | grep fun-app
