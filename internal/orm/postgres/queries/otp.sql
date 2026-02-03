-- name: CreateOTP :one
INSERT INTO otps (
    target,
    otp_code,
    purpose,
    max_attempt,
    expires_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetActiveOTP :one
SELECT * FROM otps
WHERE target = $1 
  AND purpose = $2 
  AND status = 'active'
  AND expires_at > now()
ORDER BY created_at DESC
LIMIT 1;

-- name: MarkOTPAsUsed :exec
UPDATE otps
SET status = 'used',
    used_at = now(),
    updated_at = now()
WHERE id = $1;

-- name: MarkOTPAsUsedWithCount :exec
UPDATE otps
SET status = 'used',
    used_at = now(),
    attempt_count = $2,
    updated_at = now()
WHERE id = $1;

-- name: IncrementOTPAttempt :exec
UPDATE otps
SET attempt_count = attempt_count + 1,
    updated_at = now()
WHERE id = $1;

-- name: LockOTP :exec
UPDATE otps
SET status = 'locked',
    updated_at = now()
WHERE id = $1;

-- name: LockOTPWithCount :exec
UPDATE otps
SET status = 'locked',
    attempt_count = $2,
    updated_at = now()
WHERE id = $1;

-- name: ExpireOldOTPs :exec
UPDATE otps
SET status = 'expired',
    updated_at = now()
WHERE expires_at < now()
  AND status = 'active';
