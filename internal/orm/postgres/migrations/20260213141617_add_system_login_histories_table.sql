-- +goose Up
-- +goose StatementBegin

-- =====================
-- System Login Histories table (cho admin login)
-- =====================
CREATE TABLE IF NOT EXISTS system_login_histories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    admin_id UUID NOT NULL,
    
    login_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    result VARCHAR(50) NOT NULL,   -- success, failed_password, failed_inactive, failed_not_found, ...
    failure_reason TEXT,           -- Chi tiết lý do nếu failed
    
    -- Thông tin bảo mật
    ip_address VARCHAR(45),
    user_agent TEXT,
    location VARCHAR(255),         -- Optional: city, country from IP
    
    metadata JSONB DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Foreign key (không CASCADE để giữ lại lịch sử khi xóa admin)
    CONSTRAINT fk_system_login_histories_admin
        FOREIGN KEY (admin_id)
        REFERENCES system_admins(id)
        ON DELETE SET NULL
);

-- =====================
-- Indexes (quan trọng cho query lịch sử)
-- =====================
CREATE INDEX IF NOT EXISTS idx_system_login_histories_admin_id
ON system_login_histories(admin_id);

CREATE INDEX IF NOT EXISTS idx_system_login_histories_result
ON system_login_histories(result);

CREATE INDEX IF NOT EXISTS idx_system_login_histories_login_at
ON system_login_histories(login_at DESC); -- DESC vì thường query gần nhất trước

CREATE INDEX IF NOT EXISTS idx_system_login_histories_ip_address
ON system_login_histories(ip_address);

-- Composite index cho query phổ biến: lấy history của 1 admin trong khoảng thời gian
CREATE INDEX IF NOT EXISTS idx_system_login_histories_admin_login_at
ON system_login_histories(admin_id, login_at DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_system_login_histories_admin_login_at;
DROP INDEX IF EXISTS idx_system_login_histories_ip_address;
DROP INDEX IF EXISTS idx_system_login_histories_login_at;
DROP INDEX IF EXISTS idx_system_login_histories_result;
DROP INDEX IF EXISTS idx_system_login_histories_admin_id;

DROP TABLE IF EXISTS system_login_histories;

-- +goose StatementEnd
