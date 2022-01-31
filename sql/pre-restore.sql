ALTER SYSTEM SET synchronous_commit=OFF;
CREATE USER insights WITH PASSWORD 'insights' CREATEDB;
CREATE DATABASE census WITH OWNER insights;
