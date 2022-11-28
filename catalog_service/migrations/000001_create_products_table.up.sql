CREATE TABLE IF NOT EXISTS products (
    "id" text NOT NULL PRIMARY KEY,
    "name" text NOT NULL,
    "description" text NOT NULL,
    "category" text NOT NULL,
    "price" bigint NOT NULL,
    "image" text NOT NULL
);
