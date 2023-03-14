CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    product TEXT NOT NULL,
    price INTEGER NOT NULL
);