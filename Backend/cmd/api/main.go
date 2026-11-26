package main

import (
	"errors"
	"log"

	"football/internal/config"
	"football/internal/database"
	"football/internal/handlers"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	log.Println("database connected")

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			code, msg := fiber.StatusInternalServerError, "internal server error"
			var fe *fiber.Error
			if errors.As(err, &fe) {
				code, msg = fe.Code, fe.Message
			} else {
				log.Printf("error: %v", err)
			}
			return c.Status(code).JSON(fiber.Map{"error": msg})
		},
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	api := app.Group("/api/v1")

	// 1. Health check
	api.Get("/health", func(c fiber.Ctx) error {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			return fiber.NewError(fiber.StatusServiceUnavailable, "database unavailable")
		}
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// 2. Seasons (สำหรับ Dropdown เลือกปี)
	matches := handlers.NewMatchHandler(db)
	api.Get("/seasons", matches.Seasons)

	// 3. Leagues (อ่านข้อมูลลีก)
	leagues := handlers.NewLeagueHandler(db)
	lg := api.Group("/leagues")
	lg.Get("/", leagues.List)
	lg.Get("/:id", leagues.Get)

	// 4. Teams (ข้อมูลทีมและสโมสร)
	teams := handlers.NewTeamHandler(db)
	tg := api.Group("/teams")
	tg.Get("/", teams.List)
	tg.Get("/:id", teams.Get)

	// 5. Matches (ตารางแข่งและผลบอล พร้อม filter ปี/วัน/สถานะ)
	mg := api.Group("/matches")
	mg.Get("/", matches.List)
	mg.Get("/:id", matches.Get)

	// 6. Sync on-demand (สั่งดูดข้อมูลจาก API-Football ลง Database)
	syncHandler := handlers.NewSyncHandler(db, cfg)
	api.Post("/sync/fixtures", syncHandler.SyncFixtures)

	log.Fatal(app.Listen(":" + cfg.Port))
}
