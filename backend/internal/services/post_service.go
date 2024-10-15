package services

import (
	"errors"
	"multiaura/internal/models"
	"multiaura/internal/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostService interface {
	GetPostByID(id string) (*models.Post, error)
	CreatePost(post *models.CreatePostRequest) error
	UpdatePost(id string, updates *map[string]interface{}) error
	DeletePost(id string) error
	GetRecentPosts(userID string, limit int64, page int64) ([]*models.Post, error)
}

type postService struct {
	repo     repositories.PostRepository
	userRepo repositories.UserRepository
}

func NewPostService(repo *repositories.PostRepository, userRepo *repositories.UserRepository) PostService {
	return &postService{repo: *repo, userRepo: *userRepo}
}

func (s *postService) GetPostByID(id string) (*models.Post, error) {
	post, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("failed to get post: " + err.Error())
	}
	return post, nil
}

func (s *postService) CreatePost(post *models.CreatePostRequest) error {
	user, err := s.userRepo.GetUserSummaryByID(post.UserID)
	if err != nil {
		return errors.New("failed to get user: " + err.Error())
	}

	images := make([]models.Image, len(post.Images))
	for i, img := range post.Images {
		images[i] = models.Image{
			URL: img.URL,
			ID:  primitive.NewObjectID(),
		}
	}

	newPost := &models.Post{
		ID:          primitive.NewObjectID(),
		Description: post.Description,
		Images:      images,
		CreatedAt:   time.Now(),
		CreatedBy:   *user,
		LikedBy:     nil,
		SharedBy:    nil,
		UpdatedAt:   time.Now(),
	}
	err = s.repo.Create(*newPost)
	if err != nil {
		return errors.New("failed to create post: " + err.Error())
	}

	return nil
}

func (s *postService) UpdatePost(id string, updates *map[string]interface{}) error {
	(*updates)["postID"] = id
	err := s.repo.Update(updates)
	if err != nil {
		return errors.New("failed to update post: " + err.Error())
	}
	return nil
}

func (s *postService) DeletePost(id string) error {
	err := s.repo.Delete(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("post not found")
		}
		return errors.New("failed to delete post: " + err.Error())
	}
	return nil
}

func (s *postService) GetRecentPosts(userID string, limit int64, page int64) ([]*models.Post, error) {
	followings, err := s.userRepo.GetFollowings(userID)
	if err != nil {
		return nil, errors.New("failed to fetch followings: " + err.Error())
	}

	userIDs := []string{userID}
	for _, following := range followings {
		userIDs = append(userIDs, following.ID)
	}

	posts, err := s.repo.GetRecentPosts(userIDs, limit, page)
	if err != nil {
		return nil, errors.New("failed to get recent posts: " + err.Error())
	}
	return posts, nil
}
