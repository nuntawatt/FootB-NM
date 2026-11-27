# Football API Documentation for Frontend

เอกสารประกอบการใช้งาน RESTful API สำหรับระบบฟุตบอลไทยและข้อมูลการแข่งขัน  
เหมาะสำหรับทีมนักพัฒนา Frontend (Web / Mobile App)

* **Base URL (Local):** `http://localhost:3000/api/v1`
* **Base URL (LAN / Device Test):** `http://10.5.2.103:3000/api/v1`
* **Content-Type:** `application/json` (หรือ `multipart/form-data` เมื่ออัปโหลดไฟล์)
* **Authentication Header:** `Authorization: Bearer <JWT_TOKEN>`

> สำหรับขั้นตอนการติดตั้ง Dependencies, การรัน Database Migration, และคำสั่งรันระบบ โปรดดูที่ [RUN.md](file:///Users/morgorn/Desktop/project/football/Backend/RUN.md)

---

## สรุป Endpoints ทั้งหมดในระบบ (API Overview)

| หมวดหมู่ | Method | Endpoint | สิทธิ์เข้าถึง | คำอธิบาย / พารามิเตอร์ |
|:---|:---:|:---|:---:|:---|
| **1. Auth & Profile** | `POST` | `/register` | Public | สมัครสมาชิกใหม่ (ส่งคืน Token และ User Object ทันที) |
| | `POST` | `/login` | Public | เข้าสู่ระบบด้วย username หรือ email |
| | `GET` | `/profile` | Bearer Token | ดูข้อมูลโปรไฟล์ของผู้ใช้ปัจจุบัน |
| | `PUT` | `/profile` | Bearer Token | แก้ไขข้อมูลโปรไฟล์ (รองรับทั้ง JSON และอัปโหลดไฟล์รูปภาพ) |
| | `POST` | `/profile/avatar` | Bearer Token | อัปโหลดรูปโปรไฟล์ขึ้น MinIO / S3 Storage โดยเฉพาะ |
| **2. Seasons** | `GET` | `/seasons` | Public | รายการปี/ฤดูกาลทั้งหมด (2022-2026) พร้อมจำนวนแมตช์ |
| **3. Leagues** | `GET` | `/leagues` | Public | รายชื่อลีกและฟุตบอลถ้วยทั้งหมด |
| | `GET` | `/leagues/:id` | Public | ข้อมูลลีกตาม ID (UUID) |
| **4. Teams** | `GET` | `/teams` | Public | รายชื่อสโมสรทั้งหมด (รองรับ `?league_id=`) |
| | `GET` | `/teams/:id` | Public | ข้อมูลสโมสรตาม ID (UUID) |
| | `GET` | `/teams/:id/form` | Public | ผลงาน 5 นัดล่าสุดของทีม (W/D/L) คำนวณจากผลแข่งจริง |
| | `GET` | `/teams/:id/coach` | Public | ข้อมูลหัวหน้าผู้ฝึกสอน (Coach) รูปถ่าย สัญชาติ |
| | `GET` | `/teams/:id/squad` | Public | รายชื่อผู้เล่นทุกคนในทีม ตำแหน่ง เบอร์เสื้อ รูปถ่าย |
| **5. Matches** | `GET` | `/matches` | Public | ตารางแข่งและผลบอล (รองรับ filter: year, season, team_id, group_by) |
| | `GET` | `/matches/:id` | Public | รายละเอียดแมตช์เฉพาะคู่ตาม ID |
| | `GET` | `/matches/:id/lineups` | Public | รายชื่อ 11 ตัวจริง แผนการเล่น (Formation) และตัวสำรอง |
| | `GET` | `/matches/:id/venue` | Public | ข้อมูลสนามแข่งขัน ความจุ เมือง ช่องถ่ายทอดสด |
| **6. Standings** | `GET` | `/standings` | Public | ตารางคะแนนลีก คำนวณผลจริงใน DB (แข่ง ชนะ เสมอ แพ้ แต้ม) |
| **7. Favorites** | `GET` | `/favorites` | Optional Auth | รายชื่อทีมโปรดของผู้ใช้ (ผูกกับ User ID อัตโนมัติ) |
| | `POST` | `/favorites/:team_id` | Optional Auth | เพิ่มทีมเข้าเป็นทีมโปรด |
| | `DELETE` | `/favorites/:team_id` | Optional Auth | ลบทีมออกจากทีมโปรด |
| **8. News** | `GET` | `/news` | Public | ข่าวฟุตบอลไทยล่าสุดจากสำนักข่าวจริง (รองรับ `?limit=10`) |
| | `GET` | `/heroes` | Public | ไฮไลต์ข่าวเด่น / แบนเนอร์หน้าแรก (Alias ของ `/news`) |
| **9. System** | `GET` | `/health` | Public | ตรวจสอบสถานะ Server และการเชื่อมต่อ Database |
| **10. Sync** | `POST` | `/sync/fixtures` | Public | ดึงตารางแข่งและผลบอลจาก API ภายนอกลง DB |
| | `POST` | `/sync/squads` | Public | ดูดข้อมูลโค้ชและนักเตะ 20 ทีมลง DB อัตโนมัติ |

---

## 1. Authentication & Profile (ระบบสมาชิกและโปรไฟล์)

### กฎการตรวจสอบข้อมูล (Validation Rules)
* **username:** ต้องยาว 3-30 ตัวอักษร อนุญาตเฉพาะตัวอักษรภาษาอังกฤษ, ตัวเลข และขีดล่าง `_` (ห้ามมีช่องว่างหรืออักขระพิเศษ)
* **email:** รูปแบบอีเมลมาตรฐาน ความยาวไม่เกิน 255 ตัวอักษร
* **password:** ต้องยาว 6-72 ตัวอักษร และต้องมีทั้งตัวอักษรและตัวเลขผสมกันอย่างน้อยอย่างละ 1 ตัว
* **confirm_password:** ต้องกรอก และต้องตรงกับ `password` ทุกตัวอักษร
* **avatar_url:** (ไม่บังคับ) หากส่งมา ต้องเป็น URL ที่ขึ้นต้นด้วย `http://` หรือ `https://`

---

### POST /register
สมัครสมาชิกใหม่ ระบบทำการแฮชรหัสผ่านด้วย bcrypt และส่งคืน JWT Token พร้อมข้อมูลผู้ใช้ทันที

* **Headers:** `Content-Type: application/json`
* **Request Body:**
```json
{
  "username": "somchai99",
  "email": "somchai99@example.com",
  "password": "Password123",
  "confirm_password": "Password123",
  "avatar_url": ""
}
```
* **curl:**
```bash
curl -i -X POST http://localhost:3000/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"username":"somchai99","email":"somchai99@example.com","password":"Password123","confirm_password":"Password123"}'
```
* **Response (201 Created):**
```json
{
  "status": "success",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "e65ba61b-858a-49fe-b9db-e311945ca184",
    "username": "somchai99",
    "email": "somchai99@example.com",
    "avatar_url": "https://ui-avatars.com/api/?name=somchai99&background=0F172A&color=38BDF8&size=256",
    "role": "user",
    "created_at": "2026-11-26T15:36:47.484+07:00"
  }
}
```
* **Error Responses:**
  * `400 Bad Request`: ข้อมูลไม่ผ่านเกณฑ์ เช่น `{"error": "password must contain at least one letter and one number"}`
  * `409 Conflict`: ชื่อผู้ใช้หรืออีเมลซ้ำ `{"error": "username is already taken"}`

---

### POST /login
เข้าสู่ระบบด้วยอีเมลหรือชื่อผู้ใช้

* **Headers:** `Content-Type: application/json`
* **Request Body:**
```json
{
  "email": "somchai99@example.com",
  "password": "Password123"
}
```
* **curl:**
```bash
curl -i -X POST http://localhost:3000/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"somchai99@example.com","password":"Password123"}'
```
* **Response (200 OK):**
```json
{
  "status": "success",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "e65ba61b-858a-49fe-b9db-e311945ca184",
    "username": "somchai99",
    "email": "somchai99@example.com",
    "avatar_url": "https://ui-avatars.com/api/?name=somchai99&background=0F172A&color=38BDF8&size=256",
    "role": "user",
    "created_at": "2026-11-26T15:36:47.484+07:00"
  }
}
```
* **Error Response (401 Unauthorized):**
```json
{
  "error": "invalid email/username or password"
}
```

---

### GET /profile
ดึงข้อมูลโปรไฟล์ของผู้ใช้ที่ล็อกอินอยู่

* **Headers:** `Authorization: Bearer <TOKEN>`
* **curl:**
```bash
curl -i http://localhost:3000/api/v1/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```
* **Response (200 OK):**
```json
{
  "id": "e65ba61b-858a-49fe-b9db-e311945ca184",
  "username": "somchai99",
  "email": "somchai99@example.com",
  "avatar_url": "https://ui-avatars.com/api/?name=somchai99&background=0F172A&color=38BDF8&size=256",
  "role": "user",
  "created_at": "2026-11-26T15:36:47.484+07:00"
}
```
* **Error Response (401 Unauthorized):**
```json
{
  "error": "missing authorization header"
}
```

---

### PUT /profile
แก้ไขข้อมูลโปรไฟล์ของผู้ใช้ (รองรับทั้ง JSON Body และ Multipart Form สำหรับอัปโหลดรูปภาพ)

* **Headers:** `Authorization: Bearer <TOKEN>`

**รูปแบบที่ 1: JSON Body (เปลี่ยนชื่อผู้ใช้ หรือส่ง avatar_url)**
```bash
curl -i -X PUT http://localhost:3000/api/v1/profile \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "somchai_pro",
    "avatar_url": "https://example.com/avatar.jpg"
  }'
```

**รูปแบบที่ 2: Multipart Form (อัปโหลดไฟล์รูปภาพพร้อมเปลี่ยนชื่อ)**
```bash
curl -i -X PUT http://localhost:3000/api/v1/profile \
  -H "Authorization: Bearer <TOKEN>" \
  -F "username=somchai_pro" \
  -F "avatar=@/path/to/my-photo.jpg"
```

* **Response (200 OK):**
```json
{
  "status": "success",
  "message": "profile updated successfully",
  "user": {
    "id": "e65ba61b-858a-49fe-b9db-e311945ca184",
    "username": "somchai_pro",
    "email": "somchai99@example.com",
    "avatar_url": "http://localhost:9000/football-avatars/avatars/e65ba61b-1795747973.jpg",
    "role": "user",
    "created_at": "2026-11-26T15:36:47.484+07:00"
  }
}
```

---

### POST /profile/avatar
อัปโหลดไฟล์รูปภาพโปรไฟล์ขึ้น Object Storage (MinIO / S3) โดยเฉพาะ  
*จำกัดขนาดไฟล์ไม่เกิน 5MB, รองรับสกุลไฟล์ `.jpg`, `.jpeg`, `.png`, `.webp`, `.gif`*

* **Headers:** `Authorization: Bearer <TOKEN>`
* **curl:**
```bash
curl -i -X POST http://localhost:3000/api/v1/profile/avatar \
  -H "Authorization: Bearer <TOKEN>" \
  -F "avatar=@/path/to/my-photo.png"
```
* **Response (200 OK):**
```json
{
  "status": "success",
  "message": "avatar uploaded successfully",
  "avatar_url": "http://localhost:9000/football-avatars/avatars/e65ba61b-1795747973.png",
  "user": {
    "id": "e65ba61b-858a-49fe-b9db-e311945ca184",
    "username": "somchai_pro",
    "email": "somchai99@example.com",
    "avatar_url": "http://localhost:9000/football-avatars/avatars/e65ba61b-1795747973.png",
    "role": "user",
    "created_at": "2026-11-26T15:36:47.484+07:00"
  }
}
```

---

## 2. Seasons (ปี / ฤดูกาล)

### GET /seasons
ดึงรายชื่อปี/ฤดูกาลทั้งหมด (2022 ถึง 2026) พร้อมจำนวนแมตช์ในระบบ สำหรับนำไปสร้าง Dropdown ให้ผู้ใช้เลือก

* **curl:**
```bash
curl -i http://localhost:3000/api/v1/seasons
```
* **Response (200 OK):**
```json
[
  {
    "season": 2022,
    "label": "2022/2023",
    "matches_count": 240
  },
  {
    "season": 2023,
    "label": "2023/2024",
    "matches_count": 240
  },
  {
    "season": 2024,
    "label": "2024/2025",
    "matches_count": 240
  },
  {
    "season": 2025,
    "label": "2025/2026",
    "matches_count": 0
  },
  {
    "season": 2026,
    "label": "2026/2027",
    "matches_count": 0
  }
]
```

---

## 3. Leagues (ลีกฟุตบอล)

### GET /leagues
ดึงรายชื่อลีกและฟุตบอลถ้วยทั้งหมด

* **curl:**
```bash
curl -i http://localhost:3000/api/v1/leagues
```
* **Response (200 OK):**
```json
[
  {
    "id": "0aab480b-9e10-4fea-a6fc-884160e684d5",
    "name": "Thai League 1",
    "short_name": "T1",
    "country": "Thailand",
    "logo_url": "https://media.api-sports.io/football/leagues/296.png",
    "type": "league",
    "sort_order": 0
  }
]
```

---

### GET /leagues/:id
ดึงข้อมูลลีกตาม ID (UUID)

* **curl:**
```bash
curl -i http://localhost:3000/api/v1/leagues/0aab480b-9e10-4fea-a6fc-884160e684d5
```

---

## 4. Teams (สโมสรและทีมฟุตบอล)

### GET /teams
ดึงรายชื่อทีมฟุตบอลทั้งหมด

* **Query Parameters (Optional):**
  * `league_id`: กรองเฉพาะทีมในลีกที่ระบุ (UUID)
* **curl:**
```bash
curl -i http://localhost:3000/api/v1/teams
```

---

### GET /teams/:id
ดึงข้อมูลสโมสรตาม ID

* **curl:**
```bash
curl -i http://localhost:3000/api/v1/teams/e1c4d9dc-03a5-427c-b055-57d106c81928
```

---

### GET /teams/:id/form
ดึงผลงาน 5 นัดล่าสุดของทีม (คำนวณจากผลการแข่งขันจริงในฐานข้อมูล)

* **curl:**
```bash
curl -i http://localhost:3000/api/v1/teams/e1c4d9dc-03a5-427c-b055-57d106c81928/form
```
* **Response (200 OK):**
```json
[
  {
    "match_id": "8fa85b88-1234-5678-90ab-cdef12345678",
    "date": "2025-04-20",
    "opponent": "Buriram United",
    "score": "2 - 1",
    "result": "W",
    "is_home": true
  },
  {
    "match_id": "7eb74a77-1234-5678-90ab-cdef12345678",
    "date": "2025-04-06",
    "opponent": "Port FC",
    "score": "1 - 1",
    "result": "D",
    "is_home": false
  }
]
```

---

### GET /teams/:id/coach
ดึงข้อมูลหัวหน้าผู้ฝึกสอน (Coach) ของทีม

* **curl:**
```bash
curl -i http://localhost:3000/api/v1/teams/e1c4d9dc-03a5-427c-b055-57d106c81928/coach
```
* **Response (200 OK):**
```json
{
  "id": "c1a2b3c4-1234-5678-90ab-cdef12345678",
  "team_id": "e1c4d9dc-03a5-427c-b055-57d106c81928",
  "name": "T. Sripan",
  "nationality": "Thailand",
  "photo_url": "https://media.api-sports.io/football/coachs/921.png"
}
```

---

### GET /teams/:id/squad
ดึงรายชื่อผู้เล่นทั้งหมดในทีม (Squad) พร้อมตำแหน่ง เบอร์เสื้อ อายุ และรูปถ่าย

* **curl:**
```bash
curl -i http://localhost:3000/api/v1/teams/e1c4d9dc-03a5-427c-b055-57d106c81928/squad
```
* **Response (200 OK):**
```json
[
  {
    "id": "p1a2b3c4-1234-5678-90ab-cdef12345678",
    "team_id": "e1c4d9dc-03a5-427c-b055-57d106c81928",
    "name": "P. Khammai",
    "number": 1,
    "position": "Goalkeeper",
    "age": 31,
    "nationality": "Thailand",
    "photo_url": "https://media.api-sports.io/football/players/48975.png"
  }
]
```

---

## 5. Matches (ตารางแข่งและผลบอล)

### GET /matches
ดึงรายการแข่งขันและผลบอล พร้อมข้อมูลทีมเหย้า-เยือน โลโก้ สกอร์ นาที สถานะ และข้อมูลปี/ฤดูกาล

* **Query Parameters (Optional):**
  * `year`: กรองตามปีปฏิทิน ค.ศ. ที่ลงแข่ง (`2022`, `2023`, `2024`, `2025`, `2026`)
  * `season`: กรองตามฤดูกาลฟุตบอล เช่น `2024` หรือ `2024/2025`
  * `team_id`: กรองเฉพาะแมตช์ที่ทีมนี้ลงแข่งขัน (ทั้งเหย้าและเยือน)
  * `home_team_id` / `away_team_id`: กรองคู่แข่งขันเฉพาะเจาะจง
  * `group_by`: จัดกลุ่มผลลัพธ์แยกตามปี (`year`) หรือแยกตามฤดูกาล (`season`)
  * `status`: กรองสถานะ (`finished`, `live`, `scheduled`, `postponed`)
  * `league_id`: กรองเฉพาะลีกที่ระบุ (UUID)
  * `date`: กรองตามวันที่แข่ง รูปแบบ `YYYY-MM-DD`
  * `limit`: จำนวนรายการต่อหน้า (default: 50, สูงสุด: 1000)
  * `page`: หน้าที่ต้องการดึง (default: 1)

* **curl ตัวอย่างการกรองตามปี (2022 - 2026):**
```bash
# 1. กรองแมตช์ที่แข่งในปี 2024
curl -i "http://localhost:3000/api/v1/matches?year=2024&limit=10"

# 2. กรองแมตช์ของทีมใดทีมหนึ่งในปี 2024
curl -i "http://localhost:3000/api/v1/matches?year=2024&team_id=4508928c-d4bf-4978-9cc7-515ea59b83e6"

# 3. จัดกลุ่มแมตช์แยกเป็นรายปีอัตโนมัติ
curl -i "http://localhost:3000/api/v1/matches?group_by=year"
```

* **Response Example (รายแมตช์จะมีฟิลด์ year และ season ระบุชัดเจน):**
```json
[
  {
    "id": "404405b6-19bb-455c-a8ad-a5f023be01d5",
    "year": 2025,
    "season": "2024/2025",
    "season_year": 2024,
    "match_date": "2025-05-24T18:00:00+07:00",
    "round": "Final",
    "status": "finished",
    "home_score": 2,
    "away_score": 3,
    "home_team": {
      "id": "a1b5c773-e49f-410e-be26-f8b9ab1a968d",
      "name": "Muangthong United",
      "logo_url": "https://media.api-sports.io/football/teams/2786.png"
    },
    "away_team": {
      "id": "4508928c-d4bf-4978-9cc7-515ea59b83e6",
      "name": "Buriram United",
      "logo_url": "https://media.api-sports.io/football/teams/2780.png"
    }
  }
]
```

---

### GET /matches/:id
ดึงข้อมูลแมตช์เฉพาะคู่ตาม ID

* **curl:**
```bash
curl -i http://localhost:3000/api/v1/matches/f5a892b1-1234-5678-90ab-cdef12345678
```

---

### GET /matches/:id/lineups
ดึงรายชื่อ 11 ตัวจริง แผนการเล่น (Formation) และตัวสำรองของทั้งสองฝั่ง

* **curl:**
```bash
curl -i http://localhost:3000/api/v1/matches/f5a892b1-1234-5678-90ab-cdef12345678/lineups
```
* **Response (200 OK):**
```json
{
  "match_id": "f5a892b1-1234-5678-90ab-cdef12345678",
  "home": {
    "team_id": "11bb0e93-0ec7-4632-a0c5-115f5e888ab1",
    "team_name": "Buriram United",
    "formation": "4-3-3",
    "coach_name": "Home Head Coach",
    "starting": [
      { "number": 1, "name": "Goalkeeper", "position": "GK", "rating": "7.0" },
      { "number": 9, "name": "Striker", "position": "FW", "rating": "8.0" }
    ],
    "substitutes": []
  },
  "away": {
    "team_id": "22cc1f04-1fd8-4743-b1d6-226f6f999bc2",
    "team_name": "Bangkok United",
    "formation": "4-2-3-1",
    "coach_name": "Away Head Coach",
    "starting": [
      { "number": 1, "name": "Goalkeeper", "position": "GK", "rating": "6.8" },
      { "number": 10, "name": "Playmaker", "position": "MF", "rating": "7.6" }
    ],
    "substitutes": []
  }
}
```

---

### GET /matches/:id/venue
ดึงข้อมูลสนามแข่งขัน เมือง ความจุ และช่องทางการถ่ายทอดสด

* **curl:**
```bash
curl -i http://localhost:3000/api/v1/matches/f5a892b1-1234-5678-90ab-cdef12345678/venue
```
* **Response (200 OK):**
```json
{
  "match_id": "f5a892b1-1234-5678-90ab-cdef12345678",
  "stadium": "Chang Arena Stadium",
  "city": "Buriram, Thailand",
  "capacity": 32600,
  "broadcast": [
    "TrueID Sports 1",
    "AIS PLAY Football 2",
    "PPTV HD 36"
  ]
}
```

---

## 6. Standings (ตารางคะแนนลีก)

### GET /standings
ดึงตารางคะแนนที่คำนวณจากผลการแข่งขันจริงในฐานข้อมูล (แข่ง ชนะ เสมอ แพ้ ประตูได้-เสีย ผลต่าง แต้ม และฟอร์ม 5 นัด) แคชผ่าน Redis

* **Query Parameters (Optional):**
  * `season`: ฤดูกาลที่ต้องการคำนวณ (default: `2024`, รองรับ `2023`, `2022`)
  * `league_id`: กรองเฉพาะลีกที่ต้องการ (UUID)
* **curl:**
```bash
curl -i "http://localhost:3000/api/v1/standings?season=2024"
```
* **Response (200 OK):**
```json
[
  {
    "rank": 1,
    "team_id": "11bb0e93-0ec7-4632-a0c5-115f5e888ab1",
    "team_name": "Buriram United",
    "team_logo": "https://media.api-sports.io/football/teams/2275.png",
    "played": 30,
    "won": 20,
    "drawn": 9,
    "lost": 1,
    "goals_for": 70,
    "goals_against": 27,
    "goal_diff": 43,
    "points": 69,
    "form": ["W", "W", "D", "W", "W"]
  }
]
```

---

## 7. Favorites (ระบบทีมโปรด)

*หากผู้ใช้ล็อกอิน ระบบจะผูกทีมโปรดกับ User ID จาก JWT Token โดยอัตโนมัติ (หรือส่ง Header `X-User-ID: <UUID>` สำหรับการทดสอบ)*

### GET /favorites
ดึงรายการทีมโปรดทั้งหมดของผู้ใช้

* **curl:**
```bash
curl -i http://localhost:3000/api/v1/favorites \
  -H "Authorization: Bearer <TOKEN>"
```

---

### POST /favorites/:team_id
เพิ่มทีมเข้าเป็นทีมโปรด

* **curl:**
```bash
curl -i -X POST http://localhost:3000/api/v1/favorites/e1c4d9dc-03a5-427c-b055-57d106c81928 \
  -H "Authorization: Bearer <TOKEN>"
```

---

### DELETE /favorites/:team_id
ลบทีมออกจากทีมโปรด

* **curl:**
```bash
curl -i -X DELETE http://localhost:3000/api/v1/favorites/e1c4d9dc-03a5-427c-b055-57d106c81928 \
  -H "Authorization: Bearer <TOKEN>"
```

---

## 8. Football News (ข่าวสารฟุตบอล)

ระบบข่าวสารดึงข้อมูลจากสำนักข่าวฟุตบอลไทยชั้นนำ (Ballthai API) ส่งตรงภาษาไทยแท้ ข่าวอัปเดตตลอดเวลา พร้อมภาพข่าวความละเอียดสูงและลิงก์บทความจริง มีระบบแคชผ่าน Redis เป็นเวลา 30 นาที

### GET /news
ดึงรายการข่าวฟุตบอลล่าสุด

* **Query Parameters (Optional):**
  * `limit`: จำนวนข่าวที่ต้องการดึง (default: 10, สูงสุด: 30)
* **curl:**
```bash
curl -i "http://localhost:3000/api/v1/news?limit=10"
```

---

### GET /heroes
ดึงไฮไลต์ข่าวเด่นสำหรับแสดงผลในส่วน Header Carousel / Hero Banner หน้าแรก  
*(เป็น Alias ของ `/news` รองรับ `?limit=` เช่นเดียวกัน)*

* **curl:**
```bash
curl -i "http://localhost:3000/api/v1/heroes?limit=5"
```

* **Response (200 OK):**
```json
[
  {
    "id": "hero-407370",
    "title": "“สวายเรียง” สร้างประวัติศาสตร์เปิดบ้านอัดทีมเจลีก ญี่ปุ่น ครั้งแรก",
    "subtitle": "“สวายเรียง เอฟซี” สโมสรตัวแทนจากประเทศกัมพูชา สร้างผลงานยอดเยี่ยมบนเวทีเอเชีย หลังเปิดบ้านเอาชนะ มาชิดะ เซลเวีย ทีมดังจากศึกเจลีก 1 ประเทศญี่ปุ่น ไปอย่างยอดเยี่ยม 2-0 ในศึกฟุตบอลเอเอฟซี แชมเปี้ยนส์ ลีก 2...",
    "image_url": "https://media.ballthai.com/2026/11/Preah-Khan-Reach-Svay-Rieng-FC.jpg",
    "tag": "เอเอฟซี แชมเปี้ยนส์ ลีก",
    "redirect_url": "https://www.ballthai.com/สวายเรียง-สร้างประวัติศาสตร์เปิดบ้านอัดทีมเจลีก-ญี่ปุ่น-ครั้งแรก/"
  },
  {
    "id": "hero-407366",
    "title": "“โค้ชเบนซ์” ผิดหวัง “บีจี ปทุม” พ่าย ไลออน ซิตี้ฯ ต้องแก้ไขก่อนเกมชี้ชะตา",
    "subtitle": "เฮดโค้ช “เดอะ แรบบิท” ยอมรับทีมยังขาดความเฉียบคม หลังเปิดหน้าแลกก่อนจะแพ้ ไลออน ซิตี้ เซลเลอร์ส ในศึก ACL 2...",
    "image_url": "https://media.ballthai.com/2026/11/BG-Pathum-United-2.jpg",
    "tag": "เอเอฟซี แชมเปี้ยนส์ ลีก",
    "redirect_url": "https://www.ballthai.com/โค้ชเบนซ์-ผิดหวัง-บีจี-ปทุม-พ่าย-ไลออน-ซิตี้ฯ/"
  },
  {
    "id": "hero-407360",
    "title": "“บุรีรัมย์” ตั้ง “โซรัน-เอเมอร์สัน” คุมทัพชั่วคราว",
    "subtitle": "“ปราสาทสายฟ้า” บุรีรัมย์ ยูไนเต็ด ที่เพิ่งแยกทางกับ ราอูล คาเนดา กุนซือชาวสเปน ล่าสุดประกาศมอบหมายให้ โซรัน ยานโควิช และ เอเมอร์สัน เปไรร่า ทำหน้าที่คุมทัพชั่วคราว...",
    "image_url": "https://media.ballthai.com/2026/11/BURIRAM-UNITED-4.jpg",
    "tag": "ไทยลีก",
    "redirect_url": "https://www.ballthai.com/บุรีรัมย์-ตั้ง-โซรัน-เอเมอร์สัน-คุมทัพชั่วคราว/"
  }
]
```

---

## 9. System (สถานะระบบ)

### GET /health
ตรวจสอบสถานะการทำงานของ Server และการเชื่อมต่อ Database

* **curl:**
```bash
curl -i http://localhost:3000/api/v1/health
```
* **Response (200 OK):**
```json
{
  "status": "ok",
  "database": "connected"
}
```

---

## 10. Data Synchronization (ระบบซิงค์ข้อมูล)

### POST /sync/fixtures
สั่งดึงข้อมูลตารางแข่งและผลบอลจาก API ภายนอก บันทึกลงฐานข้อมูล

* **curl:**
```bash
curl -i -X POST http://localhost:3000/api/v1/sync/fixtures \
  -H "Content-Type: application/json" \
  -d '{"season": 2024}'
```

---

### POST /sync/squads
สั่งดึงข้อมูลโค้ชและนักเตะของทุกทีม (20 ทีม) ลงฐานข้อมูลอัตโนมัติ

* **curl:**
```bash
curl -i -X POST http://localhost:3000/api/v1/sync/squads
```
