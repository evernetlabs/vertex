package admin

import (
	"context"
	"time"

	"github.com/evernetlabs/vertex/internal/app/vertex/config"
	"github.com/evernetlabs/vertex/internal/pkg/exception"
	"github.com/evernetlabs/vertex/internal/pkg/secret"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Manager struct {
	collection    *mongo.Collection
	configManager *config.Manager
}

func NewManager(collection *mongo.Collection, configManager *config.Manager) *Manager {
	return &Manager{collection: collection, configManager: configManager}
}

func (m *Manager) Init(ctx context.Context, request *InitRequest) (*Admin, error) {
	count, err := m.collection.CountDocuments(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, &exception.AuthorizationError{Message: "you are not allowed to perform this action"}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	admin := &Admin{
		ID:        bson.NewObjectID(),
		Username:  request.Username,
		Password:  string(hashedPassword),
		Creator:   request.Username,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	_, err = m.collection.InsertOne(ctx, admin)

	if err != nil {
		return nil, err
	}

	jwtSigningKey, err := secret.GenerateSecret(128)

	if err != nil {
		return nil, err
	}

	err = m.configManager.Init(ctx, jwtSigningKey, request.Vertex.Endpoint)

	if err != nil {
		return nil, err
	}

	return admin, nil
}
