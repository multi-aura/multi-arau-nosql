package routes

import (
	"multiaura/internal/databases"

	"github.com/gofiber/fiber/v2"
)

var neo4jDB *databases.Neo4jDB
var mongoDB *databases.MongoDB

func SetupRoutes(app *fiber.App) {
	neo4jDB = databases.Neo4jInstance()
	mongoDB = databases.MongoInstance()
	SetupUserRoutes(app)
	SetupRelationshipRoutes(app)
}
