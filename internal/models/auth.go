package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthTokens struct {
	ID      primitive.ObjectID `json:"_id" bson:"_id"`
	Token   string             `json:"token" bson:"token"`
	User    primitive.ObjectID `json:"user" bson:"user"`
	Created primitive.DateTime `json:"created" bson:"created"`
	Expires time.Duration      `json:"expires" bson:"expires"`
}
