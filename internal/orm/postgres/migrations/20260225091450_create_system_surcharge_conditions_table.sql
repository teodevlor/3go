-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS system_surcharge_conditions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(100) UNIQUE NOT NULL,
    condition_type VARCHAR(50) NOT NULL,
    config JSONB NOT NULL,
    is_active BOOLEAN DEFAULT true,
    CONSTRAINT chk_system_surcharge_conditions_type
        CHECK (condition_type IN ('time_window', 'weather', 'traffic', 'holiday'))
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS system_surcharge_conditions;

-- +goose StatementEnd
