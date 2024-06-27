package database

import (
	"bkoiki950/go-store/api/utils"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID        string `json:"_id" bson:"_id,omitempty"`
	FirstName string `json:"firstname" bson:"firstname"`
	LastName  string `json:"lastname" bson:"lastname"`
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
}

func getUserFilterQuery(data User) (bson.M, error) {
	filter := bson.M{}
	if data.ID != "" {
		objectId, err := utils.ConvertStringToObjectId(data.ID)
		if err != nil {
			return filter, err
		}

		filter["_id"] = objectId
	}
	if data.Email != "" {
		filter["email"] = data.Email
	}
	if data.FirstName != "" {
		filter["firstname"] = data.FirstName
	}

	if data.LastName != "" {
		filter["lastname"] = data.LastName
	}

	return filter, nil
}

func getUserUpdateQuery(data User) bson.M {
	update := bson.M{}
	if data.Password != "" {
		update["password"] = data.Password
	}
	if data.Email != "" {
		update["email"] = data.Email
	}
	if data.FirstName != "" {
		update["firstname"] = data.FirstName
	}
	if data.LastName != "" {
		update["lastname"] = data.LastName
	}

	return bson.M{"$set": update}
}

type IUserDB interface {
	GetAllUsers(filter User) ([]User, error)
	GetUserByID(id string) (User, error)
	CreateUser(data User) (User, error)
	UpdateUser(id string, data User) (User, error)
	DeleteUser(id string) (User, error)
	LoginUser(email string, password string) (User, error)
	GetUserByEmail(email string) (User, error)
}

type UserDB struct{}

var ctx = context.TODO()

func UserCollection() *UserDB {
	var db IUserDB = &UserDB{}
	return db.(*UserDB)
}

func (c *UserDB) GetAllUsers(filter User) ([]User, error) {
	var users []User
	var userColl, err = GetCollection("users")
	if userColl == nil {
		return users, err
	}

	query, err := getUserFilterQuery(filter)
	if err != nil {
		return users, err
	}
	cursor, err := userColl.Find(ctx, query)
	if err != nil {
		return users, err
	}

	if err = cursor.All(ctx, &users); err != nil {
		return users, err
	}

	defer cursor.Close(ctx)

	if len(users) == 0 {
		return []User{}, nil
	}

	return users, nil
}

func (c *UserDB) GetUserByID(id string) (User, error) {
	var user User
	var userColl, err = GetCollection("users")
	if userColl == nil {
		return user, err
	}

	filter, err := getUserFilterQuery(User{ID: id})
	if err != nil {
		return User{}, err
	}

	err = userColl.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return User{}, utils.HandleError(err, "User not found")
	}

	return user, nil
}

func (c *UserDB) CreateUser(data User) (User, error) {
	var userColl, err = GetCollection("users")
	if userColl == nil {
		return User{}, err
	}

	var userExists User
	filter, err := getUserFilterQuery(User{Email: data.Email})
	if err != nil {
		return User{}, err
	}

	err = userColl.FindOne(ctx, filter).Decode(&userExists); if err == nil {
		return User{}, fmt.Errorf("User with email %s already exists", data.Email)
	}

	pwdHash := utils.HashPassword(data.Password)
	if pwdHash == "" {
		return User{}, utils.HandleError(nil, "Something went wrong. Please try again later")
	}

	data.Password = pwdHash
	res, err := userColl.InsertOne(ctx, data)
	if err != nil {
		return User{}, err
	}

	var user User
	err = userColl.FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&user)
	if err != nil {
		return User{}, err
	}

	idstring := utils.ConvertObjectIdToString(res.InsertedID.(primitive.ObjectID))
	user.ID = idstring

	return user, nil
}

func (c *UserDB) UpdateUser(id string, data User) (User, error) {
	var userColl, err = GetCollection("users")
	if userColl == nil {
		return User{}, err
	}
	var updatedUser User
	after := options.After
	opts := options.FindOneAndUpdateOptions{ReturnDocument: &after}
	filter, err := getUserFilterQuery(User{ID: id})
	if err != nil {
		return User{}, err
	}
	update := getUserUpdateQuery(data)
	err = userColl.FindOneAndUpdate(ctx, filter, update, &opts).Decode(&updatedUser)
	if err != nil {
		return User{}, utils.HandleError(err, "User not found")
	}

	return updatedUser, nil
}

func (c *UserDB) DeleteUser(id string) (User, error) {
	var userColl, err = GetCollection("users")
	if userColl == nil {
		return User{}, err
	}
	filter, err := getUserFilterQuery(User{ID: id})
	if err != nil {
		return User{}, err
	}
	var deletedUser User
	opts := options.FindOneAndDeleteOptions{}
	err = userColl.FindOneAndDelete(ctx, filter, &opts).Decode(&deletedUser)
	if err != nil {
		return User{}, utils.HandleError(err, "User not found")
	}

	return deletedUser, nil
}

func (c *UserDB) LoginUser(email string, password string) (User, error) {
	var user User
	var userColl, err = GetCollection("users")
	if userColl == nil {
		return user, err
	}
	filter, err := getUserFilterQuery(User{Email: email})
	if err != nil {
		return User{}, err
	}
	message := "Invalid email or password"
	err = userColl.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return User{}, utils.HandleError(err, message)
	}

	isValid := utils.ComparePassword(user.Password, password)
	if !isValid {
		return User{}, fmt.Errorf(message)
	}

	return user, nil
}

func (c *UserDB) GetUserByEmail(email string) (User, error) {
	var user User
	var userColl, err = GetCollection("users")
	if userColl == nil {
		return user, err
	}
	filter, err := getUserFilterQuery(User{Email: email})
	if err != nil {
		return User{}, err
	}
	err = userColl.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return User{}, utils.HandleError(err, "User not found")
	}

	return user, nil
}
