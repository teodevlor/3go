-- name: CreateDistancePricingRule :one
INSERT INTO system_distance_pricing_rules (service_id, from_km, to_km, price_per_km, is_active)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetDistancePricingRuleByID :one
SELECT * FROM system_distance_pricing_rules
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListDistancePricingRulesByServiceID :many
SELECT * FROM system_distance_pricing_rules
WHERE service_id = $1 AND deleted_at IS NULL
ORDER BY from_km ASC;

-- name: ListDistancePricingRules :many
SELECT * FROM system_distance_pricing_rules
WHERE deleted_at IS NULL
  AND (sqlc.narg('service_id')::uuid IS NULL OR service_id = sqlc.narg('service_id')::uuid)
ORDER BY service_id ASC, from_km ASC;

-- name: UpdateDistancePricingRule :one
UPDATE system_distance_pricing_rules
SET
  service_id = $2,
  from_km = $3,
  to_km = $4,
  price_per_km = $5,
  is_active = $6,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteDistancePricingRule :exec
UPDATE system_distance_pricing_rules
SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL;
