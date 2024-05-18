package handler

import (
	"errors"
	"net/mail"

	"github.com/RaphaelHebert/DailyDices-BE/db"
	"github.com/RaphaelHebert/DailyDices-BE/helper"
	"github.com/RaphaelHebert/DailyDices-BE/model"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// CheckPasswordHash compare password with hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getUserByEmail(e string) (*model.User, error) {
	// TODO: connect to DB
	// db := database.DB
	// var user model.User
	// if err := db.Where(&model.User{Email: e}).First(&user).Error; err != nil {
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		return nil, nil
	// 	}
	// 	return nil, err
	// }

	// return dummy data
	for _, value :=  range db.UsersList {
		if value.Email == e {
			return &value, nil
		}
	}
	return nil, errors.New("user not found")
}

func getUserByUsername(u string) (*model.User, error) {
	// TODO connect to db
	// db := database.DB
	// var user model.User
	// if err := db.Where(&model.User{Username: u}).First(&user).Error; err != nil {
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		return nil, nil
	// 	}
	// 	return nil, err
	// }
	// return &user, nil

	// return dummy data
	for _, value :=  range db.UsersList {
		if value.Username == u {
			return &value, nil
		}
	}
	return nil, errors.New("user not found")
}

func isEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// Login get user and password
func Login(ctx *fiber.Ctx) error {
	type LoginInput struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	var input = &LoginInput{}
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
		userModel, err = getUserByEmail(email)
	} else {
		userModel, err = getUserByUsername(email)
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