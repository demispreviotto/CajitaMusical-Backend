-- File: migrations/reset_schema_dev.sql
-- Purpose: To completely reset the database schema in a DEVELOPMENT environment.
-- WARNING: This script will DELETE ALL DATA from these tables.

-- Drop tables in reverse dependency order to avoid foreign key issues.
-- Tables with foreign keys should be dropped before the tables they reference.
-- CASCADE is essential here to also drop any dependent objects (like FK constraints).
DROP TABLE IF EXISTS sessions CASCADE;        -- Sessions references Users
DROP TABLE IF EXISTS authentications CASCADE; -- Authentications references Users
DROP TABLE IF EXISTS authentication CASCADE; -- Authentications references Users
DROP TABLE IF EXISTS songs CASCADE;           -- Songs references Users (assuming this)
DROP TABLE IF EXISTS users CASCADE;           -- Users is referenced by others, so drop last

-- Enable the pgcrypto extension for gen_random_uuid() if not already enabled.
-- This is a one-time setup for the database.
CREATE EXTENSION IF NOT EXISTS pgcrypto;