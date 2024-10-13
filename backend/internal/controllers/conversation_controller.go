package controllers

import (
	"multiaura/internal/models"
	"multiaura/internal/services"
	APIResponse "multiaura/pkg/api_response"

	"github.com/gofiber/fiber/v2"
)

type ConversationController struct {
	service services.ConversationService
}

// NewConversationController tạo một instance mới của ConversationController
func NewConversationController(service services.ConversationService) *ConversationController {
	return &ConversationController{service}
}

// CreateConversation xử lý việc tạo một cuộc trò chuyện giữa hai người dùng
func (cc *ConversationController) CreateConversation(c *fiber.Ctx) error {

	var rep struct {
		UserIDs []string `json:"user_ids"`
		Name    string   `json:"name"`
	}
	if err := c.BodyParser(&rep); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Cannot parse JSON",
			Error:   "BadRequest",
		})
	}
	conversation, err := cc.service.CreateConversation(rep.UserIDs, rep.Name)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Fail to create conversation",
			Error:   "StatusInternalServerError",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(APIResponse.SuccessResponse{
		Status:  fiber.StatusCreated,
		Message: "Create Conversation successfully",
		Data:    conversation,
	})

}
func (cc *ConversationController) GetConversationByID(c *fiber.Ctx) error {
	// Lấy conversationID từ params
	conversationID := c.Params("conversationID")

	if conversationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Missing conversationID parameter",
			Error:   "BadRequest",
		})
	}

	// Gọi service để lấy thông tin cuộc trò chuyện
	conversation, err := cc.service.GetConversationByID(conversationID)
	if err != nil {
		// Kiểm tra từng loại lỗi cụ thể và trả về phản hồi phù hợp
		switch err.Error() {
		case "conversation not found":
			return c.Status(fiber.StatusNotFound).JSON(APIResponse.ErrorResponse{
				Status:  fiber.StatusNotFound,
				Message: "The conversation was not found",
				Error:   "ConversationNotFound",
			})
		}
	}
	return c.Status(fiber.StatusOK).JSON(APIResponse.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Get Conversation successfully",
		Data:    conversation,
	})
}

func (cc *ConversationController) GetListConversation(c *fiber.Ctx) error {
	userID := c.Params("UserID")

	if userID == "" {
		return c.Status(fiber.StatusOK).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "userID is not required",
			Error:   "BadRequest",
		})
	}
	conversation, err := cc.service.GetListConversations(userID)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "No get list conversation found",
			Error:   "Internal server error",
		})
	}
	return c.Status(fiber.StatusOK).JSON(APIResponse.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Get list conversation successfully",
		Data:    conversation,
	})
}
func (cc *ConversationController) AddMember(c *fiber.Ctx) error {
	conversationID := c.Params("conversationID")

	var req struct {
		UserID []string `json:"user_id" bson:"user_id" form:"user_id"`
	}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "userID is not required",
			Error:   "BadRequest",
		})
	}

	if conversationID == "" || len(req.UserID) == 0 {
		return c.Status(fiber.StatusOK).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid conversationID or userID",
			Error:   "BadRequest",
		})
	}

	err = cc.service.AddMembers(conversationID, req.UserID)
	if err != nil {

		return c.Status(fiber.StatusOK).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "No dont add members ",
			Error:   "StatusInternalServerError",
		})
	}

	return c.Status(fiber.StatusOK).JSON(APIResponse.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Member added successfully",
		Data:    nil,
	})
}
func (cc *ConversationController) RemoveMemberConversation(c *fiber.Ctx) error {
	conversationID := c.Params("ConversationID")
	UserID := c.Params("UserID")

	if conversationID == "" || UserID == "" {

		return c.Status(fiber.StatusOK).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid conversationID or userID",
			Error:   "BadRequest",
		})
	}
	err := cc.service.RemoveMenberConversation(conversationID, UserID)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Cannot delete member of conversation",
			Error:   "BadRequest",
		})
	}
	return c.Status(fiber.StatusOK).JSON(APIResponse.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Delete Member  successfully",
		Data:    nil,
	})
}

func (cc *ConversationController) SendMessage(c *fiber.Ctx) error {
	conversationID := c.Params("conversationID")

	var messageData struct {
		UserID  string             `json:"user_id"`
		Content models.ChatContent `json:"content"`
	}

	if err := c.BodyParser(&messageData); err != nil {

		return c.Status(fiber.StatusOK).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Cannot send messages of conversation",
			Error:   "BadRequest",
		})
	}

	// Gọi service để gửi tin nhắn
	err := cc.service.SendMessage(conversationID, messageData.UserID, messageData.Content)
	if err != nil {

		return c.Status(fiber.StatusOK).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Cannot send messages of conversation",
			Error:   "StatusInternalServerError",
		})
	}
	return c.Status(fiber.StatusOK).JSON(APIResponse.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Message sent successfully.",
		Data:    nil,
	})

}
func (cc *ConversationController) GetMessages(c *fiber.Ctx) error {
	conversationID := c.Params("conversationID")

	messages, err := cc.service.GetMessages(conversationID)
	if err != nil {

		return c.Status(fiber.StatusOK).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Unable to retrieve messages for the conversation",
			Error:   "StatusInternalServerError",
		})
	}

	return c.Status(fiber.StatusOK).JSON(APIResponse.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Get Message successfully.",
		Data:    messages,
	})
}
func (cc *ConversationController) MarkMessageAsDeleted(c *fiber.Ctx) error {
	conversationID := c.Params("conversationID")
	messageID := c.Params("messageID")

	err := cc.service.MarkMessageAsDeleted(conversationID, messageID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(APIResponse.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Message marked as deleted successfully",
		Data:    nil,
	})
}
