package repositories

import (
	"context"
	"multiaura/internal/databases"
	"multiaura/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostRepository interface {
	Repository[models.Post]
}

type postRepository struct {
	db         *databases.MongoDB
	collection *mongo.Collection
}

func NewPostRepository(db *databases.MongoDB) PostRepository {
	return &postRepository{
		db:         db,
		collection: db.Database.Collection("post"),
	}
}

func (repo *postRepository) GetByID(id string) (*models.Post, error) {
	var post models.Post

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &models.Post{}, err
	}

	filter := bson.M{"_id": objectID}

	err = repo.collection.FindOne(context.Background(), filter).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &post, nil
}

func (repo *postRepository) Create(entity models.Post) error {
	_, err := repo.collection.InsertOne(context.Background(), entity)
	return err
}

func (repo *postRepository) Delete(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}

	result, err := repo.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (repo *postRepository) Update(entityMap *map[string]interface{}) error {
	objectID, err := primitive.ObjectIDFromHex((*entityMap)["postID"].(string))
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}

	updateQuery := bson.M{"$set": entityMap}

	result, err := repo.collection.UpdateOne(context.Background(), filter, updateQuery)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
