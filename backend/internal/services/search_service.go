package services

import (
	"multiaura/internal/models"
	"multiaura/internal/repositories"
)

type SearchService interface {
	SearchNews(userID, query string, page int, limit int) ([]*models.Post, error)
	SearchPeople(userID, query string, page int, limit int) ([]*models.OtherUser, error)
	GetSuggestedFriends(userID string, page int, limit int) ([]*models.OtherUser, error)
	SearchPosts(userID, query string, page int, limit int) ([]*models.Post, error)
	SearchTrending(query string, page int, limit int) ([]*models.Post, error)
	SearchForYou(userID, query string, page int, limit int) ([]*models.Post, error)
}

type searchService struct {
	userRepo repositories.UserRepository
	postRepo repositories.PostRepository
}

func NewSearchService(userRepo *repositories.UserRepository, postRepo *repositories.PostRepository) SearchService {
	return &searchService{*userRepo, *postRepo}
}

func (s *searchService) SearchForYou(userID, query string, page int, limit int) ([]*models.Post, error) {
	posts, err := s.postRepo.SearchPostsForYou(query, userID, int64(limit), int64(page))
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *searchService) SearchNews(userID, query string, page int, limit int) ([]*models.Post, error) {
	followings, err := s.userRepo.GetFollowings(userID)
	if err != nil {
		return nil, err
	}

	userIDs := []string{userID}
	for _, following := range followings {
		userIDs = append(userIDs, following.ID)
	}
	posts, err := s.postRepo.SearchNewsMixedPosts(query, userIDs, int64(limit), int64(page))
	if err != nil {
		return nil, err
	}
	return posts, nil
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

func (s *searchService) SearchPosts(userID, query string, page int, limit int) ([]*models.Post, error) {
	blockedList, err := s.userRepo.GetBlockedList(userID)
	if err != nil {
		return nil, err
	}

	if query == "" {
		posts, err := s.postRepo.SearchPostsForYou(query, userID, int64(limit), int64(page))
		if err != nil {
			return nil, err
		}
		return posts, nil
	}
	posts, err := s.postRepo.Search(query, blockedList, int64(limit), int64(page))
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *searchService) SearchTrending(query string, page int, limit int) ([]*models.Post, error) {
	posts, err := s.postRepo.SearchTrendingPosts(query, int64(limit), int64(page))
	if err != nil {
		return nil, err
	}

	return posts, nil
}
