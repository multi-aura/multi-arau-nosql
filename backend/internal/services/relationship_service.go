package services

import (
	"errors"
	"multiaura/internal/models"
	"multiaura/internal/repositories"
)

type RelationshipService interface {
	Follow(targetUserID, userID string) error
	UnFollow(targetUserID, userID string) error
	Block(targetUserID, userID string) error
	UnBlock(targetUserID, userID string) error
	GetFriends(userID string) ([]*models.User, error)
}

type relationshipService struct {
	repo repositories.UserRepository
}

func NewRelationshipService(repo repositories.UserRepository) RelationshipService {
	return &relationshipService{repo}
}

func (s *relationshipService) Follow(targetUserID, userID string) error {
	existingUser, err := s.repo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if userID != existingUser.ID {
		return errors.New("user ID does not match")
	}

	isFollowing, err := s.repo.IsFollowingOrFriend(targetUserID, userID)
	if err != nil {
		return errors.New("failed to check follow status")
	}

	if isFollowing {
		return errors.New("user already followed or friend with target user")
	}

	err = s.repo.Follow(targetUserID, userID)
	if err != nil {
		return errors.New("failed to follow")
	}

	return nil
}

func (s *relationshipService) UnFollow(targetUserID, userID string) error {
	panic("unimplemented")
}

func (s *relationshipService) Block(targetUserID, userID string) error {
	panic("unimplemented")
}

func (s *relationshipService) UnBlock(targetUserID, userID string) error {
	panic("unimplemented")
}

func (s *relationshipService) GetFriends(userID string) ([]*models.User, error) {
	friends, err := s.repo.GetFriends(userID)
	if err != nil {
		return nil, err
	}
	return friends, nil
}