package models

type UserProfile struct {
	User             *User          `json:"user"`
	MutualFollowings []*UserSummary `json:"mutualFollowings"`
	MutualFriends    []*UserSummary `json:"mutualFriends"`
}
