-- name: CreateSurchargeRule :one
INSERT INTO system_surcharge_rules (service_id, zone_id, surcharge_type, amount, unit, condition, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetSurchargeRuleByID :one
SELECT * FROM system_surcharge_rules
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListSurchargeRules :many
SELECT * FROM system_surcharge_rules
WHERE deleted_at IS NULL
  AND (sqlc.narg('service_id')::uuid IS NULL OR service_id = sqlc.narg('service_id')::uuid)
  AND (sqlc.narg('zone_id')::uuid IS NULL OR zone_id = sqlc.narg('zone_id')::uuid)
ORDER BY service_id ASC, zone_id ASC, created_at DESC;

-- name: UpdateSurchargeRule :one
UPDATE system_surcharge_rules
SET
  service_id = $2,
  zone_id = $3,
  surcharge_type = $4,
  amount = $5,
  unit = $6,
  condition = $7,
  is_active = $8,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteSurchargeRule :exec
UPDATE system_surcharge_rules
SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL;
