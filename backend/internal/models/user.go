package models

import (
	"errors"
	"multiaura/pkg/utils"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type User struct {
	ID          string    `bson:"userID" json:"userID" form:"userID"`
	FullName    string    `bson:"fullname" json:"fullname" form:"fullname"`
	Username    string    `bson:"username" json:"username" form:"username"`
	Email       string    `bson:"email" json:"email" form:"email"`
	Password    string    `bson:"password" json:"password" form:"password"`
	PhoneNumber string    `bson:"phone" json:"phone" form:"phone"`
	Birthday    time.Time `bson:"birthday" json:"birthday" form:"birthday"`
	Gender      string    `bson:"gender" json:"gender" form:"gender"`
	Nation      string    `bson:"nation" json:"nation" form:"nation"`
	Province    string    `bson:"province" json:"province" form:"province"`
	Avatar      string    `bson:"avatar" json:"avatar" form:"avatar"`
	IsAdmin     bool      `bson:"isAdmin" json:"isAdmin" form:"isAdmin"`
	IsActive    bool      `bson:"isActive" json:"isActive" form:"isActive"`
	IsPublic    bool      `bson:"isPublic" json:"isPublic" form:"isPublic"`
}

type RegisterRequest struct {
	FullName    string `bson:"fullname" json:"fullname" form:"fullname" validate:"required"`
	Username    string `bson:"username" json:"username" form:"username" validate:"required"`
	Email       string `bson:"email" json:"email" form:"email" validate:"required,email"`
	Password    string `bson:"password" json:"password" form:"password" validate:"required,min=3"`
	PhoneNumber string `bson:"phone" json:"phone" form:"phone" validate:"required"`
	Birthday    string `bson:"birthday" json:"birthday" form:"birthday" validate:"required"`
	Gender      string `bson:"gender" json:"gender" form:"gender" validate:"required"`
	Nation      string `bson:"nation" json:"nation" form:"nation"`
	Province    string `bson:"province" json:"province" form:"province"`
}

type LoginRequest struct {
	Username string `bson:"username" json:"username" form:"username"`
	Password string `bson:"password" json:"password" form:"password" validate:"required,min=3"`
}

func (u *User) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"userID":   u.ID,
		"fullname": u.FullName,
		"username": u.Username,
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
		"isPublic": u.IsPublic,
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
		ID:          utils.GetString(data, "userID"),
		FullName:    utils.GetString(data, "fullname"),
		Username:    utils.GetString(data, "username"),
		Email:       utils.GetString(data, "email"),
		Password:    utils.GetString(data, "password"),
		PhoneNumber: utils.GetString(data, "phone"),
		Birthday:    birthday,
		Gender:      utils.GetString(data, "gender"),
		Nation:      utils.GetString(data, "nation"),
		Province:    utils.GetString(data, "province"),
		Avatar:      utils.GetString(data, "avatar"),
		IsAdmin:     utils.GetBool(data, "isAdmin"),
		IsActive:    utils.GetBool(data, "isActive"),
		IsPublic:    utils.GetBool(data, "isPublic"),
	}, nil
}
