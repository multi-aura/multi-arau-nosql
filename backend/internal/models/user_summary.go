package models

import (
	"multiaura/pkg/utils"
)

type UserSummary struct {
	ID       string `bson:"userID" json:"userID" form:"userID"`
	FullName string `bson:"fullname" json:"fullname" form:"fullname"`
	Username string `bson:"username" json:"username" form:"username"`
	Avatar   string `bson:"avatar" json:"avatar" form:"avatar"`
	IsActive bool   `bson:"isActive" json:"isActive" form:"isActive"`
}

func (u *UserSummary) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"userID":   u.ID,
		"fullname": u.FullName,
		"username": u.Username,
		"avatar":   u.Avatar,
		"isActive": u.IsActive,
	}
}

func (u *UserSummary) FromMap(data map[string]interface{}) (*UserSummary, error) {
	return &UserSummary{
		ID:       utils.GetString(data, "userID"),
		FullName: utils.GetString(data, "fullname"),
		Username: utils.GetString(data, "username"),
		Avatar:   utils.GetString(data, "avatar"),
		IsActive: utils.GetBool(data, "isActive"),
	}, nil
}
