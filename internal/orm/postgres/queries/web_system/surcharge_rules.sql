-- name: CreateSurchargeRule :one
INSERT INTO system_surcharge_rules (service_id, zone_id, amount, unit, is_active, priority, created_by, updated_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, service_id, zone_id, amount, unit, is_active, created_at, updated_at, deleted_at, priority, created_by, updated_by;

-- name: GetSurchargeRuleByID :one
SELECT id, service_id, zone_id, amount, unit, is_active, created_at, updated_at, deleted_at, priority, created_by, updated_by
FROM system_surcharge_rules
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListSurchargeRules :many
SELECT id, service_id, zone_id, amount, unit, is_active, created_at, updated_at, deleted_at, priority, created_by, updated_by
FROM system_surcharge_rules
WHERE deleted_at IS NULL
  AND (sqlc.narg('service_id')::uuid IS NULL OR service_id = sqlc.narg('service_id')::uuid)
  AND (sqlc.narg('zone_id')::uuid IS NULL OR zone_id = sqlc.narg('zone_id')::uuid)
ORDER BY service_id ASC, zone_id ASC, created_at DESC;

-- name: UpdateSurchargeRule :one
UPDATE system_surcharge_rules
SET
  service_id = $2,
  zone_id = $3,
  amount = $4,
  unit = $5,
  is_active = $6,
  priority = $7,
  updated_by = $8,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING id, service_id, zone_id, amount, unit, is_active, created_at, updated_at, deleted_at, priority, created_by, updated_by;

-- name: DeleteSurchargeRule :exec
UPDATE system_surcharge_rules
SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL;
