CREATE TYPE "category" AS ENUM (
  'Clothing',
  'Accessories',
  'Footwear',
  'Beverages'
);

CREATE TABLE "product" (
  "id" uuid PRIMARY KEY,
  "name" varchar,
  "sku" varchar,
  "category" category,
  "stock" integer,
  "price" integer,
  "imageUrl" varchar,
  "notes" varchar,
  "isAvailable" boolean,
  "location" varchar,
  "createdAt" timestamp
);
