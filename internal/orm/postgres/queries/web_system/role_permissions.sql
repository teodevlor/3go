-- name: CreateRolePermission :exec
INSERT INTO system_role_permissions (role_id, permission_id)
VALUES ($1, $2);

-- name: CreateRolePermissionsBatch :exec
INSERT INTO system_role_permissions (role_id, permission_id)
SELECT $1, unnest($2::uuid[]);

-- name: DeleteRolePermissionsByRoleID :exec
DELETE FROM system_role_permissions
WHERE role_id = $1;

-- name: GetPermissionIDsByRoleID :many
SELECT permission_id FROM system_role_permissions
WHERE role_id = $1;

-- name: GetRolePermissionPairsByRoleIDs :many
SELECT role_id, permission_id FROM system_role_permissions
WHERE role_id = ANY($1::uuid[]);
