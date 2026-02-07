-- name: GetSettingByKey :one
SELECT id, account_id, key, value, type, description, is_active, metadata, created_at, updated_at, deleted_at
FROM settings
WHERE key = $1 AND is_active = true
LIMIT 1;
