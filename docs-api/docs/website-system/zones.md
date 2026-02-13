---
sidebar_position: 1
title: Zones API
---

# Zones API

Quản lý các khu vực/vùng địa lý trong hệ thống. Mỗi zone được định nghĩa bằng một polygon GeoJSON và có hệ số giá riêng.

---

## 1. `POST /api/v1/system/zones`

Tạo zone mới với polygon GeoJSON.

**Auth:** Có thể yêu cầu (tùy cấu hình)

**Request body:**

```json
{
  "code": "ZONE_HN_1",
  "name": "Khu vực Hà Nội 1",
  "price_multiplier": 1.5,
  "is_active": true,
  "polygon": {
    "type": "Polygon",
    "coordinates": [
      [
        [105.8342, 21.0278],
        [105.8542, 21.0278],
        [105.8542, 21.0478],
        [105.8342, 21.0478],
        [105.8342, 21.0278]
      ]
    ]
  }
}
```

| Field | Kiểu | Bắt buộc | Mô tả |
|-------|------|----------|-------|
| code | string | Có | Mã zone (tối đa 100 ký tự) |
| name | string | Có | Tên zone (tối đa 255 ký tự) |
| price_multiplier | number | Có | Hệ số giá (≥ 0) |
| is_active | boolean | Không | Trạng thái hoạt động (mặc định: false) |
| polygon | GeoJSONPolygon | Có | Vùng địa lý theo chuẩn GeoJSON |
| polygon.type | string | Có | Phải là "Polygon" |
| polygon.coordinates | array | Có | Mảng các ring, mỗi ring là mảng các điểm [lng, lat] |

**Lưu ý về Polygon:**
- Polygon phải có ít nhất 1 ring
- Mỗi ring phải có ít nhất 4 điểm
- Ring phải đóng kín (điểm đầu = điểm cuối)
- Mỗi điểm có format [longitude, latitude]
- Longitude: [-180, 180], Latitude: [-90, 90]

**cURL:**

```bash
curl -X POST 'http://localhost:8080/api/v1/system/zones' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer <access_token>' \
  -d '{
    "code": "ZONE_HN_1",
    "name": "Khu vực Hà Nội 1",
    "price_multiplier": 1.5,
    "is_active": true,
    "polygon": {
      "type": "Polygon",
      "coordinates": [
        [
          [105.8342, 21.0278],
          [105.8542, 21.0278],
          [105.8542, 21.0478],
          [105.8342, 21.0478],
          [105.8342, 21.0278]
        ]
      ]
    }
  }'
```

**Response 200**

```json
{
  "status_code": 200,
  "message": "Success",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "code": "ZONE_HN_1",
    "name": "Khu vực Hà Nội 1",
    "price_multiplier": 1.5,
    "is_active": true,
    "polygon": "{\"type\":\"Polygon\",\"coordinates\":[[[105.8342,21.0278],[105.8542,21.0278],[105.8542,21.0478],[105.8342,21.0478],[105.8342,21.0278]]]}"
  }
}
```

**Response 422 - Validation Error**

```json
{
  "status_code": 422,
  "message": [
    "code: không được để trống",
    "polygon: coordinates must have at least 4 points"
  ],
  "data": null
}
```

---

## 2. `GET /api/v1/system/zones`

Lấy danh sách zones với phân trang.

**Auth:** Có thể yêu cầu (tùy cấu hình)

**Query Parameters:**

| Param | Kiểu | Mặc định | Mô tả |
|-------|------|----------|-------|
| page | number | 1 | Trang hiện tại |
| limit | number | 10 | Số items mỗi trang |
| search | string | - | Tìm kiếm theo code hoặc name |
| is_active | boolean | - | Lọc theo trạng thái |

**cURL:**

```bash
curl -X GET 'http://localhost:8080/api/v1/system/zones?page=1&limit=10' \
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
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "code": "ZONE_HN_1",
        "name": "Khu vực Hà Nội 1",
        "price_multiplier": 1.5,
        "is_active": true,
        "polygon": "{\"type\":\"Polygon\",\"coordinates\":[[[105.8342,21.0278],...]]}"
      }
    ],
    "pagination": {
      "total": 50,
      "page": 1,
      "limit": 10,
      "total_pages": 5
    }
  }
}
```

---

## 3. `GET /api/v1/system/zones/:id`

Lấy chi tiết một zone theo ID.

**Auth:** Có thể yêu cầu (tùy cấu hình)

**Path Parameters:**

| Param | Kiểu | Mô tả |
|-------|------|-------|
| id | string (UUID) | ID của zone |

**cURL:**

```bash
curl -X GET 'http://localhost:8080/api/v1/system/zones/550e8400-e29b-41d4-a716-446655440000' \
  -H 'Authorization: Bearer <access_token>'
```

**Response 200**

```json
{
  "status_code": 200,
  "message": "Success",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "code": "ZONE_HN_1",
    "name": "Khu vực Hà Nội 1",
    "price_multiplier": 1.5,
    "is_active": true,
    "polygon": "{\"type\":\"Polygon\",\"coordinates\":[[[105.8342,21.0278],...]]}"
  }
}
```

**Response 404**

```json
{
  "status_code": 404,
  "message": "Zone not found",
  "data": null
}
```

---

## 4. `PUT /api/v1/system/zones/:id`

Cập nhật thông tin zone.

**Auth:** Có thể yêu cầu (tùy cấu hình)

**Path Parameters:**

| Param | Kiểu | Mô tả |
|-------|------|-------|
| id | string (UUID) | ID của zone |

**Request body:**

```json
{
  "code": "ZONE_HN_1_UPDATED",
  "name": "Khu vực Hà Nội 1 (Cập nhật)",
  "price_multiplier": 2.0,
  "is_active": false,
  "polygon": {
    "type": "Polygon",
    "coordinates": [
      [
        [105.8342, 21.0278],
        [105.8542, 21.0278],
        [105.8542, 21.0478],
        [105.8342, 21.0478],
        [105.8342, 21.0278]
      ]
    ]
  }
}
```

**cURL:**

```bash
curl -X PUT 'http://localhost:8080/api/v1/system/zones/550e8400-e29b-41d4-a716-446655440000' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer <access_token>' \
  -d '{
    "code": "ZONE_HN_1_UPDATED",
    "name": "Khu vực Hà Nội 1 (Cập nhật)",
    "price_multiplier": 2.0,
    "is_active": false,
    "polygon": {
      "type": "Polygon",
      "coordinates": [...]
    }
  }'
```

**Response 200**

```json
{
  "status_code": 200,
  "message": "Success",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "code": "ZONE_HN_1_UPDATED",
    "name": "Khu vực Hà Nội 1 (Cập nhật)",
    "price_multiplier": 2.0,
    "is_active": false,
    "polygon": "{...}"
  }
}
```

**Response 404**

```json
{
  "status_code": 404,
  "message": "Zone not found",
  "data": null
}
```

---

## 5. `DELETE /api/v1/system/zones/:id`

Xóa một zone.

**Auth:** Có thể yêu cầu (tùy cấu hình)

**Path Parameters:**

| Param | Kiểu | Mô tả |
|-------|------|-------|
| id | string (UUID) | ID của zone |

**cURL:**

```bash
curl -X DELETE 'http://localhost:8080/api/v1/system/zones/550e8400-e29b-41d4-a716-446655440000' \
  -H 'Authorization: Bearer <access_token>'
```

**Response 200**

```json
{
  "status_code": 200,
  "message": "Zone deleted successfully",
  "data": null
}
```

**Response 404**

```json
{
  "status_code": 404,
  "message": "Zone not found",
  "data": null
}
```
