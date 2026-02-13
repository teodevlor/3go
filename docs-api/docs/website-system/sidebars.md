---
sidebar_position: 2
title: Sidebars API
---

# Sidebars API

Quản lý cấu hình sidebar/menu động cho các context khác nhau (system, app_user, app_driver...). Mỗi sidebar chứa cấu trúc menu dạng tree với các item có thể có children.

---

## 1. `POST /api/v1/system/sidebars`

Tạo cấu hình sidebar mới.

**Auth:** Có thể yêu cầu (tùy cấu hình)

**Request body:**

```json
{
  "context": "system",
  "version": "1.0.0",
  "generated_at": "2026-02-11T10:00:00Z",
  "items": [
    {
      "id": "dashboard",
      "label": "Dashboard",
      "icon": "dashboard-icon",
      "path": "/dashboard",
      "order": 1,
      "visible": true,
      "permission_required": ["view_dashboard"],
      "badge": {
        "type": "info",
        "value": "New"
      },
      "feature_flag": "dashboard_v2"
    },
    {
      "id": "users",
      "label": "Quản lý người dùng",
      "icon": "users-icon",
      "order": 2,
      "visible": true,
      "children": [
        {
          "id": "users-list",
          "label": "Danh sách",
          "path": "/users/list",
          "order": 1,
          "visible": true
        },
        {
          "id": "users-create",
          "label": "Tạo mới",
          "path": "/users/create",
          "order": 2,
          "visible": true
        }
      ]
    }
  ]
}
```

| Field | Kiểu | Bắt buộc | Mô tả |
|-------|------|----------|-------|
| context | string | Có | Context của sidebar (system, app_user...), max 100 ký tự |
| version | string | Không | Phiên bản cấu hình, max 50 ký tự |
| generated_at | datetime | Không | Thời điểm tạo cấu hình |
| items | array | Có | Danh sách items trong sidebar |
| items[].id | string | Có | ID unique của item |
| items[].label | string | Có | Nhãn hiển thị |
| items[].icon | string | Không | Tên icon |
| items[].path | string | Không | Đường dẫn routing |
| items[].order | number | Không | Thứ tự sắp xếp |
| items[].visible | boolean | Không | Hiển thị hay ẩn |
| items[].children | array | Không | Menu con (cấu trúc giống parent) |
| items[].permission_required | array | Không | Danh sách permissions cần có |
| items[].badge | object | Không | Badge hiển thị (type, value) |
| items[].feature_flag | string | Không | Feature flag để bật/tắt |

**cURL:**

```bash
curl -X POST 'http://localhost:8080/api/v1/system/sidebars' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer <access_token>' \
  -d '{
    "context": "system",
    "version": "1.0.0",
    "items": [...]
  }'
```

**Response 200**

```json
{
  "status_code": 200,
  "message": "Success",
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440000",
    "context": "system",
    "version": "1.0.0",
    "generated_at": "2026-02-11T10:00:00Z",
    "items": [...],
    "created_at": "2026-02-11T10:00:00Z",
    "updated_at": "2026-02-11T10:00:00Z"
  }
}
```

**Response 422 - Validation Error**

```json
{
  "status_code": 422,
  "message": [
    "context: không được để trống",
    "items: phải có ít nhất 1 item"
  ],
  "data": null
}
```

---

## 2. `GET /api/v1/system/sidebars`

Lấy danh sách sidebars với phân trang.

**Auth:** Có thể yêu cầu (tùy cấu hình)

**Query Parameters:**

| Param | Kiểu | Mặc định | Mô tả |
|-------|------|----------|-------|
| page | number | 1 | Trang hiện tại |
| limit | number | 10 | Số items mỗi trang |
| context | string | - | Lọc theo context |
| version | string | - | Lọc theo version |

**cURL:**

```bash
curl -X GET 'http://localhost:8080/api/v1/system/sidebars?page=1&limit=10&context=system' \
  -H 'Authorization: Bearer <access_token>'
```

**Response 200**

```json
{
  "status_code": 200,
  "message": "Success",
  "data": {
    "items": [
      {
        "id": "660e8400-e29b-41d4-a716-446655440000",
        "context": "system",
        "version": "1.0.0",
        "generated_at": "2026-02-11T10:00:00Z",
        "items": [...],
        "created_at": "2026-02-11T10:00:00Z",
        "updated_at": "2026-02-11T10:00:00Z"
      }
    ],
    "pagination": {
      "total": 10,
      "page": 1,
      "limit": 10,
      "total_pages": 1
    }
  }
}
```

---

## 3. `GET /api/v1/system/sidebars/:id`

Lấy chi tiết một sidebar theo ID.

**Auth:** Có thể yêu cầu (tùy cấu hình)

**Path Parameters:**

| Param | Kiểu | Mô tả |
|-------|------|-------|
| id | string (UUID) | ID của sidebar |

**cURL:**

```bash
curl -X GET 'http://localhost:8080/api/v1/system/sidebars/660e8400-e29b-41d4-a716-446655440000' \
  -H 'Authorization: Bearer <access_token>'
```

**Response 200**

```json
{
  "status_code": 200,
  "message": "Success",
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440000",
    "context": "system",
    "version": "1.0.0",
    "generated_at": "2026-02-11T10:00:00Z",
    "items": [
      {
        "id": "dashboard",
        "label": "Dashboard",
        "icon": "dashboard-icon",
        "path": "/dashboard",
        "order": 1,
        "visible": true,
        "permission_required": ["view_dashboard"],
        "badge": {
          "type": "info",
          "value": "New"
        }
      }
    ],
    "created_at": "2026-02-11T10:00:00Z",
    "updated_at": "2026-02-11T10:00:00Z"
  }
}
```

**Response 404**

```json
{
  "status_code": 404,
  "message": "Sidebar not found",
  "data": null
}
```

---

## 4. `PUT /api/v1/system/sidebars/:id`

Cập nhật cấu hình sidebar.

**Auth:** Có thể yêu cầu (tùy cấu hình)

**Path Parameters:**

| Param | Kiểu | Mô tả |
|-------|------|-------|
| id | string (UUID) | ID của sidebar |

**Request body:** (Giống CREATE)

```json
{
  "context": "system",
  "version": "1.0.1",
  "generated_at": "2026-02-11T11:00:00Z",
  "items": [...]
}
```

**cURL:**

```bash
curl -X PUT 'http://localhost:8080/api/v1/system/sidebars/660e8400-e29b-41d4-a716-446655440000' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer <access_token>' \
  -d '{
    "context": "system",
    "version": "1.0.1",
    "items": [...]
  }'
```

**Response 200**

```json
{
  "status_code": 200,
  "message": "Success",
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440000",
    "context": "system",
    "version": "1.0.1",
    "items": [...],
    "created_at": "2026-02-11T10:00:00Z",
    "updated_at": "2026-02-11T11:00:00Z"
  }
}
```

**Response 404**

```json
{
  "status_code": 404,
  "message": "Sidebar not found",
  "data": null
}
```

---

## 5. `DELETE /api/v1/system/sidebars/:id`

Xóa một sidebar.

**Auth:** Có thể yêu cầu (tùy cấu hình)

**Path Parameters:**

| Param | Kiểu | Mô tả |
|-------|------|-------|
| id | string (UUID) | ID của sidebar |

**cURL:**

```bash
curl -X DELETE 'http://localhost:8080/api/v1/system/sidebars/660e8400-e29b-41d4-a716-446655440000' \
  -H 'Authorization: Bearer <access_token>'
```

**Response 200**

```json
{
  "status_code": 200,
  "message": "Sidebar deleted successfully",
  "data": null
}
```

**Response 404**

```json
{
  "status_code": 404,
  "message": "Sidebar not found",
  "data": null
}
```
