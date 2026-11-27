# Football API Documentation for Frontend

Base URL: `http://localhost:3000/api/v1`  
LAN Base URL: `http://10.5.2.103:3000/api/v1`  
Content-Type: `application/json`

---

## สรุป Endpoints ทั้งหมดในระบบ (API Overview)

| หมวดหมู่ | Method | Endpoint | สิทธิ์เข้าถึง | หน้าที่ |
|:---|:---:|:---|:---:|:---|
| **1. Auth & Profile** | `POST` | `/register` | Public | สมัครสมาชิกใหม่ (ส่งคืน JWT Token) |
| | `POST` | `/login` | Public | เข้าสู่ระบบด้วยอีเมลหรือชื่อผู้ใช้ |
| | `GET` | `/profile` | Bearer Token | ดูข้อมูลโปรไฟล์ส่วนตัว |
| **2. Seasons** | `GET` | `/seasons` | Public | รายการปี/ฤดูกาลทั้งหมด (2022-2026) |
| **3. Leagues** | `GET` | `/leagues` | Public | รายชื่อลีกและบอลถ้วยทั้งหมด |
| | `GET` | `/leagues/:id` | Public | ข้อมูลลีกตาม ID |
| **4. Teams** | `GET` | `/teams` | Public | รายชื่อทีมทั้งหมด (รองรับ `?league_id=`) |
| | `GET` | `/teams/:id` | Public | ข้อมูลทีมตาม ID |
| | `GET` | `/teams/:id/form` | Public | ผลงาน 5 นัดล่าสุดของทีม (W/D/L) |
| | `GET` | `/teams/:id/coach` | Public | ข้อมูลหัวหน้าผู้ฝึกสอน |
| | `GET` | `/teams/:id/squad` | Public | รายชื่อนักเตะทั้งหมดของทีม |
| **5. Matches** | `GET` | `/matches` | Public | ตารางแข่งและผลบอล (รองรับ filter) |
| | `GET` | `/matches/:id` | Public | ข้อมูลแมตช์ตาม ID |
| | `GET` | `/matches/:id/lineups` | Public | รายชื่อ 11 ตัวจริง แผนการเล่น ตัวสำรอง |
| | `GET` | `/matches/:id/venue` | Public | ข้อมูลสนาม เมือง ความจุ ช่องถ่ายทอดสด |
| **6. Standings** | `GET` | `/standings` | Public | ตารางคะแนนคำนวณจากผลแข่งจริงใน DB มีการใช้ Redis |
| **7. Favorites** | `GET` | `/favorites` | Optional Auth | รายชื่อทีมโปรดของผู้ใช้ |
| | `POST` | `/favorites/:team_id` | Optional Auth | เพิ่มทีมเข้าทีมโปรด |
| | `DELETE` | `/favorites/:team_id` | Optional Auth | ลบทีมออกจากทีมโปรด |
| **8. Highlights** | `GET` | `/heroes` | Public | ข่าวเด่นและแบนเนอร์ไฮไลท์หน้าแรก |
| **9. System** | `GET` | `/health` | Public | ตรวจสอบสถานะ Server และ Database |
| **10. Sync** | `POST` | `/sync/fixtures` | Public | ดึงแมตช์จาก API-Football บันทึกลง DB และมีใช้ Redis |
| | `POST` | `/sync/squads` | Public | ดึงโค้ชและนักเตะ 20 ทีมลง DB อัตโนมัติ |

---

## 1. Authentication & Profile (ระบบสมาชิกและเข้าสู่ระบบ)

### กฎการตรวจสอบข้อมูล (Validation Rules)
* **username:** ต้องยาว 3-30 ตัวอักษร อนุญาตเฉพาะตัวอักษรภาษาอังกฤษ, ตัวเลข และขีดล่าง `_` (ห้ามมีช่องว่างหรืออักขระพิเศษ)
* **email:** รูปแบบอีเมลมาตรฐาน ความยาวไม่เกิน 255 ตัวอักษร
* **password:** ต้องยาว 6-72 ตัวอักษร และต้องมีทั้งตัวอักษรและตัวเลขผสมกันอย่างน้อยอย่างละ 1 ตัว
* **confirm_password:** ต้องกรอก และต้องตรงกับ `password` ทุกตัวอักษร
* **avatar_url:** (ไม่บังคับ) หากส่งมา ต้องเป็น URL ที่ขึ้นต้นด้วย `http://` หรือ `https://`

---

### POST /register
สมัครสมาชิกใหม่ ระบบจะทำการแฮชรหัสผ่านด้วย bcrypt และส่งคืน JWT Token พร้อมข้อมูลผู้ใช้ทันที

- **Headers:** `Content-Type: application/json`
- **Request Body:**
```json
{
  "username": "somchai99",
  "email": "somchai99@example.com",
  "password": "Password123",
  "confirm_password": "Password123",
  "avatar_url": ""
}
```
- **curl:**
```bash
curl -i -X POST http://localhost:3000/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"username":"somchai99","email":"somchai99@example.com","password":"Password123","confirm_password":"Password123"}'
```
- **Response (201 Created):**
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
- **Error Responses:**
  - `400 Bad Request`: ข้อมูลไม่ผ่านเกณฑ์ เช่น `{"error": "password must contain at least one letter and one number"}`
  - `409 Conflict`: ชื่อผู้ใช้หรืออีเมลซ้ำ `{"error": "username is already taken"}`

---

### POST /login
เข้าสู่ระบบด้วยอีเมลหรือชื่อผู้ใช้

- **Headers:** `Content-Type: application/json`
- **Request Body:**
```json
{
  "email": "somchai99@example.com",
  "password": "Password123"
}
```
- **curl:**
```bash
curl -i -X POST http://localhost:3000/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"somchai99@example.com","password":"Password123"}'
```
- **Response (200 OK):**
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
- **Error Response (401 Unauthorized):**
```json
{
  "error": "invalid email/username or password"
}
```

---

### GET /profile
ดึงข้อมูลโปรไฟล์ของผู้ใช้ที่ล็อกอินอยู่ (ต้องส่ง Bearer Token)

- **Headers:** `Authorization: Bearer <TOKEN>`
- **curl:**
```bash
curl -i http://localhost:3000/api/v1/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```
- **Response (200 OK):**
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
- **Error Response (401 Unauthorized):**
```json
{
  "error": "missing authorization header"
}
```

---

## 2. Seasons (ปี / ฤดูกาล)

### GET /seasons
ดึงรายชื่อปี/ฤดูกาลทั้งหมด (2022 ถึง 2026) พร้อมจำนวนแมตช์ในระบบ สำหรับนำไปทำ Dropdown ใน Frontend

- **curl:**
```bash
curl -i http://localhost:3000/api/v1/seasons
```
- **Response (200 OK):**
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
ดึงรายชื่อลีกและบอลถ้วยทั้งหมด

- **curl:**
```bash
curl -i http://localhost:3000/api/v1/leagues
```
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
    "sort_order": 0
  }
]
```

---

### GET /leagues/:id
ดึงข้อมูลลีกตาม ID

- **curl:**
```bash
curl -i http://localhost:3000/api/v1/leagues/0aab480b-9e10-4fea-a6fc-884160e684d5
```

---

## 4. Teams (สโมสร/ทีมฟุตบอล)

### GET /teams
ดึงรายชื่อทีมฟุตบอลทั้งหมด

- **Query Parameters (Optional):**
  - `league_id`: กรองเฉพาะทีมที่สังกัดลีกนี้ (UUID)
- **curl:**
```bash
curl -i http://localhost:3000/api/v1/teams
```

---

### GET /teams/:id
ดึงข้อมูลทีมตาม ID

- **curl:**
```bash
curl -i http://localhost:3000/api/v1/teams/e1c4d9dc-03a5-427c-b055-57d106c81928
```

---

### GET /teams/:id/form
ดึงผลงาน 5 นัดล่าสุดของทีม (คำนวณจากผลการแข่งขันจริงในฐานข้อมูล)

- **curl:**
```bash
curl -i http://localhost:3000/api/v1/teams/e1c4d9dc-03a5-427c-b055-57d106c81928/form
```
- **Response (200 OK):**
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

- **curl:**
```bash
curl -i http://localhost:3000/api/v1/teams/e1c4d9dc-03a5-427c-b055-57d106c81928/coach
```
- **Response (200 OK):**
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

- **curl:**
```bash
curl -i http://localhost:3000/api/v1/teams/e1c4d9dc-03a5-427c-b055-57d106c81928/squad
```
- **Response (200 OK):**
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
ดึงรายการแข่งขันและผลบอล พร้อมข้อมูลทีมเหย้า-เยือน โลโก้ สกอร์ นาที และสถานะ

- **Query Parameters (Optional):**
  - `season`: กรองตามฤดูกาล (`2024`, `2023`, `2022`)
  - `year`: กรองตามปีปฏิทินที่แข่ง (`2024`)
  - `status`: กรองสถานะ (`finished`, `live`, `scheduled`)
  - `league_id`: กรองเฉพาะลีกที่ระบุ (UUID)
  - `date`: กรองตามวันที่แข่ง รูปแบบ `YYYY-MM-DD`
  - `limit`: จำนวนรายการต่อหน้า (default: 50, สูงสุด: 200)
  - `page`: หน้าที่ต้องการดึง (default: 1)

- **curl:**
```bash
curl -i "http://localhost:3000/api/v1/matches?season=2024&limit=20"
```

---

### GET /matches/:id
ดึงข้อมูลแมตช์เฉพาะคู่ตาม ID

- **curl:**
```bash
curl -i http://localhost:3000/api/v1/matches/f5a892b1-1234-5678-90ab-cdef12345678
```

---

### GET /matches/:id/lineups
ดึงรายชื่อ 11 ตัวจริง แผนการเล่น และตัวสำรองของทั้งสองฝั่ง

- **curl:**
```bash
curl -i http://localhost:3000/api/v1/matches/f5a892b1-1234-5678-90ab-cdef12345678/lineups
```
- **Response (200 OK):**
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

- **curl:**
```bash
curl -i http://localhost:3000/api/v1/matches/f5a892b1-1234-5678-90ab-cdef12345678/venue
```
- **Response (200 OK):**
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
ดึงตารางคะแนนที่คำนวณจากผลการแข่งขันจริงในฐานข้อมูล (แข่ง, ชนะ, เสมอ, แพ้, ประตูได้-เสีย, คะแนน และฟอร์ม 5 นัด)

- **Query Parameters (Optional):**
  - `season`: ฤดูกาลที่ต้องการคำนวณ (default: `2024`, รองรับ `2023`, `2022`)
  - `league_id`: กรองเฉพาะลีกที่ต้องการ (UUID)
- **curl:**
```bash
curl -i "http://localhost:3000/api/v1/standings?season=2024"
```
- **Response (200 OK):**
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

*หากผู้ใช้ล็อกอิน ระบบจะผูกทีมโปรดกับ User ID จาก JWT Token โดยอัตโนมัติ (หรือส่ง `X-User-ID: <UUID>` เพื่อทดสอบ)*

### GET /favorites
ดึงรายการทีมโปรดทั้งหมดของผู้ใช้

- **curl:**
```bash
curl -i http://localhost:3000/api/v1/favorites \
  -H "Authorization: Bearer <TOKEN>"
```

---

### POST /favorites/:team_id
เพิ่มทีมเข้าเป็นทีมโปรด

- **curl:**
```bash
curl -i -X POST http://localhost:3000/api/v1/favorites/e1c4d9dc-03a5-427c-b055-57d106c81928 \
  -H "Authorization: Bearer <TOKEN>"
```

---

### DELETE /favorites/:team_id
ลบทีมออกจากทีมโปรด

- **curl:**
```bash
curl -i -X DELETE http://localhost:3000/api/v1/favorites/e1c4d9dc-03a5-427c-b055-57d106c81928 \
  -H "Authorization: Bearer <TOKEN>"
```

---

## 8. Highlights (แบนเนอร์ข่าวเด่นหน้าแรก)

### GET /heroes
ดึงข่าวเด่น บิ๊กแมตช์ประจำสัปดาห์ และแบนเนอร์ภาพความละเอียดสูงสำหรับหน้าแรก

- **curl:**
```bash
curl -i http://localhost:3000/api/v1/heroes
```
- **Response (200 OK):**
```json
[
  {
    "id": "hero-1",
    "title": "ศึกบิ๊กแมตช์ไทยลีก นัดหยุดโลกสุดสัปดาห์นี้!",
    "subtitle": "ติดตามการถ่ายทอดสดแบบ Real-time พร้อมผลการวิเคราะห์เจาะลึก",
    "image_url": "https://images.unsplash.com/photo-1518091043644-c1d4457512c6?auto=format&fit=crop&w=1200&q=80",
    "tag": "BIG MATCH",
    "match_id": "f5a892b1-1234-5678-90ab-cdef12345678",
    "redirect_url": "/matches"
  }
]
```

---

## 9. System (สถานะระบบ)

### GET /health
ตรวจสอบสถานะการทำงานของ Server และการเชื่อมต่อ Database

- **curl:**
```bash
curl -i http://localhost:3000/api/v1/health
```

---

## 10. Data Synchronization (ระบบดึงข้อมูลลงฐานข้อมูล)

### 10.1 รันผ่าน HTTP API

#### POST /sync/fixtures
ดึงแมตช์การแข่งขันจาก API-Football บันทึกลงฐานข้อมูล

- **curl:**
```bash
curl -i -X POST http://localhost:3000/api/v1/sync/fixtures \
  -H "Content-Type: application/json" \
  -d '{"season": 2024}'
```

#### POST /sync/squads
ดึงข้อมูลโค้ชและนักเตะของทุกทีม (20 ทีม) ลงฐานข้อมูลอัตโนมัติ

- **curl:**
```bash
curl -i -X POST http://localhost:3000/api/v1/sync/squads
```

---

### 10.2 รันผ่าน Terminal CLI

```bash
# 1. ดูดข้อมูลโค้ชและนักเตะทั้ง 20 ทีม (มีระบบเว้นจังหวะกัน rate limit และข้ามทีมที่เคยดูดแล้ว)
go run ./cmd/sync squads

# 2. ดูดโปรแกรมและผลการแข่งขันแบบระบุช่วงปี
go run ./cmd/sync fixtures -season 2022-2024

# 3. ดูรายการลีกที่ผูกไว้
go run ./cmd/sync leagues
```

---

## 11. Database Migration

หากต้องการสร้างตาราง `coaches` และ `players` ใน PostgreSQL:

```bash
PGPASSWORD=nuntawat838 psql -h localhost -U postgres -d football -f migrations/000003_create_coaches_and_players.up.sql
```
