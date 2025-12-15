#!/bin/bash

echo "Starting SustainWear API..."

# CHECK IF CONFIG EXISTS
if [ ! -f "./configs/config.toml" ]; then
    echo "WARNING: Config file not found. It will be created on first run."
fi

# CHECK IF DATABASE EXISTS
if [ ! -f "./data/sustainwear.db" ]; then
    echo "WARNING: Database not found!"
    read -p "Do you want to initialize the database now? (Y/n): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Nn]$ ]]; then
        ./scripts/init_db.sh
    else
        echo "Cannot start without database"
        exit 1
    fi
fi

# BUILD AND RUN
echo "Building application..."
go build -o bin/sustainwear cmd/api/main.go

if [ $? -eq 0 ]; then
    echo "Build successful"
    echo "Starting server..."
    ./bin/sustainwear
else
    echo "Build failed"
    exit 1
fi