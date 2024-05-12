package handler

import (
	"strconv"

	"github.com/RaphaelHebert/DailyDices-BE/helper"

	"github.com/gofiber/fiber/v2"
)

func Dices(ctx *fiber.Ctx) error {
	// get uid from token
	// user := ctx.Locals("user").(*jwt.Token)
	// claims := user.Claims.(jwt.MapClaims)
	// uid := claims["uid"].(string)

	// check user passed a number of dice
	n, err := strconv.Atoi(ctx.Query("n"))
	if err != nil {
		n = 3
	}

	dices := helper.RollDices(n)

        // Convert byte slice to slice of integers
       
	// TODO add info such as dateTime
	// if uid != "" {
	// 	Scores[uuid.NewString()] = model.Score{
	// 		Uid: uid,
	// 		Score: dices,
	// 	}
	// }

	return ctx.Status(fiber.StatusOK).JSON(dices)
}