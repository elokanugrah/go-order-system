-- migration/000001_init_schema.up.sql
CREATE TABLE "products" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "price" decimal(10, 2) NOT NULL,
  "quantity" integer NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "orders" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "total_amount" decimal(10, 2) NOT NULL,
  "status" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "order_items" (
  "id" bigserial PRIMARY KEY,
  "order_id" bigint NOT NULL REFERENCES "orders" ("id"),
  "product_id" bigint NOT NULL REFERENCES "products" ("id"),
  "quantity" integer NOT NULL,
  "price_at_order" decimal(10, 2) NOT NULL
);