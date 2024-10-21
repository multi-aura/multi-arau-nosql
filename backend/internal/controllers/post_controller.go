package controllers

import (
	"multiaura/internal/models"
	"multiaura/internal/services"
	APIResponse "multiaura/pkg/api_response"

	"github.com/gofiber/fiber/v2"
)

type PostController struct {
	service services.PostService
}

func NewPostController(service services.PostService) *PostController {
	return &PostController{service}
}

func (pc *PostController) CreatePost(c *fiber.Ctx) error {
	var post models.CreatePostRequest
	if err := c.BodyParser(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid input",
			Error:   "StatusBadRequest",
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

	post.UserID = userID
	if err := pc.service.CreatePost(&post); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
			Error:   "StatusInternalServerError",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(APIResponse.SuccessResponse{
		Status:  fiber.StatusCreated,
		Message: "Post created successfully",
		Data:    post,
	})
}

func (pc *PostController) GetPosts(c *fiber.Ctx) error {
	// posts, err := pc.service.GetPostByID("123")
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(APIResponse.ErrorResponse{
	// 		Status:  fiber.StatusInternalServerError,
	// 		Message: err.Error(),
	// 		Error:   "StatusInternalServerError",
	// 	})
	// }

	// return c.Status(fiber.StatusOK).JSON(APIResponse.SuccessResponse{
	// 	Status:  fiber.StatusOK,
	// 	Message: "Posts retrieved successfully",
	// 	Data:    posts,
	// })

	return c.Status(fiber.StatusOK).JSON(APIResponse.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Did not have logic handle this api",
		Data:    nil,
	})
}

func (pc *PostController) UpdatePost(c *fiber.Ctx) error {
	postID := c.Params("postID")
	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid input",
			Error:   "StatusBadRequest",
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

	// Cập nhật bài viết qua dịch vụ
	if err := pc.service.UpdatePost(postID, &updates); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(), // Sử dụng thông điệp lỗi từ service
			Error:   "StatusInternalServerError",
		})
	}

	return c.Status(fiber.StatusOK).JSON(APIResponse.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Post updated successfully",
		Data:    nil,
	})
}

// DeletePost xóa bài viết
func (pc *PostController) DeletePost(c *fiber.Ctx) error {
	postID := c.Params("postID")

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "StatusUnauthorized",
		})
	}

	if err := pc.service.DeletePost(postID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(), // Sử dụng thông điệp lỗi từ service
			Error:   "StatusInternalServerError",
		})
	}

	return c.Status(fiber.StatusOK).JSON(APIResponse.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Post deleted successfully",
		Data:    nil,
	})
}

func (pc *PostController) GetPostByID(c *fiber.Ctx) error {
	postID := c.Params("id")

	post, err := pc.service.GetPostByID(postID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
			Error:   "StatusInternalServerError",
		})
	}

	if post == nil {
		return c.Status(fiber.StatusNotFound).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusNotFound,
			Message: "Post not found",
			Error:   "StatusNotFound",
		})
	}

	return c.Status(fiber.StatusOK).JSON(APIResponse.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Post retrieved successfully",
		Data:    post,
	})
}

func (pc *PostController) GetRecentPosts(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized access",
			Error:   "StatusUnauthorized",
		})
	}

	var req models.PagingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
			Error:   "StatusBadRequest",
		})
	}

	if req.Limit <= 0 {
		req.Limit = 10 // Mặc định limit là 10 nếu không cung cấp hoặc không hợp lệ
	}
	if req.Page <= 0 {
		req.Page = 1 // Mặc định page là 1 nếu không cung cấp hoặc không hợp lệ
	}

	posts, err := pc.service.GetRecentPosts(userID, req.Limit, req.Page)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
			Error:   "StatusInternalServerError",
		})
	}

	return c.Status(fiber.StatusOK).JSON(APIResponse.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Posts retrieved successfully",
		Data:    posts,
	})
}
