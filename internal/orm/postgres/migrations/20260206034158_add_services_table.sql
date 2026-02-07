-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS services (
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

CREATE UNIQUE INDEX IF NOT EXISTS idx_services_code_unique
ON services(code)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_services_is_active
ON services(is_active);

CREATE INDEX IF NOT EXISTS idx_services_deleted_at
ON services(deleted_at)
WHERE deleted_at IS NULL;

DROP TRIGGER IF EXISTS update_services_updated_at ON services;
CREATE TRIGGER update_services_updated_at
BEFORE UPDATE ON services
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS update_services_updated_at ON services;

DROP INDEX IF EXISTS idx_services_deleted_at;
DROP INDEX IF EXISTS idx_services_is_active;
DROP INDEX IF EXISTS idx_services_code_unique;

DROP TABLE IF EXISTS services;

-- +goose StatementEnd
