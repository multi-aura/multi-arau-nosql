package repositories

import (
	"context"
	"log"
	"math/rand"
	"multiaura/internal/databases"
	"multiaura/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostRepository interface {
	Repository[models.Post]
	GetRecentPosts(userIDs []string, limit, page int64) ([]*models.Post, error)
	SearchTrendingPosts(query string, limit, page int64) ([]*models.Post, error)
	SearchNewsMixedPosts(query string, userIDs []string, limit, page int64) ([]*models.Post, error)
	SearchPostsForYou(query, userID string, limit, page int64) ([]*models.Post, error)
	Search(query string, blockedUserIDs []string, limit, page int64) ([]*models.Post, error)
	UploadPhotos(id string, url []string) (bool, error)
}

type postRepository struct {
	db         *databases.MongoDB
	collection *mongo.Collection
}

func NewPostRepository(db *databases.MongoDB) PostRepository {
	return &postRepository{
		db:         db,
		collection: db.Database.Collection("posts"),
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

func (repo *postRepository) GetRecentPosts(userIDs []string, limit, page int64) ([]*models.Post, error) {
	var posts []*models.Post
	sort := bson.D{{Key: "createdAt", Value: -1}}
	skip := (page - 1) * limit

	findOptions := options.Find()
	findOptions.SetSort(sort)
	findOptions.SetLimit(limit)
	findOptions.SetSkip(skip)

	filter := bson.M{"createdBy.userID": bson.M{"$in": userIDs}}

	cursor, err := repo.collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var data map[string]interface{}
		if err := cursor.Decode(&data); err != nil {
			return nil, err
		}

		post, err := new(models.Post).FromMap(data)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (repo *postRepository) SearchTrendingPosts(query string, limit, page int64) ([]*models.Post, error) {
	var posts []*models.Post
	skip := (page - 1) * limit

	// Tạo pipeline
	pipeline := mongo.Pipeline{}

	// Nếu query không rỗng, thêm điều kiện tìm kiếm bằng regex
	if query != "" {
		pipeline = append(pipeline, bson.D{
			{Key: "$match", Value: bson.M{"description": bson.M{"$regex": query, "$options": "i"}}},
		})
	}

	// Tính tổng số lượng likes và shares
	pipeline = append(pipeline, bson.D{
		{Key: "$addFields", Value: bson.D{
			{Key: "totalLikes", Value: bson.M{"$size": bson.M{"$ifNull": []interface{}{"$likedBy", []interface{}{}}}}},
			{Key: "totalShares", Value: bson.M{"$size": bson.M{"$ifNull": []interface{}{"$sharedBy", []interface{}{}}}}},
		}},
	})

	// Tạo điểm (score) dựa trên lượt like, share và thời gian tạo (createdAt)
	pipeline = append(pipeline, bson.D{
		{Key: "$addFields", Value: bson.D{
			{Key: "score", Value: bson.M{
				"$add": []interface{}{
					bson.M{"$multiply": []interface{}{"$totalLikes", 2}},  // Mỗi like nhân với 2 điểm
					bson.M{"$multiply": []interface{}{"$totalShares", 1}}, // Mỗi share nhân với 1 điểm
					bson.M{"$divide": []interface{}{
						1,
						bson.M{"$subtract": []interface{}{time.Now(), "$createdAt"}},
					}}, // Thời gian càng gần hiện tại, điểm càng cao
				},
			}},
		}},
	})

	// Sắp xếp theo điểm
	pipeline = append(pipeline, bson.D{
		{Key: "$sort", Value: bson.D{{Key: "score", Value: -1}}},
	})

	// Giới hạn số lượng kết quả
	if limit > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$limit", Value: limit}})
	}

	if skip > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$skip", Value: skip}})
	}

	cursor, err := repo.collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var data map[string]interface{}
		if err := cursor.Decode(&data); err != nil {
			return nil, err
		}

		post, err := new(models.Post).FromMap(data)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (repo *postRepository) SearchNewsMixedPosts(query string, userIDs []string, limit, page int64) ([]*models.Post, error) {
	var friendPosts []*models.Post
	var otherPosts []*models.Post

	sort := bson.D{{Key: "createdAt", Value: -1}}
	skip := (page - 1) * limit
	// Bước 1: Lấy bài viết từ bạn bè
	friendFilter := bson.M{"createdBy.userID": bson.M{"$in": userIDs}}
	if query != "" {
		// Sử dụng $regex để tìm kiếm description chứa query
		friendFilter["description"] = bson.M{"$regex": query, "$options": "i"} // "i" để tìm kiếm không phân biệt hoa thường
	}
	friendOptions := options.Find().SetSort(sort).SetLimit(limit).SetSkip(skip)

	friendCursor, err := repo.collection.Find(context.Background(), friendFilter, friendOptions)
	if err != nil {
		return nil, err
	}
	defer friendCursor.Close(context.Background())

	for friendCursor.Next(context.Background()) {
		var data map[string]interface{}
		if err := friendCursor.Decode(&data); err != nil {
			return nil, err
		}

		post, err := new(models.Post).FromMap(data)
		if err != nil {
			return nil, err
		}

		friendPosts = append(friendPosts, post)
	}
	// Bước 2: Lấy bài viết từ những người khác
	otherFilter := bson.M{"createdBy.userID": bson.M{"$nin": userIDs}}
	if query != "" {
		// Sử dụng $regex để tìm kiếm description chứa query
		otherFilter["description"] = bson.M{"$regex": query, "$options": "i"}
	}
	otherOptions := options.Find().SetSort(sort).SetLimit(limit).SetSkip(skip)

	otherCursor, err := repo.collection.Find(context.Background(), otherFilter, otherOptions)
	if err != nil {
		return nil, err
	}
	defer otherCursor.Close(context.Background())

	for otherCursor.Next(context.Background()) {
		var data map[string]interface{}
		if err := otherCursor.Decode(&data); err != nil {
			return nil, err
		}

		post, err := new(models.Post).FromMap(data)
		if err != nil {
			return nil, err
		}

		otherPosts = append(otherPosts, post)
	}
	// Bước 3: Trộn bài viết theo quy tắc đã nêu
	var mixedPosts []*models.Post
	friendCount := 0
	otherCount := 0

	for len(mixedPosts) < int(limit) {
		// Kiểm tra nếu đã hết bài viết từ bạn bè và người khác
		if friendCount >= len(friendPosts) && otherCount >= len(otherPosts) {
			break // Thoát vòng lặp nếu không còn bài viết nào
		}

		// Lấy bài viết từ bạn bè
		if friendCount < len(friendPosts) {
			// 3 bài đầu tiên từ bạn bè
			if len(mixedPosts) < 3 {
				mixedPosts = append(mixedPosts, friendPosts[friendCount])
				friendCount++
			} else if (len(mixedPosts)-3)%2 == 0 && len(mixedPosts) < int(limit) {
				// Lấy bài từ bạn bè theo quy tắc xen kẽ
				mixedPosts = append(mixedPosts, friendPosts[friendCount])
				friendCount++
			}
		}

		// Lấy bài viết từ những người khác
		if otherCount < len(otherPosts) {
			// 2 bài từ những người khác
			if len(mixedPosts) < 2 || (len(mixedPosts)%2 == 0 && len(mixedPosts) < int(limit)) {
				mixedPosts = append(mixedPosts, otherPosts[otherCount])
				otherCount++
			}
		}
	}
	// Bước 4: Random hóa kết quả
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(mixedPosts), func(i, j int) {
		mixedPosts[i], mixedPosts[j] = mixedPosts[j], mixedPosts[i]
	})

	// Giới hạn kết quả về số lượng bài viết tối đa
	if int64(len(mixedPosts)) > limit {
		mixedPosts = mixedPosts[:limit]
	}
	return mixedPosts, nil
}

func (repo *postRepository) SearchPostsForYou(query, userID string, limit, page int64) ([]*models.Post, error) {
	var posts []*models.Post

	// Tính toán số lượng bài viết cần bỏ qua
	skip := (page - 1) * limit

	// Tiêu chí: Bài viết từ bạn bè
	friendFilter := bson.M{"createdBy.userID": bson.M{"$in": []string{userID}}} // Bài viết từ bạn bè
	if query != "" {
		// Sử dụng $regex để tìm kiếm description chứa query
		friendFilter["description"] = bson.M{"$regex": query, "$options": "i"} // "i" để tìm kiếm không phân biệt hoa thường
	}
	friendSort := bson.D{{Key: "likesCount", Value: -1}} // Sắp xếp theo lượt thích

	// Kết hợp các truy vấn
	findOptions := options.Find()
	findOptions.SetSort(friendSort) // Sắp xếp theo lượt thích
	findOptions.SetLimit(limit)
	findOptions.SetSkip(skip) // Bỏ qua số bài viết đã tính toán

	// Lấy bài viết từ bạn bè
	cursor, err := repo.collection.Find(context.Background(), friendFilter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var data map[string]interface{}
		if err := cursor.Decode(&data); err != nil {
			return nil, err
		}

		post, err := new(models.Post).FromMap(data)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	// Tiêu chí: Bài viết từ những người khác
	otherFilter := bson.M{"createdBy.userID": bson.M{"$nin": []string{userID}}} // Bài viết từ những người khác
	if query != "" {
		// Sử dụng $regex để tìm kiếm description chứa query
		otherFilter["description"] = bson.M{"$regex": query, "$options": "i"}
	}

	// Lấy bài viết từ những người khác
	cursor, err = repo.collection.Find(context.Background(), otherFilter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var data map[string]interface{}
		if err := cursor.Decode(&data); err != nil {
			return nil, err
		}

		post, err := new(models.Post).FromMap(data)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	// Giới hạn số lượng bài viết về limit
	if int64(len(posts)) > limit {
		posts = posts[:limit]
	}

	return posts, nil
}

func (repo *postRepository) Search(query string, blockedUserIDs []string, limit int64, page int64) ([]*models.Post, error) {
	var posts []*models.Post
	sort := bson.D{{Key: "createdAt", Value: -1}}
	skip := (page - 1) * limit

	findOptions := options.Find()
	findOptions.SetSort(sort)
	findOptions.SetLimit(limit)
	findOptions.SetSkip(skip)

	log.Println(blockedUserIDs)

	// Tạo bộ lọc cho truy vấn
	filter := bson.M{
		"$and": []bson.M{
			{
				"description": bson.M{"$regex": query, "$options": "i"},
			},
		},
	}

	if len(blockedUserIDs) > 0 {
		filter["createdBy.userID"] = bson.M{"$nin": blockedUserIDs}
	}

	cursor, err := repo.collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var data map[string]interface{}
		if err := cursor.Decode(&data); err != nil {
			return nil, err
		}

		post, err := new(models.Post).FromMap(data)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (repo *postRepository) UploadPhotos(id string, urls []string) (bool, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	var images []models.Image
	for _, url := range urls {
		images = append(images, models.Image{
			URL: url,
			ID:  primitive.NewObjectID(),
		})
	}

	filter := bson.M{
		"_id": objID,
	}

	update := bson.M{
		"$set": bson.M{
			"images": images,
		},
	}

	result, err := repo.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return false, err
	}

	if result.MatchedCount == 0 {
		return false, nil
	}

	return true, nil
}
