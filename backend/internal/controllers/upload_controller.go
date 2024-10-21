package controllers

import (
	"mime/multipart"
	"multiaura/internal/services"
	APIResponse "multiaura/pkg/api_response"

	"github.com/gofiber/fiber/v2"
)

type UploadController struct {
	service services.UploadService
}

func NewUploadController(service services.UploadService) *UploadController {
	return &UploadController{service}
}

func (uc *UploadController) UploadProfilePhoto(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized access",
			Error:   "StatusUnauthorized",
		})
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid file",
			Error:   err.Error(),
		})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Unable to open file",
			Error:   err.Error(),
		})
	}
	defer file.Close()

	url, err := uc.service.UploadProfilePhoto(userID, file, fileHeader)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to upload profile picture",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(APIResponse.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Profile picture uploaded successfully",
		Data:    fiber.Map{"url": url},
	})
}

func (uc *UploadController) UploadPostPhotos(c *fiber.Ctx) error {
	postID := c.Params("postID")
	if postID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Missing post ID",
			Error:   "PostIDMissing",
		})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid form data",
			Error:   err.Error(),
		})
	}

	// Lấy tất cả các file từ form
	files := form.File["photos"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "No files provided",
			Error:   "NoFilesProvided",
		})
	}

	var openedFiles []multipart.File
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(APIResponse.ErrorResponse{
				Status:  fiber.StatusInternalServerError,
				Message: "Unable to open file",
				Error:   err.Error(),
			})
		}
		defer file.Close()
		openedFiles = append(openedFiles, file)
	}

	fileURLs, err := uc.service.UploadPostPhotos(postID, openedFiles, files)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse.ErrorResponse{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to upload post photos",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(APIResponse.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Post photos uploaded successfully",
		Data:    fiber.Map{"urls": fileURLs},
	})
}
