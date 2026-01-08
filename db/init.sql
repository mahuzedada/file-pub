-- File Pub Database Schema
-- This script initializes the database schema for the File Pub application
-- Run this against your RDS MySQL instance

-- Create database if it doesn't exist
CREATE DATABASE IF NOT EXISTS filepub
    DEFAULT CHARACTER SET utf8mb4
    DEFAULT COLLATE utf8mb4_unicode_ci;

USE filepub;

-- Create images table
CREATE TABLE IF NOT EXISTS images (
    id VARCHAR(36) PRIMARY KEY,
    filename VARCHAR(255) NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    s3_key VARCHAR(512) NOT NULL,
    s3_url VARCHAR(1024) NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    size BIGINT NOT NULL,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_uploaded_at (uploaded_at DESC),
    INDEX idx_filename (filename),
    INDEX idx_content_type (content_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Display table structure
DESCRIBE images;

-- Display row count
SELECT COUNT(*) as total_images FROM images;
