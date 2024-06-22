#!/bin/bash

# Define the pod name and namespace
PG_POD="pg-primary-0"  
NAMESPACE="default"

# Define where to store the backups on the host machine
BACKUP_DIR=$(dirname "$(realpath "$0")")

# Define the filename for the backup
BACKUP_FILE="pg_backup.sql"

# PostgreSQL credentials
PG_USER="postgres" 
PG_PASSWORD="root"

# Execute pg_dumpall inside the PostgreSQL pod
echo "Starting backup of all PostgreSQL databases from pod $PG_POD..."
kubectl exec $PG_POD -n $NAMESPACE -- env PGPASSWORD=$PG_PASSWORD pg_dumpall -U $PG_USER > "$BACKUP_DIR/$BACKUP_FILE"

# Check if the backup was successful
if [ $? -eq 0 ]; then
    echo "Backup completed successfully."
    echo "Backup file is located at: $BACKUP_DIR/$BACKUP_FILE"
else
    echo "Backup failed."
fi

### Help
# CREATE DATABASE metabase;
# GRANT ALL PRIVILEGES ON DATABASE metabase TO aman WITH GRANT OPTION;
# GRANT ALL ON SCHEMA public TO aman; (Need to Connect to Right DB eg. metabase)