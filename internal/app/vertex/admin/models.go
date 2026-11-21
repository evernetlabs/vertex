package admin

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Admin struct {
	ID        bson.ObjectID `bson:"_id" json:"id"`
	Username  string        `bson:"username" json:"username"`
	Password  string        `bson:"password" json:"-"`
	Creator   string        `bson:"creator" json:"creator"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}
