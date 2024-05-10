-- Alter the table to modify the phoneNumber column to be of type string and unique
ALTER TABLE staff
ALTER COLUMN "phoneNumber" TYPE VARCHAR(255),
ADD UNIQUE ("phoneNumber");

-- Add an index for the phoneNumber column
CREATE INDEX idx_phoneNumber ON staff ("phoneNumber");