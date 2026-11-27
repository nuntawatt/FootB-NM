package main

import (
	"errors"
	"log"

	"football/internal/redis"
	"football/internal/config"
	"football/internal/database"
	"football/internal/handlers"
	"football/internal/middleware"

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

	// start redis
	rCache := cache.New(cfg.RedisURL)

	app := fiber.New(fiber.Config{
		ErrorHandler: jsonErrorHandler,
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	registerRoutes(app, db, cfg, rCache)

	log.Fatal(app.Listen(":" + cfg.Port))
}

func registerRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config, rCache *cache.Cache) {
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
	sync := handlers.NewSyncHandler(db, cfg, rCache)
	api.Post("/sync/fixtures", sync.SyncFixtures)
	api.Post("/sync/squads", sync.SyncSquads)

	// Standings (ตารางคะแนน)
	standings := handlers.NewStandingsHandler(db, rCache)
	api.Get("/standings", standings.Get)

	// User & Profile
	users := handlers.NewUserHandler(db, cfg)
	protected := middleware.Protected(cfg.JWTSecret)
	optAuth := middleware.OptionalAuth(cfg.JWTSecret)

	api.Post("/register", users.Register)
	api.Post("/login", users.Login)
	api.Get("/profile", protected, users.Profile)

	// Favorites (ทีมโปรด - ดึง user จาก JWT อัตโนมัติ หรือ fallback เป็น X-User-ID)
	favorites := handlers.NewFavoriteHandler(db)
	fav := api.Group("/favorites", optAuth)
	fav.Get("/", favorites.List)
	fav.Post("/:team_id", favorites.Add)
	fav.Delete("/:team_id", favorites.Remove)

	// Extra Features (Form, Coach, Squad, Lineups, Venue, Heroes)
	extra := handlers.NewExtraHandler(db)
	api.Get("/heroes", extra.Heroes)
	api.Get("/teams/:id/form", extra.TeamForm)
	api.Get("/teams/:id/coach", extra.TeamCoach)
	api.Get("/teams/:id/squad", extra.TeamSquad)
	api.Get("/matches/:id/lineups", extra.MatchLineups)
	api.Get("/matches/:id/venue", extra.MatchVenue)
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
