#!/bin/bash

# Load environment variables if .env exists
if [ -f .env ]; then
  export $(cat .env | grep -v '^#' | xargs)
fi

# Set defaults
DB_USER=${POSTGRES_USER:-cyclingstream}
DB_NAME=${POSTGRES_DB:-cyclingstream}
CONTAINER_NAME="cyclingstream_postgres"

echo "Validating database schema..."

# Check if container is running
if ! docker ps | grep -q "$CONTAINER_NAME"; then
    echo "Error: Container $CONTAINER_NAME is not running."
    exit 1
fi

# Function to check column existence
check_column() {
    local table=$1
    local column=$2
    
    count=$(docker exec "$CONTAINER_NAME" psql -U "$DB_USER" -d "$DB_NAME" -t -c \
        "SELECT count(*) FROM information_schema.columns WHERE table_name='$table' AND column_name='$column';")
    
    if [ "$(echo $count | tr -d ' ')" == "1" ]; then
        echo "✅ $table.$column exists"
        return 0
    else
        echo "❌ $table.$column MISSING"
        return 1
    fi
}

# Function to check table existence
check_table() {
    local table=$1
    
    count=$(docker exec "$CONTAINER_NAME" psql -U "$DB_USER" -d "$DB_NAME" -t -c \
        "SELECT count(*) FROM information_schema.tables WHERE table_name='$table';")
    
    if [ "$(echo $count | tr -d ' ')" == "1" ]; then
        echo "✅ Table $table exists"
        return 0
    else
        echo "❌ Table $table MISSING"
        return 1
    fi
}

ERRORS=0

echo "Checking 'races' table schema..."
# Base fields
check_column "races" "id" || ERRORS=$((ERRORS+1))
check_column "races" "name" || ERRORS=$((ERRORS+1))
check_column "races" "is_free" || ERRORS=$((ERRORS+1))
check_column "races" "price_cents" || ERRORS=$((ERRORS+1))

# Newer fields (often missing in legacy DBs)
check_column "races" "requires_login" || ERRORS=$((ERRORS+1))
check_column "races" "stage_name" || ERRORS=$((ERRORS+1))
check_column "races" "stage_type" || ERRORS=$((ERRORS+1))
check_column "races" "elevation_meters" || ERRORS=$((ERRORS+1))
check_column "races" "estimated_finish_time" || ERRORS=$((ERRORS+1))
check_column "races" "stage_length_km" || ERRORS=$((ERRORS+1))

echo ""
echo "Checking for newer tables..."
check_table "achievements" || ERRORS=$((ERRORS+1))
check_table "user_preferences" || ERRORS=$((ERRORS+1))
check_table "chat_messages" || ERRORS=$((ERRORS+1))
# Check if chat_messages has new columns
if check_table "chat_messages"; then
    check_column "chat_messages" "user_role" || ERRORS=$((ERRORS+1))
    check_column "chat_messages" "badges" || ERRORS=$((ERRORS+1))
    check_column "chat_messages" "special_emote" || ERRORS=$((ERRORS+1))
fi

echo ""
if [ $ERRORS -eq 0 ]; then
    echo "✅ Schema validation passed! All checked columns and tables exist."
    exit 0
else
    echo "❌ Schema validation failed with $ERRORS errors."
    echo "Run 'make fix-schema' to attempt automatic repairs."
    exit 1
fi

