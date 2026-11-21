package config

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Config struct {
	ID        bson.ObjectID `bson:"_id" json:"id"`
	Key       string        `bson:"key" json:"key"`
	Value     string        `bson:"value" json:"value"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}
