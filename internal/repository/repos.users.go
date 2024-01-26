package repos

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	models "github.com/samasno/object-store/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func NewUserRepo(uri string) (*UsersRepo, error) {
	fn := "NewUserRepo"
	u := &UsersRepo{}
	u.name = "UserRepo"
	var err error

	u.client, err = GetMongoDBClient(uri)
	if err != nil {
		log.Printf("%s: %s\n", fn, err.Error())
		return nil, err
	}
	u.db = u.client.Database(defaultDatabase)

	log.Println("NewUserRepo connected")

	return u, nil
}

func (u *UsersRepo) CreateUser(user models.User) (primitive.ObjectID, error) { // create or update user
	fn := u.name + ".CreateUser"

	coll := u.db.Collection(usersColl)

	err := validateCreateUserInput(&user)
	if err != nil {
		log.Printf("%s: Failed to validate input: %s\n", fn, err.Error())
		return primitive.ObjectID{}, err
	}

	user.ID = primitive.ObjectID{}
	user.Created = primitive.NewDateTimeFromTime(time.Now())
	res, err := coll.InsertOne(context.Background(), user)
	if err != nil {
		log.Printf("%s: %s\n", fn, err.Error())
		return primitive.ObjectID{}, err
	}

	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Printf("%s: failed to retrieve id for new user document for %s\n", fn, user.Email)
		return user.ID, nil
	}

	log.Printf("%s: created user with email id %s\n", fn, id.String())

	return id, nil
}

func validateCreateUserInput(u *models.User) error {
	errs := []string{}

	emRgx := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	u.ID = primitive.ObjectID{}
	validEmail := emRgx.Match([]byte(u.Email))

	if !validEmail {
		msg := fmt.Sprintf("invalid email provided.")
		errs = append(errs, msg)
	}

	if len(u.Password) < 6 {
		msg := fmt.Sprintf("password must be at least 6 characters long.")
		errs = append(errs, msg)
	}

	if len(errs) > 0 {
		msg := strings.Join(errs, " ")
		return fmt.Errorf(msg)
	}

	h, err := hashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = h

	return nil
}

func hashPassword(str string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(str), 10)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (u *UsersRepo) DB() *mongo.Database {
	return u.db
}

func (u *UsersRepo) FetchUserById(_id primitive.ObjectID) (*models.User, error) {
	fn := u.name + ".FetchUserById"
	coll := u.db.Collection(usersColl)

	filter := bson.D{{"_id", _id}}

	user := models.User{}

	err := coll.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		log.Printf("%s: failed to fetch by id %s\n", fn, err.Error())
		return nil, err
	}

	return nil, nil
}

func (u *UsersRepo) FetchUserByEmail(email string) (*models.User, error) {
	fn := u.name + ".FetchUserByEmail"
	coll := u.db.Collection(usersColl)

	user := models.User{}

	filter := bson.D{{"email", email}, {"deleted", false}}
	err := coll.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		log.Printf("%s: failed to fetch user %s\n", fn, err.Error())
		return nil, err
	}

	return &user, nil
}

func (u *UsersRepo) DeleteUserByEmail(email string) (int, error) { // mark user as deleted, but don`t remove record
	fn := u.name + ".DeleteUserByEmail"

	update := bson.D{{"$set", bson.D{{"deleted", true}}}}
	d, err := updateUser(u.db, "email", email, update)
	if err != nil {
		log.Printf("%s: failed to delete user %s", fn, err.Error())
		return 0, err
	}

	log.Printf("%s: deleted %d user", fn, d)
	return d, nil
}

func (u *UsersRepo) DeleteUserByID(id primitive.ObjectID) (int, error) {
	fn := u.name + ".DeleteUserByID"

	update := bson.D{{"$set", bson.D{{"deleted", true}}}}
	d, err := updateUser(u.db, "_id", id, update)
	if err != nil {
		log.Printf("%s: failed to delete user %s", fn, err.Error())
		return 0, err
	}

	log.Printf("%s: deleted %d user", fn, d)
	return d, nil
}

func (u *UsersRepo) UpdateUser(_id primitive.ObjectID, info models.User) (*models.User, error) {
	return nil, nil
}

func updateUser[T string | primitive.ObjectID](db *mongo.Database, key string, val T, update interface{}) (int, error) {
	coll := db.Collection(usersColl)

	filter := bson.D{{key, val}}

	res, err := coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return 0, err
	}

	return int(res.ModifiedCount), nil
}

func cleanupTestUsersFromDB() {

}
