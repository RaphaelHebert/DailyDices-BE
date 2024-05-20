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
	"golang.org/x/crypto/bcrypt"
)

func GetUser(ctx *fiber.Ctx) error {
	uid := ctx.Query("id")
	query := bson.D{{"uid", uid}}
	cursor, err := config.Mg.Db.Collection("users").Find(ctx.Context(), query)
		if err != nil {
			fmt.Println("here")
			return ctx.Status(500).SendString(err.Error())
		}

	var users []model.User = make([]model.User, 0)

	if err := cursor.All(ctx.Context(), &users); err != nil {
		fmt.Println("here2")

		return ctx.Status(500).SendString(err.Error())

	}
		// return employees list in JSON format
	return ctx.JSON(users)
	// uid := ctx.Query("id")
	// // TODO connect to db and drop dummy data
	// err := uuid.Validate(uid)
	// if err != nil {
	// 	return ctx.SendStatus(fiber.StatusBadRequest)
	// }
	// u := db.UsersList[uid]
	// up := model.PublicUser{
	// 	Username: u.Username,
	// 	Email: u.Email,
	// 	UID: u.UID,
	// }
	// res, err := json.Marshal(up)
	// if err != nil {
	// 	return ctx.SendStatus(fiber.StatusBadRequest)
	// }
	
	// s := fmt.Sprintf("%s", res)
	// return ctx.Status(fiber.StatusOK).JSON(s)
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