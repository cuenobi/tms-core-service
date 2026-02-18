-- Make email column nullable to support OAuth users (LINE) who may not have an email
-- Also update unique index to allow multiple NULL values (PostgreSQL treats NULLs as distinct)

-- Drop the existing NOT NULL constraint
ALTER TABLE users ALTER COLUMN email DROP NOT NULL;

-- Drop old unique index and recreate as partial index (only enforce uniqueness for non-null emails)
DROP INDEX IF EXISTS idx_users_email;
CREATE UNIQUE INDEX idx_users_email ON users (email) WHERE email IS NOT NULL AND email != '';

-- Convert any existing empty string emails to NULL
UPDATE users SET email = NULL WHERE email = '';
