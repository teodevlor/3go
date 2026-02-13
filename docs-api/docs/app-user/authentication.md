---
sidebar_position: 1
title: Authentication
---

# Authentication API

Quản lý xác thực người dùng: đăng ký, đăng nhập, kích hoạt tài khoản, làm mới token và đăng xuất.

---

## 1. `POST /api/v1/auth/user/register`

Đăng ký tài khoản người dùng mới. Sau khi đăng ký thành công, hệ thống sẽ gửi mã OTP về số điện thoại.

**Auth:** Không

**Request body:**

```json
{
  "phone": "+84901234567",
  "full_name": "Nguyễn Văn A",
  "password": "Pass@123"
}
```

| Field | Kiểu | Bắt buộc | Mô tả |
|-------|------|----------|-------|
| phone | string | Có | Số điện thoại (format: +84...) |
| full_name | string | Có | Họ và tên |
| password | string | Có | Mật khẩu (tối thiểu 8 ký tự, có chữ hoa, số và ký tự đặc biệt) |

**cURL:**

```bash
curl -X POST 'http://localhost:8080/api/v1/auth/user/register' \
  -H 'Content-Type: application/json' \
  -d '{
    "phone": "+84901234567",
    "full_name": "Nguyễn Văn A",
    "password": "Pass@123"
  }'
```

**Response 200**

```json
{
  "status_code": 200,
  "message": "Success",
  "data": {
    "user_message": "Đăng ký tài khoản thành công, vui lòng kiểm tra điện thoại để nhận mã OTP"
  }
}
```

**Response 422 - Validation Error**

```json
{
  "status_code": 422,
  "message": [
    "phone: không được để trống",
    "password: Mật khẩu phải có ít nhất 8 ký tự, bao gồm chữ hoa, chữ thường, số và ký tự đặc biệt"
  ],
  "data": null
}
```

**Response 500 - User Already Exists**

```json
{
  "status_code": 500,
  "message": "Số điện thoại đã được đăng ký",
  "data": null
}
```

---

## 2. `POST /api/v1/auth/user/active`

Kích hoạt tài khoản bằng mã OTP nhận được từ SMS.

**Auth:** Không

**Request body:**

```json
{
  "phone": "+84901234567",
  "code": "123456"
}
```

| Field | Kiểu | Bắt buộc | Mô tả |
|-------|------|----------|-------|
| phone | string | Có | Số điện thoại đã đăng ký |
| code | string | Có | Mã OTP (6 chữ số) |

**cURL:**

```bash
curl -X POST 'http://localhost:8080/api/v1/auth/user/active' \
  -H 'Content-Type: application/json' \
  -d '{
    "phone": "+84901234567",
    "code": "123456"
  }'
```

**Response 200**

```json
{
  "status_code": 200,
  "message": "Success",
  "data": {
    "user_message": "Kích hoạt tài khoản thành công"
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

**Response 400 - Already Active**

```json
{
  "status_code": 400,
  "message": "Tài khoản đã được kích hoạt",
  "data": null
}
```

---

## 3. `POST /api/v1/auth/user/login`

Đăng nhập vào hệ thống. Trả về access token và refresh token.

**Auth:** Không

**Request body:**

```json
{
  "phone": "+84901234567",
  "password": "Pass@123",
  "device": {
    "device_uid": "unique-device-id-123",
    "platform": "iOS",
    "device_name": "iPhone 13",
    "os_version": "15.0",
    "app_version": "1.0.0",
    "fcm_token": "fcm_token_here"
  }
}
```

| Field | Kiểu | Bắt buộc | Mô tả |
|-------|------|----------|-------|
| phone | string | Có | Số điện thoại |
| password | string | Có | Mật khẩu |
| device | object | Có | Thông tin thiết bị |
| device.device_uid | string | Có | ID unique của thiết bị |
| device.platform | string | Có | Nền tảng (iOS, Android, Web) |
| device.device_name | string | Có | Tên thiết bị |
| device.os_version | string | Có | Phiên bản OS |
| device.app_version | string | Có | Phiên bản app |
| device.fcm_token | string | Không | Token Firebase Cloud Messaging |

**cURL:**

```bash
curl -X POST 'http://localhost:8080/api/v1/auth/user/login' \
  -H 'Content-Type: application/json' \
  -d '{
    "phone": "+84901234567",
    "password": "Pass@123",
    "device": {
      "device_uid": "unique-device-id-123",
      "platform": "iOS",
      "device_name": "iPhone 13",
      "os_version": "15.0",
      "app_version": "1.0.0",
      "fcm_token": ""
    }
  }'
```

**Response 200**

```json
{
  "status_code": 200,
  "message": "Success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user_profile": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "full_name": "Nguyễn Văn A",
      "avatar_url": "https://example.com/avatar.jpg",
      "is_active": true,
      "phone": "+84901234567",
      "email": "user@example.com"
    }
  }
}
```

**Response 401 - Invalid Credentials**

```json
{
  "status_code": 401,
  "message": "Số điện thoại hoặc mật khẩu không đúng",
  "data": null
}
```

**Response 403 - Not Active**

```json
{
  "status_code": 403,
  "message": "Tài khoản chưa được kích hoạt. Vui lòng kiểm tra mã OTP",
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

## 4. `POST /api/v1/auth/user/refresh-token`

Làm mới access token khi token cũ hết hạn.

**Auth:** Không (sử dụng refresh token)

**Request body:**

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

| Field | Kiểu | Bắt buộc | Mô tả |
|-------|------|----------|-------|
| refresh_token | string | Có | Refresh token nhận từ login |

**cURL:**

```bash
curl -X POST 'http://localhost:8080/api/v1/auth/user/refresh-token' \
  -H 'Content-Type: application/json' \
  -d '{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }'
```

**Response 200**

```json
{
  "status_code": 200,
  "message": "Success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Response 401 - Invalid Token**

```json
{
  "status_code": 401,
  "message": "Refresh token không hợp lệ hoặc đã hết hạn",
  "data": null
}
```

---

## 5. `POST /api/v1/auth/user/logout`

Đăng xuất khỏi hệ thống. Hủy access token và refresh token.

**Auth:** Bearer Token (bắt buộc)

**Request body:** Không có

**cURL:**

```bash
curl -X POST 'http://localhost:8080/api/v1/auth/user/logout' \
  -H 'Authorization: Bearer <access_token>'
```

**Response 200**

```json
{
  "status_code": 200,
  "message": "Success",
  "data": {
    "user_message": "Đăng xuất thành công"
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
