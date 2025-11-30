#!/bin/bash

# Load environment variables if .env exists
if [ -f .env ]; then
  export $(cat .env | grep -v '^#' | xargs)
fi

# Set defaults
DB_USER=${POSTGRES_USER:-cyclingstream}
DB_NAME=${POSTGRES_DB:-cyclingstream}
CONTAINER_NAME="cyclingstream_postgres"

echo "Inspecting database schema..."

# Check if container is running
if ! docker ps | grep -q "$CONTAINER_NAME"; then
    echo "Error: Container $CONTAINER_NAME is not running."
    exit 1
fi

echo "=== Tables and Row Counts ==="
docker exec "$CONTAINER_NAME" psql -U "$DB_USER" -d "$DB_NAME" -c "
SELECT 
    schemaname as schema, 
    relname as table, 
    n_live_tup as row_count
FROM pg_stat_user_tables 
ORDER BY n_live_tup DESC;
"

echo ""
echo "=== Races Table Schema ==="
docker exec "$CONTAINER_NAME" psql -U "$DB_USER" -d "$DB_NAME" -c "
SELECT 
    column_name, 
    data_type, 
    is_nullable,
    column_default
FROM information_schema.columns 
WHERE table_name = 'races'
ORDER BY ordinal_position;
"

echo ""
echo "=== Chat Messages Table Schema ==="
docker exec "$CONTAINER_NAME" psql -U "$DB_USER" -d "$DB_NAME" -c "
SELECT 
    column_name, 
    data_type, 
    is_nullable
FROM information_schema.columns 
WHERE table_name = 'chat_messages'
ORDER BY ordinal_position;
"

echo ""
echo "=== Migration Status (schema_migrations table) ==="
# Check if schema_migrations exists (golang-migrate usually creates this)
MIGRATION_TABLE_EXISTS=$(docker exec "$CONTAINER_NAME" psql -U "$DB_USER" -d "$DB_NAME" -t -c "SELECT to_regclass('schema_migrations');")

if [ "$MIGRATION_TABLE_EXISTS" != "" ]; then
    docker exec "$CONTAINER_NAME" psql -U "$DB_USER" -d "$DB_NAME" -c "SELECT * FROM schema_migrations ORDER BY version DESC LIMIT 5;"
else
    echo "No schema_migrations table found."
fi

