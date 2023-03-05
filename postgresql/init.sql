CREATE DATABASE store_api;

\c store_api;

CREATE TABLE products (
  id INT,
  name VARCHAR(255),
  description VARCHAR(1024),
  price INT
);

insert into products values(1, 'iPhone', 'iPhone 14', 100);
insert into products values(2, 'iPhone', 'iPhone 14 PRO MAX', 200);

CREATE USER "user-api" WITH PASSWORD 'qwe123';
GRANT ALL PRIVILEGES ON TABLE products TO "user-api";
