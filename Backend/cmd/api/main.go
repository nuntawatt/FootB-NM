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
	"gorm.io/gorm"
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
		ErrorHandler: jsonErrorHandler,
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	registerRoutes(app, db, cfg)

	log.Fatal(app.Listen(":" + cfg.Port))
}

func registerRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config) {
	api := app.Group("/api/v1")

	// Health check
	api.Get("/health", func(c fiber.Ctx) error {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			return fiber.NewError(fiber.StatusServiceUnavailable, "database unavailable")
		}
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Seasons
	matches := handlers.NewMatchHandler(db)
	api.Get("/seasons", matches.Seasons)

	// Leagues
	leagues := handlers.NewLeagueHandler(db)
	lg := api.Group("/leagues")
	lg.Get("/", leagues.List)
	lg.Get("/:id", leagues.Get)

	// Teams
	teams := handlers.NewTeamHandler(db)
	tg := api.Group("/teams")
	tg.Get("/", teams.List)
	tg.Get("/:id", teams.Get)

	// Matches
	mg := api.Group("/matches")
	mg.Get("/", matches.List)
	mg.Get("/:id", matches.Get)

	// Sync
	sync := handlers.NewSyncHandler(db, cfg)
	api.Post("/sync/fixtures", sync.SyncFixtures)
}

func jsonErrorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := "internal server error"

	var fe *fiber.Error
	if errors.As(err, &fe) {
		code = fe.Code
		msg = fe.Message
	} else {
		log.Printf("unhandled error: %v", err)
	}

	return c.Status(code).JSON(fiber.Map{"error": msg})
}
