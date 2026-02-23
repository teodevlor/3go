# Storage Upload API

Endpoint upload **một hoặc nhiều file**, luôn dùng chung một API.

- **Method:** `POST`
- **URL:** `/api/v1/storage/upload`
- **Content-Type:** `multipart/form-data`

## Form fields

| Key        | Bắt buộc | Mô tả |
|-----------|----------|--------|
| `file` hoặc `file[]` | Có | File cần upload. Gửi 1 part = 1 file, nhiều part = nhiều file (tối đa 20). |
| `path`    | Không | Thư mục lưu (vd: `avatars`, `avatar/user`). Key lưu sẽ dạng: `upload/{path}/{YYYYMMDD}/{uuid}.ext` |
| `visibility` | Không | `public` (mặc định) hoặc `private` |

---

## Postman – Upload nhiều file

1. Chọn **POST**, URL: `http://localhost:8080/api/v1/storage/upload`
2. **Không** set header `Content-Type` thủ công — để Postman tự gửi `multipart/form-data; boundary=...` khi chọn form-data.
3. Vào tab **Body** → chọn **form-data**.
3. Thêm các row:

   | Key          | Type | Value |
   |--------------|------|--------|
   | `path`       | Text | `avatars` (hoặc path bạn muốn) |
   | `visibility` | Text | `public` |
   | `file[]`     | File | Chọn file 1 (click "Select Files" chọn 1 hoặc nhiều) |
   | `file[]`     | File | Chọn file 2 (thêm row mới cùng key `file[]`) |
   | `file[]`     | File | Chọn file 3 … |

   **Lưu ý:** Một số phiên bản Postman cho phép chọn **nhiều file trong một ô** cho cùng key `file[]`; nếu không thì thêm nhiều row cùng key `file[]`, mỗi row một file.

4. Gửi request.

**Cách khác (dùng key `file`):** Có thể dùng key `file` thay cho `file[]`, mỗi row một file:

| Key    | Type | Value   |
|--------|------|---------|
| `path` | Text | `avatars` |
| `file` | File | file 1  |
| `file` | File | file 2  |

---

## Postman – Upload một file

Cùng endpoint, chỉ cần gửi **một** part file:

| Key          | Type | Value   |
|--------------|------|---------|
| `path`       | Text | `avatars` |
| `visibility` | Text | `public` |
| `file` hoặc `file[]` | File | Chọn 1 file |

---

## cURL – Một file

```bash
curl -X POST 'http://localhost:8080/api/v1/storage/upload' \
  -F 'path=avatars' \
  -F 'visibility=public' \
  -F 'file=@/path/to/hat-oc-cho.png'
```

## cURL – Nhiều file (multi file)

**Quan trọng:** Mỗi file phải có **một option `-F` riêng**. Không gộp nhiều file vào một `-F`.

**Đúng** — 4 file, 4 lần `-F 'file=@...'`:

```bash
curl -X POST 'http://localhost:8080/api/v1/storage/upload' \
  -F 'path=avatars' \
  -F 'visibility=public' \
  -F 'file=@/Users/kegiaumatvideptraigmail.com/hatvang.vn/images/hat-hanh-nhan.png' \
  -F 'file=@/Users/kegiaumatvideptraigmail.com/hatvang.vn/images/hat-hanh-nhan.png' \
  -F 'file=@/Users/kegiaumatvideptraigmail.com/hatvang.vn/images/hat-ho-dao.png' \
  -F 'file=@/Users/kegiaumatvideptraigmail.com/hatvang.vn/images/hat-macca.png'
```

**Sai** — thiếu `-F` trước mỗi file (curl sẽ không gửi đủ part):

```bash
# SAI: chỉ có 2 form field được gửi (path, visibility), các file không đi kèm -F nên không gửi đúng
curl -X POST '...' --form 'path=avatars' --form 'visibility=public' \
  'file=@"/path/1"' 'file=@"/path/2"' 'file=@"/path/3"'
```

Dùng key `file[]` cũng được, mỗi file một `-F`:

```bash
curl -X POST 'http://localhost:8080/api/v1/storage/upload' \
  -F 'path=avatars' \
  -F 'visibility=public' \
  -F 'file[]=@/path/to/image1.png' \
  -F 'file[]=@/path/to/image2.png'
```

---

## Response – luôn là mảng `uploads`

Dù gửi 1 hay nhiều file, response luôn có dạng `data.uploads` là **mảng** (1 file → 1 phần tử, nhiều file → nhiều phần tử). Nếu thấy `data` chỉ là một object (không có `uploads`) hoặc chỉ 1 phần tử trong khi gửi nhiều file → kiểm tra lại cách gửi (curl phải có **-F** trước từng file).

### 1 file

```json
{
  "status": true,
  "code": 200,
  "message": "Success",
  "data": {
    "uploads": [
      {
        "key": "upload/avatars/20260222/xxx.png",
        "path": "avatars",
        "size": 1308675,
        "original_filename": "hat-dieu.png",
        "visibility": "public",
        "bucket": "public-assets"
      }
    ]
  }
}
```

### Nhiều file

```json
{
  "status": true,
  "code": 200,
  "message": "Success",
  "data": {
    "uploads": [
    {
      "key": "upload/avatars/20250222/abc-123.png",
      "path": "avatars",
      "size": 12345,
      "original_filename": "hat-oc-cho.png",
      "visibility": "public",
      "bucket": "your-public-bucket"
    },
    {
      "key": "upload/avatars/20250222/def-456.png",
      "path": "avatars",
      "size": 67890,
      "original_filename": "hat-ho-dao.png",
      "visibility": "public",
      "bucket": "your-public-bucket"
    }
  ]
}
```

**cURL đúng để nhận đủ 2 (hoặc n) file trong `uploads`:** mỗi file một `-F`:

```bash
curl -X POST 'http://localhost:8080/api/v1/storage/upload' \
  -F 'path=avatars' \
  -F 'visibility=public' \
  -F 'file=@/Users/.../hat-dieu.png' \
  -F 'file=@/Users/.../hat-oc-cho.png'
```

---

## Lỗi thường gặp

- **"file is required"** / **"file is required (form-data key: ...)"**  
  - Chưa gửi part file, hoặc key form sai. Trong Postman: dùng key **`file`** hoặc **`file[]`**, type cột value chọn **File** rồi chọn file.  
  - **Quan trọng:** Không set header `Content-Type: multipart/form-data` tay — phải để Postman tự set (kèm `boundary`), nếu không server không parse được file.
- **"too many files, max 20"** → Vượt quá 20 file trong một request.
- **"multipart form required"** → Request không phải `multipart/form-data` (vd đang gửi JSON hoặc thiếu boundary).
- **Gửi 2 file nhưng `data.uploads` chỉ có 1 phần tử** → Trong curl phải dùng **-F** trước từng file (ví dụ: `-F 'file=@path1'` và `-F 'file=@path2'`). Thiếu -F thì chỉ part đầu được gửi.
