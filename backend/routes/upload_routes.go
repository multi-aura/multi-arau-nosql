package routes

import (
	"multiaura/internal/controllers"
	"multiaura/internal/middlewares"
	"multiaura/internal/repositories"
	"multiaura/internal/services"

	"github.com/gofiber/fiber/v2"
)

func SetupUploadRoutes(app *fiber.App) {
	storageRepository := repositories.NewStorageRepository()
	userRepository := repositories.NewUserRepository(neo4jDB)
	postRepository := repositories.NewPostRepository(mongoDB)
	service := services.NewUploadService(&userRepository, &postRepository, &storageRepository)
	controller := controllers.NewUploadController(service)

	uploadGroup := app.Group("/upload")

	uploadGroup.Post("/profile-photo", middlewares.AuthMiddleware(), controller.UploadProfilePhoto)
	uploadGroup.Post("/post-photos/:postID", middlewares.AuthMiddleware(), controller.UploadPostPhotos)
}
