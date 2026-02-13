-- name: CreateZone :one
INSERT INTO system_zones (
  code,
  name,
  polygon,
  price_multiplier,
  is_active
)
VALUES (
  $1,
  $2,
  ST_SetSRID(
    ST_GeomFromGeoJSON($3),
    4326
  ),
  $4,
  $5
)
RETURNING *;

-- name: GetZoneByID :one
SELECT id, code, name, ST_AsGeoJSON(polygon)::text AS polygon_geojson, price_multiplier, is_active, created_at, updated_at, deleted_at
FROM system_zones
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetZoneByCode :one
SELECT id, code, name, ST_AsGeoJSON(polygon)::text AS polygon_geojson, price_multiplier, is_active, created_at, updated_at, deleted_at
FROM system_zones
WHERE code = $1 AND deleted_at IS NULL;

-- name: ListZones :many
SELECT id, code, name, ST_AsGeoJSON(polygon)::text AS polygon_geojson, price_multiplier, is_active, created_at, updated_at, deleted_at
FROM system_zones
WHERE deleted_at IS NULL
  AND ($1 = '' OR name ILIKE '%' || $1 || '%' OR code ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountZones :one
SELECT COUNT(*) FROM system_zones
WHERE deleted_at IS NULL
  AND ($1 = '' OR name ILIKE '%' || $1 || '%' OR code ILIKE '%' || $1 || '%');

-- name: UpdateZone :one
UPDATE system_zones
SET
  code = $2,
  name = $3,
  polygon = ST_SetSRID(ST_GeomFromGeoJSON($4), 4326),
  price_multiplier = $5,
  is_active = $6,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING id, code, name, polygon, price_multiplier, is_active, created_at, updated_at, deleted_at;

-- name: DeleteZone :exec
UPDATE system_zones
SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL;
