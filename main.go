package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/RaphaelHebert/DailyDices-BE/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/uuid"
)

func dices(ctx *fiber.Ctx) error {
	// get uid from token
	// user := ctx.Locals("user").(*jwt.Token)
	// claims := user.Claims.(jwt.MapClaims)
	// uid := claims["uid"].(string)

	// check user passed a number of dice
	n, err := strconv.Atoi(ctx.Query("n"))
	if err != nil {
		n = 3
	}

	dices := RollDices(n)
	// TODO add info such as dateTime
	// if uid != "" {
	// 	Scores[uuid.NewString()] = models.Score{
	// 		Uid: uid,
	// 		Score: dices,
	// 	}
	// }
	res := fmt.Sprintf("%v", dices)
	return ctx.Status(fiber.StatusOK).JSON(res)
}

func login(ctx *fiber.Ctx) error {
	var u models.User

	ctx.Set("Content-Type", "application/json")
	ctx.BodyParser(&u)

  	fmt.Printf("login: user %v requesting login", u)
  
	// TODO retrieve user from DB
  	if u.Email == "joe@mymail.com" && u.Password == "somehash" {
		tokenString, err := CreateToken(u.Username, u.Email, u.Uid)
		if err != nil {
			ctx.SendStatus(fiber.StatusInternalServerError)
		}
		ctx.Status(fiber.StatusOK)
		ctx.SendString(tokenString)
		return nil
	} else {
		ctx.SendStatus(http.StatusUnauthorized)
		return nil
	}
}


func getUser(ctx *fiber.Ctx) error {
	uid := ctx.Query("id")
	// TODO connect to db and drop dummy data
	err := uuid.Validate(uid)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	res, err := json.Marshal(UsersList[uid])
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	s := fmt.Sprintf("%s", res)
	return ctx.Status(fiber.StatusOK).JSON(s)
}

func updateUser(ctx *fiber.Ctx) error {
	uid := ctx.Query("id")
	// TODO connect to db and drop dummy data
	err := uuid.Validate(uid)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	// check is user exists
	if _, ok := UsersList[uid]; !ok {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	newUser := models.User{}
	
	err = ctx.BodyParser(&newUser)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	
	UsersList[uid] = newUser

	return ctx.Status(fiber.StatusAccepted).JSON(newUser)
}

func deleteUser(ctx *fiber.Ctx) error {
	uid := ctx.Query("id")
	// TODO connect to db and drop dummy data
	err := uuid.Validate(uid)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	delete(UsersList, uid)
	return ctx.SendStatus(fiber.StatusNoContent)
}

func getAllUsers(ctx *fiber.Ctx) error {
	// TODO connect to db and drop dummy data
	res, err := json.Marshal(UsersList)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	s := fmt.Sprintf("%s", res)
	return ctx.Status(fiber.StatusOK).JSON(s)
}

func hello(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).SendString("you accessed a protected route")
}

func createUser(ctx *fiber.Ctx) error {
	newUser := models.User{}
	
	err := ctx.BodyParser(&newUser)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	
	user := models.User{
		Username: newUser.Username,
		Email: newUser.Email,
		Password: newUser.Password,
	}
	newUuid := uuid.NewString()
	UsersList[newUuid] = user
	// TODO connect to db and drop dummy data
	// res, err := json.Marshal(UsersList[newUuid])
	// if err != nil {
	// 	return ctx.SendStatus(fiber.StatusBadRequest)
	// }
	return ctx.Status(fiber.StatusAccepted).JSON(user)
}

func main(){
	app := fiber.New()
    app.Use(cors.New())

	app.Use(logger.New())
	app.Use(requestid.New())
	app.Get("/", Protected(), hello)


	user := app.Group("/user")
	user.Get("/all", getAllUsers)
	user.Get("/", getUser)
	user.Post("/", createUser)
	user.Delete("/", deleteUser)
	user.Put("/", updateUser)

	app.Get("/login", login)
	app.Get("/roll-dices", dices)

	// to display dummy data
	fmt.Println("mockUUID:", MockUUID)

	err := app.Listen(":8080")
	if err != nil {
		log.Fatal(err)
	}
}