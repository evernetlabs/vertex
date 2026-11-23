package profile

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/evernetlabs/vertex/internal/app/vertex/auth"
	"github.com/evernetlabs/vertex/internal/pkg/exception"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Manager struct {
	collection    *mongo.Collection
	authenticator *auth.Authenticator
}

func NewManager(collection *mongo.Collection, authenticator *auth.Authenticator) *Manager {
	return &Manager{collection: collection, authenticator: authenticator}
}

func (m *Manager) Create(ctx context.Context, request *CreationRequest, accountUsername string) (*Profile, error) {
	count, err := m.collection.CountDocuments(ctx, bson.M{
		"handle":           request.Handle,
		"account_username": accountUsername,
	})

	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, &exception.ClientError{Message: fmt.Sprintf("profile %s already exists for account %s", request.Handle, accountUsername)}
	}

	profile := &Profile{
		ID:              bson.NewObjectID(),
		AccountUsername: accountUsername,
		Handle:          request.Handle,
		DisplayName:     request.DisplayName,
		Type:            request.Type,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}

	_, err = m.collection.InsertOne(ctx, profile)

	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (m *Manager) List(ctx context.Context, accountUsername string, page int64, size int64) ([]*Profile, error) {
	profiles := make([]*Profile, 0)

	cur, err := m.collection.Find(ctx, bson.M{
		"account_username": accountUsername,
	}, options.Find().SetSkip(page*size).SetLimit(size))

	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &profiles)

	if err != nil {
		return nil, err
	}

	return profiles, nil
}

func (m *Manager) Get(ctx context.Context, handle string, accountUsername string) (*Profile, error) {
	profile := &Profile{}

	err := m.collection.FindOne(ctx, bson.M{
		"handle":           handle,
		"account_username": accountUsername,
	}).Decode(profile)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, &exception.NotFoundError{Message: fmt.Sprintf("profile %s not found for account %s", handle, accountUsername)}
	}

	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (m *Manager) Update(ctx context.Context, handle string, request *UpdateRequest, accountUsername string) error {

	fields := bson.M{
		"updated_at": time.Now().UTC(),
	}

	if request.DisplayName != "" {
		fields["display_name"] = request.DisplayName
	}

	fields["type"] = request.Type

	res, err := m.collection.UpdateOne(ctx, bson.M{
		"handle":           handle,
		"account_username": accountUsername,
	}, bson.M{
		"$set": fields,
	})

	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return &exception.NotFoundError{Message: fmt.Sprintf("profile %s not found for account %s", handle, accountUsername)}
	}

	return nil
}

func (m *Manager) Delete(ctx context.Context, handle string, accountUsername string) error {
	res, err := m.collection.DeleteOne(ctx, bson.M{
		"handle":           handle,
		"account_username": accountUsername,
	})

	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return &exception.NotFoundError{Message: fmt.Sprintf("profile %s not found for account %s", handle, accountUsername)}
	}

	return nil
}

func (m *Manager) GetToken(ctx context.Context, handle string, accountUsername string) (*TokenResponse, error) {
	count, err := m.collection.CountDocuments(ctx, bson.M{
		"handle":           handle,
		"account_username": accountUsername,
	})

	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, &exception.NotFoundError{Message: fmt.Sprintf("profile %s not found for account %s", handle, accountUsername)}
	}

	token, err := m.authenticator.GenerateProfileToken(ctx, handle, accountUsername)

	if err != nil {
		return nil, err
	}

	return &TokenResponse{Token: token}, nil
}
