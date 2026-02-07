-- name: CreateAccountAppDevice :one
INSERT INTO account_app_devices (
    account_id,
    device_id,
    app_type,
    fcm_token,
    is_active,
    metadata
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetAccountAppDeviceByID :one
SELECT * FROM account_app_devices
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetAccountAppDevice :one
SELECT * FROM account_app_devices
WHERE account_id = $1
  AND device_id = $2
  AND app_type = $3
  AND deleted_at IS NULL;

-- name: ListAccountAppDevices :many
SELECT * FROM account_app_devices
WHERE account_id = $1 AND deleted_at IS NULL
ORDER BY last_used_at DESC NULLS LAST, created_at DESC;

-- name: UpdateAccountAppDevice :one
UPDATE account_app_devices
SET
    fcm_token = COALESCE(sqlc.narg('fcm_token'), fcm_token),
    is_active = COALESCE(sqlc.narg('is_active'), is_active),
    last_used_at = COALESCE(sqlc.narg('last_used_at'), last_used_at),
    metadata = COALESCE(sqlc.narg('metadata'), metadata),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteAccountAppDevice :exec
UPDATE account_app_devices
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1;
