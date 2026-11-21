package config

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Manager struct {
	collection *mongo.Collection
}

func NewManager(collection *mongo.Collection) *Manager {
	return &Manager{collection: collection}
}

const (
	KeyJwtSigningKey  = "jwt_signing_key"
	KeyVertexEndpoint = "vertex_endpoint"
)

func (m *Manager) Init(ctx context.Context, jwtSigningKey string, vertexEndpoint string) error {
	if err := m.SetIfMissing(ctx, KeyJwtSigningKey, jwtSigningKey); err != nil {
		return err
	}

	if err := m.SetIfMissing(ctx, KeyVertexEndpoint, vertexEndpoint); err != nil {
		return err
	}

	return nil
}

func (m *Manager) SetVertexEndpoint(ctx context.Context, vertexEndpoint string) error {
	return m.Set(ctx, KeyVertexEndpoint, vertexEndpoint)
}

func (m *Manager) SetJwtSigningKey(ctx context.Context, jwtSigningKey string) error {
	return m.Set(ctx, KeyJwtSigningKey, jwtSigningKey)
}

func (m *Manager) GetVertexEndpoint(ctx context.Context) string {
	return m.Get(ctx, KeyVertexEndpoint, "localhost:8080")
}

func (m *Manager) GetJwtSigningKey(ctx context.Context) string {
	return m.Get(ctx, KeyJwtSigningKey, "secret")
}

func (m *Manager) SetIfMissing(ctx context.Context, key string, value string) error {
	_, err := m.collection.UpdateOne(ctx, bson.M{
		"key": key,
	}, bson.M{
		"$setOnInsert": bson.M{
			"_id":        bson.NewObjectID(),
			"key":        key,
			"value":      value,
			"created_at": time.Now().UTC(),
			"updated_at": time.Now().UTC(),
		},
	}, options.UpdateOne().SetUpsert(true))

	return err
}

func (m *Manager) Set(ctx context.Context, key string, value string) error {
	_, err := m.collection.UpdateOne(ctx, bson.M{
		"key": key,
	}, bson.M{
		"$set": bson.M{
			"value":      value,
			"updated_at": time.Now().UTC(),
		},
		"$setOnInsert": bson.M{
			"_id":        bson.NewObjectID(),
			"key":        key,
			"created_at": time.Now().UTC(),
		},
	}, options.UpdateOne().SetUpsert(true))

	return err
}

func (m *Manager) Get(ctx context.Context, key string, defaultValue string) string {
	var config Config

	err := m.collection.FindOne(ctx, bson.M{
		"key": key,
	}).Decode(&config)

	if err != nil {
		return defaultValue
	}

	return config.Value
}
