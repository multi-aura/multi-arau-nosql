package repositories

import (
	"context"
	"multiaura/internal/databases"
	"multiaura/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ConservationRepository interface {
	Repository[models.User]
	// GetUsersByName(name string) ([]models.User, error)
	// GetUserByEmail(email string) (*models.User, error)
}

type conservationRepository struct {
	db         *databases.MongoDB
	collection *mongo.Collection
}

func NewConservationRepository(db *databases.MongoDB) ConservationRepository {
	return &conservationRepository{
		db:         db,
		collection: db.Database.Collection("chats"),
	}
}

// GetByID implements ConservationRepository.
func (repo *conservationRepository) GetByID(id string) (*models.User, error) {
	var user models.User

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &models.User{}, err
	}

	filter := bson.M{"_id": objectID}

	err = repo.collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// Create implements ConservationRepository.
func (repo *conservationRepository) Create(entity models.User) error {
	_, err := repo.collection.InsertOne(context.Background(), entity)
	return err
}

// Delete implements ConservationRepository.
func (repo *conservationRepository) Delete(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}

	result, err := repo.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// Update implements ConservationRepository.
func (repo *conservationRepository) Update(entityMap *map[string]interface{}) error {
	filter := bson.M{"_id": (*entityMap)["userID"].(string)}

	updateQuery := bson.M{"$set": entityMap}

	result, err := repo.collection.UpdateOne(context.Background(), filter, updateQuery)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
