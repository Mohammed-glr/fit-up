#!/bin/bash

# Build script for Azure VM deployment
echo "Building Lornian Backend..."

# Install dependencies
go mod download
go mod verify

# Build the application
go build -o bin/lornian-backend main.go

echo "Build completed successfully!"
