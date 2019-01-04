#!/usr/bin/env bash

# export VAULT_ADDR=http://docker:8200

# Enable vault version 1 API's
vault login root-token && vault kv get secret &&
vault secrets disable secret &&
vault secrets enable -version=1 -path=secret kv

# Enable Transit
vault secrets enable transit

#Encrypt Decrypt Data
#vault write -f transit/keys/my-key #Generate New Key

#vault write transit/encrypt/my-key plaintext=$(base64 <<< "my secret data") #Encrypt
#vault write transit/decrypt/my-key ciphertext=<CIPHER> #Decrypt to Base64
#base64 --decode <<< <BASE_64_TEXT #Get Plain Text

#vault write -f transit/keys/my-key/rotate #Rotate Key
#vault write transit/rewrap/my-key ciphertext=<CIPHER> #Re-encrypt data using rotate key