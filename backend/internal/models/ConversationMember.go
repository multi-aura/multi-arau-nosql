package models

import "time"

type Users struct {
	UserID   string    `json:"userID" bson:"userID" form:"userID"`
	Fullname string    `json:"fullname" bson:"fullname" form:"fullname"`
	Username string    `bson:"username" json:"username" form:"username"`
	Avatar   string    `json:"avatar" bson:"avatar" form:"avatar"`
	Added_at time.Time `json:"added_at" bson:"added_at" form:"added_at"`
}
