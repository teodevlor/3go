-- +goose Up
-- +goose StatementBegin
INSERT INTO public.system_permissions (resource, action, name, description) VALUES
('DRIVER_DOCUMENT_TYPE', 'CREATE', 'Tạo loại giấy tờ tài xế', 'Cho phép tạo loại giấy tờ (catalog) cho tài xế.'),
('DRIVER_DOCUMENT_TYPE', 'LIST', 'Danh sách loại giấy tờ tài xế', 'Cho phép xem danh sách loại giấy tờ.'),
('DRIVER_DOCUMENT_TYPE', 'READ', 'Chi tiết loại giấy tờ tài xế', 'Cho phép xem chi tiết một loại giấy tờ.'),
('DRIVER_DOCUMENT_TYPE', 'UPDATE', 'Cập nhật loại giấy tờ tài xế', 'Cho phép cập nhật loại giấy tờ.'),
('DRIVER_DOCUMENT_TYPE', 'DELETE', 'Xóa loại giấy tờ tài xế', 'Cho phép xóa loại giấy tờ.')
ON CONFLICT (resource, action) DO UPDATE SET
    name        = EXCLUDED.name,
    description = EXCLUDED.description,
    updated_at  = CURRENT_TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM public.system_permissions
WHERE resource = 'DRIVER_DOCUMENT_TYPE' AND action IN ('CREATE', 'LIST', 'READ', 'UPDATE', 'DELETE');
-- +goose StatementEnd
