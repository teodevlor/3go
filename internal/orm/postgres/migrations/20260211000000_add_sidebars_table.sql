-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS system_sidebars (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    context VARCHAR(100) NOT NULL,
    version VARCHAR(50) NOT NULL DEFAULT '1.0.0',
    generated_at TIMESTAMPTZ,
    items JSONB NOT NULL DEFAULT '[]'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_system_sidebars_context
ON system_sidebars(context)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_system_sidebars_deleted_at
ON system_sidebars(deleted_at)
WHERE deleted_at IS NULL;

DROP TRIGGER IF EXISTS update_system_sidebars_updated_at ON system_sidebars;
CREATE TRIGGER update_system_sidebars_updated_at
BEFORE UPDATE ON system_sidebars
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_system_sidebars_updated_at ON system_sidebars;
DROP INDEX IF EXISTS idx_system_sidebars_deleted_at;
DROP INDEX IF EXISTS idx_system_sidebars_context;
DROP TABLE IF EXISTS system_sidebars;
-- +goose StatementEnd
