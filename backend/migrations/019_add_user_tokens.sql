-- Add tokens field to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS tokens INT DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_daily_login TIMESTAMP;

-- Give existing users starting tokens
UPDATE users SET tokens = 1000 WHERE tokens = 0;
