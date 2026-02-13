-- name: CreateSystemLoginHistory :one
INSERT INTO system_login_histories (
    admin_id,
    result,
    failure_reason,
    ip_address,
    user_agent,
    location,
    metadata
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;
