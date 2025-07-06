CREATE TABLE "Users" (
  "user_id" bigserial PRIMARY KEY,
  "username" varchar(50) UNIQUE NOT NULL,
  "password_hash" varchar(255) NOT NULL,
  "email" varchar(100) UNIQUE NOT NULL,
  "role" varchar NOT NULL DEFAULT 'customer',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "Categories" (
  "category_id" bigserial PRIMARY KEY,
  "name" varchar(100) UNIQUE NOT NULL
);

CREATE TABLE "Products" (
  "product_id" bigserial PRIMARY KEY,
  "name" varchar(100) NOT NULL,
  "description" text NOT NULL,
  "product_image" text NOT NULL,
  "price" double precision NOT NULL,
  "stock_quantity" bigint NOT NULL,
  "category_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "Orders" (
  "order_id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "total_amount" double precision NOT NULL,
  "status" varchar(50) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "CartItems" (
  "cart_item_id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "product_id" bigint NOT NULL,
  "quantity" bigint NOT NULL,
  "price" double precision NOT NULL
);

CREATE TABLE "Reviews" (
  "review_id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "product_id" bigint NOT NULL,
  "rating" bigint NOT NULL,
  "comment" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "Payments" (
  "payment_id" bigserial PRIMARY KEY,
  "order_id" bigint NOT NULL,
  "payment_method" varchar(50) NOT NULL,
  "amount" double precision NOT NULL,
  "status" varchar(50) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX "idx_users_email" ON "Users" ("email");

CREATE INDEX "idx_products_category_id" ON "Products" ("category_id");

CREATE INDEX "idx_orders_user_id" ON "Orders" ("user_id");

COMMENT ON COLUMN "Orders"."status" IS 'CHECK (status IN (''pending'', ''completed''))';

COMMENT ON COLUMN "Reviews"."rating" IS 'CHECK (rating >= 1 AND rating <= 5)';

COMMENT ON COLUMN "Payments"."status" IS 'CHECK (status IN (''pending'', ''completed''))';

ALTER TABLE "Products" ADD FOREIGN KEY ("category_id") REFERENCES "Categories" ("category_id");

ALTER TABLE "Orders" ADD FOREIGN KEY ("user_id") REFERENCES "Users" ("user_id");

ALTER TABLE "CartItems" ADD FOREIGN KEY ("user_id") REFERENCES "Users" ("user_id");

ALTER TABLE "CartItems" ADD FOREIGN KEY ("product_id") REFERENCES "Products" ("product_id");

ALTER TABLE "Reviews" ADD FOREIGN KEY ("user_id") REFERENCES "Users" ("user_id");

ALTER TABLE "Reviews" ADD FOREIGN KEY ("product_id") REFERENCES "Products" ("product_id");

ALTER TABLE "Payments" ADD FOREIGN KEY ("order_id") REFERENCES "Orders" ("order_id");

-- Insert fixed categories if they don't already exist
INSERT INTO "Categories" (name) VALUES
  ('Clothing'),
  ('Electronics'),
  ('Beauty'),
  ('Media'),
  ('Accesories')
ON CONFLICT (name) DO NOTHING;

