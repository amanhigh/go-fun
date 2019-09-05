vault kv put secret/aman aman="preet"

#http://docker:8200/ui/vault/secrets/secret/list
vault kv get secret/aman

# Vault Mysql Management
vault secrets enable database
vault write database/config/aman-mysql plugin_name=mysql-database-plugin connection_url="root:root@tcp(compose_mysql_1:3306)/" allowed_roles="aman-mysql-role" username="root" password="root"
vault write database/roles/aman-mysql-role db_name=aman-mysql creation_statements="CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';GRANT SELECT ON *.* TO '{{name}}'@'%';"     default_ttl="1h"     max_ttl="24h"
vault read database/creds/aman-mysql-role

# Revocation
#vault lease revoke database/creds/aman-mysql-role/OstydME0HqTS7QmmSB5MQVqN

# Vault Mongo Management
#vault write database/config/aman-mongo plugin_name=mongodb-database-plugin allowed_roles="aman-mongo-role" connection_url="mongodb://{{username}}:{{password}}@compose_mongo_1:27017/admin" username="root" password="root"
#vault write database/roles/aman-mongo-role db_name=aman-mongo creation_statements='{ "db": "admin", "roles": [{ "role": "readWrite" }, {"role": "read", "db": "compute"}] }' default_ttl="1h" max_ttl="24h"
#vault read database/creds/aman-mongo-role