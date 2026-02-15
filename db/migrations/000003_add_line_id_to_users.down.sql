-- Drop index and column
DROP INDEX IF EXISTS idx_users_line_id;
ALTER TABLE users DROP COLUMN IF EXISTS line_id;
