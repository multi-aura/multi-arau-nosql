package repositories

import (
	"context"
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
	UpdateRemoveruser(conversation *models.Conversation) error
	AddMemberToConversation(user []models.Users, id_conversation string) error
	AddMessageToConversation(message models.Chat, conversationID string) error
	GetMessagesByConversationID(conversationID string) ([]models.Chat, error)
	MarkMessageAsDeleted(conversationID string, messageID string) error
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
		return nil, err
	}

	// Truy vấn MongoDB dựa trên ObjectID
	filter := bson.M{"_id": objectID}
	err = repo.collection.FindOne(context.Background(), filter).Decode(&conversation)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.Conversation{}, nil
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
	filter := bson.M{"_id": (*entityMap)["_id"].(string)}

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
func (repo *conversationRepository) UpdateRemoveruser(conversation *models.Conversation) error {
	filter := bson.M{"_id": conversation.ID}

	update := bson.M{
		"$set": bson.M{
			"users":     conversation.Users,
			"updatedat": conversation.UpdatedAt,
		},
	}

	result, err := repo.collection.UpdateOne(context.Background(), filter, update)
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

	filter := bson.M{"users.userID": userID}

	cursor, err := repo.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &conversations); err != nil {
		return nil, err
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
			"users": bson.M{
				"$each": users, // Thêm từng phần tử trong mảng users
			},
		},
		"$set": bson.M{
			"updatedat": time.Now().UTC(), // Cập nhật thời gian sửa đổi
		},
	}

	_, err = repo.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
func (repo *conversationRepository) AddMessageToConversation(message models.Chat, conversationID string) error {
	conversationObjectID, err := primitive.ObjectIDFromHex(conversationID)
	if err != nil {
		return err
	}

	// Thêm tin nhắn vào mảng "chats" trong cuộc trò chuyện
	filter := bson.M{"_id": conversationObjectID}
	update := bson.M{
		"$push": bson.M{"chats": message},
		"$set":  bson.M{"updatedat": message.UpdatedAt},
	}

	_, err = repo.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
func (r *conversationRepository) GetMessagesByConversationID(conversationID string) ([]models.Chat, error) {
	objectID, err := primitive.ObjectIDFromHex(conversationID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	var conversation models.Conversation

	err = r.collection.FindOne(context.Background(), filter).Decode(&conversation)
	if err != nil {
		return nil, err
	}

	return conversation.Chats, nil
}

func (r *conversationRepository) MarkMessageAsDeleted(conversationID string, messageID string) error {
	conversationObjectID, err := primitive.ObjectIDFromHex(conversationID)
	if err != nil {
		return err
	}

	messageObjectID, err := primitive.ObjectIDFromHex(messageID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": conversationObjectID, "chats.id_chat": messageObjectID}
	update := bson.M{"$set": bson.M{"chats.$.status": "deleted"}}

	result, err := r.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err

	}
	if result.ModifiedCount == 0 {
		log.Println("No document was updated.")
	}
	return nil
}
