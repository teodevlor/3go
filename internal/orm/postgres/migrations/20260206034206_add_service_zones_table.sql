-- +goose Up
-- +goose StatementBegin

-- Bảng pivot: service – zone (nhiều-nhiều)
CREATE TABLE IF NOT EXISTS service_zones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    zone_id UUID NOT NULL,
    service_id UUID NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_service_zones_zone
        FOREIGN KEY (zone_id)
        REFERENCES zones(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_service_zones_service
        FOREIGN KEY (service_id)
        REFERENCES services(id)
        ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_service_zones_zone_service_unique
ON service_zones(zone_id, service_id)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_service_zones_zone_id
ON service_zones(zone_id);

CREATE INDEX IF NOT EXISTS idx_service_zones_service_id
ON service_zones(service_id);

CREATE INDEX IF NOT EXISTS idx_service_zones_deleted_at
ON service_zones(deleted_at)
WHERE deleted_at IS NULL;

DROP TRIGGER IF EXISTS update_service_zones_updated_at ON service_zones;
CREATE TRIGGER update_service_zones_updated_at
BEFORE UPDATE ON service_zones
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS update_service_zones_updated_at ON service_zones;

DROP INDEX IF EXISTS idx_service_zones_deleted_at;
DROP INDEX IF EXISTS idx_service_zones_service_id;
DROP INDEX IF EXISTS idx_service_zones_zone_id;
DROP INDEX IF EXISTS idx_service_zones_zone_service_unique;

DROP TABLE IF EXISTS service_zones;

-- +goose StatementEnd
