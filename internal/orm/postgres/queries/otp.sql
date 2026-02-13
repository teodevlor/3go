-- name: CreateOTP :one
INSERT INTO system_otps (
    target,
    otp_code,
    purpose,
    max_attempt,
    expires_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetActiveOTP :one
SELECT * FROM system_otps
WHERE target = $1 
  AND purpose = $2 
  AND status = 'active'
  AND expires_at > now()
ORDER BY created_at DESC
LIMIT 1;

-- name: MarkOTPAsUsed :exec
UPDATE system_otps
SET status = 'used',
    used_at = now(),
    updated_at = now()
WHERE id = $1;

-- name: MarkOTPAsUsedWithCount :exec
UPDATE system_otps
SET status = 'used',
    used_at = now(),
    attempt_count = $2,
    updated_at = now()
WHERE id = $1;

-- name: IncrementOTPAttempt :exec
UPDATE system_otps
SET attempt_count = attempt_count + 1,
    updated_at = now()
WHERE id = $1;

-- name: LockOTP :exec
UPDATE system_otps
SET status = 'locked',
    updated_at = now()
WHERE id = $1;

-- name: LockOTPWithCount :exec
UPDATE system_otps
SET status = 'locked',
    attempt_count = $2,
    updated_at = now()
WHERE id = $1;

-- name: ExpireOldOTPs :exec
UPDATE system_otps
SET status = 'expired',
    updated_at = now()
WHERE expires_at < now()
  AND status = 'active';

-- name: GetLastOTPCreatedAt :one
SELECT created_at FROM system_otps
WHERE target = $1 AND purpose = $2
ORDER BY created_at DESC
LIMIT 1;

-- name: CountOTPsCreatedSince :one
SELECT COUNT(*)::int FROM system_otps
WHERE target = $1 AND purpose = $2 AND created_at >= $3;

-- name: GetOldestOTPCreatedAtSince :one
SELECT created_at FROM system_otps
WHERE target = $1 AND purpose = $2 AND created_at >= $3
ORDER BY created_at ASC
LIMIT 1;
