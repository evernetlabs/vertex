package account

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/evernetlabs/vertex/internal/app/vertex/auth"
	"github.com/evernetlabs/vertex/internal/pkg/exception"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Manager struct {
	collection    *mongo.Collection
	authenticator *auth.Authenticator
}

func NewManager(collection *mongo.Collection, authenticator *auth.Authenticator) *Manager {
	return &Manager{collection: collection, authenticator: authenticator}
}

func (m *Manager) SignUp(ctx context.Context, request *SignUpRequest) (*Account, error) {
	count, err := m.collection.CountDocuments(ctx, bson.M{
		"username": request.Username,
	})

	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, &exception.ClientError{Message: fmt.Sprintf("account with username %s already exists", request.Username)}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	account := &Account{
		ID:        bson.NewObjectID(),
		Username:  request.Username,
		Password:  string(hashedPassword),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	_, err = m.collection.InsertOne(ctx, account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *Manager) GetToken(ctx context.Context, request *TokenRequest) (*TokenResponse, error) {
	account := &Account{}

	err := m.collection.FindOne(ctx, bson.M{
		"username": request.Username,
	}).Decode(account)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, &exception.AuthenticationError{Message: "invalid username or password"}
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(request.Password))

	if err != nil {
		return nil, &exception.AuthenticationError{Message: "invalid username or password"}
	}

	token, err := m.authenticator.GenerateAccountToken(ctx, account.Username)

	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		Token: token,
	}, nil
}

func (m *Manager) Get() {

}

func (m *Manager) ChangePassword() {

}

func (m *Manager) Delete() {

}
