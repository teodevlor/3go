-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS system_package_size_pricing (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    service_id UUID NOT NULL,

    package_size VARCHAR(100) NOT NULL,
    extra_price DECIMAL(10, 2) NOT NULL,

    is_active BOOLEAN NOT NULL DEFAULT true,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_system_package_size_pricing_service
        FOREIGN KEY (service_id)
        REFERENCES system_services(id)
        ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_system_package_size_pricing_service_size_unique
ON system_package_size_pricing(service_id, package_size)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_system_package_size_pricing_service_id
ON system_package_size_pricing(service_id);

CREATE INDEX IF NOT EXISTS idx_system_package_size_pricing_is_active
ON system_package_size_pricing(is_active);

CREATE INDEX IF NOT EXISTS idx_system_package_size_pricing_deleted_at
ON system_package_size_pricing(deleted_at)
WHERE deleted_at IS NULL;

DROP TRIGGER IF EXISTS update_system_package_size_pricing_updated_at ON system_package_size_pricing;
CREATE TRIGGER update_system_package_size_pricing_updated_at
BEFORE UPDATE ON system_package_size_pricing
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS update_system_package_size_pricing_updated_at ON system_package_size_pricing;

DROP INDEX IF EXISTS idx_system_package_size_pricing_deleted_at;
DROP INDEX IF EXISTS idx_system_package_size_pricing_is_active;
DROP INDEX IF EXISTS idx_system_package_size_pricing_service_id;
DROP INDEX IF EXISTS idx_system_package_size_pricing_service_size_unique;

DROP TABLE IF EXISTS system_package_size_pricing;

-- +goose StatementEnd
