package services

import (
	"errors"
	"fmt"
	"multiaura/internal/models"
	"multiaura/internal/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ConversationService interface {
	CreateConversation(userIDs []string, name string) (*models.Conversation, error)
	GetConversationByID(id string) (*models.Conversation, error)
	GetListConversations(id string) ([]models.Conversation, error)
	AddMember(conversationID string, userID []string) error
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

	if len(userIDs) < 2 {
		return nil, errors.New("at least two users are required to create a conversation")
	} else if len(userIDs) > 2 {
		ConversationType = "Group"
	}

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

	newConversation := models.Conversation{
		ID:               primitive.NewObjectID(),
		Name:             name,
		ConversationType: ConversationType,
		Users:            users,
		Chats:            []models.Chat{},
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

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

func (c *conversationService) GetListConversations(id string) ([]models.Conversation, error) {
	if id == "" {
		return nil, errors.New("id not found")
	}
	listConversation, err := c.repo.GetListConversations(id)
	if err != nil {
		return nil, errors.New("error getting list of conversations")
	}
	if len(listConversation) == 0 {
		return nil, errors.New("no conversations")
	}
	return listConversation, nil

}

func (c *conversationService) AddMember(conversationID string, userIDs []string) error {

	var users []models.Users
	for _, userID := range userIDs {

		user, err := c.userRepo.GetByID(userID)
		if err != nil {
			return err
		}

		if user == nil {
			return fmt.Errorf("user with ID %s not found", userID)
		}

		users = append(users, models.Users{
			UserID:   user.ID,
			Fullname: user.FullName,
			Avatar:   user.Avatar,
			LastSeen: time.Now(),
		})
	}

	err := c.repo.AddMemberToConversation(users, conversationID)
	if err != nil {
		return err
	}

	return nil
}
