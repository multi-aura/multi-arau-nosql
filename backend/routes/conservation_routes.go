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

	relationships.Get("/conversation-by-id/:conversationID", controller.GetConversationByID)
	relationships.Post("/create-messages", controller.CreateConversation)
	relationships.Get("/list-messages/:userID", controller.GetListConversation)
	relationships.Post("/add-member-conversation/:conversationID", controller.AddMember)
	relationships.Delete("/remove-member-conversation/:conversationID/Menber/:userID", controller.RemoveMenberConversation)
	relationships.Post("/send-message/:conversationID", controller.SendMessage)
	relationships.Get("/messages-chat/:conversationID", controller.GetMessages)
	relationships.Put("/messages-delete/:conversationID/:messageID", controller.MarkMessageAsDeleted)
}
