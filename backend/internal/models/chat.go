package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	ChatContent struct {
		Text     string `json:"text" bson:"text,omitempty" form:"text,omitempty"`
		Image    string `json:"image" bson:"image,omitempty" form:"image,omitempty"`
		VoiceURL string `json:"voice_url" bson:"voice_url,omitempty" form:"voice_url,omitempty"`
	}

	Chat struct {
		ID        primitive.ObjectID `json:"id_chat" bson:"id_chat,omitempty" form:"id_chat,omitempty"`
		Sender    Users              `json:"sender" bson:"sender" form:"sender"`
		Content   ChatContent        `json:"content" bson:"content" form:"content"`
		Emotion   []string           `json:"emotion" bson:"emotion,omitempty" form:"emotion,omitempty"`
		CreatedAt time.Time          `json:"createdat" bson:"createdat" form:"createdat"`
		UpdatedAt time.Time          `json:"updatedat" bson:"updatedat" form:"updatedat"`
		Status    string             `json:"status" bson:"status" form:"status"`
	}

	Conversation struct {
		ID               primitive.ObjectID `json:"_id" bson:"_id,omitempty" form:"_id,omitempty"`
		Name             string             `json:"name_conversation" bson:"name_conversation" form:"name_conversation"`
		ConversationType string             `json:"conversation_type" bson:"conversation_type" form:"conversation_type"`
		Users            []Users            `json:"users" bson:"users" form:"users"`
		Chats            []Chat             `json:"chats" bson:"chats" form:"chats"`
		SeenBy           []SeenBy           `json:"seen_by" bson:"seen_by" form:"seen_by"`
		CreatedAt        time.Time          `json:"createdat" bson:"createdat" form:"createdat"`
		UpdatedAt        time.Time          `json:"updatedat" bson:"updatedat" form:"updatedat"`
	}

	SeenBy struct {
		UserID primitive.ObjectID `json:"user_id" bson:"user_id" form:"user_id"`
		SeenAt time.Time          `json:"seen_at" bson:"seen_at" form:"seen_at"`
	}
)
