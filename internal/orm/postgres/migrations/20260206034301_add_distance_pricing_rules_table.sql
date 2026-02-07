-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS distance_pricing_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    service_id UUID NOT NULL,

    from_km NUMERIC(10, 2) NOT NULL,
    to_km NUMERIC(10, 2) NOT NULL,
    price_per_km DECIMAL(10, 2) NOT NULL,

    is_active BOOLEAN NOT NULL DEFAULT true,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_distance_pricing_rules_service
        FOREIGN KEY (service_id)
        REFERENCES services(id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_distance_pricing_rules_service_id
ON distance_pricing_rules(service_id);

CREATE INDEX IF NOT EXISTS idx_distance_pricing_rules_is_active
ON distance_pricing_rules(is_active);

CREATE INDEX IF NOT EXISTS idx_distance_pricing_rules_deleted_at
ON distance_pricing_rules(deleted_at)
WHERE deleted_at IS NULL;

DROP TRIGGER IF EXISTS update_distance_pricing_rules_updated_at ON distance_pricing_rules;
CREATE TRIGGER update_distance_pricing_rules_updated_at
BEFORE UPDATE ON distance_pricing_rules
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS update_distance_pricing_rules_updated_at ON distance_pricing_rules;

DROP INDEX IF EXISTS idx_distance_pricing_rules_deleted_at;
DROP INDEX IF EXISTS idx_distance_pricing_rules_is_active;
DROP INDEX IF EXISTS idx_distance_pricing_rules_service_id;

DROP TABLE IF EXISTS distance_pricing_rules;

-- +goose StatementEnd
