package handler

import (
	"fmt"
	"strconv"

	"github.com/RaphaelHebert/DailyDices-BE/db"
	"github.com/RaphaelHebert/DailyDices-BE/helper"
	"github.com/RaphaelHebert/DailyDices-BE/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
)

func Dices(ctx *fiber.Ctx) error {
	// get uid from token
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	uid := claims["uid"].(string)

	// check user passed a number of dice
	n, err := strconv.Atoi(ctx.Query("n"))
	if err != nil {
		n = 3
	}

	dices := helper.RollDices(n)

    // Convert byte slice to slice of integers
       
	// TODO add info such as dateTime
	s := []model.Score{}
	if(len(db.Scores[uid]) == 0 ){
		db.Scores[uid] = s
	}
	if uid != "" {
		db.Scores[uid] = append(db.Scores[uid], model.Score{
			UID: uuid.NewString(),
			Score: dices,
		})
	}
	fmt.Println(db.Scores)
	return ctx.Status(fiber.StatusOK).JSON(dices)
}