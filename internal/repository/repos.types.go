package repos

import "go.mongodb.org/mongo-driver/mongo"

type Repo struct {
	client      *mongo.Client
	db          *mongo.Database
	collections map[string]*mongo.Collection
	name        string
}

type UsersRepo Repo

type BlocksRepo Repo
