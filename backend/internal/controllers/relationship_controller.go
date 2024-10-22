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

func (uc *RelationshipController) GetFollowers(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "StatusUnauthorized",
		})
	}

	followers, err := uc.service.GetFollowers(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Fail to get followers",
			Error:   "StatusInternalServerError",
		})
	}

	if len(followers) == 0 {
		return c.Status(fiber.StatusOK).JSON(APIResponse.SuccessResponse{
			Status:  fiber.StatusOK,
			Message: "No followers found",
			Data:    followers,
		})
	}

	return c.Status(fiber.StatusOK).JSON(APIResponse.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Followers retrieved successfully",
		Data:    followers,
	})
}

func (uc *RelationshipController) GetFollowings(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "StatusUnauthorized",
		})
	}

	followings, err := uc.service.GetFollowings(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Fail to get followings",
			Error:   "StatusInternalServerError",
		})
	}

	if len(followings) == 0 {
		return c.Status(fiber.StatusOK).JSON(APIResponse.SuccessResponse{
			Status:  fiber.StatusOK,
			Message: "No followings found",
			Data:    followings,
		})
	}

	return c.Status(fiber.StatusOK).JSON(APIResponse.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Followings retrieved successfully",
		Data:    followings,
	})
}

func (uc *RelationshipController) GetBlockedUsers(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "StatusUnauthorized",
		})
	}

	followings, err := uc.service.GetBlockedUsers(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Fail to get blocked users",
			Error:   "StatusInternalServerError",
		})
	}

	if len(followings) == 0 {
		return c.Status(fiber.StatusOK).JSON(APIResponse.SuccessResponse{
			Status:  fiber.StatusOK,
			Message: "No blocked users found",
			Data:    followings,
		})
	}

	return c.Status(fiber.StatusOK).JSON(APIResponse.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Blocked users retrieved successfully",
		Data:    followings,
	})
}

func (uc *RelationshipController) GetProfile(c *fiber.Ctx) error {
	username := c.Params("username")
	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Username is required",
			Error:   "StatusBadRequest",
		})
	}

	userID := c.Locals("userID").(string)

	userProfile, err := uc.service.GetProfile(userID, username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
			Error:   "StatusInternalServerError",
		})
	}

	return c.Status(fiber.StatusOK).JSON(APIResponse.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "User profile retrieved successfully",
		Data:    userProfile,
	})
}

func (uc *RelationshipController) GetRelationshipStatus(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "StatusUnauthorized",
		})
	}

	targetUserID := c.Params("userID")
	if targetUserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Target userID is required",
			Error:   "StatusBadRequest",
		})
	}

	relationshipStatus, err := uc.service.GetRelationship(userID, targetUserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to get relationship status",
			Error:   "StatusInternalServerError",
		})
	}

	// Trả về trạng thái quan hệ
	return c.Status(fiber.StatusOK).JSON(APIResponse.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Relationship status retrieved successfully",
		Data:    relationshipStatus,
	})
}
