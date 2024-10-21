package services

import (
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"multiaura/internal/repositories"
	"multiaura/pkg/utils"
)

type UploadService interface {
	UploadProfilePhoto(userID string, file multipart.File, fileHeader *multipart.FileHeader) (string, error)
	UploadPostPhotos(postID string, files []multipart.File, fileHeaders []*multipart.FileHeader) ([]string, error)
}

type uploadService struct {
	userRepo    *repositories.UserRepository
	postRepo    *repositories.PostRepository
	storageRepo *repositories.StorageRepository
}

func NewUploadService(userRepo *repositories.UserRepository, postRepo *repositories.PostRepository, storageRepo *repositories.StorageRepository) UploadService {
	return &uploadService{userRepo, postRepo, storageRepo}
}

func (s *uploadService) UploadProfilePhoto(userID string, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	url, err := (*s.storageRepo).UploadFile(file, fileHeader, "profile-photos")
	if err != nil {
		return "", err
	}

	result, err := (*s.userRepo).UploadProfilePhoto(userID, url)
	if err != nil {
		return "", err
	}

	if !result {
		fileName, err := utils.ExtractFileName(url)
		if err != nil {
			return "", err
		} else {
			fmt.Println("File name:", fileName)
		}
		if deleteErr := (*s.storageRepo).DeleteFile(fileName); deleteErr != nil {
			return "", errors.New("failed to upload profile photo, and unable to delete file")
		}
		return "", errors.New("failed to update your profile photo")
	}

	return url, nil
}

func (s *uploadService) UploadPostPhotos(postID string, files []multipart.File, fileHeaders []*multipart.FileHeader) ([]string, error) {
	var fileURLs []string
	folder := fmt.Sprintf("posts/%s", postID)

	// Upload từng file và lưu URL của file
	for i, file := range files {
		log.Println("file:", file)
		fileURL, err := (*s.storageRepo).UploadFile(file, fileHeaders[i], folder)
		if err != nil {
			// Nếu xảy ra lỗi, xóa các file đã upload trước đó
			s.cleanupUploadedFiles(fileURLs)
			return nil, errors.New("failed to upload files")
		}
		fileURLs = append(fileURLs, fileURL)
	}
	log.Println("urls:", fileURLs)

	// Cập nhật thông tin URL của ảnh vào database
	result, err := (*s.postRepo).UploadPhotos(postID, fileURLs)
	if err != nil {
		s.cleanupUploadedFiles(fileURLs)
		return nil, err
	}

	if !result {
		s.cleanupUploadedFiles(fileURLs)
		return nil, errors.New("failed to update post with uploaded photos")
	}

	return fileURLs, nil
}

func (s *uploadService) cleanupUploadedFiles(fileURLs []string) {
	for _, fileURL := range fileURLs {
		// Tách tên file từ URL
		fileName, err := utils.ExtractFileName(fileURL)
		if err != nil {
			fmt.Println("Error extracting file name:", err)
			continue
		}
		
		// Xóa file từ storage
		if err := (*s.storageRepo).DeleteFile(fileName); err != nil {
			fmt.Println("Error deleting file:", fileName, err)
		}
	}
}

