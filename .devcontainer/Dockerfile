FROM mcr.microsoft.com/devcontainers/go:1-1.23-bookworm

# Install additional Go tools
RUN go install github.com/remyoudompheng/go-misc/deadcode@latest
RUN go install golang.org/x/tools/cmd/goimports@latest