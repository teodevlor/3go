-- +goose Up
-- +goose StatementBegin
-- Seed superadmin. Yêu cầu: chạy make pg-up (migrations) trước, sau đó make pg-seed-up.
-- Nếu bảng system_admins chưa tồn tại thì bỏ qua (tránh lỗi khi chạy seed trước migration).

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.tables
        WHERE table_schema = 'public' AND table_name = 'system_admins'
    ) THEN
        IF NOT EXISTS (SELECT 1 FROM system_admins WHERE email = 'supperadmin@gogogo.com') THEN
            INSERT INTO system_admins (
                email,
                password_hash,
                full_name,
                department,
                is_active,
                created_at,
                updated_at
            ) VALUES (
                'supperadmin@gogogo.com',
                crypt('password', gen_salt('bf', 10)),
                'Super Admin',
                'admin',
                TRUE,
                NOW(),
                NOW()
            );
        END IF;
    END IF;
END
$$;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.tables
        WHERE table_schema = 'public' AND table_name = 'system_admins'
    ) THEN
        DELETE FROM system_admins WHERE email = 'supperadmin@gogogo.com';
    END IF;
END
$$;
-- +goose StatementEnd
