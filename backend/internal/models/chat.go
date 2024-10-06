package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	ChatContent struct {
		Text  string `json:"text" form:"text" bson:"text,omitempty"`
		Image string `json:"image" form:"image" bson:"image,omitempty"`
		Voice string `json:"voice" form:"voice" bson:"voice,omitempty"`
	}

	Chat struct {
		ID       primitive.ObjectID `json:"_id" form:"_id" bson:"_id,omitempty"`
		Sender   User               `json:"sender" form:"sender" bson:"sender"`
		Content  ChatContent        `json:"content" form:"content" bson:"content"`
		Emotion  []string           `json:"emotion" form:"emotion" bson:"emotion,omitempty"`
		CreateAt time.Time          `json:"created_at" form:"created_at" bson:"created_at"`
		UpdateAt time.Time          `json:"updated_at" form:"updated_at" bson:"updated_at"`
		Status   int                `json:"status" form:"status" bson:"status"`
	}

	Conversation struct {
		ID               primitive.ObjectID `json:"_id" form:"_id" bson:"_id,omitempty"`
		Name             string             `json:"name_conversation" form:"name_conversation" bson:"name_conversation"`
		ConversationType string             `json:"conversation_type" form:"conversation_type" bson:"conversation_type"`
		Users            []Users            `json:"users" form:"users" bson:"users"`
		Chats            []Chat             `json:"chats" form:"chats" bson:"chats"`
		SeenBy           []SeenBy           `json:"seen_by" form:"seen_by" bson:"seen_by"`
		CreateAt         time.Time          `json:"created_at" form:"created_at" bson:"created_at"`
		UpdateAt         time.Time          `json:"updated_at" form:"updated_at" bson:"updated_at"`
	}

	SeenBy struct {
		UserID primitive.ObjectID `json:"user_id" form:"user_id" bson:"user_id"`
		SeenAt time.Time          `json:"seen_at" form:"seen_at" bson:"seen_at"`
	}

	Users struct {
		UserID   primitive.ObjectID `json:"user_id" form:"user_id" bson:"user_id"`
		Fullname string             `json:"fullname" form:"fullname" bson:"fullname"`
		Avatar   string             `json:"avatar" form:"avatar" bson:"avatar"`
		LastSeen time.Time          `json:"last_seen" form:"last_seen" bson:"last_seen"`
	}
)
