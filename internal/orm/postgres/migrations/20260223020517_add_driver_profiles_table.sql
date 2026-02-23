-- +goose Up
-- +goose StatementBegin
CREATE TYPE driver_profile_status AS ENUM (
    'PENDING_PROFILE',
    'DOCUMENT_INCOMPLETE',
    'PENDING_VERIFICATION',
    'ACTIVE',
    'SUSPENDED',
    'REJECTED'
);

CREATE TABLE IF NOT EXISTS driver_profiles (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id uuid NOT NULL UNIQUE,
    full_name varchar(255) NOT NULL,
    date_of_birth date,
    gender varchar(20),
    address text,

    global_status driver_profile_status NOT NULL DEFAULT 'PENDING_PROFILE',

    rating numeric(3,2) DEFAULT 5.0,
    total_completed_orders integer DEFAULT 0,

    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamptz
);

ALTER TABLE driver_profiles
ADD CONSTRAINT fk_driver_profiles_account
FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS driver_profiles;
DROP TYPE IF EXISTS driver_profile_status;
-- +goose StatementEnd
