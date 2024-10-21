package routes

import (
	"multiaura/internal/controllers"
	"multiaura/internal/middlewares"
	"multiaura/internal/repositories"
	"multiaura/internal/services"

	"github.com/gofiber/fiber/v2"
)

func SetupSearchRoutes(app *fiber.App) {
	userRepository := repositories.NewUserRepository(neo4jDB)
	postRepository := repositories.NewPostRepository(mongoDB)
	service := services.NewSearchService(&userRepository, &postRepository)
	controller := controllers.NewSearchController(service)

	search := app.Group("/search")
	search.Post("/people", middlewares.AuthMiddleware(), controller.SearchPeople)
	search.Post("/people/:q", middlewares.AuthMiddleware(), controller.SearchPeople)

	search.Post("/posts", middlewares.AuthMiddleware(), controller.SearchPosts)
	search.Post("/posts/:q", middlewares.AuthMiddleware(), controller.SearchPosts)

	search.Post("/trending", middlewares.AuthMiddleware(), controller.SearchTrending)
	search.Post("/trending/:q", middlewares.AuthMiddleware(), controller.SearchTrending)

	search.Post("/for-you", middlewares.AuthMiddleware(), controller.SearchForYou)
	search.Post("/for-you/:q", middlewares.AuthMiddleware(), controller.SearchForYou)

	search.Post("/news", middlewares.AuthMiddleware(), controller.SearchNews)
	search.Post("/news/:q", middlewares.AuthMiddleware(), controller.SearchNews)
}
