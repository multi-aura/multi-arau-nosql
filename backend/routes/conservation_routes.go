package routes

import (
	"multiaura/internal/controllers"
	"multiaura/internal/middlewares"
	"multiaura/internal/repositories"
	"multiaura/internal/services"

	"github.com/gofiber/fiber/v2"
)

func SetupConversationRoutes(app *fiber.App) {
	repository := repositories.NewConversationRepository(mongoDB)
	Userrepository := repositories.NewUserRepository(neo4jDB)

	service := services.NewConversationService(repository, Userrepository)
	controller := controllers.NewConversationController(service)

	relationships := app.Group("/conservations")

	relationships.Post("/chat/:userID", middlewares.AuthMiddleware(), controller.CreateConversation)
	relationships.Get("/messages/:conversationID", controller.GetConversationByID)

}
