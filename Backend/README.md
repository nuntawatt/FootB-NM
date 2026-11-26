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
| **Teams** | `GET` | `/teams` | รายชื่อทีมและสโมสรทั้งหมด (รองรับ filter ตามลีก) |
| | `GET` | `/teams/:id` | ข้อมูลทีมตาม ID |
| **Matches** | `GET` | `/matches` | ตารางแข่งและผลบอล (รองรับ filter ปี/ฤดูกาล/สถานะ/วัน) |
| | `GET` | `/matches/:id` | ข้อมูลแมตช์ตาม ID |
| **Sync** | `POST` | `/sync/fixtures` | สั่งดึงข้อมูลจาก API-Football บันทึกลง Database ทันที |

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
  },
  {
    "id": "a33ef81e-4ff1-4ef0-bfc8-9d2ce5e640d2",
    "name": "Thai League 2",
    "short_name": "T2",
    "country": "Thailand",
    "logo_url": "https://media.api-sports.io/football/leagues/297.png",
    "type": "league",
    "sort_order": 1
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
- **Response (200 OK):** คืน Object ข้อมูลลีก
- **Error (404 Not Found):** `{"error": "league not found"}`

---

## 4. Teams (สโมสร/ทีมฟุตบอล)

### GET /teams
ดึงรายชื่อทีมฟุตบอลทั้งหมด

- **Query Parameters (Optional):**
  - `league_id`: กรองเฉพาะทีมที่สังกัดลีกนี้ (UUID)
- **curl (ดึงทั้งหมด):**
```bash
curl -i http://localhost:3000/api/v1/teams
```
- **curl (กรองตามลีก):**
```bash
curl -i "http://localhost:3000/api/v1/teams?league_id=0aab480b-9e10-4fea-a6fc-884160e684d5"
```
- **Response (200 OK):**
```json
[
  {
    "id": "11bb0e93-0ec7-4632-a0c5-115f5e888ab1",
    "league_id": "0aab480b-9e10-4fea-a6fc-884160e684d5",
    "name": "Buriram United",
    "short_name": "",
    "logo_url": "https://media.api-sports.io/football/teams/2275.png"
  },
  {
    "id": "22cc1f04-1fd8-4743-b1d6-226f6f999bc2",
    "league_id": "0aab480b-9e10-4fea-a6fc-884160e684d5",
    "name": "Bangkok United",
    "short_name": "",
    "logo_url": "https://media.api-sports.io/football/teams/2276.png"
  }
]
```

---

### GET /teams/:id
ดึงข้อมูลทีมตาม ID

- **curl:**
```bash
curl -i http://localhost:3000/api/v1/teams/11bb0e93-0ec7-4632-a0c5-115f5e888ab1
```
- **Response (200 OK):**
```json
{
  "id": "11bb0e93-0ec7-4632-a0c5-115f5e888ab1",
  "league_id": "0aab480b-9e10-4fea-a6fc-884160e684d5",
  "name": "Buriram United",
  "short_name": "",
  "logo_url": "https://media.api-sports.io/football/teams/2275.png"
}
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

- **curl (ดึงแมตช์ของฤดูกาล 2024):**
```bash
curl -i "http://localhost:3000/api/v1/matches?season=2024&limit=20"
```
- **curl (ดึงแมตช์ของฤดูกาล 2023):**
```bash
curl -i "http://localhost:3000/api/v1/matches?season=2023&limit=20"
```
- **curl (ดึงแมตช์ของฤดูกาล 2022):**
```bash
curl -i "http://localhost:3000/api/v1/matches?season=2022&limit=20"
```
- **curl (ดูเฉพาะผลบอลที่จบแล้ว):**
```bash
curl -i "http://localhost:3000/api/v1/matches?status=finished&limit=20"
```
- **curl (ดูบอลสด):**
```bash
curl -i "http://localhost:3000/api/v1/matches?status=live"
```
- **curl (ดูแมตช์ตามวันที่):**
```bash
curl -i "http://localhost:3000/api/v1/matches?date=2025-04-30"
```
- **Response (200 OK):**
```json
[
  {
    "id": "f5a892b1-1234-5678-90ab-cdef12345678",
    "league_id": "0aab480b-9e10-4fea-a6fc-884160e684d5",
    "home_team_id": "11bb0e93-0ec7-4632-a0c5-115f5e888ab1",
    "away_team_id": "22cc1f04-1fd8-4743-b1d6-226f6f999bc2",
    "match_date": "2025-04-30T18:00:00+07:00",
    "home_score": 7,
    "away_score": 0,
    "home_penalty_score": null,
    "away_penalty_score": null,
    "minute": null,
    "round": "Regular Season - 30",
    "status": "finished",
    "home_team": {
      "id": "11bb0e93-0ec7-4632-a0c5-115f5e888ab1",
      "name": "Buriram United",
      "logo_url": "https://media.api-sports.io/football/teams/2275.png"
    },
    "away_team": {
      "id": "22cc1f04-1fd8-4743-b1d6-226f6f999bc2",
      "name": "Nong Bua Pitchaya",
      "logo_url": "https://media.api-sports.io/football/teams/2280.png"
    },
    "league": {
      "id": "0aab480b-9e10-4fea-a6fc-884160e684d5",
      "name": "Thai League 1",
      "type": "league"
    }
  }
]
```

---

### GET /matches/:id
ดึงข้อมูลแมตช์เฉพาะคู่ตาม ID

- **curl:**
```bash
curl -i http://localhost:3000/api/v1/matches/f5a892b1-1234-5678-90ab-cdef12345678
```
- **Response (200 OK):** คืน Object แมตช์พร้อมข้อมูลทีมเหย้าและเยือน
- **Error (404 Not Found):** `{"error": "match not found"}`

---

## 6. Sync On-Demand (ซิงก์ข้อมูลลงฐานข้อมูล)

### POST /sync/fixtures
สั่งให้ Backend วิ่งไปดึงข้อมูลจาก API-Football ตามปีที่ระบุ แล้วนำมา Upsert บันทึกลง PostgreSQL ทันที

- **Headers:** `Content-Type: application/json`
- **Request Body:**
```json
{
  "season": 2024    // int, required (เช่น 2022, 2023, 2024)
}
```
- **curl:**
```bash
curl -i -X POST http://localhost:3000/api/v1/sync/fixtures \
  -H "Content-Type: application/json" \
  -d '{"season": 2024}'
```
- **Response (200 OK):**
```json
{
  "status": "success",
  "season": 2024,
  "leagues_synced": 1,
  "teams_created": 1,
  "matches_created": 240,
  "matches_updated": 0,
  "matches_skipped": 0,
  "quota_remaining": {
    "day": 91,
    "minute": 7
  }
}
```
