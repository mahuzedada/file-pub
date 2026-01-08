#!/bin/bash

# Production Environment Setup Script
# This script validates production configuration

set -e

echo "==================================="
echo "File Pub - Production Setup"
echo "==================================="
echo ""

# Check if .env.prod exists
if [ ! -f .env.prod ]; then
    echo "Error: .env.prod file not found"
    echo "Please copy .env.prod template and configure it:"
    echo "  cp .env.prod.example .env.prod"
    exit 1
fi

# Source the prod environment
source .env.prod

echo "Step 1: Validating configuration..."
if [ -z "$S3_BUCKET" ]; then
    echo "Error: S3_BUCKET not set in .env.prod"
    exit 1
fi

if [ -z "$DB_HOST" ]; then
    echo "Error: DB_HOST not set in .env.prod"
    exit 1
fi

echo ""
echo "==================================="
echo "Production Configuration"
echo "==================================="
echo ""
echo "S3 Bucket: ${S3_BUCKET} (pre-existing)"
echo "Database: ${DB_HOST}:${DB_PORT}"
echo "Region: ${S3_REGION}"
echo ""
echo "Prerequisites (should already exist):"
echo "  ✓ S3 bucket: ${S3_BUCKET}"
echo "  ✓ RDS MySQL instance: ${DB_HOST}"
echo "  ✓ VPC with proper security groups"
echo "  ✓ EC2 instance with IAM role for S3 access"
echo ""
echo "Next steps:"
echo "1. Initialize database (if not done):"
echo "   mysql -h ${DB_HOST} -u ${DB_USER} -p < db/init.sql"
echo ""
echo "2. Build and deploy:"
echo "   make prod-build"
echo "   make prod-deploy SSH_HOST=ec2-user@your-ec2-ip"
echo ""
