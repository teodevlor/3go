-- +goose Up
-- +goose StatementBegin
DO $$
DECLARE
    v_account_id uuid;
BEGIN
    -- Tìm hoặc tạo system account dùng cho cấu hình hệ thống (settings.account_id NOT NULL)
    SELECT id INTO v_account_id FROM accounts WHERE email = 'system@internal' LIMIT 1;

    IF v_account_id IS NULL THEN
        -- Tạo account hệ thống giả nếu chưa tồn tại (phone NOT NULL + UNIQUE nên dùng placeholder)
        INSERT INTO accounts (email, password_hash, phone, created_at, updated_at)
        VALUES ('system@internal', '$2a$10$placeholderhashvaluewithsufficientlength', '00000000000', NOW(), NOW())
        RETURNING id INTO v_account_id;
    END IF;

    -- Insert config_re_forgot_password setting (rate-limit cho request forgot password)
    IF NOT EXISTS (SELECT 1 FROM settings WHERE key = 'config_re_forgot_password' AND account_id = v_account_id) THEN
        INSERT INTO settings (
            account_id,
            key,
            value,
            type,
            description,
            is_active,
            metadata,
            created_at,
            updated_at
        ) VALUES (
            v_account_id,
            'config_re_forgot_password',
            '{
                "maxCount": 3,
                "timeOutExpired": 300,
                "timeOutResent": 60,
                "blockDurations": {
                    "violation1": 300,
                    "violation2": 900,
                    "violation3": 3600,
                    "violation4Plus": 86400
                },
                "trackingTTL": 90000
            }'::jsonb,
            'web_system',
            'Config re-request forgot password with progressive blocking mechanism',
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
DELETE FROM settings WHERE key = 'config_re_forgot_password';
-- Optional: We generally don't delete the system account in down migration 
-- as it might be used by other things or created manually.
-- +goose StatementEnd
