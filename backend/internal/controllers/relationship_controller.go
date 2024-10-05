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

func (uc *RelationshipController) Follow(c *fiber.Ctx) error {
	userID := c.Params("userID")
	targetUserID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "StatusUnauthorized",
		})
	}

	err := uc.service.Follow(targetUserID, userID)
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
		case "failed to check follow status":
			return c.Status(fiber.StatusInternalServerError).JSON(APIResponse.ErrorResponse{
				Status:  fiber.StatusInternalServerError,
				Message: "Failed to check follow status",
				Error:   "StatusInternalServerError",
			})
		case "user already followed or friend with target user":
			return c.Status(fiber.StatusConflict).JSON(APIResponse.ErrorResponse{
				Status:  fiber.StatusConflict,
				Message: "User already followed or friend with target user",
				Error:   "StatusConflict",
			})
		case "failed to follow":
			return c.Status(fiber.StatusInternalServerError).JSON(APIResponse.ErrorResponse{
				Status:  fiber.StatusInternalServerError,
				Message: "Failed to follow",
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
		Message: "Successfully followed the user",
		Data:    nil,
	})
}

func (uc *RelationshipController) UnFollow(c *fiber.Ctx) error {
	userID := c.Params("userID")
	targetUserID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "StatusUnauthorized",
		})
	}

	err := uc.service.UnFollow(targetUserID, userID)
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
		case "failed to check follow status":
			return c.Status(fiber.StatusInternalServerError).JSON(APIResponse.ErrorResponse{
				Status:  fiber.StatusInternalServerError,
				Message: "Failed to check follow status",
				Error:   "StatusInternalServerError",
			})
		case "user is not following or friend with target user":
			return c.Status(fiber.StatusConflict).JSON(APIResponse.ErrorResponse{
				Status:  fiber.StatusConflict,
				Message: "User is not following or friend with target user",
				Error:   "StatusConflict",
			})
		case "failed to unfollow":
			return c.Status(fiber.StatusInternalServerError).JSON(APIResponse.ErrorResponse{
				Status:  fiber.StatusInternalServerError,
				Message: "Failed to unfollow",
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
		Message: "Successfully unfollowed the user",
		Data:    nil,
	})
}


func (uc *RelationshipController) Block(c *fiber.Ctx) error {
	userID := c.Params("userID")
	targetUserID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "StatusUnauthorized",
		})
	}

	err := uc.service.Block(targetUserID, userID)
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
		case "user already blocked":
			return c.Status(fiber.StatusConflict).JSON(APIResponse.ErrorResponse{
				Status:  fiber.StatusConflict,
				Message: "User already blocked",
				Error:   "StatusConflict",
			})
		case "failed to block user":
			return c.Status(fiber.StatusInternalServerError).JSON(APIResponse.ErrorResponse{
				Status:  fiber.StatusInternalServerError,
				Message: "Failed to block user",
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
		Message: "Successfully blocked the user",
		Data:    nil,
	})
}

func (uc *RelationshipController) UnBlock(c *fiber.Ctx) error {
	userID := c.Params("userID")
	targetUserID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "StatusUnauthorized",
		})
	}

	err := uc.service.UnBlock(targetUserID, userID)
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
