# Football API Documentation for Frontend

Base URL: `http://localhost:3000/api/v1`  
LAN Base URL: `http://10.5.2.103:3000/api/v1`  
Content-Type: `application/json`

---

### 1. Health Check
ตรวจสอบสถานะการทำงานของ Server และ Database

- **Endpoint:** `GET /health`
- **Request:** ไม่มี
- **Response (200 OK):**
```json
{
  "status": "ok"
}
```

---

### 2. Get All Leagues
ดึงรายชื่อลีกและฟุตบอลถ้วยทั้งหมดในระบบ เรียงตามชื่อลีก

- **Endpoint:** `GET /leagues`
- **Request:** ไม่มี
- **Response (200 OK):**
```json
[
  {
    "id": "0aab480b-9e10-4fea-a6fc-884160e684d5",
    "name": "Thai League 1",
    "short_name": "T1",
    "country": "Thailand",
    "logo_url": "https://media.api-sports.io/football/leagues/296.png",
    "type": "league",
    "sort_order": 0,
    "created_at": "2026-11-26T10:00:00Z",
    "updated_at": "2026-11-26T10:00:00Z"
  }
]
```

---

### 3. Create League
สร้างข้อมูลลีกใหม่เข้าสู่ระบบ

- **Endpoint:** `POST /leagues`
- **Request Body:**
```json
{
  "name": "Premier League",
  "short_name": "EPL",
  "country": "England",
  "logo_url": "https://example.com/logo.png"
}
```
| ฟิลด์ | ชนิดข้อมูล | จำเป็น | รายละเอียด |
|:---|:---:|:---:|:---|
| `name` | string | ใช่ | ชื่อของลีก |
| `short_name` | string | ไม่ใช่ | ตัวย่อของลีก |
| `country` | string | ไม่ใช่ | ประเทศ |
| `logo_url` | string | ไม่ใช่ | URL รูปโลโก้ |

- **Response (201 Created):**
```json
{
  "id": "e8432a51-3b7c-4866-8f3e-289cf0b31e9a",
  "name": "Premier League",
  "short_name": "EPL",
  "country": "England",
  "logo_url": "https://example.com/logo.png",
  "type": "league",
  "sort_order": 0,
  "created_at": "2026-11-26T11:20:00Z",
  "updated_at": "2026-11-26T11:20:00Z"
}
```
- **Error (400 Bad Request):**
```json
{
  "error": "name is required"
}
```

---

### 4. Get League by ID
ดึงข้อมูลลีกเฉพาะรายการตาม UUID

- **Endpoint:** `GET /leagues/:id`
- **URL Parameter:** `:id` (UUID ของลีก)
- **Request:** ไม่มี
- **Response (200 OK):**
```json
{
  "id": "0aab480b-9e10-4fea-a6fc-884160e684d5",
  "name": "Thai League 1",
  "short_name": "T1",
  "country": "Thailand",
  "logo_url": "https://media.api-sports.io/football/leagues/296.png",
  "type": "league",
  "sort_order": 0,
  "created_at": "2026-11-26T10:00:00Z",
  "updated_at": "2026-11-26T10:00:00Z"
}
```
- **Error (404 Not Found):**
```json
{
  "error": "league not found"
}
```

---

### 5. Update League
แก้ไขข้อมูลลีกตาม UUID

- **Endpoint:** `PUT /leagues/:id`
- **URL Parameter:** `:id` (UUID ของลีก)
- **Request Body:**
```json
{
  "name": "Thai League 1 (Hilux Revo)",
  "short_name": "T1",
  "country": "Thailand",
  "logo_url": "https://example.com/new_logo.png"
}
```
- **Response (200 OK):**
```json
{
  "id": "0aab480b-9e10-4fea-a6fc-884160e684d5",
  "name": "Thai League 1 (Hilux Revo)",
  "short_name": "T1",
  "country": "Thailand",
  "logo_url": "https://example.com/new_logo.png",
  "type": "league",
  "sort_order": 0,
  "created_at": "2026-11-26T10:00:00Z",
  "updated_at": "2026-11-26T11:30:00Z"
}
```

---

### 6. Delete League
ลบลีกออกจากระบบ (Soft Delete)

- **Endpoint:** `DELETE /leagues/:id`
- **URL Parameter:** `:id` (UUID ของลีก)
- **Request:** ไม่มี
- **Response (204 No Content):** ไม่มี Body
