-- +goose Up
-- +goose StatementBegin

-- Đổi tên bảng system_login_histories thành app_login_histories
ALTER TABLE system_login_histories RENAME TO app_login_histories;

-- Đổi tên các indexes
ALTER INDEX idx_system_login_histories_account_id RENAME TO idx_app_login_histories_account_id;
ALTER INDEX idx_system_login_histories_device_id RENAME TO idx_app_login_histories_device_id;
ALTER INDEX idx_system_login_histories_app_type RENAME TO idx_app_login_histories_app_type;
ALTER INDEX idx_system_login_histories_result RENAME TO idx_app_login_histories_result;
ALTER INDEX idx_system_login_histories_login_at RENAME TO idx_app_login_histories_login_at;
ALTER INDEX idx_system_login_histories_ip_address RENAME TO idx_app_login_histories_ip_address;
ALTER INDEX idx_system_login_histories_account_login_at RENAME TO idx_app_login_histories_account_login_at;

-- Đổi tên các constraints
ALTER TABLE app_login_histories RENAME CONSTRAINT fk_system_login_histories_account TO fk_app_login_histories_account;
ALTER TABLE app_login_histories RENAME CONSTRAINT fk_system_login_histories_device TO fk_app_login_histories_device;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Đổi tên lại về system_login_histories
ALTER TABLE app_login_histories RENAME TO system_login_histories;

-- Đổi tên lại các indexes
ALTER INDEX idx_app_login_histories_account_id RENAME TO idx_system_login_histories_account_id;
ALTER INDEX idx_app_login_histories_device_id RENAME TO idx_system_login_histories_device_id;
ALTER INDEX idx_app_login_histories_app_type RENAME TO idx_system_login_histories_app_type;
ALTER INDEX idx_app_login_histories_result RENAME TO idx_system_login_histories_result;
ALTER INDEX idx_app_login_histories_login_at RENAME TO idx_system_login_histories_login_at;
ALTER INDEX idx_app_login_histories_ip_address RENAME TO idx_system_login_histories_ip_address;
ALTER INDEX idx_app_login_histories_account_login_at RENAME TO idx_system_login_histories_account_login_at;

-- Đổi tên lại các constraints
ALTER TABLE system_login_histories RENAME CONSTRAINT fk_app_login_histories_account TO fk_system_login_histories_account;
ALTER TABLE system_login_histories RENAME CONSTRAINT fk_app_login_histories_device TO fk_system_login_histories_device;

-- +goose StatementEnd
