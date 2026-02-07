-- name: CreateDevice :one
INSERT INTO devices (
    device_uid,
    platform,
    device_name,
    os_version,
    app_version,
    metadata
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetDeviceByID :one
SELECT * FROM devices
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetDeviceByUID :one
SELECT * FROM devices
WHERE device_uid = $1 AND deleted_at IS NULL;

-- name: UpdateDevice :one
UPDATE devices
SET
    device_name = COALESCE(sqlc.narg('device_name'), device_name),
    os_version = COALESCE(sqlc.narg('os_version'), os_version),
    app_version = COALESCE(sqlc.narg('app_version'), app_version),
    metadata = COALESCE(sqlc.narg('metadata'), metadata),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteDevice :exec
UPDATE devices
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1;
