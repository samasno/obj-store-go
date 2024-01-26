package repos

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	models "github.com/samasno/object-store/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type testConfig struct {
	URI string `json:"uri"`
}

func TestCreateUser(t *testing.T) {
	config := getTestConfig()

	u, err := NewUserRepo(config.URI)
	if err != nil {
		t.Fatal(err.Error())
	}

	tu := models.User{
		Email:    "email@email.com",
		Password: "newpassword",
	}

	id, err := u.CreateUser(tu)
	if err != nil {
		t.Fatal(err.Error())
	}

	println(string(id.Hex()))
}

func TestFetchUserByEmail(t *testing.T) {
	config := getTestConfig()

	u, err := NewUserRepo(config.URI)
	if err != nil {
		t.Fatal(err.Error())
	}

	email := "email@email.com"

	user, err := u.FetchUserByEmail(email)
	if err != nil {
		t.Fatal(err.Error())
	}

	println("Got user")
	println(user.Email)
	println(user.Password)
}

func TestUniqueEmailConstraint(t *testing.T) {

}

func TestUpdateUser(t *testing.T) {

}

func TestDeleteUserByID(t *testing.T) {
	config := getTestConfig()

	u, err := NewUserRepo(config.URI)
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

	u, err := NewUserRepo(config.URI)
	if err != nil {
		t.Fatal(err.Error())
	}

	d, err := u.DeleteUserByEmail("email@email.com")
	if err != nil {
		t.Fatal(err.Error())
	}

	fmt.Printf("Deleted %d user records\n", d)
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
