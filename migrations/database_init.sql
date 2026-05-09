CREATE DATABASE shortlink_db;
CREATE USER shortlinkuser WITH PASSWORD 'shortlinkpassword';
GRANT ALL PRIVILEGES ON DATABASE shortlink_db to shortlinkuser;
CREATE SCHEMA IF NOT EXISTS app AUTHORIZATION shortlinkuser;
GRANT ALL ON SCHEMA app TO shortlinkuser;
ALTER DEFAULT PRIVILEGES IN SCHEMA app GRANT ALL ON TABLES TO shortlinkuser;
ALTER ROLE shortlinkuser SET search_path TO app, public;