-- name: CreateOTPAudit :one
INSERT INTO otp_audits (
    otp_id,
    target,
    purpose,
    attempt_number,
    result,
    failure_reason,
    ip_address,
    user_agent,
    metadata
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;