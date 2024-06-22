#!/bin/bash

# Define the pod name and namespace if not default
PG_POD="pg-primary-0"  # Replace with your actual PostgreSQL pod name
NAMESPACE="default"  # Replace with the namespace of the PostgreSQL pod if different

# Define the backup file location on your host machine
BACKUP_DIR=$(dirname "$(realpath "$0")")
BACKUP_FILE="$BACKUP_DIR/pg_backup.sql"  # Adjust this if your backup file has a different name or path

# PostgreSQL credentials
PG_USER="postgres"  # Replace with your PostgreSQL user
PG_PASSWORD="root"  # Replace with your actual PostgreSQL password

# Check if the backup file exists
if [ ! -f "$BACKUP_FILE" ]; then
    echo "Backup file does not exist: $BACKUP_FILE"
    exit 1
fi

# Restore the backup into the PostgreSQL database
echo "Starting restoration of PostgreSQL databases into pod $PG_POD..."
cat "$BACKUP_FILE" | kubectl exec -i $PG_POD -n $NAMESPACE -- env PGPASSWORD=$PG_PASSWORD psql -U $PG_USER -d postgres -q --set ON_ERROR_STOP=0

# Check if the restoration was successful
if [ $? -eq 0 ]; then
    echo "Restoration completed successfully."
else
    echo "Restoration failed."
fi
