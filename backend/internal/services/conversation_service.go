package services

import (
	"errors"
	"multiaura/internal/models"
	"multiaura/internal/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ConversationService interface {
	CreateConversation(userIDs []string, name string) (*models.Conversation, error)
	GetConversationByID(id string) (*models.Conversation, error)
}

type conversationService struct {
	repo     repositories.ConversationRepository
	userRepo repositories.UserRepository
}

func NewConversationService(repo repositories.ConversationRepository, userRepo repositories.UserRepository) ConversationService {
	return &conversationService{
		repo:     repo,
		userRepo: userRepo, // Gán userRepo vào service
	}
}

// CreateConversation implements ConversationService.
func (c *conversationService) CreateConversation(userIDs []string, name string) (*models.Conversation, error) {
	ConversationType := "Private"

	// Kiểm tra số lượng userIDs để xác định loại cuộc trò chuyện
	if len(userIDs) < 2 {
		return nil, errors.New("at least two users are required to create a conversation")
	} else if len(userIDs) > 2 {
		ConversationType = "Group"
	}

	// Tạo slice chứa các đối tượng models.Users
	var users []models.Users
	for _, id := range userIDs {
		user, err := c.userRepo.GetByID(id)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, errors.New("user not found")
		}

		users = append(users, models.Users{
			UserID:   user.ID,
			Fullname: user.FullName,
			Avatar:   user.Avatar,
			LastSeen: time.Now(),
		})
	}

	// Tạo cuộc trò chuyện mới
	newConversation := models.Conversation{
		ID:               primitive.NewObjectID(),
		Name:             name,
		ConversationType: ConversationType,
		Users:            users,
		Chats:            []models.Chat{},
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Lưu cuộc trò chuyện vào cơ sở dữ liệu
	err := c.repo.Create(newConversation)
	if err != nil {
		return nil, errors.New("failed to create conversation")
	}

	return &newConversation, nil
}

func (c *conversationService) GetConversationByID(id string) (*models.Conversation, error) {
	if id == "" {
		return nil, errors.New("ID not found")
	}

	conversation, err := c.repo.GetByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return conversation, nil
}
