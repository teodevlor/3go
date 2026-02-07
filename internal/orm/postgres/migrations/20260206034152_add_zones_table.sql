-- +goose Up
-- +goose StatementBegin

-- Yêu cầu: cài PostGIS trước (macOS: brew install postgis)
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS zones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    code VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    polygon GEOMETRY(Polygon, 4326),

    price_multiplier DECIMAL(10, 2) NOT NULL DEFAULT 1.0,
    is_active BOOLEAN NOT NULL DEFAULT true,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_zones_code_unique
ON zones(code)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_zones_is_active
ON zones(is_active);

CREATE INDEX IF NOT EXISTS idx_zones_deleted_at
ON zones(deleted_at)
WHERE deleted_at IS NULL;

DROP TRIGGER IF EXISTS update_zones_updated_at ON zones;
CREATE TRIGGER update_zones_updated_at
BEFORE UPDATE ON zones
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS update_zones_updated_at ON zones;

DROP INDEX IF EXISTS idx_zones_deleted_at;
DROP INDEX IF EXISTS idx_zones_is_active;
DROP INDEX IF EXISTS idx_zones_code_unique;

DROP TABLE IF EXISTS zones;

-- +goose StatementEnd
