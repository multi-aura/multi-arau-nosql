package repositories

import (
	"context"
	"errors"
	"multiaura/internal/databases"
	"multiaura/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type UserRepository interface {
	Repository[models.User]
	GetUsersByName(name string) ([]models.User, error)
	GetUserSummaryByID(id string) (*models.UserSummary, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByPhone(phone string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	Follow(targetUserID, userID string) error
	UnFollow(targetUserID, userID string) error
	Block(targetUserID, userID string) error
	UnBlock(targetUserID, userID string) error
	IsFollowingOrFriend(targetUserID, userID string) (bool, error)
	IsBlockedBy(targetUserID, userID string) (bool, error)
	IsBlocking(targetUserID, userID string) (bool, error)
	IsFollowing(targetUserID, userID string) (bool, error)
	IsFollowedBy(targetUserID, userID string) (bool, error)
	IsFriend(targetUserID, userID string) (bool, error)
	GetFriends(userID string) ([]*models.User, error)
	GetFollowers(userID string) ([]*models.UserSummary, error)
	GetFollowings(userID string) ([]*models.UserSummary, error)
	GetBlockedList(userID string) ([]string, error)
	GetBlockedUsers(userID string) ([]*models.UserSummary, error)
	GetRelationship(targetUserID, userID string) (models.RelationshipStatus, error)
	Search(userID, query string, page, limit int) ([]*models.OtherUser, error)
	GetSuggestedFriends(userID string, page, limit int) ([]*models.OtherUser, error)
	UploadProfilePhoto(userID, url string) (bool, error)
	GetMutualFollowings(targetUserID, userID string) ([]*models.UserSummary, error)
	GetMutualFriends(targetUserID, userID string) ([]*models.UserSummary, error)
}

type userRepository struct {
	db *databases.Neo4jDB
}

func NewUserRepository(db *databases.Neo4jDB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (repo *userRepository) GetByID(id string) (*models.User, error) {
	ctx := context.Background()
	session := repo.db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.Run(ctx, "MATCH (u:User {userID: $userID, isActive: true}) RETURN u", map[string]interface{}{
		"userID": id,
	})
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		node, found := record.Get("u")
		if !found {
			return nil, errors.New("user not found")
		}

		userNode := node.(neo4j.Node)
		user := &models.User{}
		user, err := user.FromMap(userNode.Props)
		if err != nil {
			return nil, errors.New("error converting map to User")
		}
		return user, nil
	}

	return nil, errors.New("user with id " + id + " not found")
}

func (repo *userRepository) GetUserSummaryByID(id string) (*models.UserSummary, error) {
	ctx := context.Background()
	session := repo.db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.Run(ctx, "MATCH (u:User {userID: $userID, isActive: true}) RETURN u.userID, u.fullname, u.username, u.avatar, u.isActive", map[string]interface{}{
		"userID": id,
	})
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		userSummary := &models.UserSummary{}
		if userIDVal, ok := record.Get("u.userID"); ok {
			userSummary.ID = userIDVal.(string)
		}
		if fullnameVal, ok := record.Get("u.fullname"); ok {
			userSummary.FullName = fullnameVal.(string)
		}
		if usernameVal, ok := record.Get("u.username"); ok {
			userSummary.Username = usernameVal.(string)
		}
		if avatarVal, ok := record.Get("u.avatar"); ok {
			userSummary.Avatar = avatarVal.(string)
		}
		if isActive, ok := record.Get("u.isActive"); ok {
			userSummary.IsActive = isActive.(bool)
		}
		return userSummary, nil
	}

	return nil, errors.New("user with id " + id + " not found")
}

func (repo *userRepository) Create(user models.User) error {
	ctx := context.Background()
	session := repo.db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	tx, err := session.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	user.ID = uuid.NewString()
	user.IsActive = true
	_, err = tx.Run(ctx,
		"CREATE (u:User) SET u = $userProps",
		map[string]interface{}{
			"userProps": user.ToMap(),
		},
	)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (repo *userRepository) Update(entityMap *map[string]interface{}) error {
	ctx := context.Background()
	session := repo.db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	tx, err := session.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	userID := (*entityMap)["userID"].(string)

	userProps := make(map[string]interface{})

	for key, value := range *entityMap {
		if key != "userID" {
			userProps[key] = value
		}
	}

	result, err := tx.Run(ctx,
		"MATCH (u:User {userID: $userID}) SET u += $userProps RETURN u",
		map[string]interface{}{
			"userID":    userID,
			"userProps": userProps,
		},
	)

	if err != nil {
		return err
	}

	if !result.Next(ctx) {
		return errors.New("user with id " + userID + " not found")
	}

	return tx.Commit(ctx)
}

func (repo *userRepository) Delete(id string) error {
	ctx := context.Background()
	session := repo.db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	tx, err := session.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	result, err := tx.Run(ctx,
		"MATCH (u:User {userID: $userID}) SET u.isActive = false RETURN u",
		map[string]interface{}{
			"userID": id,
		},
	)
	if err != nil {
		return err
	}

	summary, err := result.Consume(ctx)
	if err != nil {
		return err
	}

	if summary.Counters().PropertiesSet() == 0 {
		return errors.New("user with id " + id + " not found")
	}

	return tx.Commit(ctx)
}

func (repo *userRepository) GetUsersByName(name string) ([]models.User, error) {
	ctx := context.Background()
	session := repo.db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.Run(ctx,
		"MATCH (u:User) WHERE u.username CONTAINS $name RETURN u",
		map[string]interface{}{
			"name": name,
		},
	)
	if err != nil {
		return nil, err
	}

	var users []models.User
	for result.Next(ctx) {
		record := result.Record()
		node, found := record.Get("u")
		if !found {
			continue
		}

		userNode := node.(neo4j.Node)
		user := &models.User{}
		user, err := user.FromMap(userNode.Props)
		if err != nil {
			return nil, errors.New("error converting map to User")
		}
		users = append(users, *user)
	}

	if len(users) == 0 {
		return nil, errors.New("no users found with the name " + name)
	}

	return users, nil
}

func (repo *userRepository) GetUserByEmail(email string) (*models.User, error) {
	ctx := context.Background()
	session := repo.db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.Run(ctx,
		"MATCH (u:User {email: $email}) RETURN u",
		map[string]interface{}{
			"email": email,
		},
	)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		node, found := record.Get("u")
		if !found {
			return nil, errors.New("user not found")
		}

		userNode := node.(neo4j.Node)
		user := &models.User{}
		user, err := user.FromMap(userNode.Props)
		if err != nil {
			return nil, errors.New("error converting map to User")
		}
		return user, nil
	}

	return nil, errors.New("user with email " + email + " not found")
}

func (repo *userRepository) GetUserByPhone(phone string) (*models.User, error) {
	ctx := context.Background()
	session := repo.db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.Run(ctx,
		"MATCH (u:User {phone: $phone}) RETURN u",
		map[string]interface{}{
			"phone": phone,
		},
	)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		node, found := record.Get("u")
		if !found {
			return nil, errors.New("user not found")
		}

		userNode := node.(neo4j.Node)
		user := &models.User{}
		user, err := user.FromMap(userNode.Props)
		if err != nil {
			return nil, errors.New("error converting map to User")
		}
		return user, nil
	}

	return nil, errors.New("user with phone " + phone + " not found")
}

func (repo *userRepository) GetUserByUsername(username string) (*models.User, error) {
	ctx := context.Background()
	session := repo.db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.Run(ctx,
		"MATCH (u:User {username: $username}) RETURN u",
		map[string]interface{}{
			"username": username,
		},
	)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		node, found := record.Get("u")
		if !found {
			return nil, errors.New("user not found")
		}

		userNode := node.(neo4j.Node)
		user := &models.User{}
		user, err := user.FromMap(userNode.Props)
		if err != nil {
			return nil, errors.New("error converting map to User")
		}
		return user, nil
	}

	return nil, errors.New("user with username " + username + " not found")
}

func (repo *userRepository) Follow(targetUserID, userID string) error {
	ctx := context.Background()
	session := repo.db.Driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(ctx, `
			MATCH (u1:User {userID: $targetUserID}), (u2:User {userID: $userID})
			MERGE (u1)-[f1:FOLLOWS {since: timestamp()}]->(u2)
			WITH u1, u2, f1
			OPTIONAL MATCH (u2)-[f2:FOLLOWS]->(u1)
			WITH u1, u2, f1, f2
			WHERE f2 IS NOT NULL
			MERGE (u1)-[:FRIEND_WITH {since: f1.since}]->(u2)
			MERGE (u2)-[:FRIEND_WITH {since: f2.since}]->(u1)
			DELETE f1, f2
			RETURN u1, u2
		`, map[string]interface{}{
			"targetUserID": targetUserID,
			"userID":       userID,
		})
		if err != nil {
			return nil, err
		}
		return nil, nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (repo *userRepository) UnFollow(targetUserID, userID string) error {
	ctx := context.Background()
	session := repo.db.Driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(ctx, `
			MATCH (u1:User {userID: $targetUserID})-[r:FOLLOWS|FRIEND_WITH]->(u2:User {userID: $userID})
			OPTIONAL MATCH (u2)-[f2:FRIEND_WITH]->(u1)
			WITH u1, u2, r, f2
			DELETE r
			WITH u1, u2, f2
			WHERE f2 IS NOT NULL
			CREATE (u2)-[:FOLLOWS {since: f2.since}]->(u1)
			DELETE f2
			RETURN u1, u2
		`, map[string]interface{}{
			"targetUserID": targetUserID,
			"userID":       userID,
		})
		if err != nil {
			return nil, err
		}
		return nil, nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (repo *userRepository) Block(targetUserID, userID string) error {
	ctx := context.Background()
	session := repo.db.Driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(ctx, `
			MATCH (u1:User {userID: $targetUserID}), (u2:User {userID: $userID})
			MERGE (u1)-[:BLOCKED]->(u2)
			RETURN u1, u2
		`, map[string]interface{}{
			"targetUserID": targetUserID,
			"userID":       userID,
		})
		if err != nil {
			return nil, err
		}
		return nil, nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (repo *userRepository) UnBlock(targetUserID, userID string) error {
	ctx := context.Background()
	session := repo.db.Driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(ctx, `
			MATCH (u1:User {userID: $targetUserID})-[r:BLOCKED]->(u2:User {userID: $userID})
			DELETE r
			RETURN u1, u2
		`, map[string]interface{}{
			"targetUserID": targetUserID,
			"userID":       userID,
		})
		if err != nil {
			return nil, err
		}
		return nil, nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (repo *userRepository) IsFollowingOrFriend(targetUserID, userID string) (bool, error) {
	ctx := context.Background()
	session := repo.db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
			MATCH (u1:User {userID: $targetUserID})-[r:FOLLOWS|FRIEND_WITH]->(u2:User {userID: $userID})
			RETURN COUNT(r) > 0 AS isFollowingOrFriend
		`

		record, err := tx.Run(ctx, query, map[string]interface{}{
			"targetUserID": targetUserID,
			"userID":       userID,
		})

		if err != nil {
			return false, err
		}

		if record.Next(ctx) {
			isFollowingOrFriend, _ := record.Record().Get("isFollowingOrFriend")
			return isFollowingOrFriend.(bool), nil
		}

		return false, errors.New("unexpected result")
	})

	if err != nil {
		return false, err
	}

	return result.(bool), nil
}

func (repo *userRepository) IsBlockedBy(targetUserID, userID string) (bool, error) {
	ctx := context.Background()
	session := repo.db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
			MATCH (u1:User {userID: $targetUserID})<-[r:BLOCKED]-(u2:User {userID: $userID})
			RETURN COUNT(r) > 0 AS isBlocked
		`

		record, err := tx.Run(ctx, query, map[string]interface{}{
			"targetUserID": targetUserID,
			"userID":       userID,
		})

		if err != nil {
			return false, err
		}

		if record.Next(ctx) {
			isBlocked, _ := record.Record().Get("isBlocked")
			return isBlocked.(bool), nil
		}

		return false, errors.New("unexpected result")
	})

	if err != nil {
		return false, err
	}

	return result.(bool), nil
}

func (repo *userRepository) IsBlocking(targetUserID, userID string) (bool, error) {
	ctx := context.Background()
	session := repo.db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
			MATCH (u1:User {userID: $targetUserID})-[r:BLOCKED]->(u2:User {userID: $userID})
			RETURN COUNT(r) > 0 AS isBlocking
		`

		record, err := tx.Run(ctx, query, map[string]interface{}{
			"userID":       userID,
			"targetUserID": targetUserID,
		})

		if err != nil {
			return false, err
		}

		if record.Next(ctx) {
			isBlocking, _ := record.Record().Get("isBlocking")
			return isBlocking.(bool), nil
		}

		return false, errors.New("unexpected result")
	})

	if err != nil {
		return false, err
	}

	return result.(bool), nil
}

func (repo *userRepository) IsFollowedBy(targetUserID, userID string) (bool, error) {
	ctx := context.Background()
	session := repo.db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
			MATCH (u1:User {userID: $targetUserID})<-[r:FOLLOWS]-(u2:User {userID: $userID})
			RETURN COUNT(r) > 0 AS isFollowedBy
		`

		record, err := tx.Run(ctx, query, map[string]interface{}{
			"targetUserID": targetUserID,
			"userID":       userID,
		})

		if err != nil {
			return false, err
		}

		if record.Next(ctx) {
			isFollowedBy, _ := record.Record().Get("isFollowedBy")
			return isFollowedBy.(bool), nil
		}

		return false, errors.New("unexpected result")
	})

	if err != nil {
		return false, err
	}

	return result.(bool), nil
}

func (repo *userRepository) IsFollowing(targetUserID, userID string) (bool, error) {
	ctx := context.Background()
	session := repo.db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
			MATCH (u1:User {userID: $targetUserID})-[r:FOLLOWS]->(u2:User {userID: $userID})
			RETURN COUNT(r) > 0 AS isFollowing
		`

		record, err := tx.Run(ctx, query, map[string]interface{}{
			"userID":       userID,
			"targetUserID": targetUserID,
		})

		if err != nil {
			return false, err
		}

		if record.Next(ctx) {
			isFollowing, _ := record.Record().Get("isFollowing")
			return isFollowing.(bool), nil
		}

		return false, errors.New("unexpected result")
	})

	if err != nil {
		return false, err
	}

	return result.(bool), nil
}

func (repo *userRepository) IsFriend(targetUserID, userID string) (bool, error) {
	ctx := context.Background()
	session := repo.db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
			MATCH (u1:User {userID: $userID})-[r:FRIEND_WITH]-(u2:User {userID: $targetUserID})
			RETURN COUNT(r) > 0 AS isFriend
		`

		record, err := tx.Run(ctx, query, map[string]interface{}{
			"userID":       userID,
			"targetUserID": targetUserID,
		})

		if err != nil {
			return false, err
		}

		if record.Next(ctx) {
			isFriend, _ := record.Record().Get("isFriend")
			return isFriend.(bool), nil
		}

		return false, errors.New("unexpected result")
	})

	if err != nil {
		return false, err
	}

	return result.(bool), nil
}

func (repo *userRepository) GetFriends(userID string) ([]*models.User, error) {
	ctx := context.Background()
	session := repo.db.Driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		records, err := tx.Run(ctx, `
			MATCH (u:User {userID: $userID})-[:FRIEND_WITH]->(friend:User)
			RETURN friend
		`, map[string]interface{}{
			"userID": userID,
		})
		if err != nil {
			return nil, err
		}

		// Collect friend models into a slice
		var friends []*models.User
		for records.Next(ctx) {
			record := records.Record()
			friendNode, _ := record.Get("friend")
			friendUser := &models.User{}

			// Convert Neo4j node properties to User model
			friendNodeProps := friendNode.(neo4j.Node).Props
			friendUser, err = friendUser.FromMap(friendNodeProps)
			if err != nil {
				return nil, errors.New("error converting map to User")
			}

			friends = append(friends, friendUser)
		}

		// Check if there were any errors during the record fetching
		if err = records.Err(); err != nil {
			return nil, err
		}
		return friends, nil
	})

	if err != nil {
		return nil, err
	}

	friendList, ok := result.([]*models.User)
	if !ok {
		return nil, errors.New("failed to cast result to []*models.User")
	}

	return friendList, nil
}

func (repo *userRepository) GetFollowers(userID string) ([]*models.UserSummary, error) {
	ctx := context.Background()
	session := repo.db.Driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		records, err := tx.Run(ctx, `
			MATCH (u:User{userID: $userID})<-[:FOLLOWS|FRIEND_WITH]-(f:User)
			RETURN f.userID AS userID, f.fullname AS fullname, f.username AS username, f.avatar AS avatar
		`, map[string]interface{}{
			"userID": userID,
		})
		if err != nil {
			return nil, err
		}

		var followers []*models.UserSummary
		for records.Next(ctx) {
			record := records.Record()
			followerUser := &models.UserSummary{}

			if userIDVal, ok := record.Get("userID"); ok {
				followerUser.ID = userIDVal.(string)
			}
			if fullnameVal, ok := record.Get("fullname"); ok {
				followerUser.FullName = fullnameVal.(string)
			}
			if usernameVal, ok := record.Get("username"); ok {
				followerUser.Username = usernameVal.(string)
			}
			if avatarVal, ok := record.Get("avatar"); ok {
				followerUser.Avatar = avatarVal.(string)
			}
			if isActive, ok := record.Get("isActive"); ok {
				followerUser.IsActive = isActive.(bool)
			}

			followers = append(followers, followerUser)
		}

		if err = records.Err(); err != nil {
			return nil, err
		}
		return followers, nil
	})

	if err != nil {
		return nil, err
	}

	followerList, ok := result.([]*models.UserSummary)
	if !ok {
		return nil, errors.New("failed to cast result to []*models.UserSummary")
	}

	return followerList, nil
}

func (repo *userRepository) GetFollowings(userID string) ([]*models.UserSummary, error) {
	ctx := context.Background()
	session := repo.db.Driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		records, err := tx.Run(ctx, `
			MATCH (u:User{userID: $userID})-[:FOLLOWS|FRIEND_WITH]->(f:User)
			WHERE u.userID = $userID
			RETURN f.userID AS userID, f.fullname AS fullname, f.username AS username, f.avatar AS avatar
		`, map[string]interface{}{
			"userID": userID,
		})
		if err != nil {
			return nil, err
		}

		var followings []*models.UserSummary
		for records.Next(ctx) {
			record := records.Record()
			followingUser := &models.UserSummary{}

			if userIDVal, ok := record.Get("userID"); ok {
				followingUser.ID = userIDVal.(string)
			}
			if fullnameVal, ok := record.Get("fullname"); ok {
				followingUser.FullName = fullnameVal.(string)
			}
			if usernameVal, ok := record.Get("username"); ok {
				followingUser.Username = usernameVal.(string)
			}
			if avatarVal, ok := record.Get("avatar"); ok {
				followingUser.Avatar = avatarVal.(string)
			}
			if isActive, ok := record.Get("isActive"); ok {
				followingUser.IsActive = isActive.(bool)
			}

			followings = append(followings, followingUser)
		}

		if err = records.Err(); err != nil {
			return nil, err
		}
		return followings, nil
	})

	if err != nil {
		return nil, err
	}

	followingList, ok := result.([]*models.UserSummary)
	if !ok {
		return nil, errors.New("failed to cast result to []*models.UserSummary")
	}

	return followingList, nil
}

func (repo *userRepository) GetRelationship(targetUserID, userID string) (models.RelationshipStatus, error) {
	ctx := context.Background()
	session := repo.db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
            MATCH (u1:User {userID: $targetUserID}), (u2:User {userID: $userID})
            OPTIONAL MATCH (u1)-[r1:FOLLOWS]->(u2)
            OPTIONAL MATCH (u2)-[r2:FOLLOWS]->(u1)
            OPTIONAL MATCH (u1)-[r3:BLOCKED]->(u2)
            OPTIONAL MATCH (u2)-[r4:BLOCKED]->(u1)
            OPTIONAL MATCH (u1)-[r5:FRIEND_WITH]->(u2)
            OPTIONAL MATCH (u2)-[r6:FRIEND_WITH]->(u1)
            RETURN 
                CASE 
                    WHEN COUNT(r3) > 0 THEN 'BLOCKING'
                    WHEN COUNT(r4) > 0 THEN 'BLOCKED_BY'
                    WHEN COUNT(r5) > 0 THEN 'FRIEND'
                    WHEN COUNT(r1) > 0 THEN 'FOLLOWING'
                    WHEN COUNT(r2) > 0 THEN 'FOLLOWED_BY'
                    ELSE 'NO_RELATIONSHIP' 
                END AS status,
                COALESCE(
                    r3.since, 
                    r4.since, 
                    r5.since, 
                    r1.since, 
                    r2.since
                ) AS since
        `

		// Thực hiện truy vấn
		record, err := tx.Run(ctx, query, map[string]interface{}{
			"userID":       userID,
			"targetUserID": targetUserID,
		})

		if err != nil {
			return nil, err
		}

		if record.Next(ctx) {
			recordData := record.Record()

			relationshipStatus, _ := recordData.Get("status")
			since, _ := recordData.Get("since")

			var status models.RelationshipStatusType

			switch relationshipStatus.(string) {
			case "BLOCKING":
				status = models.Blocking
			case "BLOCKED_BY":
				status = models.Blocked
			case "FRIEND":
				status = models.Friend
			case "FOLLOWING":
				status = models.Following
			case "FOLLOWED_BY":
				status = models.Follower
			default:
				status = models.NoRelationship
			}

			var sinceTime *time.Time
			if since != nil {
				if sinceTimeValue, ok := since.(time.Time); ok {
					sinceTime = &sinceTimeValue
				}
			}

			return models.RelationshipStatus{
				Status: status,
				Since:  sinceTime,
			}, nil
		}

		return models.RelationshipStatus{}, errors.New("no relationship data found")
	})

	if err != nil {
		return models.RelationshipStatus{}, err
	}

	return result.(models.RelationshipStatus), nil
}

func (repo *userRepository) Search(userID, query string, page, limit int) ([]*models.OtherUser, error) {
	ctx := context.Background()
	session := repo.db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		// Calculate the number of records to skip based on the page and limit
		skip := (page - 1) * limit

		searchQuery := `
			MATCH (u:User)
			WHERE (u.fullname CONTAINS $query OR u.username CONTAINS $query)
				AND u.userID <> $currentUserID
				AND NOT EXISTS {
					MATCH (u)-[:BLOCKED]-(b:User {userID: $currentUserID})
				}
			OPTIONAL MATCH (u1:User {userID: $currentUserID})
			OPTIONAL MATCH (u1)-[r1:FOLLOWS]->(u)
			OPTIONAL MATCH (u)-[r2:FOLLOWS]->(u1)
			OPTIONAL MATCH (u1)-[r3:FRIEND_WITH]-(u)
			OPTIONAL MATCH (u)-[r4:FRIEND_WITH]-(u1)
			RETURN u,
			CASE 
				WHEN COUNT(r3) > 0 THEN 'FRIEND'
				WHEN COUNT(r4) > 0 THEN 'FRIEND'
				WHEN COUNT(r1) > 0 THEN 'FOLLOWING'
				WHEN COUNT(r2) > 0 THEN 'FOLLOWED_BY'
				ELSE 'NO_RELATIONSHIP'
			END AS relationshipStatus,
			COALESCE(
				MAX(CASE WHEN r3 IS NOT NULL THEN r3.since END), 
				MAX(CASE WHEN r4 IS NOT NULL THEN r4.since END), 
				MAX(CASE WHEN r1 IS NOT NULL THEN r1.since END), 
				MAX(CASE WHEN r2 IS NOT NULL THEN r2.since END)
			) AS since
			SKIP $skip
			LIMIT $limit
		`

		// Execute the query with pagination parameters
		records, err := tx.Run(ctx, searchQuery, map[string]interface{}{
			"query":         query,
			"currentUserID": userID,
			"skip":          skip,
			"limit":         limit,
		})

		if err != nil {
			return nil, err
		}

		var otherUsers []*models.OtherUser

		// Loop through the result set
		for records.Next(ctx) {
			record := records.Record()

			// Create a map to combine all data
			userData := make(map[string]interface{})

			// Get the user node
			otherUserNode, _ := record.Get("u")
			otherUserNodeProps := otherUserNode.(neo4j.Node).Props

			// Get the relationship status and since timestamp
			relationshipStatus, _ := record.Get("relationshipStatus")
			since, _ := record.Get("since")

			// Add user properties to the map
			for key, value := range otherUserNodeProps {
				userData[key] = value
			}

			// Add relationship status and since to the map
			userData["relationshipStatus"] = relationshipStatus
			userData["since"] = since

			// Convert the map to OtherUser model
			otherUser := &models.OtherUser{}
			otherUser, err = otherUser.FromMap(userData)
			if err != nil {
				return nil, errors.New("error converting map to User")
			}

			otherUsers = append(otherUsers, otherUser)
		}

		return otherUsers, nil
	})

	if err != nil {
		return nil, err
	}

	users := result.([]*models.OtherUser)
	return users, nil
}

func (repo *userRepository) GetSuggestedFriends(userID string, page, limit int) ([]*models.OtherUser, error) {
	ctx := context.Background()
	session := repo.db.Driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
			MATCH (me:User {userID: $userID})
			OPTIONAL MATCH (me)-[:FOLLOWS|FRIEND_WITH]-(friend:User)
			WITH me, COUNT(friend) AS friendCount
			MATCH (fof:User)
			WHERE fof.userID <> me.userID
				AND NOT (me)-[:BLOCKED]-(fof)
				AND NOT (me)-[:FOLLOWS|FRIEND_WITH]->(fof)
			RETURN fof, friendCount
			ORDER BY
				friendCount DESC,
				CASE WHEN fof.province = me.province THEN 0 ELSE 1 END,
				CASE WHEN fof.nation = me.nation THEN 0 ELSE 1 END
			SKIP $skip
            LIMIT $limit
        `
		skip := (page - 1) * limit

		records, err := tx.Run(ctx, query, map[string]interface{}{
			"userID": userID,
			"skip":   skip,
			"limit":  limit,
		})
		if err != nil {
			return nil, err
		}

		var suggestedUsers []*models.OtherUser

		for records.Next(ctx) {
			record := records.Record()
			otherUserNode, _ := record.Get("fof")
			otherUser := &models.OtherUser{}
			otherUserNodeProps := otherUserNode.(neo4j.Node).Props

			otherUser, err = otherUser.FromMap(otherUserNodeProps)
			if err != nil {
				return nil, errors.New("error converting map to UserSummary")
			}

			suggestedUsers = append(suggestedUsers, otherUser)
		}

		return suggestedUsers, nil
	})

	if err != nil {
		return nil, err
	}

	users := result.([]*models.OtherUser)
	return users, nil
}

func (repo *userRepository) GetBlockedList(userID string) ([]string, error) {
	ctx := context.Background()
	session := repo.db.Driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		records, err := tx.Run(ctx, `
			MATCH (u:User{userID: $userID})-[r:BLOCKED]-(blocked:User)
			RETURN blocked.userID AS userID
		`, map[string]interface{}{
			"userID": userID,
		})
		if err != nil {
			return nil, err
		}

		var blockeds []string
		for records.Next(ctx) {
			record := records.Record()

			if userIDVal, ok := record.Get("userID"); ok {
				blockeds = append(blockeds, userIDVal.(string))
			}
		}

		if err = records.Err(); err != nil {
			return nil, err
		}
		return blockeds, nil
	})

	if err != nil {
		return nil, err
	}

	blockedList, ok := result.([]string)
	if !ok {
		return nil, errors.New("failed to cast result to []*string")
	}

	return blockedList, nil
}

func (repo *userRepository) GetBlockedUsers(userID string) ([]*models.UserSummary, error) {
	ctx := context.Background()
	session := repo.db.Driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		records, err := tx.Run(ctx, `
			MATCH (u:User{userID: $userID})-[:BLOCKED]->(f:User)
			RETURN f.userID AS userID, f.fullname AS fullname, f.username AS username, f.avatar AS avatar
		`, map[string]interface{}{
			"userID": userID,
		})
		if err != nil {
			return nil, err
		}

		var followings []*models.UserSummary
		for records.Next(ctx) {
			record := records.Record()
			followingUser := &models.UserSummary{}

			if userIDVal, ok := record.Get("userID"); ok {
				followingUser.ID = userIDVal.(string)
			}
			if fullnameVal, ok := record.Get("fullname"); ok {
				followingUser.FullName = fullnameVal.(string)
			}
			if usernameVal, ok := record.Get("username"); ok {
				followingUser.Username = usernameVal.(string)
			}
			if avatarVal, ok := record.Get("avatar"); ok {
				followingUser.Avatar = avatarVal.(string)
			}
			if isActive, ok := record.Get("isActive"); ok {
				followingUser.IsActive = isActive.(bool)
			}

			followings = append(followings, followingUser)
		}

		if err = records.Err(); err != nil {
			return nil, err
		}
		return followings, nil
	})

	if err != nil {
		return nil, err
	}

	followingList, ok := result.([]*models.UserSummary)
	if !ok {
		return nil, errors.New("failed to cast result to []*models.UserSummary")
	}

	return followingList, nil
}

func (repo *userRepository) UploadProfilePhoto(userID, url string) (bool, error) {
	ctx := context.Background()
	session := repo.db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	tx, err := session.BeginTransaction(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Close(ctx)

	_, err = tx.Run(ctx,
		"MATCH (u:User {userID: $userID}) SET u.avatar = $url RETURN u",
		map[string]interface{}{
			"userID": userID,
			"url":    url,
		},
	)
	if err != nil {
		return false, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (repo *userRepository) GetMutualFollowings(targetUserID, userID string) ([]*models.UserSummary, error) {
	ctx := context.Background()
	session := repo.db.Driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		records, err := tx.Run(ctx, `
			MATCH (u1:User {userID: $userID})-[:FOLLOWS]->(mutual:User)<-[:FOLLOWS]-(u2:User {userID: $targetUserID})
			RETURN mutual.userID AS userID, mutual.fullname AS fullname, mutual.username AS username, mutual.avatar AS avatar
		`, map[string]interface{}{
			"userID":       userID,
			"targetUserID": targetUserID,
		})
		if err != nil {
			return nil, err
		}

		var mutualFollowings []*models.UserSummary
		for records.Next(ctx) {
			record := records.Record()
			mutualUser := &models.UserSummary{}

			if userIDVal, ok := record.Get("userID"); ok {
				mutualUser.ID = userIDVal.(string)
			}
			if fullnameVal, ok := record.Get("fullname"); ok {
				mutualUser.FullName = fullnameVal.(string)
			}
			if usernameVal, ok := record.Get("username"); ok {
				mutualUser.Username = usernameVal.(string)
			}
			if avatarVal, ok := record.Get("avatar"); ok {
				mutualUser.Avatar = avatarVal.(string)
			}

			mutualFollowings = append(mutualFollowings, mutualUser)
		}

		if err = records.Err(); err != nil {
			return nil, err
		}
		return mutualFollowings, nil
	})

	if err != nil {
		return nil, err
	}

	mutualFollowingList, ok := result.([]*models.UserSummary)
	if !ok {
		return nil, errors.New("failed to cast result to []*models.UserSummary")
	}

	return mutualFollowingList, nil
}

func (repo *userRepository) GetMutualFriends(targetUserID, userID string) ([]*models.UserSummary, error) {
	ctx := context.Background()
	session := repo.db.Driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		records, err := tx.Run(ctx, `
			MATCH (u1:User {userID: $userID})-[:FRIEND_WITH]->(mutual:User)<-[:FRIEND_WITH]-(u2:User {userID: $targetUserID})
			RETURN mutual.userID AS userID, mutual.fullname AS fullname, mutual.username AS username, mutual.avatar AS avatar
		`, map[string]interface{}{
			"userID":       userID,
			"targetUserID": targetUserID,
		})
		if err != nil {
			return nil, err
		}

		var mutualFriends []*models.UserSummary
		for records.Next(ctx) {
			record := records.Record()
			mutualUser := &models.UserSummary{}

			if userIDVal, ok := record.Get("userID"); ok {
				mutualUser.ID = userIDVal.(string)
			}
			if fullnameVal, ok := record.Get("fullname"); ok {
				mutualUser.FullName = fullnameVal.(string)
			}
			if usernameVal, ok := record.Get("username"); ok {
				mutualUser.Username = usernameVal.(string)
			}
			if avatarVal, ok := record.Get("avatar"); ok {
				mutualUser.Avatar = avatarVal.(string)
			}

			mutualFriends = append(mutualFriends, mutualUser)
		}

		if err = records.Err(); err != nil {
			return nil, err
		}
		return mutualFriends, nil
	})

	if err != nil {
		return nil, err
	}

	mutualFriendList, ok := result.([]*models.UserSummary)
	if !ok {
		return nil, errors.New("failed to cast result to []*models.UserSummary")
	}

	return mutualFriendList, nil
}
