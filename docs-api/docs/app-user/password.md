---
sidebar_position: 3
title: Password Management
---

# Password Management API

Quản lý mật khẩu: quên mật khẩu và đặt lại mật khẩu thông qua OTP.

---

## 1. `POST /api/v1/auth/user/forgot-password`

Yêu cầu đặt lại mật khẩu. Hệ thống sẽ gửi mã OTP về số điện thoại.

**Auth:** Không

**Request body:**

```json
{
  "phone": "+84901234567"
}
```

| Field | Kiểu | Bắt buộc | Mô tả |
|-------|------|----------|-------|
| phone | string | Có | Số điện thoại đã đăng ký |

**cURL:**

```bash
curl -X POST 'http://localhost:8080/api/v1/auth/user/forgot-password' \
  -H 'Content-Type: application/json' \
  -d '{
    "phone": "+84901234567"
  }'
```

**Response 200**

```json
{
  "status_code": 200,
  "message": "Success",
  "data": {
    "user_message": "Vui lòng kiểm tra điện thoại để nhận mã OTP đặt lại mật khẩu"
  }
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

**Response 429 - Too Many Requests**

```json
{
  "status_code": 429,
  "message": "Bạn đã yêu cầu quá nhiều lần. Vui lòng thử lại sau 60 giây",
  "data": null
}
```

**Response 422 - Validation Error**

```json
{
  "status_code": 422,
  "message": [
    "phone: không được để trống",
    "phone: không đúng định dạng"
  ],
  "data": null
}
```

---

## 2. `POST /api/v1/auth/user/reset-password`

Đặt lại mật khẩu bằng mã OTP nhận được từ forgot-password.

**Auth:** Không

**Request body:**

```json
{
  "phone": "+84901234567",
  "code": "123456",
  "new_password": "NewPass@123",
  "confirm_password": "NewPass@123"
}
```

| Field | Kiểu | Bắt buộc | Mô tả |
|-------|------|----------|-------|
| phone | string | Có | Số điện thoại |
| code | string | Có | Mã OTP (6 chữ số) |
| new_password | string | Có | Mật khẩu mới (tối thiểu 8 ký tự, có chữ hoa, số và ký tự đặc biệt) |
| confirm_password | string | Có | Xác nhận mật khẩu mới (phải khớp với new_password) |

**cURL:**

```bash
curl -X POST 'http://localhost:8080/api/v1/auth/user/reset-password' \
  -H 'Content-Type: application/json' \
  -d '{
    "phone": "+84901234567",
    "code": "123456",
    "new_password": "NewPass@123",
    "confirm_password": "NewPass@123"
  }'
```

**Response 200**

```json
{
  "status_code": 200,
  "message": "Success",
  "data": {
    "user_message": "Đặt lại mật khẩu thành công"
  }
}
```

**Response 400 - Invalid OTP**

```json
{
  "status_code": 400,
  "message": "Mã OTP không hợp lệ hoặc đã hết hạn",
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
    "new_password: Mật khẩu phải có ít nhất 8 ký tự, bao gồm chữ hoa, chữ thường, số và ký tự đặc biệt",
    "confirm_password: Mật khẩu xác nhận không khớp"
  ],
  "data": null
}
```
