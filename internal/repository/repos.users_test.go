package repos

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/samasno/object-store/internal/models"
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

func TestDeleteUserById(t *testing.T) {

}

func TestDeleteUserByEmail(t *testing.T) {

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
