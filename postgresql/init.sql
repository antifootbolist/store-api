CREATE DATABASE store_api;
CREATE USER "user-api" WITH PASSWORD 'qwe123';
GRANT ALL PRIVILEGES ON DATABASE store_api TO "user-api";
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "user-api";
