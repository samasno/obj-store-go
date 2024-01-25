package repos

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	defaultDatabase = "object-store"
	usersColl       = "users"
)

var existingClients map[string]*mongo.Client

var clientsLock *sync.Mutex = &sync.Mutex{}

func GetMongoDBClient(uri string) (*mongo.Client, error) {
	fn := "GetMongodbClient"
	clientsLock.Lock()
	defer clientsLock.Unlock()
	c, ok := existingClients[uri]
	if ok {
		return c, nil
	}
	options := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var client *mongo.Client
	errChan := make(chan error)

	go func(errChan chan error) {
		var err error
		client, err = mongo.Connect(ctx, options)
		if err != nil {
			log.Printf("%s: %s\n", fn, err.Error())
			errChan <- err
			return
		}
		errChan <- nil
	}(errChan)

	select {
	case <-ctx.Done():
		log.Printf("%s: %s\n", fn, ctx.Err().Error())
		return nil, ctx.Err()
	case err := <-errChan:
		if err != nil {
			log.Printf("%s: %s", fn, err.Error())
			return nil, err
		}
	}

	return client, nil
}
