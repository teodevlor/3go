-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS public.system_roles (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    code        varchar(64)  NOT NULL,
    name        varchar(255) NOT NULL,
    description text,
    is_active   boolean      NOT NULL DEFAULT true,

    created_at  timestamptz  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  timestamptz  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at  timestamptz,

    CONSTRAINT system_roles_code_key UNIQUE (code)
);

CREATE INDEX IF NOT EXISTS idx_system_roles_is_active
ON public.system_roles(is_active);

CREATE INDEX IF NOT EXISTS idx_system_roles_deleted_at
ON public.system_roles(deleted_at)
WHERE deleted_at IS NULL;

DROP TRIGGER IF EXISTS update_system_roles_updated_at ON public.system_roles;
CREATE TRIGGER update_system_roles_updated_at
BEFORE UPDATE ON public.system_roles
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS update_system_roles_updated_at ON public.system_roles;

DROP INDEX IF EXISTS idx_system_roles_deleted_at;
DROP INDEX IF EXISTS idx_system_roles_is_active;

DROP TABLE IF EXISTS public.system_roles;

-- +goose StatementEnd
