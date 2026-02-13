-- name: GetSystemAdminByEmail :one
SELECT * FROM system_admins
WHERE email = $1 LIMIT 1;

-- name: GetSystemAdminByID :one
SELECT * FROM system_admins
WHERE id = $1 LIMIT 1;

-- name: UpdateSystemAdminLastLoginAt :exec
UPDATE system_admins
SET last_login_at = $2
WHERE id = $1;
