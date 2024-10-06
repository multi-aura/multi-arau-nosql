package repositories

import (
	"context"
	"errors"
	"log"
	"multiaura/internal/databases"
	"multiaura/internal/models"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type UserRepository interface {
	Repository[models.User]
	GetUsersByName(name string) ([]models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByPhone(phone string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	Follow(targetUserID, userID string) error
	UnFollow(targetUserID, userID string) error
	Block(targetUserID, userID string) error
	UnBlock(targetUserID, userID string) error
	IsFollowingOrFriend(targetUserID, userID string) (bool, error)
	IsBlocked(targetUserID, userID string) (bool, error)
	GetFriends(userID string) ([]*models.User, error)
	GetFollowers(userID string) ([]*models.UserSummary, error)
	GetFollowings(userID string) ([]*models.UserSummary, error)
}

type userRepository struct {
	db *databases.Neo4jDB
}

func NewUserRepository(db *databases.Neo4jDB) UserRepository {
	return &userRepository{db: db}
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
		if key != "userID" { // Không thêm userID vào userProps
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
		log.Println(err.Error())
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
			MERGE (u1)-[f1:FOLLOWS]->(u2)
			WITH u1, u2, f1
			OPTIONAL MATCH (u2)-[f2:FOLLOWS]->(u1)
			WITH u1, u2, f1, f2
			WHERE f2 IS NOT NULL
			MERGE (u1)-[:FRIEND_WITH]->(u2)
			MERGE (u2)-[:FRIEND_WITH]->(u1)
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
			CREATE (u2)-[:FOLLOWS]->(u1)
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

func (repo *userRepository) IsBlocked(targetUserID, userID string) (bool, error) {
	ctx := context.Background()
	session := repo.db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
			MATCH (u1:User {userID: $targetUserID})-[r:BLOCKED]->(u2:User {userID: $userID})
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
			MATCH (u:User)-[:FOLLOWS|FRIEND_WITH]->(f:User)
			WHERE u.userID = $userID
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
			MATCH (u:User)<-[:FOLLOWS|FRIEND_WITH]-(f:User)
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
