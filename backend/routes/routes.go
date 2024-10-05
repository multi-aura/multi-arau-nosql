package routes

import (
	"multiaura/internal/databases"

	"github.com/gofiber/fiber/v2"
)

var neo4jDB *databases.Neo4jDB

func SetupRoutes(app *fiber.App) {
	neo4jDB = databases.Instance()
	SetupUserRoutes(app)
	SetupRelationshipRoutes(app)
}
