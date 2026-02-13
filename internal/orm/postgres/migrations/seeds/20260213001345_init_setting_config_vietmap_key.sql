-- +goose Up
-- +goose StatementBegin
DO $$
DECLARE
    v_account_id uuid;
BEGIN
    SELECT id INTO v_account_id FROM accounts WHERE email = 'system@internal' LIMIT 1;

    IF v_account_id IS NULL THEN
        INSERT INTO accounts (email, password_hash, phone, created_at, updated_at)
        VALUES ('system@internal', '$2a$10$placeholderhashvaluewithsufficientlength', '00000000000', NOW(), NOW())
        RETURNING id INTO v_account_id;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM system_settings WHERE key = 'config_vietmap_key' AND account_id = v_account_id) THEN
        INSERT INTO system_settings (
            account_id, key, value, type, description, is_active, metadata, created_at, updated_at
        ) VALUES (
            v_account_id,
            'config_vietmap_key',
            '{"apiKey": "api_key_vietmap"}'::jsonb,
            'web_system',
            'Vietmap API key for map/geocoding services',
            true,
            '{}'::jsonb,
            NOW(),
            NOW()
        );
    END IF;
END
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM system_settings WHERE key = 'config_vietmap_key';
-- +goose StatementEnd
