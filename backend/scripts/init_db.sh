#!/bin/bash

# DATABASE INITIALIZATION SCRIPT
# Creates SQLite database and runs schema

DB_PATH="./data/sustainwear.db"
SCHEMA_FILE="./scripts/schema.sql"
SEED_FILE="./scripts/seed.sql"

echo "Initializing SustainWear Database..."

# CREATE DATA DIRECTORY IF IT DOESN'T EXIST
mkdir -p ./data

# CHECK IF DATABASE EXISTS
if [ -f "$DB_PATH" ]; then
    echo "WARNING: Database already exists at $DB_PATH"
    read -p "Do you want to delete and recreate it? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        rm "$DB_PATH"
        echo "Deleted existing database"
    else
        echo "Aborted"
        exit 0
    fi
fi

# CREATE DATABASE AND RUN SCHEMA
echo "Creating database schema..."
sqlite3 "$DB_PATH" < "$SCHEMA_FILE"

if [ $? -eq 0 ]; then
    echo "Schema created successfully"
else
    echo "Failed to create schema"
    exit 1
fi

# ASK IF USER WANTS TO SEED DATA
read -p "Do you want to insert seed data? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "Inserting seed data..."
    sqlite3 "$DB_PATH" < "$SEED_FILE"
    if [ $? -eq 0 ]; then
        echo "Seed data inserted successfully"
    else
        echo "Failed to insert seed data"
        exit 1
    fi
fi

echo "Database initialization complete!"
echo "Database location: $DB_PATH"