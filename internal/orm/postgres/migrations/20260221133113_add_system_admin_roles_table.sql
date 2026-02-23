-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS public.system_admin_roles (
    admin_id    uuid        NOT NULL,
    role_id     uuid        NOT NULL,

    assigned_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    assigned_by uuid,

    PRIMARY KEY (admin_id, role_id),

    CONSTRAINT fk_system_admin_roles_admin
        FOREIGN KEY (admin_id)
        REFERENCES public.system_admins(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_system_admin_roles_role
        FOREIGN KEY (role_id)
        REFERENCES public.system_roles(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_system_admin_roles_assigned_by
        FOREIGN KEY (assigned_by)
        REFERENCES public.system_admins(id)
        ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_system_admin_roles_admin_id
ON public.system_admin_roles(admin_id);

CREATE INDEX IF NOT EXISTS idx_system_admin_roles_role_id
ON public.system_admin_roles(role_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_system_admin_roles_role_id;
DROP INDEX IF EXISTS idx_system_admin_roles_admin_id;

DROP TABLE IF EXISTS public.system_admin_roles;

-- +goose StatementEnd
