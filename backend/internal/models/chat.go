package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	ChatContent struct {
		Text     string `json:"text" bson:"text,omitempty"`
		Image    string `json:"image" bson:"image,omitempty"`
		VoiceURL string `json:"voice_url" bson:"voice_url,omitempty"`
	}

	Chat struct {
		ID        primitive.ObjectID `json:"id_chat" bson:"id_chat,omitempty"`
		Sender    User               `json:"sender" bson:"sender"`
		Content   ChatContent        `json:"content" bson:"content"`
		Emotion   []string           `json:"emotion" bson:"emotion,omitempty"`
		CreatedAt time.Time          `json:"created_at" bson:"created_at"`
		UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
		Status    string             `json:"status" bson:"status"`
	}

	Conversation struct {
		ID               primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
		Name             string             `json:"name_conversation" bson:"name_conversation"`
		ConversationType string             `json:"conversation_type" bson:"conversation_type"`
		Users            []Users            `json:"users" bson:"users"`
		Chats            []Chat             `json:"chats" bson:"chats"`
		SeenBy           []SeenBy           `json:"seen_by" bson:"seen_by"`
		CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
		UpdatedAt        time.Time          `json:"updated_at" bson:"updated_at"`
	}

	SeenBy struct {
		UserID primitive.ObjectID `json:"user_id" bson:"user_id"`
		SeenAt time.Time          `json:"seen_at" bson:"seen_at"`
	}

	Users struct {
		UserID   string    `json:"user_id" bson:"user_id"`
		Fullname string    `json:"fullname" bson:"fullname"`
		Avatar   string    `json:"avatar" bson:"avatar"`
		LastSeen time.Time `json:"last_seen" bson:"last_seen"`
	}
)
