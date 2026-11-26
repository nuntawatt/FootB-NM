// Command api starts the Football API HTTP server.
//
// All configuration is read from environment variables (see internal/config).
// Run with:
//
//	go run ./cmd/api
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
	// 1. Load configuration from environment / .env file.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	// 2. Open the PostgreSQL connection pool.
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	log.Println("database connected")

	// 3. Create the Fiber application.
	//    The custom error handler ensures every error is returned as JSON,
	//    including panics caught by the recover middleware.
	app := fiber.New(fiber.Config{
		ErrorHandler: jsonErrorHandler,
	})

	// 4. Global middleware (applied in order):
	//    recover -> logger -> cors
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	// 5. Register all routes.
	registerRoutes(app, db, cfg)

	// 6. Start listening.
	log.Fatal(app.Listen(":" + cfg.Port))
}

// registerRoutes mounts every route group under /api/v1.
func registerRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config) {
	api := app.Group("/api/v1")

	// Health check — also verifies the database is reachable.
	api.Get("/health", func(c fiber.Ctx) error {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			return fiber.NewError(fiber.StatusServiceUnavailable, "database unavailable")
		}
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Seasons — dropdown data for the frontend.
	matches := handlers.NewMatchHandler(db)
	api.Get("/seasons", matches.Seasons)

	// Leagues — read-only competition list.
	leagues := handlers.NewLeagueHandler(db)
	lg := api.Group("/leagues")
	lg.Get("/", leagues.List)
	lg.Get("/:id", leagues.Get)

	// Teams — clubs filtered optionally by league.
	teams := handlers.NewTeamHandler(db)
	tg := api.Group("/teams")
	tg.Get("/", teams.List)
	tg.Get("/:id", teams.Get)

	// Matches — fixture list with rich filtering, plus single match detail.
	mg := api.Group("/matches")
	mg.Get("/", matches.List)
	mg.Get("/:id", matches.Get)

	// Sync — on-demand pull from API-Football into the database.
	sync := handlers.NewSyncHandler(db, cfg)
	api.Post("/sync/fixtures", sync.SyncFixtures)
}

// jsonErrorHandler converts any error into a consistent JSON response.
// fiber.Error values carry an explicit HTTP status code; all other errors
// are treated as 500 Internal Server Error and logged.
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
