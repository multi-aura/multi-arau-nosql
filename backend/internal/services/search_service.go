package services

import (
	"multiaura/internal/models"
	"multiaura/internal/repositories"
)

type SearchService interface {
	SearchNews(userID, query string) ([]*models.Post, error)
	SearchPeople(userID, query string, page int, limit int) ([]*models.OtherUser, error)
	GetSuggestedFriends(userID string, page int, limit int) ([]*models.OtherUser, error)
	SearchPosts(userID, query string) ([]*models.Post, error)
	SearchTrending(userID, query string) ([]*models.Post, error)
	SearchForYou(userID, query string) ([]*models.Post, error)
}

type searchService struct {
	userRepo repositories.UserRepository
	postRepo repositories.PostRepository
}

func NewSearchService(userRepo *repositories.UserRepository, postRepo *repositories.PostRepository) SearchService {
	return &searchService{*userRepo, *postRepo}
}

func (s *searchService) SearchForYou(userID, query string) ([]*models.Post, error) {
	panic("unimplemented")
}

func (s *searchService) SearchNews(userID, query string) ([]*models.Post, error) {
	panic("unimplemented")
}

func (s *searchService) SearchPeople(userID, query string, page int, limit int) ([]*models.OtherUser, error) {
	otherUsers, err := s.userRepo.Search(userID, query, page, limit)
	if err != nil {
		return nil, err
	}

	return otherUsers, nil
}

func (s *searchService) GetSuggestedFriends(userID string, page int, limit int) ([]*models.OtherUser, error) {
	otherUsers, err := s.userRepo.GetSuggestedFriends(userID, page, limit)
	if err != nil {
		return nil, err
	}

	return otherUsers, nil
}

func (s *searchService) SearchPosts(userID, query string) ([]*models.Post, error) {
	panic("unimplemented")
}

func (s *searchService) SearchTrending(userID, query string) ([]*models.Post, error) {
	panic("unimplemented")
}
