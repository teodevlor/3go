-- name: CreateLoginHistory :one
INSERT INTO app_login_histories (
    account_id,
    device_id,
    app_type,
    result,
    failure_reason,
    ip_address,
    user_agent,
    location,
    metadata
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;
