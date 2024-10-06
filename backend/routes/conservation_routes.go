package routes

import (
	"multiaura/internal/controllers"
	"multiaura/internal/middlewares"
	"multiaura/internal/repositories"
	"multiaura/internal/services"

	"github.com/gofiber/fiber/v2"
)

func SetupConservationRoutes(app *fiber.App) {
	repository := repositories.NewConservationRepository(mongoDB)
	service := services.NewConservationService(repository)
	controller := controllers.NewConservationController(service)

	relationships := app.Group("/conservations")

	relationships.Post("/chat/:userID", middlewares.AuthMiddleware(), controller.CreateConservation)
}
