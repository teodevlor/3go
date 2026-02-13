-- +goose Up
-- +goose StatementBegin

-- Ensure pgcrypto extension is available
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

DO $$
BEGIN
    -- Insert superadmin if not exists
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
END
$$;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM system_admins WHERE email = 'supperadmin@gogogo.com';

-- +goose StatementEnd
