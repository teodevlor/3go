-- name: CreatePermission :one
INSERT INTO system_permissions (resource, action, name, description)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetPermissionByID :one
SELECT * FROM system_permissions
WHERE id = $1;

-- name: GetPermissionByCode :one
SELECT * FROM system_permissions
WHERE code = $1;

-- name: GetPermissionsByIDs :many
SELECT * FROM system_permissions
WHERE id = ANY($1::uuid[]);

-- name: UpdatePermission :one
UPDATE system_permissions
SET
    resource = $2,
    action = $3,
    name = $4,
    description = $5,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: ListPermissions :many
SELECT * FROM system_permissions
WHERE ($1::text = '' OR resource ILIKE '%' || $1 || '%' OR action ILIKE '%' || $1 || '%' OR name ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountPermissions :one
SELECT COUNT(*) FROM system_permissions
WHERE ($1::text = '' OR resource ILIKE '%' || $1 || '%' OR action ILIKE '%' || $1 || '%' OR name ILIKE '%' || $1 || '%');

-- name: DeletePermission :exec
DELETE FROM system_permissions
WHERE id = $1;
