package models

import "time"

type RelationshipStatusType string

const (
	NoRelationship RelationshipStatusType = "noRelationship"
	IsFollowing    RelationshipStatusType = "isFollowing"
	IsFollowedBy   RelationshipStatusType = "isFollowedBy"
	IsBlocking     RelationshipStatusType = "isBlocking"
	IsBlockedBy    RelationshipStatusType = "isBlocked"
	IsFriend       RelationshipStatusType = "isFriend"
)

type RelationshipStatus struct {
	Status RelationshipStatusType `json:"status"`
	Since  *time.Time             `json:"since,omitempty"`
}
