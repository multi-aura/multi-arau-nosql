package services

import (
	"errors"
	"multiaura/internal/models"
	"multiaura/internal/repositories"
	"multiaura/internal/websocket/group"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ConversationService interface {
	CreateConversation(userIDs []string, name string) (*models.Conversation, error)
	GetConversationByID(id string) (*models.Conversation, error)
	GetListConversations(id string) ([]models.Conversation, error)
	RemoveMenberConversation(ConversationID string, UserID string) error
	AddMembers(conversationID string, userIDs []string) error
	SendMessage(conversationID, userID string, content models.ChatContent) error
	GetMessages(conversationID string) ([]models.Chat, error)
	MarkMessageAsDeleted(conversationID string, messageID string) error
}

type conversationService struct {
	repo            repositories.ConversationRepository
	userRepo        repositories.UserRepository
	websocketGroups map[string]*group.Group
}

func NewConversationService(repo repositories.ConversationRepository, userRepo repositories.UserRepository) ConversationService {
	return &conversationService{
		repo:            repo,
		userRepo:        userRepo,
		websocketGroups: make(map[string]*group.Group),
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
			Username: user.Username,
			Added_at: time.Now(),
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

func (c *conversationService) AddMembers(conversationID string, userIDs []string) error {
	if conversationID == "" {
		return errors.New("conversation ID is required")
	}

	conversation, err := c.repo.GetByID(conversationID)
	if err != nil {
		return err
	}
	if conversation == nil {
		return err

	}

	existingUsers := conversation.Users
	existingUserMap := make(map[string]bool)

	for _, user := range existingUsers {
		existingUserMap[user.UserID] = true
	}

	var newUsers []models.Users

	for _, userID := range userIDs {
		user, err := c.userRepo.GetByID(userID)
		if err != nil {
			return err

		}
		if user == nil {
			return err

		}

		if !existingUserMap[userID] {
			newUser := models.Users{
				UserID:   user.ID,
				Fullname: user.FullName,
				Avatar:   user.Avatar,
				Username: user.Username,
				Added_at: time.Now(),
			}
			newUsers = append(newUsers, newUser)
		}
	}

	if len(newUsers) == 0 {
		return nil
	}

	err = c.repo.AddMemberToConversation(newUsers, conversationID)
	if err != nil {
		return err

	}

	return nil
}
func (c *conversationService) RemoveMenberConversation(ConversationID string, UserID string) error {
	if ConversationID == "" {
		return errors.New("no conversation ID specified")
	}
	if UserID == "" {
		return errors.New("no user ID specified")
	}
	conversation, err := c.repo.GetByID(ConversationID)
	if err != nil {
		return errors.New("No conversation with ID " + ConversationID)
	}
	var UsersUpdate []models.Users
	userFound := false
	for _, user := range conversation.Users {
		if user.UserID == UserID {
			userFound = true
			continue
		}
		UsersUpdate = append(UsersUpdate, user)

	}
	if !userFound {
		return errors.New("User not found in conversation list.")

	}
	conversation.Users = UsersUpdate
	conversation.UpdatedAt = time.Now()
	err = c.repo.UpdateRemoveruser(conversation)
	if err != nil {
		return errors.New("Failed to update conversation")
	}
	return nil
}
func (cs *conversationService) SendMessage(conversationID, userID string, content models.ChatContent) error {
	user, err := cs.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	newMessage := models.Chat{
		ID: primitive.NewObjectID(),
		Sender: models.Users{
			UserID:   user.ID,
			Fullname: user.FullName,
			Avatar:   user.Avatar,
		},
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    "sent",
	}

	// Lưu tin nhắn vào database
	err = cs.repo.AddMessageToConversation(newMessage, conversationID)
	if err != nil {
		return err
	}

	// // Phát tin nhắn qua WebSocket tới tất cả các client trong Group
	// messageData, err := json.Marshal(newMessage)
	// if err != nil {
	// 	return err
	// }

	// if group, ok := cs.websocketGroups[conversationID]; ok {
	// 	log.Println("Broadcasting message to WebSocket group:", conversationID)
	// 	group.BroadcastMessage(messageData)
	// } else {
	// 	log.Println("No WebSocket group found for conversationID:", conversationID)
	// }

	return nil
}
func (s *conversationService) GetMessages(conversationID string) ([]models.Chat, error) {
	return s.repo.GetMessagesByConversationID(conversationID)
}

func (s *conversationService) MarkMessageAsDeleted(conversationID string, messageID string) error {
	return s.repo.MarkMessageAsDeleted(conversationID, messageID)
}
