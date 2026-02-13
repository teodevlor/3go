-- name: CreatePackageSizePricing :one
INSERT INTO system_package_size_pricing (service_id, package_size, extra_price, is_active)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetPackageSizePricingByID :one
SELECT * FROM system_package_size_pricing
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListPackageSizePricings :many
SELECT * FROM system_package_size_pricing
WHERE deleted_at IS NULL
  AND (sqlc.narg('service_id')::uuid IS NULL OR service_id = sqlc.narg('service_id')::uuid)
ORDER BY service_id ASC, package_size ASC, created_at DESC;

-- name: UpdatePackageSizePricing :one
UPDATE system_package_size_pricing
SET
  service_id = $2,
  package_size = $3,
  extra_price = $4,
  is_active = $5,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeletePackageSizePricing :exec
UPDATE system_package_size_pricing
SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL;
