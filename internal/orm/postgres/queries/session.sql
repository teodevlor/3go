-- name: CreateSession :one
INSERT INTO sessions (
    account_app_device_id,
    refresh_token_hash,
    expires_at,
    ip_address,
    user_agent,
    metadata
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetSessionByID :one
SELECT * FROM sessions
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetSessionByRefreshTokenHash :one
SELECT * FROM sessions
WHERE refresh_token_hash = $1
  AND is_revoked = false
  AND deleted_at IS NULL;

-- name: ListActiveSessions :many
SELECT s.*, aad.account_id, aad.device_id, aad.app_type
FROM sessions s
JOIN account_app_devices aad ON s.account_app_device_id = aad.id
WHERE aad.account_id = $1
  AND s.is_revoked = false
  AND s.expires_at > CURRENT_TIMESTAMP
  AND s.deleted_at IS NULL
ORDER BY s.last_active_at DESC;

-- name: UpdateSessionActivity :exec
UPDATE sessions
SET last_active_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: RevokeSession :exec
UPDATE sessions
SET
    is_revoked = true,
    revoked_at = CURRENT_TIMESTAMP,
    revoked_reason = $2
WHERE id = $1;

-- name: RevokeSessionByRefreshToken :exec
UPDATE sessions
SET
    is_revoked = true,
    revoked_at = CURRENT_TIMESTAMP,
    revoked_reason = $2
WHERE refresh_token_hash = $1;

-- name: RevokeAllSessionsByAccount :exec
UPDATE sessions
SET
    is_revoked = true,
    revoked_at = CURRENT_TIMESTAMP,
    revoked_reason = $2
WHERE account_app_device_id IN (
    SELECT id FROM account_app_devices WHERE account_id = $1
);

-- name: RevokeAllSessionsByAccountAppDevice :exec
UPDATE sessions
SET
    is_revoked = true,
    revoked_at = CURRENT_TIMESTAMP,
    revoked_reason = $2
WHERE account_app_device_id = $1;

-- name: DeleteExpiredSessions :exec
UPDATE sessions
SET deleted_at = CURRENT_TIMESTAMP
WHERE expires_at < CURRENT_TIMESTAMP
  AND deleted_at IS NULL;
