-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS public.system_permissions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),

    resource varchar(100) NOT NULL,
    action   varchar(100) NOT NULL,

    code varchar(200)
        GENERATED ALWAYS AS (resource || '.' || action) STORED,

    name        varchar(255) NOT NULL,
    description text,

    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamptz,

    UNIQUE(resource, action),
    UNIQUE(code)
);

CREATE INDEX IF NOT EXISTS idx_system_permissions_deleted_at
ON public.system_permissions(deleted_at)
WHERE deleted_at IS NULL;

DROP TRIGGER IF EXISTS update_system_permissions_updated_at ON public.system_permissions;
CREATE TRIGGER update_system_permissions_updated_at
BEFORE UPDATE ON public.system_permissions
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS update_system_permissions_updated_at ON public.system_permissions;

DROP INDEX IF EXISTS idx_system_permissions_deleted_at;
DROP TABLE IF EXISTS public.system_permissions;

-- +goose StatementEnd
