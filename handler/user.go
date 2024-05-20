package handler

import (
	"encoding/json"
	"fmt"

	"github.com/RaphaelHebert/DailyDices-BE/config"
	"github.com/RaphaelHebert/DailyDices-BE/db"
	"github.com/RaphaelHebert/DailyDices-BE/model"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func GetUser(ctx *fiber.Ctx) error {
	var user model.User

	uid := ctx.Query("id")

	err := uuid.Validate(uid)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	err = config.Mg.Db.Collection("users").FindOne(ctx.Context(), bson.M{"uid": uid}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Handle no documents found
			return ctx.Status(404).SendString("User not found")
		}
		// Handle other errors
		return ctx.Status(500).SendString(err.Error())
	}

	return ctx.JSON(user)
}

func UpdateUser(ctx *fiber.Ctx) error {
	uid := ctx.Query("id")
	// TODO connect to db and drop dummy data
	err := uuid.Validate(uid)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	// check is user exists
	if _, ok := db.UsersList[uid]; !ok {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	newUser := model.User{}
	
	err = ctx.BodyParser(&newUser)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	newUser.Password = db.UsersList[uid].Password
	db.UsersList[uid] = newUser

	return ctx.Status(fiber.StatusAccepted).JSON(newUser)
}

func DeleteUser(ctx *fiber.Ctx) error {
	uid := ctx.Query("id")
	// TODO connect to db and drop dummy data
	err := uuid.Validate(uid)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	delete(db.UsersList, uid)
	return ctx.SendStatus(fiber.StatusNoContent)
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
		UID: uuid.NewString(),
	}
	db.UsersList[user.UID] = user
	db.Scores[user.UID] = []model.Score{}
	// TODO connect to db and drop dummy data
	// res, err := json.Marshal(UsersList[newUuid])
	// if err != nil {
	// 	return ctx.SendStatus(fiber.StatusBadRequest)
	// }
	return ctx.Status(fiber.StatusAccepted).JSON(model.UserPublic{
		Username: newUser.Username,
		Email: newUser.Email,
		UID: user.UID,
	})
}