FROM mcr.microsoft.com/devcontainers/go:1-1.23-bookworm

# Add user configuration
ARG USER=vscode

# Go environment variables
ENV GOPATH=/home/vscode/go
ENV GOBIN=/home/vscode/go/bin
ENV PATH=${PATH}:${GOBIN}

# Install additional Go tools
RUN go install github.com/remyoudompheng/go-misc/deadcode@latest && \
    go install golang.org/x/tools/cmd/goimports@latest && \
    go install github.com/vektra/mockery/v2@v2.51.1

# Install X11 dependencies for clipboard support
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    libx11-dev \
    libxcb1-dev \
    libxau-dev \
    libxdmcp-dev && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*