package controllers

import (
	"multiaura/internal/services"
	APIResponse "multiaura/pkg/api_response"

	"github.com/gofiber/fiber/v2"
)

type RelationshipController struct {
	service services.RelationshipService
}

func NewRelationshipController(service services.RelationshipService) *RelationshipController {
	return &RelationshipController{service}
}

func (uc *RelationshipController) GetFriends(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "StatusUnauthorized",
		})
	}

	friends, err := uc.service.GetFriends(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Fail to get friends",
			Error:   "StatusInternalServerError",
		})
	}

	if len(friends) == 0 {
		return c.Status(fiber.StatusOK).JSON(APIResponse.SuccessResponse{
			Status:  fiber.StatusOK,
			Message: "No friends found",
			Data:    friends,
		})
	}

	return c.Status(fiber.StatusOK).JSON(APIResponse.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Friends retrieved successfully",
		Data:    friends,
	})
}
