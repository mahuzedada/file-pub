#!/bin/bash

# Development Environment Setup Script
# This script starts Docker MySQL for local development

set -e

echo "==================================="
echo "File Pub - Development Setup"
echo "==================================="
echo ""

# Check if .env.dev exists
if [ ! -f .env.dev ]; then
    echo "Error: .env.dev file not found"
    echo "Please copy .env.dev template and configure it:"
    echo "  cp .env.dev.example .env.dev"
    exit 1
fi

# Source the dev environment
source .env.dev

echo "Step 1: Checking configuration..."
if [ -z "$S3_BUCKET" ]; then
    echo "Error: S3_BUCKET not set in .env.dev"
    echo "Please set S3_BUCKET to your existing S3 bucket name"
    exit 1
fi

echo "Configuration:"
echo "  - S3 Bucket: ${S3_BUCKET}"
echo "  - Database: localhost:3306"
echo ""

echo "Step 2: Starting Docker MySQL..."
docker-compose up -d mysql

echo ""
echo "Step 3: Waiting for MySQL to be ready..."
sleep 10

echo ""
echo "==================================="
echo "Development environment ready!"
echo "==================================="
echo ""
echo "Database: MySQL running on localhost:3306"
echo "  - User: devuser"
echo "  - Password: devpassword"
echo "  - Database: filepub"
echo ""
echo "S3 Bucket: ${S3_BUCKET} (pre-existing)"
echo ""
echo "To start the application:"
echo "  make dev-run"
echo ""
echo "To stop MySQL:"
echo "  make dev-down"
echo ""
