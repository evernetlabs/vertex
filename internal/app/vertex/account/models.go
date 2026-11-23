package account

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Account struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string        `json:"username" bson:"username"`
	Password  string        `json:"-" bson:"password"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}
