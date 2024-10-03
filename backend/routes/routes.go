package routes

import (
	"multiaura/internal/databases"

	"github.com/gofiber/fiber/v2"
)

var mongoDB *databases.MongoDB

func SetupRoutes(app *fiber.App) {
	mongoDB = databases.Instance()
	setupUserRoutes(app)

}
