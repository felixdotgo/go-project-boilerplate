-- Drop refresh_tokens table
DROP TABLE IF EXISTS refresh_tokens;

-- Drop user_oauth_providers table
DROP TABLE IF EXISTS user_oauth_providers;

-- Drop indexes from users
DROP INDEX IF EXISTS idx_users_email_verified;

-- Remove profile columns from users
ALTER TABLE users
    DROP COLUMN IF EXISTS name,
    DROP COLUMN IF EXISTS first_name,
    DROP COLUMN IF EXISTS last_name,
    DROP COLUMN IF EXISTS avatar_url,
    DROP COLUMN IF EXISTS email_verified,
    DROP COLUMN IF EXISTS locale;

-- Restore password constraint (NOTE: This will fail if there are OAuth-only users without passwords)
-- ALTER TABLE users ALTER COLUMN password SET NOT NULL;
