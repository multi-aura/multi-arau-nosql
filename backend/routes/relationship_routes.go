package routes

import (
	"multiaura/internal/controllers"
	"multiaura/internal/middlewares"
	"multiaura/internal/repositories"
	"multiaura/internal/services"

	"github.com/gofiber/fiber/v2"
)

func SetupRelationshipRoutes(app *fiber.App) {
	repository := repositories.NewUserRepository(neo4jDB)
	service := services.NewRelationshipService(repository)
	controller := controllers.NewRelationshipController(service)

	relationships := app.Group("/relationships")

	// relationships.Post("/follow/:userID", controller.Follow)
	// relationships.Delete("/unfollow/:userID", controller.Unfollow)
	// relationships.Post("/block/:userID", controller.Block)
	// relationships.Delete("/unblock/:userID", controller.Unblock)
	relationships.Get("/friends", middlewares.AuthMiddleware(), controller.GetFriends)
}
