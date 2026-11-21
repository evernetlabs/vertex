package admin

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/evernetlabs/vertex/internal/app/vertex/auth"
	"github.com/evernetlabs/vertex/internal/app/vertex/config"
	"github.com/evernetlabs/vertex/internal/pkg/exception"
	"github.com/evernetlabs/vertex/internal/pkg/secret"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type Manager struct {
	collection    *mongo.Collection
	configManager *config.Manager
	authenticator *auth.Authenticator
}

func NewManager(collection *mongo.Collection, configManager *config.Manager, authenticator *auth.Authenticator) *Manager {
	return &Manager{collection: collection, configManager: configManager, authenticator: authenticator}
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

func (m *Manager) GetToken(ctx context.Context, request *TokenRequest) (*TokenResponse, error) {
	admin := &Admin{}

	err := m.collection.FindOne(ctx, bson.M{"username": request.Username}).Decode(admin)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, &exception.AuthenticationError{Message: "invalid username and password combination"}
	}

	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(request.Password))

	if err != nil {
		return nil, &exception.AuthenticationError{Message: "invalid username and password combination"}
	}

	token, err := m.authenticator.GenerateAdminToken(ctx, admin.Username)

	if err != nil {
		return nil, err
	}

	return &TokenResponse{Token: token}, nil
}

func (m *Manager) Get(ctx context.Context, username string) (*Admin, error) {
	admin := &Admin{}

	err := m.collection.FindOne(ctx, bson.M{"username": username}).Decode(admin)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, &exception.NotFoundError{Message: fmt.Sprintf("admin %s not found", username)}
	}

	if err != nil {
		return nil, err
	}

	return admin, nil
}

func (m *Manager) ChangePassword(ctx context.Context, username string, request *PasswordChangeRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"password":   string(hashedPassword),
			"updated_at": time.Now().UTC(),
		},
	}

	result, err := m.collection.UpdateOne(ctx, bson.M{"username": username}, update)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return &exception.NotFoundError{Message: fmt.Sprintf("admin %s not found", username)}
	}

	return nil
}

func (m *Manager) Add(ctx context.Context, request *AdditionRequest, creator string) (*AdditionResponse, error) {
	count, err := m.collection.CountDocuments(ctx, bson.M{"username": request.Username})

	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, &exception.ClientError{Message: fmt.Sprintf("admin %s already exists", request.Username)}
	}

	password, err := secret.GenerateSecret(12)

	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	admin := &Admin{
		ID:        bson.NewObjectID(),
		Username:  request.Username,
		Password:  string(hashedPassword),
		Creator:   creator,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	_, err = m.collection.InsertOne(ctx, admin)

	if err != nil {
		return nil, err
	}

	return &AdditionResponse{Password: password, Admin: admin}, nil
}

func (m *Manager) List(ctx context.Context, page int64, size int64) ([]*Admin, error) {
	admins := make([]*Admin, 0)

	cur, err := m.collection.Find(ctx, bson.M{}, options.Find().SetSkip(page*size).SetLimit(size))

	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &admins)

	if err != nil {
		return nil, err
	}

	return admins, nil
}

func (m *Manager) ResetPassword(ctx context.Context, username string) (*PasswordResetResponse, error) {

	password, err := secret.GenerateSecret(12)

	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$set": bson.M{
			"password":   string(hashedPassword),
			"updated_at": time.Now().UTC(),
		},
	}

	result, err := m.collection.UpdateOne(ctx, bson.M{"username": username}, update)

	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, &exception.NotFoundError{Message: fmt.Sprintf("admin %s not found", username)}
	}

	return &PasswordResetResponse{Password: password}, nil
}

func (m *Manager) Delete(ctx context.Context, username string) error {
	result, err := m.collection.DeleteOne(ctx, bson.M{"username": username})

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return &exception.NotFoundError{Message: fmt.Sprintf("admin %s not found", username)}
	}

	return nil
}
