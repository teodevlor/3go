-- +goose Up
-- +goose StatementBegin

-- =====================
-- Sessions table (phiên đăng nhập cụ thể: account + device + app)
-- =====================
CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    account_app_device_id UUID NOT NULL,
    
    refresh_token_hash VARCHAR(255) NOT NULL, -- Hash của refresh token (không lưu plain text)
    
    expires_at TIMESTAMPTZ NOT NULL,
    is_revoked BOOLEAN NOT NULL DEFAULT false,
    revoked_at TIMESTAMPTZ,
    revoked_reason VARCHAR(255),      -- logout, security, password_reset, ...
    
    last_active_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Thông tin bảo mật bổ sung
    ip_address VARCHAR(45),           -- IPv4 hoặc IPv6
    user_agent TEXT,
    
    metadata JSONB DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,

    -- Foreign keys
    CONSTRAINT fk_sessions_account_app_device
        FOREIGN KEY (account_app_device_id)
        REFERENCES account_app_devices(id)
        ON DELETE CASCADE
);

-- =====================
-- Indexes
-- =====================
CREATE UNIQUE INDEX IF NOT EXISTS idx_sessions_refresh_token_hash_unique
ON sessions(refresh_token_hash)
WHERE is_revoked = false AND deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_sessions_account_app_device_id
ON sessions(account_app_device_id);

CREATE INDEX IF NOT EXISTS idx_sessions_is_revoked
ON sessions(is_revoked);

CREATE INDEX IF NOT EXISTS idx_sessions_expires_at
ON sessions(expires_at);

CREATE INDEX IF NOT EXISTS idx_sessions_last_active_at
ON sessions(last_active_at);

CREATE INDEX IF NOT EXISTS idx_sessions_deleted_at
ON sessions(deleted_at)
WHERE deleted_at IS NOT NULL;

-- =====================
-- Trigger: auto update updated_at
-- =====================
DROP TRIGGER IF EXISTS update_sessions_updated_at ON sessions;
CREATE TRIGGER update_sessions_updated_at
BEFORE UPDATE ON sessions
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS update_sessions_updated_at ON sessions;

DROP INDEX IF EXISTS idx_sessions_deleted_at;
DROP INDEX IF EXISTS idx_sessions_last_active_at;
DROP INDEX IF EXISTS idx_sessions_expires_at;
DROP INDEX IF EXISTS idx_sessions_is_revoked;
DROP INDEX IF EXISTS idx_sessions_account_app_device_id;
DROP INDEX IF EXISTS idx_sessions_refresh_token_hash_unique;

DROP TABLE IF EXISTS sessions;

-- +goose StatementEnd
