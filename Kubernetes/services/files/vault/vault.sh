#export VAULT_ADDR=http://localhost:8200

echo -en "\033[1;32m \n Secrets \033[0m \n"
#http://localhost:8200/ui/vault/secrets/secret/list
# vault login root-token
# vault secrets disable secret
# vault secrets enable -version=1 -path=secret kv
vault kv put secret/aman aman="preet"
vault kv get secret/aman

echo -en "\033[1;32m \n DB Management \033[0m \n"
#Increase lease Time
vault secrets enable database
vault write sys/mounts/database/tune max_lease_ttl="87600h"

# Vault Mysql Management
# Use plugin_name=mysql-legacy-database-plugin for mysql < 3.7
echo -en "\nMysql\n"
vault write database/config/aman-mysql plugin_name=mysql-database-plugin connection_url="{{username}}:{{password}}@tcp(mysql-primary:3306)/" allowed_roles="aman-mysql-role" username="root" password="root"
vault write database/roles/aman-mysql-role db_name=aman-mysql creation_statements="CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';GRANT SELECT ON *.* TO '{{name}}'@'%';"     default_ttl="1h"     max_ttl="24h"
vault read database/creds/aman-mysql-role

# Revocation
#vault lease revoke database/creds/aman-mysql-role/OstydME0HqTS7QmmSB5MQVqN

# Vault Mongo Management
# echo -en "\nMongo\n"
# vault write database/config/aman-mongo plugin_name=mongodb-database-plugin allowed_roles="aman-mongo-role" connection_url="mongodb://{{username}}:{{password}}@mongo-mongodb:27017/admin" username="root" password="root"
# vault write database/roles/aman-mongo-role db_name=aman-mongo creation_statements='{ "db": "admin", "roles": [{ "role": "readWrite" }, {"role": "read", "db": "compute"}] }' default_ttl="1h" max_ttl="24h"
# vault read database/creds/aman-mongo-role

# Enable Transit
echo -en "\033[1;32m \n Transit Engine \033[0m \n"
vault secrets enable transit

vault write -f transit/keys/my-key #Generate New Key
vault write transit/encrypt/my-key plaintext=$(echo "my secret data" | base64) #Encrypt
# vault write transit/decrypt/my-key ciphertext=<CIPHER> #Decrypt to Base64
# echo <BASE_64_TEXT> | base64 -d #Get Plain Text

vault write -f transit/keys/my-key/rotate #Rotate Key
# vault write transit/rewrap/my-key ciphertext=<CIPHER> #Re-encrypt data using rotate key