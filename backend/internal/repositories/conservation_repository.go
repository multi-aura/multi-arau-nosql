package repositories

import (
	"context"
	"fmt"
	"log"
	"multiaura/internal/databases"
	"multiaura/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ConversationRepository interface {
	Repository[models.Conversation]
	GetListConversations(userID string) ([]models.Conversation, error)
	AddMemberToConversation(user []models.Users, id_conversation string) error
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

func (repo *conversationRepository) GetListConversations(userID string) ([]models.Conversation, error) {
	var conversations []models.Conversation

	filter := bson.M{"users.user_id": userID}

	cursor, err := repo.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("error finding conversations: %v", err)
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &conversations); err != nil {
		return nil, fmt.Errorf("error decoding conversations: %v", err)
	}

	return conversations, nil
}
func (repo *conversationRepository) AddMemberToConversation(users []models.Users, id_conversation string) error {
	id_conversationRepository, err := primitive.ObjectIDFromHex(id_conversation)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": id_conversationRepository}

	update := bson.M{
		"$push": bson.M{
			"users": users,
		},
		"$set": bson.M{
			"updatedat": time.Now().UTC(),
		},
	}
	_, err = repo.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}
