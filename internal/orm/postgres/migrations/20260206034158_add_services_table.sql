-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS system_services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    code VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    base_price DECIMAL(10, 2) NOT NULL,
    min_price DECIMAL(10, 2) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_system_services_code_unique
ON system_services(code)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_system_services_is_active
ON system_services(is_active);

CREATE INDEX IF NOT EXISTS idx_system_services_deleted_at
ON system_services(deleted_at)
WHERE deleted_at IS NULL;

DROP TRIGGER IF EXISTS update_system_services_updated_at ON system_services;
CREATE TRIGGER update_system_services_updated_at
BEFORE UPDATE ON system_services
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS update_system_services_updated_at ON system_services;

DROP INDEX IF EXISTS idx_system_services_deleted_at;
DROP INDEX IF EXISTS idx_system_services_is_active;
DROP INDEX IF EXISTS idx_system_services_code_unique;

DROP TABLE IF EXISTS system_services;

-- +goose StatementEnd
