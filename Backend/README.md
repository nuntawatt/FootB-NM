# Football API Documentation for Frontend

Base URL: `http://localhost:3000/api/v1`  
LAN Base URL: `http://10.5.2.103:3000/api/v1`  
Content-Type: `application/json`

---

## สรุป Endpoints ทั้งหมดในระบบ

| หมวดหมู่ | Method | Endpoint | หน้าที่ |
|:---|:---:|:---|:---|
| **System** | `GET` | `/health` | ตรวจสอบสถานะ Server และ Database |
| **Seasons** | `GET` | `/seasons` | รายการปี/ฤดูกาลทั้งหมด (2022-2026) สำหรับ Dropdown |
| **Leagues** | `GET` | `/leagues` | รายชื่อลีกและบอลถ้วยทั้งหมด |
| | `GET` | `/leagues/:id` | ข้อมูลลีกตาม ID |
| **Teams** | `GET` | `/teams` | รายชื่อทีมทั้งหมด (รองรับ ?league_id=) |
| | `GET` | `/teams/:id` | ข้อมูลทีมตาม ID |
| | `GET` | `/teams/:id/form` | ผลงาน 5 นัดล่าสุดของทีม (W/D/L) |
| | `GET` | `/teams/:id/coach` | ข้อมูลโค้ชของทีม |
| | `GET` | `/teams/:id/squad` | รายชื่อนักเตะทั้งหมดของทีม |
| **Matches** | `GET` | `/matches` | ตารางแข่งและผลบอล (รองรับ filter ปี/ฤดูกาล/สถานะ/วัน) |
| | `GET` | `/matches/:id` | ข้อมูลแมตช์ตาม ID |
| | `GET` | `/matches/:id/lineups` | รายชื่อ 11 ตัวจริงและแผนการเล่น |
| | `GET` | `/matches/:id/venue` | ข้อมูลสนามและช่องถ่ายทอดสด |
| **Standings** | `GET` | `/standings` | ตารางคะแนนคำนวณจากผลแข่งจริงใน DB |
| **Favorites** | `GET` | `/favorites` | รายชื่อทีมโปรดของผู้ใช้ |
| | `POST` | `/favorites/:team_id` | เพิ่มทีมเข้าทีมโปรด |
| | `DELETE` | `/favorites/:team_id` | ลบทีมออกจากทีมโปรด |
| **Highlights**| `GET` | `/heroes` | ข่าวเด่นและแบนเนอร์ไฮไลท์หน้าแรก |
| **Sync** | `POST` | `/sync/fixtures` | ดึงแมตช์จาก API-Football บันทึกลง Database |
| | `POST` | `/sync/squads` | ดึงโค้ชและนักเตะ 20 ทีมลง Database อัตโนมัติ |

---

## 1. System

### GET /health
ตรวจสอบสถานะการทำงานของ Server และ Database

- **curl:**
```bash
curl -i http://localhost:3000/api/v1/health
```
- **Response (200 OK):**
```json
{
  "status": "ok"
}
```

---

## 2. Seasons (ปี / ฤดูกาล)

### GET /seasons
ดึงรายชื่อปี/ฤดูกาลทั้งหมด (2022 ถึง 2026) พร้อมจำนวนแมตช์ในระบบ สำหรับนำไปทำ Dropdown เลือกปีใน Frontend

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
  },
  {
    "id": "p2b3c4d5-1234-5678-90ab-cdef12345678",
    "team_id": "e1c4d9dc-03a5-427c-b055-57d106c81928",
    "name": "E. Fernandes",
    "number": 4,
    "position": "Defender",
    "age": 30,
    "nationality": "Brazil",
    "photo_url": "https://media.api-sports.io/football/players/12345.png"
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
  - `status`: กรองสถานะ (`finished` จบแล้ว, `live` กำลังเตะ, `scheduled` ยังไม่แข่ง)
  - `league_id`: กรองเฉพาะลีกที่ระบุ (UUID)
  - `date`: กรองตามวันที่แข่ง รูปแบบ `YYYY-MM-DD` (เช่น `2025-04-30`)
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
  },
  {
    "rank": 2,
    "team_id": "22cc1f04-1fd8-4743-b1d6-226f6f999bc2",
    "team_name": "Bangkok United",
    "team_logo": "https://media.api-sports.io/football/teams/2276.png",
    "played": 30,
    "won": 17,
    "drawn": 10,
    "lost": 3,
    "goals_for": 53,
    "goals_against": 24,
    "goal_diff": 29,
    "points": 61,
    "form": ["W", "D", "W", "W", "L"]
  }
]
```

---

## 7. Favorites (ระบบทีมโปรด)

### GET /favorites
ดึงรายการทีมโปรดทั้งหมดของผู้ใช้

- **Headers (Optional):**
  - `X-User-ID`: UUID ของผู้ใช้ (หากไม่ส่ง ระบบจะใช้ Default Test User ID ให้โดยอัตโนมัติ)
- **curl:**
```bash
curl -i http://localhost:3000/api/v1/favorites
```

---

### POST /favorites/:team_id
เพิ่มทีมเข้าเป็นทีมโปรด

- **curl:**
```bash
curl -i -X POST http://localhost:3000/api/v1/favorites/e1c4d9dc-03a5-427c-b055-57d106c81928
```
- **Response (201 Created):**
```json
{
  "status": "success",
  "message": "team added to favorites",
  "team_id": "e1c4d9dc-03a5-427c-b055-57d106c81928"
}
```

---

### DELETE /favorites/:team_id
ลบทีมออกจากทีมโปรด

- **curl:**
```bash
curl -i -X DELETE http://localhost:3000/api/v1/favorites/e1c4d9dc-03a5-427c-b055-57d106c81928
```
- **Response (200 OK):**
```json
{
  "status": "success",
  "message": "team removed from favorites",
  "team_id": "e1c4d9dc-03a5-427c-b055-57d106c81928"
}
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
  },
  {
    "id": "hero-2",
    "title": "สรุปตารางคะแนนล่าสุด ใครเต็งแชมป์ประจำฤดูกาล?",
    "subtitle": "เจาะลึกฟอร์ม 5 นัดหลังสุดของหัวตาราง",
    "image_url": "https://images.unsplash.com/photo-1579952363873-27f3bade9f55?auto=format&fit=crop&w=1200&q=80",
    "tag": "ANALYSIS",
    "redirect_url": "/standings"
  }
]
```

---

## 9. Data Synchronization (ระบบดึงข้อมูลลงฐานข้อมูล)

### 9.1 รันผ่าน HTTP API

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

### 9.2 รันผ่าน Terminal CLI

```bash
# 1. ดูดข้อมูลโค้ชและนักเตะทั้ง 20 ทีม (รันคำสั่งเดียวจบ มีระบบเว้นจังหวะกัน rate limit อัตโนมัติ)
go run ./cmd/sync squads

# 2. ดูดโปรแกรมและผลการแข่งขันแบบระบุช่วงปี
go run ./cmd/sync fixtures -season 2022-2024

# 3. ดูรายการลีกที่ผูกไว้
go run ./cmd/sync leagues
```

---

## 10. Database Migration

หากต้องการสร้างตาราง `coaches` และ `players` ใน PostgreSQL:

```bash
PGPASSWORD=nuntawat838 psql -h localhost -U postgres -d football -f migrations/000003_create_coaches_and_players.up.sql
```
