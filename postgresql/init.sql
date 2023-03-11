CREATE DATABASE store_api;

\c store_api;

CREATE TABLE products (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255),
  description VARCHAR(1024),
  price INT
);

CREATE USER "user-api" WITH PASSWORD 'qwe123';
GRANT ALL PRIVILEGES ON TABLE products TO "user-api";
ALTER TABLE products OWNER TO "user-api";
