package models

import (
	"errors"
	"multiaura/pkg/utils"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type OtherUser struct {
	ID                 string              `bson:"userID" json:"userID" form:"userID"`
	FullName           string              `bson:"fullname" json:"fullname" form:"fullname"`
	Username           string              `bson:"username" json:"username" form:"username"`
	Birthday           time.Time           `bson:"birthday" json:"birthday" form:"birthday"`
	Gender             string              `bson:"gender" json:"gender" form:"gender"`
	Nation             string              `bson:"nation" json:"nation" form:"nation"`
	Province           string              `bson:"province" json:"province" form:"province"`
	Avatar             string              `bson:"avatar" json:"avatar" form:"avatar"`
	IsAdmin            bool                `bson:"isAdmin" json:"isAdmin" form:"isAdmin"`
	IsActive           bool                `bson:"isActive" json:"isActive" form:"isActive"`
	IsPublic           bool                `bson:"isPublic" json:"isPublic" form:"isPublic"`
	RelationshipStatus *RelationshipStatus `bson:"relationshipStatus,omitempty" json:"relationshipStatus,omitempty"`
}

func (u *OtherUser) ToMap() map[string]interface{} {
	data := map[string]interface{}{
		"userID":   u.ID,
		"fullname": u.FullName,
		"username": u.Username,
		"birthday": u.Birthday,
		"gender":   u.Gender,
		"nation":   u.Nation,
		"province": u.Province,
		"avatar":   u.Avatar,
		"isAdmin":  u.IsAdmin,
		"isActive": u.IsActive,
		"isPublic": u.IsPublic,
	}

	if u.RelationshipStatus != nil {
		data["relationshipStatus"] = map[string]interface{}{
			"status": u.RelationshipStatus.Status,
			"since":  u.RelationshipStatus.Since,
		}
	}

	return data
}

func (u *OtherUser) FromMap(data map[string]interface{}) (*OtherUser, error) {
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

	// Xử lý relationshipStatus và since
	var relationshipStatus *RelationshipStatus
	if val, ok := data["relationshipStatus"]; ok {
		statusStr, ok := val.(string)
		if ok {
			// Chuyển đổi statusStr thành RelationshipStatusType
			var status RelationshipStatusType
			switch statusStr {
			case "BLOCKING":
				status = Blocking
			case "BLOCKED_BY":
				status = Blocked
			case "FRIEND":
				status = Friend
			case "FOLLOWING":
				status = Following
			case "FOLLOWED_BY":
				status = Follower
			default:
				status = NoRelationship
			}

			// Lấy timestamp từ "since"
			var since *time.Time
			if sinceVal, ok := data["since"]; ok {
				if timestamp, ok := sinceVal.(int64); ok {
					parsedSince := time.UnixMilli(timestamp) // Chuyển timestamp thành thời gian
					since = &parsedSince
				}
			}

			// Gán lại giá trị cho relationshipStatus
			relationshipStatus = &RelationshipStatus{
				Status: status,
				Since:  since,
			}
		}
	}

	return &OtherUser{
		ID:                 utils.GetString(data, "userID"),
		FullName:           utils.GetString(data, "fullname"),
		Username:           utils.GetString(data, "username"),
		Birthday:           birthday,
		Gender:             utils.GetString(data, "gender"),
		Nation:             utils.GetString(data, "nation"),
		Province:           utils.GetString(data, "province"),
		Avatar:             utils.GetString(data, "avatar"),
		IsAdmin:            utils.GetBool(data, "isAdmin"),
		IsActive:           utils.GetBool(data, "isActive"),
		IsPublic:           utils.GetBool(data, "isPublic"),
		RelationshipStatus: relationshipStatus,
	}, nil
}
