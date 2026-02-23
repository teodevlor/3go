-- +goose Up
-- +goose StatementBegin
CREATE TYPE driver_document_status AS ENUM (
    'PENDING',
    'APPROVED',
    'REJECTED'
);

CREATE TABLE IF NOT EXISTS driver_documents (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    driver_id uuid NOT NULL,
    document_type_id uuid NOT NULL,

    file_url text NOT NULL,
    expire_at date,

    status driver_document_status NOT NULL DEFAULT 'PENDING',

    reject_reason text,
    verified_at timestamptz,
    verified_by uuid,

    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamptz
);

ALTER TABLE driver_documents
ADD CONSTRAINT fk_driver_documents_driver
FOREIGN KEY (driver_id) REFERENCES driver_profiles(id) ON DELETE CASCADE;

ALTER TABLE driver_documents
ADD CONSTRAINT fk_driver_documents_type
FOREIGN KEY (document_type_id) REFERENCES driver_document_types(id);

ALTER TABLE driver_documents
ADD CONSTRAINT fk_driver_documents_admin
FOREIGN KEY (verified_by) REFERENCES system_admins(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS driver_documents;
DROP TYPE IF EXISTS driver_document_status;
-- +goose StatementEnd
