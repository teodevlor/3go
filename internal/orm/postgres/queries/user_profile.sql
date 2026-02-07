-- name: CreateUserProfile :one
INSERT INTO user_profiles (
    account_id, full_name, avatar_url, is_active, metadata
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING id, account_id, full_name, avatar_url, is_active, metadata, created_at, updated_at, deleted_at;

-- name: GetUserProfileByAccountId :one
SELECT id, account_id, full_name, avatar_url, is_active, metadata, created_at, updated_at, deleted_at FROM user_profiles WHERE account_id = $1 LIMIT 1;

-- name: GetUserProfileByID :one
SELECT id, account_id, full_name, avatar_url, is_active, metadata, created_at, updated_at, deleted_at FROM user_profiles WHERE id = $1 LIMIT 1;

-- name: UpdateUserProfile :one
UPDATE user_profiles
SET 
    full_name = $2,
    avatar_url = $3,
    is_active = $4,
    metadata = $5,
    updated_at = $6
WHERE id = $1
RETURNING *;