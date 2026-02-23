-- name: CreateRole :one
INSERT INTO system_roles (code, name, description, is_active)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetRoleByID :one
SELECT * FROM system_roles
WHERE id = $1;

-- name: GetRoleByCode :one
SELECT * FROM system_roles
WHERE code = $1;

-- name: GetRolesByIDs :many
SELECT * FROM system_roles
WHERE id = ANY($1::uuid[]);

-- name: UpdateRole :one
UPDATE system_roles
SET
    code = $2,
    name = $3,
    description = $4,
    is_active = $5,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: ListRoles :many
SELECT * FROM system_roles
WHERE ($1::text = '' OR code ILIKE '%' || $1 || '%' OR name ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountRoles :one
SELECT COUNT(*) FROM system_roles
WHERE ($1::text = '' OR code ILIKE '%' || $1 || '%' OR name ILIKE '%' || $1 || '%');

-- name: DeleteRole :exec
DELETE FROM system_roles
WHERE id = $1;
