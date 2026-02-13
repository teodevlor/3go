-- name: CreateSystemAdminRefreshToken :one
INSERT INTO system_admin_refresh_tokens (
    admin_id,
    refresh_token_hash,
    expires_at,
    ip_address,
    user_agent,
    metadata
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetSystemAdminRefreshTokenByHash :one
SELECT * FROM system_admin_refresh_tokens
WHERE refresh_token_hash = $1 
  AND is_revoked = false 
  AND deleted_at IS NULL
LIMIT 1;

-- name: RevokeSystemAdminRefreshTokenByHash :exec
UPDATE system_admin_refresh_tokens
SET is_revoked = true,
    revoked_at = CURRENT_TIMESTAMP,
    revoked_reason = $2
WHERE refresh_token_hash = $1
  AND is_revoked = false
  AND deleted_at IS NULL;

-- name: RevokeAllSystemAdminRefreshTokens :exec
UPDATE system_admin_refresh_tokens
SET is_revoked = true,
    revoked_at = CURRENT_TIMESTAMP,
    revoked_reason = $2
WHERE admin_id = $1
  AND is_revoked = false
  AND deleted_at IS NULL;

-- name: UpdateSystemAdminRefreshTokenActivity :exec
UPDATE system_admin_refresh_tokens
SET last_active_at = CURRENT_TIMESTAMP
WHERE id = $1
  AND is_revoked = false
  AND deleted_at IS NULL;
