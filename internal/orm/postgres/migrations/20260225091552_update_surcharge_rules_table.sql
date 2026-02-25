-- +goose Up
-- +goose StatementBegin

ALTER TABLE system_surcharge_rules
    DROP COLUMN IF EXISTS condition,
    DROP COLUMN IF EXISTS surcharge_type,
    ADD COLUMN IF NOT EXISTS priority INT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS created_by UUID NOT NULL,
    ADD COLUMN IF NOT EXISTS updated_by UUID NOT NULL;

ALTER TABLE system_surcharge_rules
ADD CONSTRAINT fk_system_surcharge_rules_created_by
FOREIGN KEY (created_by) REFERENCES system_admins(id),
ADD CONSTRAINT fk_system_surcharge_rules_updated_by
FOREIGN KEY (updated_by) REFERENCES system_admins(id);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE system_surcharge_rules
    DROP CONSTRAINT IF EXISTS fk_system_surcharge_rules_created_by,
    DROP CONSTRAINT IF EXISTS fk_system_surcharge_rules_updated_by,
    DROP COLUMN IF EXISTS created_by,
    DROP COLUMN IF EXISTS updated_by,
    DROP COLUMN IF EXISTS priority,
    ADD COLUMN IF NOT EXISTS surcharge_type VARCHAR(100),
    ADD COLUMN IF NOT EXISTS condition JSONB DEFAULT '{}'::jsonb;

-- +goose StatementEnd
