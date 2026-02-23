-- +goose Up
-- +goose StatementBegin
CREATE TYPE driver_service_status AS ENUM (
    'PENDING_DOCUMENT',
    'PENDING_APPROVAL',
    'ACTIVE',
    'SUSPENDED',
    'REJECTED'
);

CREATE TABLE IF NOT EXISTS driver_services (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    driver_id uuid NOT NULL,
    service_id uuid NOT NULL,

    status driver_service_status NOT NULL DEFAULT 'PENDING_DOCUMENT',

    approved_at timestamptz,
    approved_by uuid, -- system_admins

    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamptz,

    UNIQUE (driver_id, service_id)
);

ALTER TABLE driver_services
ADD CONSTRAINT fk_driver_services_driver
FOREIGN KEY (driver_id) REFERENCES driver_profiles(id) ON DELETE CASCADE;

ALTER TABLE driver_services
ADD CONSTRAINT fk_driver_services_service
FOREIGN KEY (service_id) REFERENCES system_services(id) ON DELETE CASCADE;

ALTER TABLE driver_services
ADD CONSTRAINT fk_driver_services_admin
FOREIGN KEY (approved_by) REFERENCES system_admins(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS driver_services;
DROP TYPE IF EXISTS driver_service_status;
-- +goose StatementEnd
