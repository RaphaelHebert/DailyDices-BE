package handler

import (
	"encoding/json"
	"fmt"

	"github.com/RaphaelHebert/DailyDices-BE/config"
	"github.com/RaphaelHebert/DailyDices-BE/db"
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
			return ctx.Status(500).SendString(err.Error())
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
	uid, err := primitive.ObjectIDFromHex(params)
	// the provided ID might be invalid ObjectID
	if err != nil {
		return ctx.SendStatus(400)
	}

	user := model.PublicUser{}
	// Parse body into struct
	err = ctx.BodyParser(&user); 
	if err != nil {
		fmt.Println("ehre")
		return ctx.Status(400).SendString(err.Error())
	}

	// Find the user and update its data
	query := bson.D{{Key: "_id", Value: uid}}
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

	if err != nil {
		return ctx.SendStatus(500)
	}

	// the employee might not exist
	if result.DeletedCount < 1 {
		return ctx.SendStatus(404)
	}

	// the record was deleted
	return ctx.SendStatus(204)
}

func GetAllUsers(ctx *fiber.Ctx) error {
	// TODO connect to db and drop dummy data
	// TODO update to return PublicUser
	res, err := json.Marshal(db.UsersList)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	s := fmt.Sprintf("%s", res)
	return ctx.Status(fiber.StatusOK).JSON(s)
}

func CreateUser(ctx *fiber.Ctx) error {
	newUser := model.User{}
	
	err := ctx.BodyParser(&newUser)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	
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

	filter := bson.D{{Key: "_id", Value: iu.InsertedID}}
	createdRecord := collection.FindOne(ctx.Context(), filter)

	// decode the Mongo record into model.User
	createdUser := &model.User{}
	createdRecord.Decode(createdUser)

	// return the created Employee in JSON format
	return ctx.Status(201).JSON(createdUser)
}