package controllers

import (
	"multiaura/internal/models"
	"multiaura/internal/services"
	APIResponse "multiaura/pkg/api_response"
	"multiaura/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{service}
}

// Register
func (uc *UserController) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	err := uc.service.Register(&req)
	if err != nil {
		if err.Error() == "email already exists" {
			return c.Status(fiber.StatusConflict).JSON(APIResponse.ErrorResponse{
				Status:  fiber.StatusConflict,
				Message: "Email already exists",
				Error:   "Conflict",
			})
		}
		// Return specific errors if any
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
			Error:   "BadRequest",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(APIResponse.SuccessResponse{
		Status:  fiber.StatusCreated,
		Message: "Register successful",
		Data:    req,
	})
}

// Login
func (uc *UserController) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Call the Login service method
	existingUser, err := uc.service.Login(req.Username, req.Password)
	if err != nil {
		if err.Error() == "invalid email" {
			return c.Status(fiber.StatusUnauthorized).JSON(APIResponse.ErrorResponse{
				Status:  fiber.StatusUnauthorized,
				Message: "Invalid email",
				Error:   "StatusUnauthorized",
			})
		} else if err.Error() == "invalid password" {
			return c.Status(fiber.StatusUnauthorized).JSON(APIResponse.ErrorResponse{
				Status:  fiber.StatusUnauthorized,
				Message: "Invalid password",
				Error:   "StatusUnauthorized",
			})
		}
		// Handle other unexpected errors
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "StatusInternalServerError",
		})
	}

	// Generate JWT token if login is successful
	token, err := jwt.GenerateToken(*existingUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Could not create token",
			Error:   "StatusInternalServerError",
		})
	}
	return c.Status(fiber.StatusOK).JSON(APIResponse.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Login successful",
		Data: fiber.Map{
			"token": token,
			"data":  existingUser,
		},
	})
}

// Delete
func (uc *UserController) DeleteUser(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "StatusUnauthorized",
		})
	}

	err := uc.service.DeleteAccount(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusNotFound,
			Message: "User not found",
			Error:   "StatusNotFound",
		})
	}
	return c.Status(fiber.StatusOK).JSON(APIResponse.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "User deleted successfully",
		Data:    nil,
	})
}

// Update
func (uc *UserController) UpdateUser(c *fiber.Ctx) error {
	// Khai báo một map để lưu dữ liệu cập nhật
	updatedData := make(map[string]interface{})

	// Phân tích JSON vào map
	if err := c.BodyParser(&updatedData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "StatusUnauthorized",
		})
	}

	// Thêm userID vào map cập nhật
	updatedData["userID"] = userID

	// Gọi hàm Update với map thay vì đối tượng User
	err := uc.service.Update(&updatedData)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(APIResponse.ErrorResponse{
				Status:  fiber.StatusNotFound,
				Message: "User not found",
				Error:   "StatusNotFound",
			})
		} else if err.Error() == "user ID does not match" {
			return c.Status(fiber.StatusForbidden).JSON(APIResponse.ErrorResponse{
				Status:  fiber.StatusForbidden,
				Message: "You do not have permission to update this user",
				Error:   "StatusForbidden",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to update user information",
			Error:   "StatusInternalServerError",
		})
	}

	return c.Status(fiber.StatusOK).JSON(APIResponse.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "User updated successfully",
		Data:    updatedData,
	})
}

// ForgotPassword
func (uc *UserController) ForgotPassword(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	err := uc.service.ForgotPassword(req.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
			Error:   "BadRequest",
		})
	}
	return c.Status(fiber.StatusOK).JSON(APIResponse.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Password reset instructions sent to your email",
		Data:    nil,
	})
}

// ChangePassword
func (uc *UserController) ChangePassword(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "StatusUnauthorized",
		})
	}

	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	err := uc.service.ChangePassword(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
			Error:   "BadRequest",
		})
	}
	return c.Status(fiber.StatusOK).JSON(APIResponse.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Password changed successfully",
		Data:    nil,
	})
}