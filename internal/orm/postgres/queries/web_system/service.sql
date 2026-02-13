-- name: CreateService :one
INSERT INTO system_services (
  code,
  name,
  base_price,
  min_price,
  is_active
)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetServiceByID :one
SELECT * FROM system_services
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetServiceByCode :one
SELECT * FROM system_services
WHERE code = $1 AND deleted_at IS NULL;

-- name: ListServices :many
SELECT * FROM system_services
WHERE deleted_at IS NULL
  AND ($1 = '' OR name ILIKE '%' || $1 || '%' OR code ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountServices :one
SELECT COUNT(*) FROM system_services
WHERE deleted_at IS NULL
  AND ($1 = '' OR name ILIKE '%' || $1 || '%' OR code ILIKE '%' || $1 || '%');

-- name: UpdateService :one
UPDATE system_services
SET
  code = $2,
  name = $3,
  base_price = $4,
  min_price = $5,
  is_active = $6,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteService :exec
UPDATE system_services
SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL;
