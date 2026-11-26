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

	// health check
	api.Get("/health", func(c fiber.Ctx) error {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			return fiber.NewError(fiber.StatusServiceUnavailable, "database unavailable")
		}
		return c.JSON(fiber.Map{"status": "ok"})
	})

	leagues := handlers.NewLeagueHandler(db)
	g := api.Group("/leagues")
	g.Get("/", leagues.List)
	g.Post("/", leagues.Create)
	g.Get("/:id", leagues.Get)
	g.Put("/:id", leagues.Update)
	g.Delete("/:id", leagues.Delete)

	log.Fatal(app.Listen(":" + cfg.Port))
}
