-- name: CreateDriverProfile :one
INSERT INTO driver_profiles (
    account_id,
    full_name,
    date_of_birth,
    gender,
    address,
    global_status
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetDriverProfileByAccountID :one
SELECT * FROM driver_profiles
WHERE account_id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: GetDriverProfileByID :one
SELECT * FROM driver_profiles
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: UpdateDriverProfile :one
UPDATE driver_profiles
SET
    full_name = COALESCE(sqlc.narg('full_name'), full_name),
    date_of_birth = COALESCE(sqlc.narg('date_of_birth'), date_of_birth),
    gender = COALESCE(sqlc.narg('gender'), gender),
    address = COALESCE(sqlc.narg('address'), address),
    global_status = COALESCE(sqlc.narg('global_status'), global_status),
    updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg('id') AND deleted_at IS NULL
RETURNING *;
