# System Run Guide (คู่มือขั้นตอนการรันระบบ)

คู่มือสำหรับนักพัฒนา Backend สำหรับการติดตั้ง Dependencies, การรัน Database Migration, และคำสั่งเริ่มต้นใช้งานระบบทั้งหมด

---

## 1. การติดตั้ง Go Dependencies & Packages (Go Get)

หากต้องการดาวน์โหลดและติดตั้ง Packages ทั้งหมดที่ระบบใช้งาน สามารถใช้คำสั่งเดียว:

```bash
cd Backend
go mod download
go mod tidy
```

### รายการ Packages หลักที่ระบบใช้งาน (`go get`):
* **Web Framework:**
  ```bash
  go get github.com/gofiber/fiber/v3
  ```
* **Database ORM & PostgreSQL Driver:**
  ```bash
  go get gorm.io/gorm
  go get gorm.io/driver/postgres
  ```
* **Redis Client:**
  ```bash
  go get github.com/redis/go-redis/v9
  ```
* **Object Storage SDK (MinIO / S3):**
  ```bash
  go get github.com/minio/minio-go/v7
  ```
* **JWT Authentication:**
  ```bash
  go get github.com/golang-jwt/jwt/v5
  ```
* **Password Hashing (bcrypt):**
  ```bash
  go get golang.org/x/crypto/bcrypt
  ```
* **UUID Generator:**
  ```bash
  go get github.com/google/uuid
  ```
* **Environment Variable (.env):**
  ```bash
  go get github.com/joho/godotenv
  ```

---

## 2. การเตรียมไฟล์คอนฟิก (.env)

คัดลอกไฟล์คอนฟิกเริ่มต้น (หากยังไม่มีไฟล์ `.env`):

```bash
cd Backend
cp .env.example .env
```

---

## 3. คำสั่งรัน Database Migration

ไฟล์ Migration ทั้งหมดจะอยู่ในโฟลเดอร์ `Backend/migrations/`:
1. `000001_init_schema.up.sql` (สร้างตารางผู้ใช้ ลีก สโมสร และแมตช์)
2. `000002_design_additions.up.sql` (เพิ่มฟิลด์ดีไซน์)
3. `000003_create_coaches_and_players.up.sql` (สร้างตารางโค้ชและนักเตะ)

### วิธีที่ 1: รันด้วยเครื่องมือ `golang-migrate` CLI (แนะนำ)

ติดตั้ง CLI (หากยังไม่ได้ติดตั้ง):
* macOS (Homebrew): `brew install golang-migrate`
* Go Install: `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`

**คำสั่งรัน Migration ขึ้นทั้งหมด (Migrate Up):**
```bash
migrate -path migrations -database "postgres://postgres:nuntawat838@localhost:5432/football?sslmode=disable" up
```

**คำสั่ง Rollback ถอยหลัง 1 สเต็ป (Migrate Down):**
```bash
migrate -path migrations -database "postgres://postgres:nuntawat838@localhost:5432/football?sslmode=disable" down 1
```

---

### วิธีที่ 2: รันผ่าน Docker Compose Exec (ไม่ต้องติดตั้ง CLI ในเครื่อง)

หากกำลังรัน PostgreSQL ผ่าน Docker สามารถสั่งรัน SQL เข้าไปในคอนเทนเนอร์ได้ทันที:

```bash
cat migrations/000001_init_schema.up.sql | docker compose exec -T postgres psql -U postgres -d football
cat migrations/000002_design_additions.up.sql | docker compose exec -T postgres psql -U postgres -d football
cat migrations/000003_create_coaches_and_players.up.sql | docker compose exec -T postgres psql -U postgres -d football
```

---

### วิธีที่ 3: รันด้วยคำสั่ง `psql` บนเครื่องโดยตรง

```bash
PGPASSWORD=nuntawat838 psql -h localhost -U postgres -d football -f migrations/000001_init_schema.up.sql
PGPASSWORD=nuntawat838 psql -h localhost -U postgres -d football -f migrations/000002_design_additions.up.sql
PGPASSWORD=nuntawat838 psql -h localhost -U postgres -d football -f migrations/000003_create_coaches_and_players.up.sql
```

---

## 4. วิธีการรันระบบ (Run Application)

### วิธีที่ 1: รันทั้งระบบด้วย Docker Compose (แนะนำ)
เปิดทำงานครบทั้ง **API + PostgreSQL + Redis + MinIO (S3 & Web Console)** ในคำสั่งเดียว:

```bash
cd Backend

# สั่งรันเซอร์วิสทั้งหมดในโหมด Background
docker compose up -d --build

# ตรวจสอบสถานะการทำงาน
docker compose ps

# ดู Logs ของ Backend API แบบเรียลไทม์
docker compose logs -f api
```

---

### วิธีที่ 2: รันแบบ Development (รัน Go ในเครื่อง)
สำหรับช่วงพัฒนาโค้ชที่ต้องการความรวดเร็ว โดยรันเฉพาะ Database, Cache และ Storage ใน Docker:

```bash
cd Backend

# 1. เปิดเซอร์วิส Infrastructure
docker compose up -d postgres redis minio

# 2. รัน Database Migration (ตามข้อ 3)
migrate -path migrations -database "postgres://postgres:nuntawat838@localhost:5432/football?sslmode=disable" up

# 3. สั่งรัน Go API Server
go run ./cmd/api/main.go
```

---

## 5. ช่องทางการเข้าใช้งานระบบ (Service Endpoints & GUI)

| เซอร์วิส | URL / Port | บัญชีเข้าสู่ระบบ | รายละเอียด |
|:---|:---|:---:|:---|
| **Backend API** | `http://localhost:3000/api/v1` | - | RESTful API ทุก Endpoint |
| **Health Check** | `http://localhost:3000/api/v1/health` | - | ตรวจสอบสถานะ Server & Database |
| **MinIO Web GUI** | `http://localhost:9001` | User: `minioadmin`<br>Pass: `minioadmin` | หน้าเว็บดูและพรีวิวรูปภาพโปรไฟล์ |
| **MinIO S3 API** | `http://localhost:9000` | - | Object Storage API |
| **PostgreSQL** | `localhost:5432` | User: `postgres`<br>Pass: `nuntawat838` | Database: `football` |
| **Redis** | `localhost:6379` | - | Memory Cache |

---

## 6. คำสั่งดึงข้อมูลฟุตบอลจริง (Data Sync - Optional)

หากต้องการดึงข้อมูลสโมสร โค้ช นักเตะ และโปรแกรมแข่งขันจริงจาก Football API เข้าสู่ฐานข้อมูล:

```bash
cd Backend

# 1. ดูดข้อมูลโค้ชและนักเตะ 20 สโมสร
go run ./cmd/sync squads

# 2. ดูดโปรแกรมและผลการแข่งขัน
go run ./cmd/sync fixtures -season 2024
```

---

## 7. คำสั่งหยุดการทำงาน (Stop Services)

```bash
cd Backend

# หยุดคอนเทนเนอร์ทั้งหมด
docker compose down

# หยุดคอนเทนเนอร์พร้อมลบ Volume ข้อมูล (ใช้เมื่อต้องการ Reset ข้อมูลใหม่หมด)
# docker compose down -v
```
