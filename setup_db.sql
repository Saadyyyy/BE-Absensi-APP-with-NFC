-- Setup script untuk database attendance system
-- Jalankan script ini sebagai superuser PostgreSQL

-- Create database
CREATE DATABASE dimasrio;

-- Create user (optional, jika ingin user khusus)
-- CREATE USER attendance_user WITH PASSWORD 'your_password';
-- GRANT ALL PRIVILEGES ON DATABASE attendance_db TO attendance_user;

-- Connect to the database
\c dimasrio;

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create function for generating UUID (alternative to uuid-ossp)
CREATE OR REPLACE FUNCTION gen_random_uuid() RETURNS uuid AS $$
SELECT uuid_generate_v4();
$$ LANGUAGE SQL;

-- Database is ready for GORM auto-migration