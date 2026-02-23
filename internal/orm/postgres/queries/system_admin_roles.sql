-- name: InsertAdminRole :exec
INSERT INTO system_admin_roles (admin_id, role_id, assigned_by)
VALUES ($1, $2, $3);

-- name: DeleteAdminRolesByAdminID :exec
DELETE FROM system_admin_roles WHERE admin_id = $1;

-- name: GetRoleIDsByAdminID :many
SELECT role_id FROM system_admin_roles WHERE admin_id = $1;
