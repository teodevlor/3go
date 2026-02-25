-- name: CreateDriverProfileStatusHistory :one
INSERT INTO driver_profile_status_histories (driver_id, from_status, to_status, changed_by, reason)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
