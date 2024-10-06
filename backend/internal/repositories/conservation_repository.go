package repositories

import (
	"context"
	"log"
	"multiaura/internal/databases"
	"multiaura/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ConversationRepository interface {
	Repository[models.Conversation]
}

type conversationRepository struct {
	db         *databases.MongoDB
	collection *mongo.Collection
}

func NewConversationRepository(db *databases.MongoDB) ConversationRepository {
	if db == nil || db.Database == nil {
		log.Fatal("MongoDB instance or database is nil")
	}

	return &conversationRepository{
		db:         db,
		collection: db.Database.Collection("chats"),
	}
}

func (repo *conversationRepository) GetByID(id string) (*models.Conversation, error) {
	var conversation models.Conversation

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &models.Conversation{}, err
	}

	filter := bson.M{"_id": objectID}

	err = repo.collection.FindOne(context.Background(), filter).Decode(&conversation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &conversation, nil
}

func (repo *conversationRepository) Create(conversation models.Conversation) error {
	_, err := repo.collection.InsertOne(context.Background(), conversation)
	return err
}

func (repo *conversationRepository) Delete(id string) error {
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

func (repo *conversationRepository) Update(entityMap *map[string]interface{}) error {
	filter := bson.M{"_id": (*entityMap)["user_id"].(string)}

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
