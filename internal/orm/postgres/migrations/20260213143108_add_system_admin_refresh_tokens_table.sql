-- +goose Up
-- +goose StatementBegin

-- =====================
-- system_admin_refresh_tokens table (quản lý refresh token cho admin)
-- =====================
CREATE TABLE IF NOT EXISTS system_admin_refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    admin_id UUID NOT NULL,
    
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

    -- Foreign key
    CONSTRAINT fk_system_admin_refresh_tokens_admin
        FOREIGN KEY (admin_id)
        REFERENCES system_admins(id)
        ON DELETE CASCADE
);

-- =====================
-- Indexes
-- =====================
CREATE UNIQUE INDEX IF NOT EXISTS idx_system_admin_refresh_tokens_token_hash_unique
ON system_admin_refresh_tokens(refresh_token_hash)
WHERE is_revoked = false AND deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_system_admin_refresh_tokens_admin_id
ON system_admin_refresh_tokens(admin_id);

CREATE INDEX IF NOT EXISTS idx_system_admin_refresh_tokens_is_revoked
ON system_admin_refresh_tokens(is_revoked);

CREATE INDEX IF NOT EXISTS idx_system_admin_refresh_tokens_expires_at
ON system_admin_refresh_tokens(expires_at);

CREATE INDEX IF NOT EXISTS idx_system_admin_refresh_tokens_last_active_at
ON system_admin_refresh_tokens(last_active_at);

CREATE INDEX IF NOT EXISTS idx_system_admin_refresh_tokens_deleted_at
ON system_admin_refresh_tokens(deleted_at)
WHERE deleted_at IS NOT NULL;

-- =====================
-- Trigger: auto update updated_at
-- =====================
DROP TRIGGER IF EXISTS update_system_admin_refresh_tokens_updated_at ON system_admin_refresh_tokens;
CREATE TRIGGER update_system_admin_refresh_tokens_updated_at
BEFORE UPDATE ON system_admin_refresh_tokens
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS update_system_admin_refresh_tokens_updated_at ON system_admin_refresh_tokens;

DROP INDEX IF EXISTS idx_system_admin_refresh_tokens_deleted_at;
DROP INDEX IF EXISTS idx_system_admin_refresh_tokens_last_active_at;
DROP INDEX IF EXISTS idx_system_admin_refresh_tokens_expires_at;
DROP INDEX IF EXISTS idx_system_admin_refresh_tokens_is_revoked;
DROP INDEX IF EXISTS idx_system_admin_refresh_tokens_admin_id;
DROP INDEX IF EXISTS idx_system_admin_refresh_tokens_token_hash_unique;

DROP TABLE IF EXISTS system_admin_refresh_tokens;

-- +goose StatementEnd
