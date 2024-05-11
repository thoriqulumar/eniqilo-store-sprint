CREATE TABLE "transaction" (
  "transactionId" uuid PRIMARY KEY,
  "customerId" uuid,
  "productDetails" JSONB,
  "paid" int,
  "change" int,
  "createdAt" timestamp
);