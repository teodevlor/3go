---
sidebar_position: 4
title: OTP Management
---

# OTP Management API

Quản lý OTP: gửi lại mã OTP khi cần thiết.

---

## 1. `POST /api/v1/otp/resend`

Gửi lại mã OTP cho các trường hợp: đăng ký tài khoản hoặc quên mật khẩu.

**Auth:** Không

**Request body:**

```json
{
  "phone": "+84901234567",
  "otp_type": "register"
}
```

| Field | Kiểu | Bắt buộc | Mô tả |
|-------|------|----------|-------|
| phone | string | Có | Số điện thoại |
| otp_type | string | Có | Loại OTP: "register" hoặc "forgot_password" |

**cURL:**

```bash
curl -X POST 'http://localhost:8080/api/v1/otp/resend' \
  -H 'Content-Type: application/json' \
  -d '{
    "phone": "+84901234567",
    "otp_type": "register"
  }'
```

**Response 200**

```json
{
  "status_code": 200,
  "message": "Success",
  "data": {
    "user_message": "Mã OTP đã được gửi lại thành công"
  }
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

**Response 404 - User Not Found**

```json
{
  "status_code": 404,
  "message": "Số điện thoại không tồn tại trong hệ thống",
  "data": null
}
```

**Response 422 - Validation Error**

```json
{
  "status_code": 422,
  "message": [
    "phone: không được để trống",
    "otp_type: phải là 'register' hoặc 'forgot_password'"
  ],
  "data": null
}
```

---

## Lưu ý

### Giới hạn gửi OTP

- Mỗi số điện thoại chỉ được gửi OTP tối đa **5 lần** trong vòng **1 giờ**
- Thời gian chờ giữa 2 lần gửi: **60 giây**
- Mã OTP có hiệu lực trong **5 phút**

### OTP Type

- **register**: Dùng cho đăng ký tài khoản mới
- **forgot_password**: Dùng cho quên mật khẩu

### Xử lý lỗi

Khi gặp lỗi 429 (Too Many Requests), ứng dụng nên:
1. Hiển thị thông báo lỗi cho người dùng
2. Disable nút "Gửi lại OTP" trong thời gian chờ
3. Hiển thị đếm ngược thời gian còn lại
