package handler

import (
	"context"
	"errors"
	"fmt"
	"net/mail"

	"github.com/RaphaelHebert/DailyDices-BE/helper"
	"github.com/RaphaelHebert/DailyDices-BE/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// CheckPasswordHash compare password with hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getUser(k string, v string) (*model.User, error) {
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

func isEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// Login get user and password
func Login(ctx *fiber.Ctx) error {

	var input = &model.LoginInput{}
	var userData model.User

	if err := ctx.BodyParser(input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on login request", "data": err})
	}

	email := input.Email
	pass := input.Password

	var userModel *model.User
	var err error

	// user can checking by email or username
	if isEmail(email) {
		userModel, err = getUser("email", email)
	} else {
		userModel, err = getUser("username", email)
	}

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Internal Server Error", "data": err})
	} else if userModel == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid identity or password", "data": err})
	} else {
		userData = model.User{
			UID:      userModel.UID,
			Username: userModel.Username,
			Email:    userModel.Email,
			Password: userModel.Password,
		}
	}
	// check password
	if !CheckPasswordHash(pass, userData.Password) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid identity or password", "data": nil})
	}
	
	token, err := helper.CreateToken(userData.Username, userData.Email, string(userData.UID))

	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.JSON(fiber.Map{"status": "success", "message": "Success login", "data": token})
}

// Login get user and password
func Token(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	uid := claims["uid"].(string)
	fmt.Println(uid)
	fmt.Println(user)
	
	u, err := getUser("_id", uid)
	if err != nil {
		fmt.Println("no user")
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	token, err := helper.CreateToken(u.Username, u.Email, uid)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.JSON(fiber.Map{"status": "success", "message": "Success login", "data": token})
}