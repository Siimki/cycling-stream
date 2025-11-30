#!/bin/bash

# Load environment variables if .env exists
if [ -f .env ]; then
  export $(cat .env | grep -v '^#' | xargs)
fi

# Set defaults
DB_USER=${POSTGRES_USER:-cyclingstream}
DB_NAME=${POSTGRES_DB:-cyclingstream}
CONTAINER_NAME="cyclingstream_postgres"

echo "Attempting to fix database schema..."

# Check if container is running
if ! docker ps | grep -q "$CONTAINER_NAME"; then
    echo "Error: Container $CONTAINER_NAME is not running."
    exit 1
fi

# Apply fixes from SQL file
# We use cat to pipe the file content to ensure it works even if the file is outside container context
cat scripts/fix-schema.sql | docker exec -i "$CONTAINER_NAME" psql -U "$DB_USER" -d "$DB_NAME"

if [ $? -eq 0 ]; then
    echo "✅ Schema fixes applied successfully."
else
    echo "❌ Failed to apply schema fixes."
    exit 1
fi
