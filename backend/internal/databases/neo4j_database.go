package databases

import (
	"context"
	"log"
	config "multiaura/internal/configs/dev"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jDB struct {
	Driver neo4j.DriverWithContext
}

var neo4jInstance *Neo4jDB

func Neo4jInstance() *Neo4jDB {
	return neo4jInstance
}

// Khởi tạo kết nối Neo4j
func NewNeo4jDB(cfg *config.Neo4jConfig) (*Neo4jDB, error) {
	if neo4jInstance != nil {
		return neo4jInstance, nil
	}

	// Tạo driver kết nối đến Neo4j
	driver, err := neo4j.NewDriverWithContext(cfg.URI, neo4j.BasicAuth(cfg.Username, cfg.Password, ""))
	if err != nil {
		return nil, err
	}

	neo4jInstance = &Neo4jDB{
		Driver: driver,
	}

	return neo4jInstance, nil
}

// Disconnect Neo4j
func (db *Neo4jDB) Disconnect() {
	if db.Driver != nil {
		err := db.Driver.Close(context.Background())
		if err != nil {
			log.Printf("Failed to disconnect Neo4j connection: %v", err)
		} else {
			log.Println("Neo4j connection disconnected!")
		}
	}
}
