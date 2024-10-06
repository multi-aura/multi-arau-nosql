package controllers

import (
	"log"
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

	log.Printf("Conversation ID: %v", conversationID)
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
