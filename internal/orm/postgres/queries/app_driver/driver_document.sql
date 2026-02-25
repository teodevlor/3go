-- name: CreateDriverDocument :one
INSERT INTO driver_documents (
  driver_id,
  document_type_id,
  file_url,
  expire_at,
  status
)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetDriverDocumentByID :one
SELECT * FROM driver_documents
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListDriverDocumentsByDriverID :many
SELECT * FROM driver_documents
WHERE driver_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: UpdateDriverDocument :one
UPDATE driver_documents
SET
  file_url = $2,
  expire_at = $3,
  status = $4,
  reject_reason = $5,
  verified_at = $6,
  verified_by = $7,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateDriverDocumentPartial :one
-- Bulk update (PATCH): chỉ cập nhật các field được truyền (khác NULL).
UPDATE driver_documents
SET
  file_url = COALESCE(sqlc.narg('file_url'), file_url),
  expire_at = COALESCE(sqlc.narg('expire_at'), expire_at),
  status = COALESCE(sqlc.narg('status'), status),
  reject_reason = COALESCE(sqlc.narg('reject_reason'), reject_reason),
  updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg('id') AND deleted_at IS NULL
RETURNING *;

-- name: DeleteDriverDocument :exec
UPDATE driver_documents
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL;
