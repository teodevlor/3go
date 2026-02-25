-- +goose Up
-- +goose StatementBegin
-- Xóa unique index cũ (cho phép trùng code khi service_id đều NULL)
DROP INDEX IF EXISTS idx_driver_document_types_code_unique;

-- Unique cho loại giấy tờ chung (service_id IS NULL): mỗi code chỉ 1 bản ghi global
CREATE UNIQUE INDEX uq_driver_doc_types_code_global
ON driver_document_types(code)
WHERE service_id IS NULL;

-- Unique cho loại giấy tờ theo từng service: (code, service_id) không trùng
CREATE UNIQUE INDEX uq_driver_doc_types_code_service
ON driver_document_types(code, service_id)
WHERE service_id IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS uq_driver_doc_types_code_global;
DROP INDEX IF EXISTS uq_driver_doc_types_code_service;

CREATE UNIQUE INDEX idx_driver_document_types_code_unique
ON driver_document_types(code, service_id);
-- +goose StatementEnd
