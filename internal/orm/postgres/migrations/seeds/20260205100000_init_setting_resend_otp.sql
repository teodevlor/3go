-- +goose Up
-- +goose StatementBegin
DO $$
DECLARE
    v_account_id uuid;
BEGIN
    -- Attempt to find a system account or create a placeholder for system settings
    -- because settings.account_id is NOT NULL.
    SELECT id INTO v_account_id FROM accounts WHERE email = 'system@internal' LIMIT 1;

    IF v_account_id IS NULL THEN
        -- Create a dummy system account if it doesn't exist (phone NOT NULL + UNIQUE nên dùng placeholder)
        INSERT INTO accounts (email, password_hash, phone, created_at, updated_at)
        VALUES ('system@internal', '$2a$10$placeholderhashvaluewithsufficientlength', '00000000000', NOW(), NOW())
        RETURNING id INTO v_account_id;
    END IF;

    -- Insert config_resent_otp setting
    IF NOT EXISTS (SELECT 1 FROM settings WHERE key = 'config_resent_otp' AND account_id = v_account_id) THEN
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
            'config_resent_otp',
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
            'Config resent OTP with progressive blocking mechanism',
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
DELETE FROM settings WHERE key = 'config_resent_otp';
-- Optional: We generally don't delete the system account in down migration 
-- as it might be used by other things or created manually.
-- +goose StatementEnd
