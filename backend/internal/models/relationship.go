package models

import "time"

type RelationshipStatusType string

const (
	NoRelationship RelationshipStatusType = "No Relationship"
	Following      RelationshipStatusType = "Following"
	Follower       RelationshipStatusType = "Followed"
	Blocking       RelationshipStatusType = "Blocking"
	Blocked        RelationshipStatusType = "Blocked"
	Friend         RelationshipStatusType = "Friend"
)

type RelationshipStatus struct {
	Status RelationshipStatusType `json:"status"`
	Since  *time.Time             `json:"since,omitempty"`
}
