package routes

import (
	"multiaura/internal/controllers"
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

	relationships.Get("/messages/:conversationID", controller.GetConversationByID)
	relationships.Post("/createmessages", controller.CreateConversation)
	relationships.Get("/ListMessages/:UserID", controller.GetListConversation)
	relationships.Post("/AddMenberConversation/:conversationID/members/:userID", controller.AddMember)

}
