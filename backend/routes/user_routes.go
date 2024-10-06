package routes

import (
	"multiaura/internal/controllers"
	"multiaura/internal/middlewares"
	"multiaura/internal/repositories"
	"multiaura/internal/services"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App) {
	repository := repositories.NewUserRepository(neo4jDB)
	service := services.NewUserService(repository)
	controller := controllers.NewUserController(service)

	userGroup := app.Group("/user")

	userGroup.Post("/register", controller.Register)
	userGroup.Post("/login", controller.Login)
	userGroup.Delete("/delete", middlewares.AuthMiddleware(), controller.DeleteUser)
	userGroup.Put("/update", middlewares.AuthMiddleware(), controller.UpdateUser)
}
