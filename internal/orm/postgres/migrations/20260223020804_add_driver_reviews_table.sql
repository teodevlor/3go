-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS driver_reviews (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    driver_id uuid NOT NULL,
    order_id uuid NOT NULL UNIQUE,
    customer_id uuid NOT NULL,
    rating smallint NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment text,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamptz
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS driver_reviews;
-- +goose StatementEnd
