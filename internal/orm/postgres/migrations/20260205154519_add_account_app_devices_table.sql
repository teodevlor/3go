-- +goose Up
-- +goose StatementBegin

-- =====================
-- Account App Devices table (bảng trung tâm: account + device + app_type)
-- =====================
CREATE TABLE IF NOT EXISTS account_app_devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    account_id UUID NOT NULL,
    device_id UUID NOT NULL,
    app_type VARCHAR(50) NOT NULL, -- user, driver, admin, ...

    -- Thông tin bổ sung (optional)
    fcm_token TEXT,                -- Firebase Cloud Messaging token cho push notification
    is_active BOOLEAN NOT NULL DEFAULT true,
    last_used_at TIMESTAMPTZ,

    metadata JSONB DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,

    -- Foreign keys
    CONSTRAINT fk_account_app_devices_account
        FOREIGN KEY (account_id)
        REFERENCES accounts(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_account_app_devices_device
        FOREIGN KEY (device_id)
        REFERENCES devices(id)
        ON DELETE CASCADE
);

-- =====================
-- Unique constraint: 1 account + 1 device + 1 app_type = duy nhất
-- =====================
CREATE UNIQUE INDEX IF NOT EXISTS idx_account_app_devices_unique
ON account_app_devices(account_id, device_id, app_type)
WHERE deleted_at IS NULL;

-- =====================
-- Indexes
-- =====================
CREATE INDEX IF NOT EXISTS idx_account_app_devices_account_id
ON account_app_devices(account_id);

CREATE INDEX IF NOT EXISTS idx_account_app_devices_device_id
ON account_app_devices(device_id);

CREATE INDEX IF NOT EXISTS idx_account_app_devices_app_type
ON account_app_devices(app_type);

CREATE INDEX IF NOT EXISTS idx_account_app_devices_is_active
ON account_app_devices(is_active);

CREATE INDEX IF NOT EXISTS idx_account_app_devices_deleted_at
ON account_app_devices(deleted_at)
WHERE deleted_at IS NOT NULL;

-- =====================
-- Trigger: auto update updated_at
-- =====================
DROP TRIGGER IF EXISTS update_account_app_devices_updated_at ON account_app_devices;
CREATE TRIGGER update_account_app_devices_updated_at
BEFORE UPDATE ON account_app_devices
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS update_account_app_devices_updated_at ON account_app_devices;

DROP INDEX IF EXISTS idx_account_app_devices_deleted_at;
DROP INDEX IF EXISTS idx_account_app_devices_is_active;
DROP INDEX IF EXISTS idx_account_app_devices_app_type;
DROP INDEX IF EXISTS idx_account_app_devices_device_id;
DROP INDEX IF EXISTS idx_account_app_devices_account_id;
DROP INDEX IF EXISTS idx_account_app_devices_unique;

DROP TABLE IF EXISTS account_app_devices;

-- +goose StatementEnd
