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
	search.Get("/people", middlewares.AuthMiddleware(), controller.SearchPeople)
	search.Get("/people/:q", middlewares.AuthMiddleware(), controller.SearchPeople)
	search.Get("/posts/:q", middlewares.AuthMiddleware(), controller.SearchPosts)
	search.Get("/trending/:q", middlewares.AuthMiddleware(), controller.SearchTrending)
	search.Get("/for-you/:q", middlewares.AuthMiddleware(), controller.SearchForYou)
	search.Get("/news/:q", middlewares.AuthMiddleware(), controller.SearchNews)
}
