-- name: CreateDriverDocumentType :one
INSERT INTO driver_document_types (
  code,
  name,
  description,
  is_required,
  require_expire_date,
  service_id,
  is_active
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetDriverDocumentTypeByID :one
SELECT * FROM driver_document_types
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListRequiredDriverDocumentTypesByServiceID :many
-- Trả về document types áp dụng cho service: theo service_id HOẶC chung (service_id IS NULL), chỉ lấy bản ghi đang active.
SELECT * FROM driver_document_types
WHERE deleted_at IS NULL
  AND is_active = true
  AND (service_id = $1 OR service_id IS NULL)
ORDER BY service_id NULLS LAST, created_at ASC;

-- name: GetDriverDocumentTypeByCodeGlobal :one
SELECT * FROM driver_document_types
WHERE code = $1 AND service_id IS NULL AND deleted_at IS NULL;

-- name: GetDriverDocumentTypeByCodeAndServiceID :one
SELECT * FROM driver_document_types
WHERE code = $1 AND service_id = $2 AND deleted_at IS NULL;

-- name: ListDriverDocumentTypes :many
SELECT * FROM driver_document_types
WHERE deleted_at IS NULL
  AND ($1 = '' OR name ILIKE '%' || $1 || '%' OR code ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListDriverDocumentTypesByServiceID :many
SELECT * FROM driver_document_types
WHERE deleted_at IS NULL AND service_id = $1
  AND ($2 = '' OR name ILIKE '%' || $2 || '%' OR code ILIKE '%' || $2 || '%')
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: CountDriverDocumentTypes :one
SELECT COUNT(*) FROM driver_document_types
WHERE deleted_at IS NULL
  AND ($1 = '' OR name ILIKE '%' || $1 || '%' OR code ILIKE '%' || $1 || '%');

-- name: CountDriverDocumentTypesByServiceID :one
SELECT COUNT(*) FROM driver_document_types
WHERE deleted_at IS NULL AND service_id = $1
  AND ($2 = '' OR name ILIKE '%' || $2 || '%' OR code ILIKE '%' || $2 || '%');

-- name: UpdateDriverDocumentType :one
UPDATE driver_document_types
SET
  code = $2,
  name = $3,
  description = $4,
  is_required = $5,
  require_expire_date = $6,
  service_id = $7,
  is_active = $8
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteDriverDocumentType :exec
UPDATE driver_document_types
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL;
