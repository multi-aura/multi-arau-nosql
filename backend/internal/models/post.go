package models

import (
	"multiaura/pkg/utils"
	"time"
)

type Image struct {
	URL string `bson:"url" json:"url" form:"url"`
	ID  string `bson:"_id,omitempty" json:"_id,omitempty" form:"_id,omitempty"`
}

type Post struct {
	ID          string       `bson:"_id,omitempty" json:"_id,omitempty" form:"_id,omitempty"`
	Description string       `bson:"description" json:"description" form:"description"`
	Images      []Image      `bson:"images" json:"images" form:"images"`
	CreatedAt   time.Time    `bson:"createdAt" json:"createdAt" form:"createdAt"`
	CreatedBy   UserSummary  `bson:"createdBy" json:"createdBy" form:"createdBy"`
	LikedBy     []UserSummary `bson:"likedBy" json:"likedBy" form:"likedBy"`
	SharedBy    []string     `bson:"sharedBy" json:"sharedBy" form:"sharedBy"`
	UpdatedAt   time.Time    `bson:"updatedAt" json:"updatedAt" form:"updatedAt"`
}

func (p *Post) ToMap() map[string]interface{} {
	images := make([]map[string]interface{}, len(p.Images))
	for i, img := range p.Images {
		images[i] = map[string]interface{}{
			"url": img.URL,
			"_id": img.ID,
		}
	}

	likedBy := make([]map[string]interface{}, len(p.LikedBy))
	for i, user := range p.LikedBy {
		likedBy[i] = user.ToMap()
	}

	return map[string]interface{}{
		"_id":         p.ID,
		"description": p.Description,
		"images":      images,
		"createdAt":   p.CreatedAt,
		"createdBy":   p.CreatedBy.ToMap(),
		"likedBy":     likedBy,
		"sharedBy":    p.SharedBy,
		"updatedAt":   p.UpdatedAt,
	}
}

func (p *Post) FromMap(data map[string]interface{}) (*Post, error) {
	imageData := utils.GetArray(data, "images")
	images := make([]Image, len(imageData))
	for i, img := range imageData {
		imgMap := img.(map[string]interface{})
		images[i] = Image{
			URL: utils.GetString(imgMap, "url"),
			ID:  utils.GetString(imgMap, "_id"),
		}
	}

	// Chuyển đổi LikedBy
	likedByData := utils.GetArray(data, "likedBy")
	likedBy := make([]UserSummary, len(likedByData))
	for i, usr := range likedByData {
		userMap := usr.(map[string]interface{})
		userSummary := UserSummary{}
		_, err := userSummary.FromMap(userMap)
		if err != nil {
			return nil, err
		}
		likedBy[i] = userSummary
	}

	// Chuyển đổi CreatedBy
	createdByData := utils.GetMap(data, "createdBy")
	createdBy := UserSummary{}
	_, err := createdBy.FromMap(createdByData)
	if err != nil {
		return nil, err
	}

	return &Post{
		ID:          utils.GetString(data, "_id"),
		Description: utils.GetString(data, "description"),
		Images:      images,
		CreatedAt:   utils.GetTime(data, "createdAt"),
		CreatedBy:   createdBy,
		LikedBy:     likedBy,
		SharedBy:    utils.GetStringArray(data, "sharedBy"),
		UpdatedAt:   utils.GetTime(data, "updatedAt"),
	}, nil
}