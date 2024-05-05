package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func dices(ctx *fiber.Ctx) error {
	res := fmt.Sprintf("%v", RollDices(3))
	return ctx.Status(fiber.StatusOK).JSON(res)
}

func getUser(ctx *fiber.Ctx) error {
	uid := ctx.Query("id")
	// TODO connect to db and drop dummy data
	res := fmt.Sprintf("%s", UsersList[uid])
	return ctx.Status(fiber.StatusOK).JSON(res)
}

func main(){
	app := fiber.New()

	app.Use(logger.New())
	app.Use(requestid.New())

	app.Get("/", dices)

	user := app.Group("/user")
	user.Get("/", getUser)

	err := app.Listen(":8080")
	if err != nil {
		log.Fatal(err)
	}
}