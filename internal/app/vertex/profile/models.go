package profile

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Profile struct {
	ID              bson.ObjectID `bson:"_id,omitempty" json:"id"`
	AccountUsername string        `bson:"account_username" json:"account_username"`
	Handle          string        `bson:"handle" json:"handle"`
	DisplayName     string        `bson:"display_name" json:"display_name"`
	Type            string        `bson:"type" json:"type"`
	CreatedAt       time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time     `bson:"updated_at" json:"updated_at"`
}
