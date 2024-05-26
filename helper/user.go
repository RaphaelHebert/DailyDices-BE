package helper

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"unicode"

	"github.com/RaphaelHebert/DailyDices-BE/config"
	"github.com/RaphaelHebert/DailyDices-BE/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Function to validate the password based on given requirements
func ValidatePassword(password string) bool {
	var hasMinLen, hasUpper, hasLower, hasNumber, hasSpecial bool
	const minLen = 8
	specialChars := "!@#$%^&*"

	if len(password) >= minLen {
		hasMinLen = true
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case strings.ContainsRune(specialChars, char):
			hasSpecial = true
		}
	}

	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

var collection = config.Mg.Db.Collection("users")

func GetUser(k string, v string) (*model.User, error) {
	var user model.User
	var _id primitive.ObjectID	
	var err error
	var filter primitive.D
	if k == "_id" {
		_id, err = primitive.ObjectIDFromHex(v)
		filter = bson.D{{Key: k, Value: _id}}
		if err != nil {
			return nil, err
		}
	} else {
		filter = bson.D{{Key: k, Value: v}}
	}

	
	err = collection.FindOne(context.TODO(), filter).Decode(&user)
	fmt.Println(err)
	if err != nil {
	if err == mongo.ErrNoDocuments {
		// Handle no documents found
		return nil, errors.New("No document found")
	}
		// Handle other errors
		return nil, err
	}

	return &user, nil
}

func IsEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}