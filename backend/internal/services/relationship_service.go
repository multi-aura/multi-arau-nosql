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
	GetFollowers(userID string) ([]*models.UserSummary, error)
	GetFollowings(userID string) ([]*models.UserSummary, error)
	GetProfile(userID, username string) (*models.UserProfile, error)
	IsBlockedBy(targetUserID, userID string) (bool, error)
	GetBlockedUsers(userID string) ([]*models.UserSummary, error)
	GetRelationship(userID, targetUserID string) (*models.RelationshipStatus, error)
}

type relationshipService struct {
	repo repositories.UserRepository
}

func NewRelationshipService(repo *repositories.UserRepository) RelationshipService {
	return &relationshipService{*repo}
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

	if !isFollowing {
		return errors.New("user is not following or friend with target user")
	}

	err = s.repo.UnFollow(targetUserID, userID)
	if err != nil {
		return errors.New("failed to unfollow")
	}

	return nil
}

func (s *relationshipService) Block(targetUserID, userID string) error {
	existingUser, err := s.repo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if userID != existingUser.ID {
		return errors.New("user ID does not match")
	}

	isBlocked, err := s.repo.IsBlocking(targetUserID, userID)
	if err != nil {
		return errors.New("failed to check block status")
	}

	if isBlocked {
		return errors.New("user already blocked")
	}

	err = s.repo.Block(targetUserID, userID)
	if err != nil {
		return errors.New("failed to block user")
	}

	return nil
}

func (s *relationshipService) UnBlock(targetUserID, userID string) error {
	existingUser, err := s.repo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if userID != existingUser.ID {
		return errors.New("user ID does not match")
	}

	isBlocked, err := s.repo.IsBlocking(targetUserID, userID)
	if err != nil {
		return errors.New("failed to check block status")
	}

	if !isBlocked {
		return errors.New("user is not blocked")
	}

	err = s.repo.UnBlock(targetUserID, userID)
	if err != nil {
		return errors.New("failed to unblock user")
	}

	return nil
}

func (s *relationshipService) GetFriends(userID string) ([]*models.User, error) {
	friends, err := s.repo.GetFriends(userID)
	if err != nil {
		return nil, err
	}
	return friends, nil
}

func (s *relationshipService) GetFollowers(userID string) ([]*models.UserSummary, error) {
	followers, err := s.repo.GetFollowers(userID)
	if err != nil {
		return nil, err
	}

	return followers, nil
}

func (s *relationshipService) GetFollowings(userID string) ([]*models.UserSummary, error) {
	followings, err := s.repo.GetFollowings(userID)
	if err != nil {
		return nil, err
	}

	return followings, nil
}

func (s *relationshipService) GetBlockedUsers(userID string) ([]*models.UserSummary, error) {
	blockedList, err := s.repo.GetBlockedUsers(userID)
	if err != nil {
		return nil, err
	}

	return blockedList, nil
}

func (s *relationshipService) GetProfile(userID, username string) (*models.UserProfile, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	isBlocked, err := s.repo.IsBlockedBy(userID, user.ID)
	if err == nil {
		if isBlocked {
			return nil, errors.New("blocked by this user")
		}
	}

	isBlocking, err := s.repo.IsBlocking(userID, user.ID)
	if err == nil {
		if isBlocking {
			return nil, errors.New("blocking this user")
		}
	}

	// mutualFollowings, err := s.repo.GetMutualFollowings(userID, user.ID)
	// if err != nil {
	// 	return nil, err
	// }

	mutualFriends, err := s.repo.GetMutualFriends(userID, user.ID)
	if err != nil {
		return nil, err
	}

	// Tạo một UserProfile mới
	userProfile := &models.UserProfile{
		User: user,
		// MutualFollowings: mutualFollowings,
		MutualFriends: mutualFriends,
	}

	return userProfile, nil
}

func (s *relationshipService) IsBlockedBy(targetUserID, userID string) (bool, error) {
	isBlocked, err := s.repo.IsBlockedBy(targetUserID, userID)
	if err != nil {
		return false, errors.New("failed to check block status")
	}
	return isBlocked, nil
}

func (s *relationshipService) GetRelationship(userID, targetUserID string) (*models.RelationshipStatus, error) {
	userExists, err := s.repo.GetByID(userID)
	if err != nil || userExists == nil {
		return &models.RelationshipStatus{}, errors.New("user not found")
	}

	targetUserExists, err := s.repo.GetByID(targetUserID)
	if err != nil || targetUserExists == nil {
		return &models.RelationshipStatus{}, errors.New("target user not found")
	}

	relationshipStatus, err := s.repo.GetRelationship(targetUserID, userID)
	if err != nil {
		return &models.RelationshipStatus{}, err
	}

	return &relationshipStatus, nil
}
