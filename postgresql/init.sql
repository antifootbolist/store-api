CREATE DATABASE store-api;

\c store-api;

CREATE TABLE products (
  id INT,
  name VARCHAR(255),
  description VARCHAR(1024),
  price INT
);

CREATE USER "user-api" WITH PASSWORD 'qwe123';
GRANT ALL PRIVILEGES ON TABLE products TO "user-api";