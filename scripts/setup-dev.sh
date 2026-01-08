#!/bin/bash

# Development Environment Setup Script
# This script starts Docker MySQL for local development

set -e

echo "==================================="
echo "File Pub - Development Setup"
echo "==================================="
echo ""

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
ENV_FILE="$PROJECT_ROOT/.env.dev"

# Check if .env.dev exists
if [ ! -f "$ENV_FILE" ]; then
    echo "Error: .env.dev file not found at $ENV_FILE"
    echo "Please copy .env.dev template and configure it:"
    echo "  cp .env.dev.example .env.dev"
    exit 1
fi

# Source the dev environment (export all variables)
echo "Loading environment from: $ENV_FILE"
set -a
source "$ENV_FILE"
set +a

echo "Step 1: Checking configuration..."
echo "DEBUG: S3_BUCKET='$S3_BUCKET'"
echo "DEBUG: Shell: $SHELL"
echo "DEBUG: Bash version: $BASH_VERSION"

if [ -z "$S3_BUCKET" ]; then
    echo ""
    echo "Error: S3_BUCKET not set in .env.dev"
    echo "Please check that .env.dev contains: S3_BUCKET=your-bucket-name"
    echo ""
    echo "Contents of .env.dev:"
    grep -E "^S3_BUCKET=" "$ENV_FILE" || echo "  (S3_BUCKET not found in file)"
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
