package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/vector-ops/ticketeer/config"
	"github.com/vector-ops/ticketeer/db"
	"github.com/vector-ops/ticketeer/handlers"
	"github.com/vector-ops/ticketeer/middleware"
	"github.com/vector-ops/ticketeer/repository"
	"github.com/vector-ops/ticketeer/services"
)

func main() {
	envConfig := config.NewEnvConfig()
	db := db.Init(envConfig, db.DBMigrator)

	app := fiber.New(fiber.Config{
		AppName:      "Ticketeer",
		ServerHeader: "Fiber",
	})

	// repositories
	eventRepo := repository.NewEventRepository(db)
	ticketRepo := repository.NewTicketRepository(db)
	authRepo := repository.NewAuthRepository(db)

	// services
	authService := services.NewAuthService(authRepo)

	// Routes
	router := app.Group("/api/v1")
	handlers.NewAuthHandler(router.Group("/auth"), authService)

	privateRoutes := router.Use(middleware.AuthProtected(db))

	// handelers
	handlers.NewEventHandler(privateRoutes.Group("/event"), eventRepo)
	handlers.NewTicketHandler(privateRoutes.Group("/ticket"), ticketRepo)

	app.Listen(fmt.Sprintf(":" + envConfig.ServerPort))
}
