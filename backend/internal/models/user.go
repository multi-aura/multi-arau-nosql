package models

import (
	"errors"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type User struct {
	ID          string    `bson:"user_id" json:"user_id" form:"user_id"`
	FullName    string    `bson:"fullname" json:"fullname" form:"fullname"`
	Email       string    `bson:"email" json:"email" form:"email"`
	Password    string    `bson:"password" json:"password" form:"password"`
	PhoneNumber string    `bson:"phone" json:"phone" form:"phone"`
	Birthday    time.Time `bson:"birthday" json:"birthday" form:"birthday"`
	Gender      string    `bson:"gender" json:"gender" form:"gender"`
	Nation      string    `bson:"nation" json:"nation" form:"nation"`
	Province    string    `bson:"province" json:"province" form:"province"`
	Avatar      string    `bson:"avatar" json:"avatar" form:"avatar"`
	IsAdmin     bool      `bson:"isAdmin" json:"isAdmin" form:"isAdmin"`
	IsActive    bool      `bson:"IsActive" json:"IsActive" form:"IsActive"`
}

type RegisterRequest struct {
	FullName    string `bson:"fullname" json:"fullname" form:"fullname" validate:"required"`
	Email       string `bson:"email" json:"email" form:"email" validate:"required,email"`
	Password    string `bson:"password" json:"password" form:"password" validate:"required,min=3"`
	PhoneNumber string `bson:"phone" json:"phone" form:"phone" validate:"required"`
	Birthday    string `bson:"birthday" json:"birthday" form:"birthday" validate:"required"`
	Gender      string `bson:"gender" json:"gender" form:"gender" validate:"required"`
	Nation      string `bson:"nation" json:"nation" form:"nation"`
	Province    string `bson:"province" json:"province" form:"province"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *User) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"user_id":  u.ID,
		"fullname": u.FullName,
		"email":    u.Email,
		"password": u.Password,
		"phone":    u.PhoneNumber,
		"birthday": u.Birthday,
		"gender":   u.Gender,
		"nation":   u.Nation,
		"province": u.Province,
		"avatar":   u.Avatar,
		"isAdmin":  u.IsAdmin,
		"isActive": u.IsActive,
	}
}

func (u *User) FromMap(data map[string]interface{}) (*User, error) {
	var birthday time.Time

	if val, ok := data["birthday"]; ok {
		if dateStr, ok := val.(string); ok {
			parsedDate, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				return nil, errors.New("invalid birthday format")
			}
			birthday = parsedDate
		} else if date, ok := val.(dbtype.Date); ok {
			parsedDate, err := time.Parse("2006-01-02", date.String())
			if err != nil {
				return nil, errors.New("invalid birthday format")
			}
			birthday = parsedDate
		} else {
			birthday = time.Now()
		}
	} else {
		birthday = time.Now()
	}

	return &User{
		ID:          getString(data, "user_id"),
		FullName:    getString(data, "fullname"),
		Email:       getString(data, "email"),
		Password:    getString(data, "password"),
		PhoneNumber: getString(data, "phone"),
		Birthday:    birthday,
		Gender:      getString(data, "gender"),
		Nation:      getString(data, "nation"),
		Province:    getString(data, "province"),
		Avatar:      getString(data, "avatar"),
		IsAdmin:     getBool(data, "isAdmin"),
		IsActive:    getBool(data, "isActive"),
	}, nil
}

func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getBool(data map[string]interface{}, key string) bool {
	if val, ok := data[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}
