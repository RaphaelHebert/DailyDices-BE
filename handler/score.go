package handler

import (
	"github.com/RaphaelHebert/DailyDices-BE/db"
	"github.com/golang-jwt/jwt/v5"

	"github.com/gofiber/fiber/v2"
)

func Scores(ctx *fiber.Ctx) error {
	// get uid from token
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	uid := claims["uid"].(string)

	r := db.Scores[uid]
	
	return ctx.Status(fiber.StatusOK).JSON(r)
}