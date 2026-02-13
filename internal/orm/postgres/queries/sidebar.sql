-- name: CreateSidebar :one
INSERT INTO system_sidebars (context, version, generated_at, items)
VALUES ($1, $2, $3, $4)
RETURNING id, context, version, generated_at, items, created_at, updated_at, deleted_at;

-- name: GetSidebarByID :one
SELECT id, context, version, generated_at, items, created_at, updated_at, deleted_at
FROM system_sidebars
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListSidebars :many
SELECT id, context, version, generated_at, items, created_at, updated_at, deleted_at
FROM system_sidebars
WHERE deleted_at IS NULL
  AND ($1 = '' OR context = $1)
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountSidebars :one
SELECT COUNT(*) FROM system_sidebars
WHERE deleted_at IS NULL
  AND ($1 = '' OR context = $1);

-- name: UpdateSidebar :one
UPDATE system_sidebars
SET
  context = $2,
  version = $3,
  generated_at = $4,
  items = $5,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING id, context, version, generated_at, items, created_at, updated_at, deleted_at;

-- name: DeleteSidebar :exec
UPDATE system_sidebars
SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL;
