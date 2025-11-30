#!/bin/bash

# Load environment variables if .env exists
if [ -f .env ]; then
  export $(cat .env | grep -v '^#' | xargs)
fi

# Set defaults
DB_USER=${POSTGRES_USER:-cyclingstream}
DB_NAME=${POSTGRES_DB:-cyclingstream}
CONTAINER_NAME="cyclingstream_postgres"
BACKUP_DIR="./backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/db_backup_$TIMESTAMP.sql"

# Create backup directory if it doesn't exist
mkdir -p "$BACKUP_DIR"

echo "Creating database backup..."
echo "Container: $CONTAINER_NAME"
echo "Database: $DB_NAME"
echo "Output file: $BACKUP_FILE"

# check if container is running
if ! docker ps | grep -q "$CONTAINER_NAME"; then
    echo "Error: Container $CONTAINER_NAME is not running."
    exit 1
fi

# Create backup
docker exec -t "$CONTAINER_NAME" pg_dump -U "$DB_USER" -d "$DB_NAME" -F p > "$BACKUP_FILE"

if [ $? -eq 0 ]; then
    echo "Backup completed successfully: $BACKUP_FILE"
    # Keep only last 10 backups
    ls -t "$BACKUP_DIR"/db_backup_*.sql | tail -n +11 | xargs -r rm
else
    echo "Backup failed!"
    rm -f "$BACKUP_FILE"
    exit 1
fi

