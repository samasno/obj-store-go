package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email"`
	Password string             `json:"password_hash" bson:"password_hash,omitempty"`
	Created  primitive.DateTime `json:"created,omitempty" bson:"created"`
	Verified bool               `json:"verified" bson:"verified"`
	Deleted  bool               `json:"deleted" bson:"deleted"`
}
