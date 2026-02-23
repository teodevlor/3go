-- +goose Up
-- +goose StatementBegin

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.tables
        WHERE table_schema = 'public' AND table_name = 'system_permissions'
    ) THEN
        INSERT INTO public.system_permissions (resource, action, name, description) VALUES
-- ZONE (khu vực)
('ZONE', 'CREATE', 'Tạo khu vực', 'Cho phép tạo khu vực giao hàng mới.'),
('ZONE', 'LIST', 'Danh sách khu vực', 'Cho phép xem danh sách khu vực.'),
('ZONE', 'READ', 'Chi tiết khu vực', 'Cho phép xem chi tiết một khu vực.'),
('ZONE', 'UPDATE', 'Cập nhật khu vực', 'Cho phép cập nhật thông tin khu vực.'),
('ZONE', 'DELETE', 'Xóa khu vực', 'Cho phép xóa khu vực.'),

-- SERVICE (dịch vụ)
('SERVICE', 'CREATE', 'Tạo dịch vụ', 'Cho phép tạo dịch vụ mới.'),
('SERVICE', 'LIST', 'Danh sách dịch vụ', 'Cho phép xem danh sách dịch vụ.'),
('SERVICE', 'READ', 'Chi tiết dịch vụ', 'Cho phép xem chi tiết một dịch vụ.'),
('SERVICE', 'UPDATE', 'Cập nhật dịch vụ', 'Cho phép cập nhật thông tin dịch vụ.'),
('SERVICE', 'DELETE', 'Xóa dịch vụ', 'Cho phép xóa dịch vụ.'),

-- DISTANCE_PRICING_RULE (quy tắc giá theo km)
('DISTANCE_PRICING_RULE', 'CREATE', 'Tạo quy tắc giá theo km', 'Cho phép tạo quy tắc giá theo khoảng cách.'),
('DISTANCE_PRICING_RULE', 'LIST', 'Danh sách quy tắc giá theo km', 'Cho phép xem danh sách quy tắc giá theo km.'),
('DISTANCE_PRICING_RULE', 'READ', 'Chi tiết quy tắc giá theo km', 'Cho phép xem chi tiết quy tắc.'),
('DISTANCE_PRICING_RULE', 'UPDATE', 'Cập nhật quy tắc giá theo km', 'Cho phép cập nhật quy tắc giá theo km.'),
('DISTANCE_PRICING_RULE', 'DELETE', 'Xóa quy tắc giá theo km', 'Cho phép xóa quy tắc giá theo km.'),

-- SURCHARGE_RULE (quy tắc phụ thu)
('SURCHARGE_RULE', 'CREATE', 'Tạo quy tắc phụ thu', 'Cho phép tạo quy tắc phụ thu.'),
('SURCHARGE_RULE', 'LIST', 'Danh sách quy tắc phụ thu', 'Cho phép xem danh sách quy tắc phụ thu.'),
('SURCHARGE_RULE', 'READ', 'Chi tiết quy tắc phụ thu', 'Cho phép xem chi tiết quy tắc phụ thu.'),
('SURCHARGE_RULE', 'UPDATE', 'Cập nhật quy tắc phụ thu', 'Cho phép cập nhật quy tắc phụ thu.'),
('SURCHARGE_RULE', 'DELETE', 'Xóa quy tắc phụ thu', 'Cho phép xóa quy tắc phụ thu.'),

-- PACKAGE_SIZE_PRICING (quy tắc giá theo kích thước gói)
('PACKAGE_SIZE_PRICING', 'CREATE', 'Tạo quy tắc giá theo kích thước gói', 'Cho phép tạo quy tắc giá theo kích thước gói.'),
('PACKAGE_SIZE_PRICING', 'LIST', 'Danh sách quy tắc giá theo kích thước gói', 'Cho phép xem danh sách quy tắc.'),
('PACKAGE_SIZE_PRICING', 'READ', 'Chi tiết quy tắc giá theo kích thước gói', 'Cho phép xem chi tiết quy tắc.'),
('PACKAGE_SIZE_PRICING', 'UPDATE', 'Cập nhật quy tắc giá theo kích thước gói', 'Cho phép cập nhật quy tắc.'),
('PACKAGE_SIZE_PRICING', 'DELETE', 'Xóa quy tắc giá theo kích thước gói', 'Cho phép xóa quy tắc.'),

-- SIDEBAR (cấu hình sidebar)
('SIDEBAR', 'CREATE', 'Tạo sidebar', 'Cho phép tạo cấu hình sidebar.'),
('SIDEBAR', 'LIST', 'Danh sách sidebar', 'Cho phép xem danh sách sidebar.'),
('SIDEBAR', 'READ', 'Chi tiết sidebar', 'Cho phép xem chi tiết sidebar.'),
('SIDEBAR', 'UPDATE', 'Cập nhật sidebar', 'Cho phép cập nhật sidebar.'),
('SIDEBAR', 'DELETE', 'Xóa sidebar', 'Cho phép xóa sidebar.'),

-- ROLE (vai trò)
('ROLE', 'CREATE', 'Tạo vai trò', 'Cho phép tạo vai trò mới.'),
('ROLE', 'LIST', 'Danh sách vai trò', 'Cho phép xem danh sách vai trò.'),
('ROLE', 'READ', 'Chi tiết vai trò', 'Cho phép xem chi tiết vai trò.'),
('ROLE', 'UPDATE', 'Cập nhật vai trò', 'Cho phép cập nhật vai trò.'),
('ROLE', 'DELETE', 'Xóa vai trò', 'Cho phép xóa vai trò.'),

-- ADMIN (quản trị viên)
('ADMIN', 'CREATE', 'Thêm quản trị viên', 'Cho phép thêm quản trị viên.'),
('ADMIN', 'LIST', 'Danh sách quản trị viên', 'Cho phép xem danh sách quản trị viên.'),
('ADMIN', 'READ', 'Chi tiết quản trị viên', 'Cho phép xem chi tiết quản trị viên.'),
('ADMIN', 'UPDATE', 'Cập nhật quản trị viên', 'Cho phép cập nhật thông tin quản trị viên.'),
('ADMIN', 'DELETE', 'Xóa quản trị viên', 'Cho phép xóa quản trị viên.'),

-- PERMISSION (quyền)
('PERMISSION', 'CREATE', 'Tạo quyền', 'Cho phép tạo quyền mới.'),
('PERMISSION', 'LIST', 'Danh sách quyền', 'Cho phép xem danh sách quyền.'),
('PERMISSION', 'READ', 'Chi tiết quyền', 'Cho phép xem chi tiết quyền.'),
('PERMISSION', 'UPDATE', 'Cập nhật quyền', 'Cho phép cập nhật quyền.'),
('PERMISSION', 'DELETE', 'Xóa quyền', 'Cho phép xóa quyền.')
        ON CONFLICT (resource, action) DO UPDATE SET
            name        = EXCLUDED.name,
            description = EXCLUDED.description,
            updated_at  = CURRENT_TIMESTAMP;
    END IF;
END
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.tables
        WHERE table_schema = 'public' AND table_name = 'system_permissions'
    ) THEN
        DELETE FROM public.system_permissions
        WHERE (resource, action) IN (
            ('ZONE', 'CREATE'), ('ZONE', 'LIST'), ('ZONE', 'READ'), ('ZONE', 'UPDATE'), ('ZONE', 'DELETE'),
            ('SERVICE', 'CREATE'), ('SERVICE', 'LIST'), ('SERVICE', 'READ'), ('SERVICE', 'UPDATE'), ('SERVICE', 'DELETE'),
            ('DISTANCE_PRICING_RULE', 'CREATE'), ('DISTANCE_PRICING_RULE', 'LIST'), ('DISTANCE_PRICING_RULE', 'READ'), ('DISTANCE_PRICING_RULE', 'UPDATE'), ('DISTANCE_PRICING_RULE', 'DELETE'),
            ('SURCHARGE_RULE', 'CREATE'), ('SURCHARGE_RULE', 'LIST'), ('SURCHARGE_RULE', 'READ'), ('SURCHARGE_RULE', 'UPDATE'), ('SURCHARGE_RULE', 'DELETE'),
            ('PACKAGE_SIZE_PRICING', 'CREATE'), ('PACKAGE_SIZE_PRICING', 'LIST'), ('PACKAGE_SIZE_PRICING', 'READ'), ('PACKAGE_SIZE_PRICING', 'UPDATE'), ('PACKAGE_SIZE_PRICING', 'DELETE'),
            ('SIDEBAR', 'CREATE'), ('SIDEBAR', 'LIST'), ('SIDEBAR', 'READ'), ('SIDEBAR', 'UPDATE'), ('SIDEBAR', 'DELETE'),
            ('ROLE', 'CREATE'), ('ROLE', 'LIST'), ('ROLE', 'READ'), ('ROLE', 'UPDATE'), ('ROLE', 'DELETE'),
            ('ADMIN', 'CREATE'), ('ADMIN', 'LIST'), ('ADMIN', 'READ'), ('ADMIN', 'UPDATE'), ('ADMIN', 'DELETE'),
            ('PERMISSION', 'CREATE'), ('PERMISSION', 'LIST'), ('PERMISSION', 'READ'), ('PERMISSION', 'UPDATE'), ('PERMISSION', 'DELETE')
        );
    END IF;
END
$$;
-- +goose StatementEnd
