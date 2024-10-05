package services

import (
	"errors"
	"multiaura/internal/models"
	"multiaura/internal/repositories"
	"multiaura/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(req *models.RegisterRequest) error
	Login(email string, password string) (*models.User, error)
	Logout(userID string) error
	DeleteAccount(userID string) error
	Update(userMap *map[string]interface{}) error
	ForgotPassword(email string) error
	ChangePassword(userID, oldPassword, newPassword string) error
	ComparePassword(hashedPassword string, plainPassword string) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

// Register a new user
func (s *userService) Register(req *models.RegisterRequest) error {
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.FullName == "" {
		return errors.New("fullname is required")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	if req.PhoneNumber == "" {
		return errors.New("phonenumber is required")
	}

	reqMap, err := utils.StructToMap(req)
	if err != nil {
		return errors.New("failed to convert request to map")
	}
	user := &models.User{}
	user, err = user.FromMap(reqMap)
	if err != nil {
		return errors.New("failed to convert to User")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}
	user.Password = string(hashedPassword)

	err = s.repo.Create(*user)
	if err != nil {
		return err
	}

	return nil
}

// Login a user by email
func (s *userService) Login(email string, password string) (*models.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}

func (s *userService) ComparePassword(hashedPassword string, plainPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		return errors.New("invalid password")
	}
	return nil
}

func (s *userService) Logout(userID string) error {
	return nil
}

// Delete a user account
func (s *userService) DeleteAccount(userID string) error {
	existingUser, err := s.repo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if userID != existingUser.ID {
		return errors.New("user ID does not match")
	}
	return s.repo.Delete(userID)
}

// Update a user's information
func (s *userService) Update(userMap *map[string]interface{}) error {
	userID := (*userMap)["user_id"].(string)
	existingUser, err := s.repo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if userID != existingUser.ID {
		return errors.New("user ID does not match")
	}

	if err := s.repo.Update(userMap); err != nil {
		return errors.New("failed to update user information")
	}

	return nil
}

// ForgotPassword
func (s *userService) ForgotPassword(email string) error {
	return nil
}

// Change a user's password
func (s *userService) ChangePassword(userID, oldPassword, newPassword string) error {
	// user, err := s.repo.GetByID(userID)
	// if err != nil {
	// 	return err
	// }

	// // Check if the old password matches
	// if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
	// 	return errors.New("invalid old password")
	// }

	// // Hash the new password
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	// if err != nil {
	// 	return err
	// }

	// // Update password in the database
	// user.Password = string(hashedPassword)
	// return s.repo.Update(*user)
	return errors.New("can not change password")
}
