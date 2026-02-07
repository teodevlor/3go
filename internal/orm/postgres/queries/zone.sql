--name: CreateZone :one
INSERT INTO zones (
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
);
