---
sidebar_position: 2
title: Profile Management
---

# Profile Management API

Quản lý thông tin profile người dùng: xem và cập nhật thông tin cá nhân.

---

## 1. `GET /api/v1/auth/user/profile`

Lấy thông tin profile của người dùng đang đăng nhập.

**Auth:** Bearer Token (bắt buộc)

**Request body:** Không có

**cURL:**

```bash
curl -X GET 'http://localhost:8080/api/v1/auth/user/profile' \
  -H 'Authorization: Bearer <access_token>'
```

**Response 200**

```json
{
  "status_code": 200,
  "message": "Success",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "full_name": "Nguyễn Văn A",
    "avatar_url": "https://example.com/avatar.jpg",
    "is_active": true,
    "phone": "+84901234567",
    "email": "user@example.com"
  }
}
```

**Response 401 - Unauthorized**

```json
{
  "status_code": 401,
  "message": "Unauthorized",
  "data": null
}
```

**Response 404 - User Not Found**

```json
{
  "status_code": 404,
  "message": "Người dùng không tồn tại",
  "data": null
}
```

---

## 2. `PUT /api/v1/auth/user/update-user-profile`

Cập nhật thông tin profile của người dùng.

**Auth:** Bearer Token (bắt buộc)

**Request body:**

```json
{
  "full_name": "Nguyễn Văn B",
  "avatar_url": "https://example.com/new-avatar.jpg"
}
```

| Field | Kiểu | Bắt buộc | Mô tả |
|-------|------|----------|-------|
| full_name | string | Có | Họ và tên mới |
| avatar_url | string | Có | URL ảnh đại diện mới |

**cURL:**

```bash
curl -X PUT 'http://localhost:8080/api/v1/auth/user/update-user-profile' \
  -H 'Authorization: Bearer <access_token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "full_name": "Nguyễn Văn B",
    "avatar_url": "https://example.com/new-avatar.jpg"
  }'
```

**Response 200**

```json
{
  "status_code": 200,
  "message": "Success",
  "data": {
    "user_message": "Cập nhật thông tin thành công",
    "user_profile": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "full_name": "Nguyễn Văn B",
      "avatar_url": "https://example.com/new-avatar.jpg",
      "is_active": true,
      "phone": "+84901234567",
      "email": "user@example.com"
    }
  }
}
```

**Response 401 - Unauthorized**

```json
{
  "status_code": 401,
  "message": "Unauthorized",
  "data": null
}
```

**Response 404 - User Not Found**

```json
{
  "status_code": 404,
  "message": "Người dùng không tồn tại",
  "data": null
}
```

**Response 422 - Validation Error**

```json
{
  "status_code": 422,
  "message": [
    "full_name: không được để trống",
    "avatar_url: không hợp lệ"
  ],
  "data": null
}
```
