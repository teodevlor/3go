---
sidebar_position: 1
slug: /
---

# Tài liệu Go Structure API

Tài liệu API cho dự án **Go Structure** - một hệ thống backend được xây dựng bằng Go với kiến trúc Clean Architecture.

<div class="api-docs-config">

**Cấu hình biến (giống Postman)**

Đặt **Base URL** và **Access token**, bấm **Lưu**. **Toàn bộ khối cURL** trên trang sẽ **cập nhật theo** (URL và token hiển thị đúng giá trị đã set). Bấm nút **Copy** có sẵn phía trên mỗi khối → copy ra **dùng được ngay**, kèm toast **"Đã sao chép!"**. Chưa set thì vẫn dùng ví dụ mặc định.

| Biến | Ví dụ mặc định | Mô tả |
|------|----------------|-------|
| Base URL | `http://localhost:8080` | Địa chỉ server (không gồm `/api/v1`) |
| Access token | `<access_token>` | JWT, dùng trong header `Authorization: Bearer ...` |

<div class="api-docs-config__form">

<label><strong>Base URL:</strong> <input id="api-base-url" type="text" placeholder="http://localhost:8080" autocomplete="off" class="api-docs-config__input" /></label>

<label><strong>Access token:</strong> <input id="api-access-token" type="text" placeholder="(tùy chọn)" autocomplete="off" class="api-docs-config__input" /></label>

<button id="api-save-config" class="api-docs-config__btn">Lưu</button>
<span id="api-config-status" class="api-docs-config__status"></span>

</div>

</div>

---

## Base URL

```
/api/v1
```

Ví dụ: `http://localhost:8080/api/v1/auth/user/login`

---

## Xác thực (JWT)

Các endpoint yêu cầu đăng nhập dùng header:

```http
Authorization: Bearer <access_token>
```

Sau khi login/refresh, dùng `accessToken` trong response. Khi `401` hoặc token hết hạn: gọi **refresh-token** hoặc chuyển người dùng về màn hình đăng nhập.

---

## Định dạng Response

Mọi API trả về JSON thống nhất với cấu trúc:

```json
{
  "status_code": 200,
  "message": "Success",
  "data": { ... }
}
```

### Response thành công

**HTTP:** `200`, `201`

```json
{
  "status_code": 200,
  "message": "Success",
  "data": {
    "result": { ... }
  }
}
```

### Response thất bại

**HTTP:** `4xx`, `5xx`

```json
{
  "status_code": 400,
  "message": "Mô tả lỗi",
  "data": null
}
```

### Response Validation Error

**HTTP:** `422`

```json
{
  "status_code": 422,
  "message": ["Lỗi validation 1", "Lỗi validation 2"],
  "data": null
}
```

---

## Các Module API

Tài liệu được tổ chức theo các module chính:

### 1. **Website System**
Quản lý hệ thống website bao gồm:
- **Zones**: Quản lý các khu vực/vùng
- **Sidebars**: Quản lý menu sidebar

### 2. **App User**
Xác thực và quản lý người dùng bao gồm:
- Đăng ký, đăng nhập
- Xác thực OTP
- Quản lý profile
- Quên/đổi mật khẩu

---

Chọn module ở sidebar để xem chi tiết từng endpoint.
