-- Add line_id to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS line_id VARCHAR(255) UNIQUE;

-- Create index on line_id
CREATE INDEX IF NOT EXISTS idx_users_line_id ON users(line_id);
