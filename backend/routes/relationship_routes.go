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
	service := services.NewRelationshipService(&repository)
	controller := controllers.NewRelationshipController(service)

	relationships := app.Group("/relationships")

	relationships.Post("/status/:userID", middlewares.AuthMiddleware(), controller.GetRelationshipStatus)
	relationships.Post("/follow/:userID", middlewares.AuthMiddleware(), controller.Follow)
	relationships.Delete("/unfollow/:userID", middlewares.AuthMiddleware(), controller.UnFollow)
	relationships.Post("/block/:userID", middlewares.AuthMiddleware(), controller.Block)
	relationships.Delete("/unblock/:userID", middlewares.AuthMiddleware(), controller.UnBlock)
	relationships.Get("/friends", middlewares.AuthMiddleware(), controller.GetFriends)
	relationships.Get("/followers", middlewares.AuthMiddleware(), controller.GetFollowers)
	relationships.Get("/followings", middlewares.AuthMiddleware(), controller.GetFollowings)
	relationships.Get("/blocked", middlewares.AuthMiddleware(), controller.GetBlockedUsers)
	app.Get("/:username", middlewares.AuthMiddleware(), controller.GetProfile)
}
