-- file: 10-create-user-and-db.sql
-- CREATE DATABASE persons;
-- CREATE ROLE program WITH PASSWORD 'test';
-- GRANT ALL PRIVILEGES ON DATABASE persons TO program;
-- ALTER ROLE program WITH LOGIN;

CREATE SCHEMA IF NOT EXISTS persons;

CREATE TABLE person
(
    id      BIGSERIAL PRIMARY KEY,
    "name"  TEXT UNIQUE NOT NULL,
    age     INT NOT NULL,
    address TEXT NOT NULL,
    "work"  TEXT NOT NULL
);