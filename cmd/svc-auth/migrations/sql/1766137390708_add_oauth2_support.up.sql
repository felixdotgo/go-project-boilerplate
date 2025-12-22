-- Add user profile fields to users table
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS name VARCHAR(255),
    ADD COLUMN IF NOT EXISTS first_name VARCHAR(100),
    ADD COLUMN IF NOT EXISTS last_name VARCHAR(100),
    ADD COLUMN IF NOT EXISTS avatar_url TEXT,
    ADD COLUMN IF NOT EXISTS email_verified BOOLEAN DEFAULT FALSE NOT NULL,
    ADD COLUMN IF NOT EXISTS locale VARCHAR(10);

-- Make password nullable for OAuth-only users
ALTER TABLE users ALTER COLUMN password DROP NOT NULL;

-- Add index for email verification status
CREATE INDEX IF NOT EXISTS idx_users_email_verified ON users(email_verified);

-- Create user_oauth_providers table for multi-provider support
CREATE TABLE IF NOT EXISTS user_oauth_providers (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    provider VARCHAR(50) NOT NULL,
    provider_user_id VARCHAR(255) NOT NULL,

    -- OAuth tokens (encrypted)
    access_token TEXT,
    refresh_token TEXT,
    token_expiry TIMESTAMP,

    -- Provider-specific profile data
    email VARCHAR(255),
    name VARCHAR(255),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    avatar_url TEXT,
    email_verified BOOLEAN DEFAULT FALSE NOT NULL,
    locale VARCHAR(10),

    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL,

    -- Constraints
    CONSTRAINT fk_user_oauth_providers_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    -- Each user can link each provider only once
    CONSTRAINT uq_user_provider UNIQUE (user_id, provider)
);

-- Create indexes for efficient lookups
CREATE UNIQUE INDEX IF NOT EXISTS idx_provider_user_id
    ON user_oauth_providers(provider, provider_user_id);
CREATE INDEX IF NOT EXISTS idx_user_oauth_providers_user_id
    ON user_oauth_providers(user_id);

-- Create refresh_tokens table
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    token VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    revoked_at TIMESTAMP,

    CONSTRAINT fk_refresh_tokens_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- Create indexes for refresh tokens
CREATE UNIQUE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens(token);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_revoked_at ON refresh_tokens(revoked_at) WHERE revoked_at IS NOT NULL;

-- Add comments
COMMENT ON TABLE user_oauth_providers IS 'Stores OAuth provider links - users can have multiple providers';
COMMENT ON COLUMN user_oauth_providers.provider IS 'OAuth provider name: google, facebook, github, etc.';
COMMENT ON COLUMN user_oauth_providers.provider_user_id IS 'Unique user ID from the OAuth provider';
COMMENT ON COLUMN user_oauth_providers.access_token IS 'Encrypted OAuth provider access token';
COMMENT ON COLUMN user_oauth_providers.refresh_token IS 'Encrypted OAuth provider refresh token';

COMMENT ON TABLE refresh_tokens IS 'Stores hashed refresh tokens for OAuth2 token refresh flow';
COMMENT ON COLUMN refresh_tokens.token IS 'SHA-256 hash of the refresh token';
COMMENT ON COLUMN refresh_tokens.revoked_at IS 'Timestamp when token was revoked (null if still valid)';
