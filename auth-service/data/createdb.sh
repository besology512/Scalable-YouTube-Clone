#!/bin/bash

set -e

DB_FILE="data/auth_demo.db"
SCHEMA_FILE="internal/db/sqlc/schema.sql"

# 1. Remove existing DB if it exists
if [ -f "$DB_FILE" ]; then
    echo "Removing existing $DB_FILE..."
    rm "$DB_FILE"
fi

# 2. Create new DB and apply schema
echo "Creating new $DB_FILE using $SCHEMA_FILE..."
sqlite3 "$DB_FILE" < "$SCHEMA_FILE"

# 3. Show the tables
echo "Database created successfully. Tables:"
sqlite3 "$DB_FILE" ".tables"

# 4. Optional: describe the users table
echo
echo "Schema for 'users':"
sqlite3 "$DB_FILE" ".schema users"
