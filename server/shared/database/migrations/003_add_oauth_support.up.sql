-- Add OAuth tables

-- OAuth State table for CSRF protection
CREATE TABLE oauth_states (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid(),
    state TEXT NOT NULL UNIQUE,
    provider TEXT NOT NULL,
    redirect_url TEXT,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Accounts table for OAuth account linking (extends the existing Account type)
CREATE TABLE accounts (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id TEXT NOT NULL,
    type TEXT NOT NULL,                   -- 'oauth'
    provider TEXT NOT NULL,               -- 'google', 'github', 'facebook'
    provider_account_id TEXT NOT NULL,    -- OAuth provider's user ID
    refresh_token TEXT,                   -- OAuth refresh token
    access_token TEXT,                    -- OAuth access token  
    expires_at INTEGER,                   -- Token expiration
    token_type TEXT,                      -- 'Bearer'
    scope TEXT,                           -- OAuth scopes
    id_token TEXT,                        -- OpenID Connect ID token
    session_state TEXT,                   -- Provider session state
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_accounts_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(provider, provider_account_id)
);

-- Indexes for performance
CREATE INDEX idx_oauth_states_expires_at ON oauth_states(expires_at);
CREATE INDEX idx_oauth_states_provider ON oauth_states(provider);
CREATE INDEX idx_accounts_user_id ON accounts(user_id);
CREATE INDEX idx_accounts_provider ON accounts(provider);
CREATE INDEX idx_accounts_provider_account_id ON accounts(provider_account_id);
