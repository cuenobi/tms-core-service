-- Re-add NOT NULL constraint to phone_number
-- WARNING: This will fail if there are any NULL values in the column
ALTER TABLE users ALTER COLUMN phone_number SET NOT NULL;
