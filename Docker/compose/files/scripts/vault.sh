# Vault Mysql Management
vault secrets enable database
vault write database/config/aman-mysql plugin_name=mysql-database-plugin connection_url="root:root@tcp(compose_mysql_1:3306)/" allowed_roles="aman-mysql-role" username="root" password="root"
vault write database/roles/aman-mysql-role db_name=aman-mysql creation_statements="CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';GRANT SELECT ON *.* TO '{{name}}'@'%';"     default_ttl="1h"     max_ttl="24h"
vault read database/creds/aman-mysql-role

#vault lease revoke database/creds/aman-mysql-role/OstydME0HqTS7QmmSB5MQVqN