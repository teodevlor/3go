-- +goose Up
-- +goose StatementBegin
-- Đảm bảo supper-admin luôn có full quyền (chạy sau khi seed permissions).

DO $$
DECLARE
    v_admin_id uuid;
    v_role_id  uuid;
BEGIN
    -- Đảm bảo các bảng cần thiết tồn tại
    IF EXISTS (
        SELECT 1 FROM information_schema.tables
        WHERE table_schema = 'public' AND table_name = 'system_admins'
    )
    AND EXISTS (
        SELECT 1 FROM information_schema.tables
        WHERE table_schema = 'public' AND table_name = 'system_roles'
    )
    AND EXISTS (
        SELECT 1 FROM information_schema.tables
        WHERE table_schema = 'public' AND table_name = 'system_admin_roles'
    )
    AND EXISTS (
        SELECT 1 FROM information_schema.tables
        WHERE table_schema = 'public' AND table_name = 'system_permissions'
    )
    AND EXISTS (
        SELECT 1 FROM information_schema.tables
        WHERE table_schema = 'public' AND table_name = 'system_role_permissions'
    ) THEN

        -- Lấy admin supperadmin
        SELECT id INTO v_admin_id
        FROM system_admins
        WHERE email = 'supperadmin@gogogo.com'
        LIMIT 1;

        IF v_admin_id IS NOT NULL THEN
            -- Lấy (hoặc tạo) role supper-admin
            SELECT id INTO v_role_id
            FROM public.system_roles
            WHERE code = 'supper-admin'
            LIMIT 1;

            IF v_role_id IS NULL THEN
                INSERT INTO public.system_roles (code, name, description, is_active)
                VALUES (
                    'supper-admin',
                    'Supper Admin',
                    'Role supper admin với toàn quyền hệ thống',
                    TRUE
                )
                RETURNING id INTO v_role_id;
            END IF;

            -- Đảm bảo supperadmin đã được gán role supper-admin
            IF NOT EXISTS (
                SELECT 1
                FROM public.system_admin_roles
                WHERE admin_id = v_admin_id AND role_id = v_role_id
            ) THEN
                INSERT INTO public.system_admin_roles (admin_id, role_id, assigned_at, assigned_by)
                VALUES (v_admin_id, v_role_id, NOW(), v_admin_id);
            END IF;

            -- Gán toàn bộ permissions hiện có cho role supper-admin
            INSERT INTO public.system_role_permissions (role_id, permission_id)
            SELECT v_role_id, p.id
            FROM public.system_permissions p
            WHERE NOT EXISTS (
                SELECT 1
                FROM public.system_role_permissions rp
                WHERE rp.role_id = v_role_id AND rp.permission_id = p.id
            );
        END IF;
    END IF;
END
$$;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DO $$
DECLARE
    v_role_id uuid;
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.tables
        WHERE table_schema = 'public' AND table_name = 'system_roles'
    )
    AND EXISTS (
        SELECT 1 FROM information_schema.tables
        WHERE table_schema = 'public' AND table_name = 'system_role_permissions'
    )
    AND EXISTS (
        SELECT 1 FROM information_schema.tables
        WHERE table_schema = 'public' AND table_name = 'system_admin_roles'
    ) THEN
        SELECT id INTO v_role_id
        FROM public.system_roles
        WHERE code = 'supper-admin'
        LIMIT 1;

        IF v_role_id IS NOT NULL THEN
            -- Chỉ xóa mapping, không xóa role để tránh ảnh hưởng chỗ khác
            DELETE FROM public.system_role_permissions
            WHERE role_id = v_role_id;

            DELETE FROM public.system_admin_roles
            WHERE role_id = v_role_id;
        END IF;
    END IF;
END
$$;

-- +goose StatementEnd

