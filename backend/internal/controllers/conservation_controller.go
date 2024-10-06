package controllers

import (
	"multiaura/internal/services"
	APIResponse "multiaura/pkg/api_response"

	"github.com/gofiber/fiber/v2"
)

type ConservationController struct {
	service services.ConservationService
}

func NewConservationController(service services.ConservationService) *ConservationController {
	return &ConservationController{service}
}

func (uc *ConservationController) CreateConservation(c *fiber.Ctx) error {
	userID := c.Params("userID")
	targetUserID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "StatusUnauthorized",
		})
	}

	err := uc.service.CreateConservation(targetUserID, userID)
	if err != nil {
		switch err.Error() {
		case "user not found":
			return c.Status(fiber.StatusNotFound).JSON(APIResponse.ErrorResponse{
				Status:  fiber.StatusNotFound,
				Message: "User not found",
				Error:   "StatusNotFound",
			})
		case "user ID does not match":
			return c.Status(fiber.StatusBadRequest).JSON(APIResponse.ErrorResponse{
				Status:  fiber.StatusBadRequest,
				Message: "User ID does not match",
				Error:   "StatusBadRequest",
			})
		case "failed to check block status":
			return c.Status(fiber.StatusInternalServerError).JSON(APIResponse.ErrorResponse{
				Status:  fiber.StatusInternalServerError,
				Message: "Failed to check block status",
				Error:   "StatusInternalServerError",
			})
		case "user is not blocked":
			return c.Status(fiber.StatusConflict).JSON(APIResponse.ErrorResponse{
				Status:  fiber.StatusConflict,
				Message: "User is not blocked",
				Error:   "StatusConflict",
			})
		case "failed to unblock user":
			return c.Status(fiber.StatusInternalServerError).JSON(APIResponse.ErrorResponse{
				Status:  fiber.StatusInternalServerError,
				Message: "Failed to unblock user",
				Error:   "StatusInternalServerError",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(APIResponse.ErrorResponse{
				Status:  fiber.StatusInternalServerError,
				Message: "An unexpected error occurred",
				Error:   "StatusInternalServerError",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(APIResponse.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Successfully unblocked the user",
		Data:    nil,
	})
}
