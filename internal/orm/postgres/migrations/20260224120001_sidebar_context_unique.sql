-- +goose Up
-- +goose StatementBegin
-- Xóa index cũ (non-unique)
DROP INDEX IF EXISTS idx_system_sidebars_context;

-- Unique cho context: mỗi context chỉ có 1 sidebar (chưa xóa)
CREATE UNIQUE INDEX uq_system_sidebars_context
ON system_sidebars(context)
WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS uq_system_sidebars_context;

CREATE INDEX idx_system_sidebars_context
ON system_sidebars(context)
WHERE deleted_at IS NULL;
-- +goose StatementEnd
