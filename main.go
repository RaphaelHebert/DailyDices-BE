package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func dices(ctx *fiber.Ctx) error {
	res := fmt.Sprintf("%v", RollDices(3))
	return ctx.Status(fiber.StatusOK).JSON(res)
}

func main(){
	app := fiber.New()
	app.Get("/", dices)
	err := app.Listen(":8080")
	if err != nil {
		log.Fatal(err)
	}
}