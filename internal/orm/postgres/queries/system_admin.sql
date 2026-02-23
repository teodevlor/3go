-- name: GetSystemAdminByEmail :one
SELECT * FROM system_admins
WHERE email = $1 LIMIT 1;

-- name: GetSystemAdminByID :one
SELECT * FROM system_admins
WHERE id = $1 LIMIT 1;

-- name: CreateSystemAdmin :one
INSERT INTO system_admins (email, password_hash, full_name, department, is_active)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateSystemAdmin :one
UPDATE system_admins
SET email = $2, full_name = $3, department = $4, is_active = $5, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: ListSystemAdmins :many
SELECT * FROM system_admins
WHERE ($1::text = '' OR email ILIKE '%' || $1 || '%' OR full_name ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountSystemAdmins :one
SELECT COUNT(*) FROM system_admins
WHERE ($1::text = '' OR email ILIKE '%' || $1 || '%' OR full_name ILIKE '%' || $1 || '%');

-- name: DeleteSystemAdmin :exec
DELETE FROM system_admins WHERE id = $1;

-- name: UpdateSystemAdminLastLoginAt :exec
UPDATE system_admins
SET last_login_at = $2
WHERE id = $1;
