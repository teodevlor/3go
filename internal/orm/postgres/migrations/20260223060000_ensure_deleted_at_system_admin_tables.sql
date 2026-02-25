-- +goose Up
-- +goose StatementBegin

-- Thêm cột deleted_at vào system_admins nếu chưa có (sửa lỗi login: column "deleted_at" does not exist)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'system_admins' AND column_name = 'deleted_at'
    ) THEN
        ALTER TABLE public.system_admins ADD COLUMN deleted_at TIMESTAMPTZ;
        CREATE INDEX IF NOT EXISTS idx_system_admins_deleted_at
            ON public.system_admins (deleted_at)
            WHERE deleted_at IS NULL;
    END IF;
END;
$$;

-- Thêm cột deleted_at vào system_admin_refresh_tokens nếu chưa có
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'system_admin_refresh_tokens' AND column_name = 'deleted_at'
    ) THEN
        ALTER TABLE public.system_admin_refresh_tokens ADD COLUMN deleted_at TIMESTAMPTZ;
        CREATE INDEX IF NOT EXISTS idx_system_admin_refresh_tokens_deleted_at
            ON public.system_admin_refresh_tokens(deleted_at)
            WHERE deleted_at IS NOT NULL;
    END IF;
END;
$$;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Không rollback tự động (giữ cột để tránh phá app). Nếu cần rollback thủ công:
-- ALTER TABLE public.system_admin_refresh_tokens DROP COLUMN IF EXISTS deleted_at;
-- ALTER TABLE public.system_admins DROP COLUMN IF EXISTS deleted_at;

-- +goose StatementEnd
