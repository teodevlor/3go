-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS driver_profile_status_histories (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    driver_id uuid NOT NULL,
    from_status driver_profile_status,
    to_status driver_profile_status NOT NULL,
    changed_by uuid, -- admin_id hoặc system
    reason text,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamptz
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS driver_profile_status_histories;
-- +goose StatementEnd
