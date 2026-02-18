-- Revert: make email NOT NULL again and restore original unique index
DROP INDEX IF EXISTS idx_users_email;
CREATE UNIQUE INDEX idx_users_email ON users (email);
UPDATE users SET email = '' WHERE email IS NULL;
ALTER TABLE users ALTER COLUMN email SET NOT NULL;
