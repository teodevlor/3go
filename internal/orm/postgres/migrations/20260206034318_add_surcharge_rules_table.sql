-- +goose Up
-- +goose StatementBegin

-- unit: 'percent' | 'fixed'
-- condition: jsonb e.g. {"time_range":["17:00","19:00"],"days":["mon","tue","wed","thu","fri"]}
CREATE TABLE IF NOT EXISTS surcharge_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    service_id UUID NOT NULL,
    zone_id UUID NOT NULL,

    surcharge_type VARCHAR(100) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    unit VARCHAR(100) NOT NULL, -- 'percent' | 'fixed'

    condition JSONB DEFAULT '{}'::jsonb,

    is_active BOOLEAN NOT NULL DEFAULT true,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_surcharge_rules_service
        FOREIGN KEY (service_id)
        REFERENCES services(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_surcharge_rules_zone
        FOREIGN KEY (zone_id)
        REFERENCES zones(id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_surcharge_rules_service_id
ON surcharge_rules(service_id);

CREATE INDEX IF NOT EXISTS idx_surcharge_rules_zone_id
ON surcharge_rules(zone_id);

CREATE INDEX IF NOT EXISTS idx_surcharge_rules_is_active
ON surcharge_rules(is_active);

CREATE INDEX IF NOT EXISTS idx_surcharge_rules_deleted_at
ON surcharge_rules(deleted_at)
WHERE deleted_at IS NULL;

DROP TRIGGER IF EXISTS update_surcharge_rules_updated_at ON surcharge_rules;
CREATE TRIGGER update_surcharge_rules_updated_at
BEFORE UPDATE ON surcharge_rules
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS update_surcharge_rules_updated_at ON surcharge_rules;

DROP INDEX IF EXISTS idx_surcharge_rules_deleted_at;
DROP INDEX IF EXISTS idx_surcharge_rules_is_active;
DROP INDEX IF EXISTS idx_surcharge_rules_zone_id;
DROP INDEX IF EXISTS idx_surcharge_rules_service_id;

DROP TABLE IF EXISTS surcharge_rules;

-- +goose StatementEnd
