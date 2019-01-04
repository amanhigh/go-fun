#!/usr/bin/env bash

# export VAULT_ADDR=http://docker:8200

# Enable vault version 1 API's
vault login root-token && vault kv get secret &&
vault secrets disable secret &&
vault secrets enable -version=1 -path=secret kv

# Enable Transit