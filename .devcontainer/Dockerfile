FROM mcr.microsoft.com/devcontainers/go:1-1.23-bookworm

# Add user configuration
ARG USER=vscode

# Go environment variables
ENV GOPATH=/home/vscode/go
ENV GOBIN=/home/vscode/go/bin
ENV PATH=${PATH}:${GOBIN}

# Install additional Go tools
RUN go install github.com/remyoudompheng/go-misc/deadcode@latest
RUN go install golang.org/x/tools/cmd/goimports@latest