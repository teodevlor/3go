# API Documentation (Docusaurus)

Tài liệu API cho **Gogogo API**, cấu trúc song song với `src/modules/`.

**`docs-api` là project Docusaurus riêng** (site tài liệu tĩnh), chạy độc lập với Nest API.

## Port

| Ứng dụng   | Port  | URL                    |
|------------|-------|------------------------|
| **Nest API** (backend) | 3000 | http://localhost:3000  |
| **Docs API** (Docusaurus) | 3001 | http://localhost:3001 |

Docusaurus dùng **3001** để tránh trùng với API đang chạy 3000.

## Chạy

```bash
# Cài phụ thuộc (từ thư mục gốc project)
pnpm docs:install

# Chạy dev (mở http://localhost:3001)
pnpm docs:start

# Build
pnpm docs:build
```

Hoặc trong `docs-api/`:

```bash
pnpm install
pnpm start     # http://localhost:3001
pnpm build
pnpm serve     # xem bản build tại http://localhost:3001
```

## Cấu trúc

```
docs-api/
├── docs/
│   ├── intro.md              # Trang chủ
│   └── modules/              # Song song src/modules/
│       ├── overview.md
│       ├── accounts.md       # accounts
│       ├── auth-profile-user.md
│       ├── auth-profile-driver.md
│       ├── profile-users.md
│       ├── profile-drivers.md
│       └── settings.md
├── sidebars.js
├── docusaurus.config.js
└── package.json
```

Khi thêm/sửa API trong từng module, cập nhật file tương ứng trong `docs/modules/`.
