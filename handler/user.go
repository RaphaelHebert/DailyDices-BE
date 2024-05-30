package handler

import (
	"fmt"

	"github.com/RaphaelHebert/DailyDices-BE/config"
	"github.com/RaphaelHebert/DailyDices-BE/helper"
	"github.com/RaphaelHebert/DailyDices-BE/model"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var collection = config.Mg.Db.Collection("users")

func GetUser(ctx *fiber.Ctx) error {
	var user model.User

	uid := ctx.Params("id")
	
	_id, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return ctx.SendStatus(400)
	}

	filter := bson.D{{Key: "_id", Value: _id}}

	err = collection.FindOne(ctx.Context(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Handle no documents found
			return ctx.Status(404).SendString("User not found")
		}
		// Handle other errors
		return ctx.Status(500).SendString(err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(user)
}


func UpdateUser(ctx *fiber.Ctx) error {
	params := ctx.Params("id")
	_id, err := primitive.ObjectIDFromHex(params)
	// the provided ID might be invalid ObjectID
	if err != nil {
		return ctx.SendStatus(400)
	}

	user := model.PublicUser{}
	// Parse body into struct
	err = ctx.BodyParser(&user); 
	if err != nil {
		return ctx.Status(400).SendString(err.Error())
	}

	// Check if email is available
	if uf, err := helper.GetUser("email", user.Email); err == nil {
		if uf.UID != params {
			return ctx.Status(fiber.StatusBadRequest).SendString("Email is not available")
		} 
	}

	// Check if username is available
	if uf, err := helper.GetUser("username", user.Username); err == nil {
		if uf.UID != params {
			return ctx.Status(fiber.StatusBadRequest).SendString("Username is not available")
		}
	}
	// Find the user and update its data
	query := bson.D{{Key: "_id", Value: _id}}
	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "username", Value: user.Username},
				{Key: "email", Value: user.Email},
			},
		},
	}
	
	err = collection.FindOneAndUpdate(ctx.Context(), query, update).Err()

	fmt.Println("nu: ", user)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return ctx.SendStatus(404)
		}
		return ctx.SendStatus(500)
	}
	user.UID = params
	// return the updated user
	return ctx.Status(fiber.StatusAccepted).JSON(user)
}

func DeleteUser(ctx *fiber.Ctx) error {
	params := ctx.Params("id")
	uid, err := primitive.ObjectIDFromHex(params)
	// the provided ID might be invalid ObjectID
	if err != nil {
		return ctx.SendStatus(400)
	}

	// find and delete the user with the given ID
	query := bson.D{{Key: "_id", Value: uid}}

	// TODO delete scores in scores collection ?

	result, err := collection.DeleteOne(ctx.Context(), &query)
	fmt.Println("not deleting ", result)
	if err != nil {
		return ctx.SendStatus(500)
	}

	// the user might not exist
	if result.DeletedCount < 1 {
		fmt.Println("not deleting ", uid)
		return ctx.SendStatus(404)
	}

	// the record was deleted
	return ctx.SendStatus(204)
}

func GetAllUsers(ctx *fiber.Ctx) error {
	// get all records as a cursor
	fmt.Println("ok user/all")
	query := bson.D{{}}
	cursor, err := collection.Find(ctx.Context(), query)
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}

	var users []model.User = make([]model.User, 0)

	// iterate the cursor and decode each item into a users
	if err := cursor.All(ctx.Context(), &users); err != nil {
		return ctx.Status(500).SendString(err.Error())

	}
	fmt.Println(users)

	// return users list in JSON format
	return ctx.JSON(users)
}

func CreateUser(ctx *fiber.Ctx) error {
	newUser := model.User{}
	err := ctx.BodyParser(&newUser)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// Validate the password
	if !helper.ValidatePassword(newUser.Password) {
		return ctx.Status(fiber.StatusBadRequest).SendString("Password does not meet the requirements")
	}

	// Validate the email
	if !helper.IsEmail(newUser.Email) {
		return ctx.Status(fiber.StatusBadRequest).SendString("Email is not an email")
	} 
	// Check if email is available
	if _, err = helper.GetUser("email", newUser.Email); err == nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("Email is not available")
	}
	// Check if username is available
	if _, err = helper.GetUser("username", newUser.Username); err == nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("Username is not available")
	}

	// Insert new user
	password, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)

	user := model.User{
		Username: newUser.Username,
		Email: newUser.Email,
		Password: string(password),
		UID: "",
	}

	iu, err := collection.InsertOne(ctx.Context(), user)
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}

	// Retrieve new user
	filter := bson.D{{Key: "_id", Value: iu.InsertedID}}
	createdRecord := collection.FindOne(ctx.Context(), filter)

	// Decode the Mongo record into model.User
	createdUser := &model.User{}
	createdRecord.Decode(createdUser)

	// Return the created user in JSON format
	return ctx.Status(201).JSON(createdUser)
}