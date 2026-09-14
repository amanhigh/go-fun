# The 1-1 track follows the latest stable Go release; the repository baseline
# is pinned to Go 1.27.1 in go.work and the module files.
FROM mcr.microsoft.com/devcontainers/go:1-1-bookworm

# Add user configuration
ARG USER=vscode

# Go environment variables
ENV GOPATH=/home/vscode/go
ENV GOBIN=/home/vscode/go/bin
ENV PATH=${PATH}:${GOBIN}

# Install additional Go tools
RUN go install github.com/remyoudompheng/go-misc/deadcode@latest && \
    go install golang.org/x/tools/cmd/goimports@latest && \
    go install github.com/vektra/mockery/v3@v3.7.4

# Install display dependencies for clipboard support in desktop and headless environments
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    libx11-dev \
    libxcb1-dev \
    libxau-dev \
    libxdmcp-dev \
    xvfb \
    xauth && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
