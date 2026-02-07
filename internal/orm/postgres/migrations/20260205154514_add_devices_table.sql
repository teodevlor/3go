-- +goose Up
-- +goose StatementBegin

-- =====================
-- Devices table (thiết bị vật lý, không biết user hay driver)
-- =====================
CREATE TABLE IF NOT EXISTS devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    device_uid VARCHAR(255) NOT NULL, -- UUID từ app (unique identifier từ mobile)
    platform VARCHAR(50) NOT NULL,    -- ios, android, web, ...
    device_name VARCHAR(255),         -- iPhone 15 Pro, Samsung Galaxy, ...
    os_version VARCHAR(100),          -- iOS 17.2, Android 14, ...
    app_version VARCHAR(100),         -- 1.0.0, 2.1.3, ...
    
    metadata JSONB DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

-- =====================
-- Indexes
-- =====================
CREATE UNIQUE INDEX IF NOT EXISTS idx_devices_device_uid_unique
ON devices(device_uid);

CREATE INDEX IF NOT EXISTS idx_devices_platform
ON devices(platform);

CREATE INDEX IF NOT EXISTS idx_devices_deleted_at
ON devices(deleted_at)
WHERE deleted_at IS NOT NULL;

-- =====================
-- Trigger: auto update updated_at
-- =====================
DROP TRIGGER IF EXISTS update_devices_updated_at ON devices;
CREATE TRIGGER update_devices_updated_at
BEFORE UPDATE ON devices
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS update_devices_updated_at ON devices;

DROP INDEX IF EXISTS idx_devices_deleted_at;
DROP INDEX IF EXISTS idx_devices_platform;
DROP INDEX IF EXISTS idx_devices_device_uid_unique;

DROP TABLE IF EXISTS devices;

-- +goose StatementEnd
