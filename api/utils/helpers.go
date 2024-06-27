package utils

import (
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func ConvertStringToObjectId(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}

func ConvertObjectIdToString(id primitive.ObjectID) string {
	return id.Hex()
}

func NewError(message string) error {
	return errors.New(message)
}

func HandleError(err error, message string) error {
	if err == nil && message != ""{
		return NewError(message)
	}

	if err == nil {
		return nil
	}

	if message == "" {
		message = "Data not found"
	}

	if err.Error() == "mongo: no documents in result" {
		return NewError(message)
	}

	return err
}

func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); if err != nil {
		return ""
	}

	return string(hash)
}

func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	fmt.Println(err)
	return err == nil
}
