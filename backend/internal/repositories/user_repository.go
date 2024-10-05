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

	result, err := session.Run(ctx, "MATCH (u:User {user_id: $user_id, isActive: true}) RETURN u", map[string]interface{}{
		"user_id": id,
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
	existsEmail, _ := repo.GetUserByEmail(user.Email)
	if existsEmail != nil {
		return errors.New("email already exists")
	}

	existsPhone, _ := repo.GetUserByPhone(user.PhoneNumber)
	if existsPhone != nil {
		return errors.New("phone already exists")
	}

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

	userID := (*entityMap)["user_id"].(string)

	userProps := make(map[string]interface{})

	for key, value := range *entityMap {
		if key != "user_id" { // Không thêm user_id vào userProps
			userProps[key] = value
		}
	}

	result, err := tx.Run(ctx,
		"MATCH (u:User {user_id: $user_id}) SET u += $userProps RETURN u",
		map[string]interface{}{
			"user_id":   userID,
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
		"MATCH (u:User {user_id: $user_id}) SET u.isActive = false RETURN u",
		map[string]interface{}{
			"user_id": id,
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
