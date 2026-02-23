-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS public.system_role_permissions (
    role_id       uuid NOT NULL,
    permission_id uuid NOT NULL,

    PRIMARY KEY (role_id, permission_id),

    CONSTRAINT fk_system_role_permissions_role
        FOREIGN KEY (role_id)
        REFERENCES public.system_roles(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_system_role_permissions_permission
        FOREIGN KEY (permission_id)
        REFERENCES public.system_permissions(id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_system_role_permissions_role_id
ON public.system_role_permissions(role_id);

CREATE INDEX IF NOT EXISTS idx_system_role_permissions_permission_id
ON public.system_role_permissions(permission_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_system_role_permissions_permission_id;
DROP INDEX IF EXISTS idx_system_role_permissions_role_id;

DROP TABLE IF EXISTS public.system_role_permissions;

-- +goose StatementEnd
