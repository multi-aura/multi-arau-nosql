package services

import (
	"errors"
	"multiaura/internal/models"
	"multiaura/internal/repositories"
	"multiaura/pkg/utils"
	"multiaura/pkg/validators"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(req *models.RegisterRequest) error
	Login(username string, password string) (*models.User, error)
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
	if req.Username == "" {
		return errors.New("username is required")
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

	existsEmail, _ := s.repo.GetUserByEmail(user.Email)
	if existsEmail != nil {
		return errors.New("email already exists")
	}

	existsPhone, _ := s.repo.GetUserByPhone(user.PhoneNumber)
	if existsPhone != nil {
		return errors.New("phone already exists")
	}

	existsUsername, _ := s.repo.GetUserByUsername(user.Username)
	if existsUsername != nil {
		return errors.New("username already exists")
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
func (s *userService) Login(username string, password string) (*models.User, error) {
	var user *models.User
	var err error

	if isValid := validators.IsValidateEmail(username); isValid {
		user, err = s.repo.GetUserByEmail(username)
		if err != nil {
			return nil, errors.New("user not found with this email")
		}
	} else {
		user, err = s.repo.GetUserByUsername(username)
		if err != nil {
			return nil, errors.New("user not found with this username")
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
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

func (s *userService) Update(userMap *map[string]interface{}) error {
	userID := (*userMap)["userID"].(string)
	existingUser, err := s.repo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if userID != existingUser.ID {
		return errors.New("user ID does not match")
	}
	existsPhone, _ := s.repo.GetUserByPhone((*userMap)["phone"].(string))
	if existsPhone != nil {
		return errors.New("phone already exists")
	}

	if err := s.repo.Update(userMap); err != nil {
		return errors.New("failed to update user information")
	}

	return nil
}

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
