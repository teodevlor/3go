-- name: CreateServiceZone :one
INSERT INTO system_service_zones (zone_id, service_id)
VALUES ($1, $2)
RETURNING *;

-- name: ListServiceZonesByServiceID :many
SELECT * FROM system_service_zones
WHERE service_id = $1 AND deleted_at IS NULL
ORDER BY created_at ASC;

-- name: DeleteServiceZonesByServiceID :exec
UPDATE system_service_zones
SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE service_id = $1 AND deleted_at IS NULL;
