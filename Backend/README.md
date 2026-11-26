# Football API Documentation for Frontend

Base URL: `http://localhost:3000/api/v1`  
LAN Base URL: `http://10.5.2.103:3000/api/v1`  
Content-Type: `application/json`

---

## สรุปภาพรวม Endpoints ทั้งหมด

| หมวดหมู่ | Method | Endpoint | หน้าที่ |
|:---|:---:|:---|:---|
| **System** | `GET` | `/health` | ตรวจสอบสถานะ Server และ Database |
| **Leagues** | `GET` | `/leagues` | รายชื่อลีกทั้งหมด |
| | `GET` | `/leagues/:id` | ข้อมูลลีกตาม ID |
| | `POST` | `/leagues` | สร้างลีกใหม่ |
| | `PUT` | `/leagues/:id` | แก้ไขข้อมูลลีก |
| | `DELETE`| `/leagues/:id` | ลบลีก (Soft Delete) |
| **Teams** | `GET` | `/teams` | รายชื่อทีมทั้งหมด (รองรับ filter ตามลีก) |
| | `GET` | `/teams/:id` | ข้อมูลทีมตาม ID |
| **Matches** | `GET` | `/matches` | ตารางแข่งและผลบอล (รองรับ filter สถานะ/วัน/ลีก) |
| | `GET` | `/matches/:id` | ข้อมูลแมตช์ตาม ID |

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

## 2. Leagues (ลีกฟุตบอล)

### GET /leagues
ดึงรายชื่อทุกลีกในระบบ

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

### POST /leagues
สร้างลีกใหม่

- **curl:**
```bash
curl -i -X POST http://localhost:3000/api/v1/leagues \
  -H "Content-Type: application/json" \
  -d '{"name": "Premier League", "short_name": "EPL", "country": "England", "logo_url": "https://example.com/logo.png"}'
```
- **Request Body:**
```json
{
  "name": "Premier League",      // string, required
  "short_name": "EPL",           // string, optional
  "country": "England",          // string, optional
  "logo_url": "https://...png"   // string, optional
}
```
- **Response (201 Created):**
```json
{
  "id": "e8432a51-3b7c-4866-8f3e-289cf0b31e9a",
  "name": "Premier League",
  "short_name": "EPL",
  "country": "England",
  "logo_url": "https://...png",
  "type": "league",
  "sort_order": 0
}
```

---

### GET /leagues/:id
ดึงข้อมูลลีกตาม UUID

- **curl:**
```bash
curl -i http://localhost:3000/api/v1/leagues/0aab480b-9e10-4fea-a6fc-884160e684d5
```
- **Response (200 OK):** คืน Object ข้อมูลลีก
- **Error (404 Not Found):** `{"error": "league not found"}`

---

### PUT /leagues/:id
แก้ไขข้อมูลลีก

- **curl:**
```bash
curl -i -X PUT http://localhost:3000/api/v1/leagues/0aab480b-9e10-4fea-a6fc-884160e684d5 \
  -H "Content-Type: application/json" \
  -d '{"name": "Thai League 1 (Hilux Revo)", "short_name": "T1", "country": "Thailand"}'
```
- **Response (200 OK):** คืน Object ข้อมูลลีกหลังแก้ไข

---

### DELETE /leagues/:id
ลบลีก (Soft Delete)

- **curl:**
```bash
curl -i -X DELETE http://localhost:3000/api/v1/leagues/0aab480b-9e10-4fea-a6fc-884160e684d5
```
- **Response (204 No Content):** ไม่มี Body

---

## 3. Teams (สโมสร/ทีมฟุตบอล)

### GET /teams
ดึงรายชื่อทีมฟุตบอลทั้งหมด พร้อมสังกัดลีกและโลโก้

- **Query Parameters (Optional):**
  - `league_id` (string UUID): กรองเฉพาะทีมที่สังกัดลีกนี้
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
ดึงข้อมูลทีมตาม UUID

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

## 4. Matches (ตารางการแข่งขันและผลบอล)

### GET /matches
ดึงรายการแข่งขันและผลบอล พร้อมข้อมูลทีมเหย้า-เยือน โลโก้ สกอร์ นาที และสถานะ

- **Query Parameters (Optional):**
  - `status`: กรองตามสถานะ (`finished` จบแล้ว, `live` กำลังเตะ, `scheduled` ยังไม่แข่ง)
  - `league_id`: กรองเฉพาะลีกที่ระบุ (UUID)
  - `date`: กรองเฉพาะวันที่แข่ง รูปแบบ `YYYY-MM-DD` (เช่น `2025-04-30`)
  - `limit`: จำนวนรายการ (default: 50, สูงสุด: 200)
  - `page`: หน้าที่ต้องการดึง (default: 1)

- **curl (ดึง 10 แมตช์ล่าสุด):**
```bash
curl -i "http://localhost:3000/api/v1/matches?limit=10"
```
- **curl (ดูเฉพาะผลบอลที่แข่งจบแล้ว):**
```bash
curl -i "http://localhost:3000/api/v1/matches?status=finished&limit=20"
```
- **curl (ดูบอลสดที่กำลังแข่ง):**
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
ดึงข้อมูลแมตช์เฉพาะคู่ตาม UUID

- **curl:**
```bash
curl -i http://localhost:3000/api/v1/matches/f5a892b1-1234-5678-90ab-cdef12345678
```
- **Response (200 OK):** คืน Object ข้อมูลแมตช์พร้อมข้อมูลทีมเหย้าและเยือน
- **Error (404 Not Found):** `{"error": "match not found"}`
