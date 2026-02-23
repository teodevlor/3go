-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS driver_document_types (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    code varchar(100) NOT NULL,
    name varchar(255) NOT NULL,
    description text,

    is_required boolean NOT NULL DEFAULT true,
    require_expire_date boolean NOT NULL DEFAULT false,

    service_id uuid NULL, -- null = áp dụng cho mọi service

    is_active boolean NOT NULL DEFAULT true,

    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamptz
);

ALTER TABLE driver_document_types
ADD CONSTRAINT fk_driver_document_types_service
FOREIGN KEY (service_id) REFERENCES system_services(id) ON DELETE CASCADE;

CREATE UNIQUE INDEX idx_driver_document_types_code_unique
ON driver_document_types(code, service_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS driver_document_types;
-- +goose StatementEnd
