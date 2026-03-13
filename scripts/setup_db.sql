-- Script untuk setup database PostgreSQL
-- Jalankan sebagai postgres superuser

-- Create database
CREATE DATABASE mikmongo;

-- Create user (optional - bisa pakai postgres default)
-- CREATE USER mikmongo WITH PASSWORD 'mikmongo';
-- GRANT ALL PRIVILEGES ON DATABASE mikmongo TO mikmongo;

-- Grant privileges to postgres user (if using default)
GRANT ALL PRIVILEGES ON DATABASE mikmongo TO postgres;
