package repos

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	models "github.com/samasno/object-store/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type testConfig struct {
	URI string `json:"uri"`
}

var testEmails []string = []string{"test1@email.com", "test2@email.com"}

func TestCreateUser(t *testing.T) {
	config := getTestConfig()

	u, err := NewUserRepo(config.URI, testDatabase)
	if err != nil {
		t.Fatal(err.Error())
	}

	tu := models.User{
		Email:    testEmails[0],
		Password: "newpassword",
	}

	id, err := u.CreateUser(tu)
	if err != nil {
		t.Fatal(err.Error())
	}

	println(string(id.Hex()))
	err = u.CloseConnection()
	if err != nil {
		t.Fatal(err.Error())
	}
	cleanupTestUsers(testEmails[0])
}

func TestFetchUserByEmail(t *testing.T) {
	config := getTestConfig()

	u, err := NewUserRepo(config.URI, testDatabase)
	if err != nil {
		t.Fatal(err.Error())
	}

	createTestUsers(testEmails...)

	user, err := u.FetchUserByEmail(testEmails[0])
	if err != nil {
		t.Fatal(err.Error())
	}

	println("Got user")
	println(user.Email)
	println(user.Password)
	cleanupTestUsers(testEmails...)
}

func TestUpdateUserPassword(t *testing.T) {
	u := getTestUsersRepo()

	np := "testchangepassword"

	ids := createTestUsers(testEmails...)

	d, err := u.UpdateUserPassword(ids[0], np)
	if err != nil {
		t.Fatal(err.Error())
	}

	fmt.Printf("modified %d document\n", d)

	cleanupTestUsers(testEmails...)
}

func TestDeleteUserByID(t *testing.T) {
	config := getTestConfig()

	u, err := NewUserRepo(config.URI, testDatabase)
	if err != nil {
		t.Fatal(err.Error())
	}

	idB := [12]byte{}
	b, err := hex.DecodeString("65b2e545bb9601c2e8c0dbcc")
	if err != nil {
		t.Fatal(err.Error())
	}
	for i, bb := range b {
		idB[i] = bb
	}

	id := primitive.ObjectID(idB)
	d, err := u.DeleteUserByID(id)
	if err != nil {
		t.Fatal(err.Error())
	}

	fmt.Printf("Deleted %d user records\n", d)
}

func TestDeleteUserByEmail(t *testing.T) {
	config := getTestConfig()

	u, err := NewUserRepo(config.URI, testDatabase)
	if err != nil {
		t.Fatal(err.Error())
	}

	createTestUsers(testEmails...)

	d, err := u.DeleteUserByEmail(testEmails[0])
	if err != nil {
		t.Fatal(err.Error())
	}

	fmt.Printf("Deleted %d user records\n", d)

	cleanupTestUsers(testEmails...)
}

func TestCreateTestUsers(t *testing.T) {
	ids := createTestUsers(testEmails...)
	for _, id := range ids {
		println("Got test id", id.String())
	}
}

func TestDeleteTestUsers(t *testing.T) {
	cleanupTestUsers(testEmails...)
}

func getTestUsersRepo() *UsersRepo {
	config := getTestConfig()
	client, err := GetMongoDBClient(config.URI)
	if err != nil {
		log.Fatal(err.Error())
	}

	db := client.Database(testDatabase)

	u := &UsersRepo{
		client: client,
		db:     db,
		name:   "TestRepo",
	}

	if err != nil {
		log.Fatal(err.Error())
	}

	return u
}

func getTestConfig() testConfig {
	b, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	config := testConfig{}
	err = json.Unmarshal(b, &config)
	if err != nil {
		log.Fatal(err.Error())
	}

	return config
}

func createTestUsers(emails ...string) []primitive.ObjectID {
	u := getTestUsersRepo()
	password := "testpassword"
	ids := []primitive.ObjectID{}

	for _, email := range emails {
		user := models.User{
			Email:    email,
			Password: password,
		}

		id, err := u.CreateUser(user)
		if err != nil {
			log.Fatal()
		}
		ids = append(ids, id)
	}

	return ids
}

func cleanupTestUsers(emails ...string) {
	u := getTestUsersRepo()
	d := 0

	for _, email := range emails {
		filter := bson.D{{"email", email}}
		coll := u.db.Collection(usersColl)
		res, err := coll.DeleteOne(context.Background(), filter)
		if err != nil {
			println("cleanupTestUsers: Error deleting " + email + "\n" + err.Error())
		}

		d += int(res.DeletedCount)
	}

	fmt.Printf("cleanupTestUsers: Deleted %d users records\n", d)
	u.CloseConnection()
}
