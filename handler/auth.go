// package handler

// import (
// 	"fmt"
// 	"net/http"

// 	"github.com/RaphaelHebert/DailyDices-BE/helper"
// 	"github.com/RaphaelHebert/DailyDices-BE/model"
// 	"github.com/gofiber/fiber/v2"
// )

// func Login(ctx *fiber.Ctx) error {
// 	var u model.User

// 	ctx.Set("Content-Type", "application/json")
// 	ctx.BodyParser(&u)

//   	fmt.Printf("login: user %v requesting login", u)

// 	// TODO retrieve user from DB
//   	if u.Email == "joe@mymail.com" && u.Password == "somehash" {
// 		tokenString, err := helper.CreateToken(u.Username, u.Email, u.Uid)
// 		if err != nil {
// 			ctx.SendStatus(fiber.StatusInternalServerError)
// 		}
// 		ctx.Status(fiber.StatusOK)
// 		ctx.SendString(tokenString)
// 		return nil
// 	} else {
// 		ctx.SendStatus(http.StatusUnauthorized)
// 		return nil
// 	}
// }

package handler

import (
	"errors"
	"fmt"
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
	type UserData struct {
		UID      string   `json:"Uid"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var input = &LoginInput{}
	var userData UserData

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

	fmt.Print(userModel)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Internal Server Error", "data": err})
	} else if userModel == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid identity or password", "data": err})
	} else {
		userData = UserData{
			UID:      userModel.UID,
			Username: userModel.Username,
			Email:    userModel.Email,
			Password: userModel.Password,
		}
	}

	fmt.Print(pass, userData.Password)
	// check password
	if !CheckPasswordHash(pass, userData.Password) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid identity or password", "data": nil})
	}
	
	token, _ := helper.CreateToken(userData.Username, userData.Email, string(userData.UID))

	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.JSON(fiber.Map{"status": "success", "message": "Success login", "data": token})
}