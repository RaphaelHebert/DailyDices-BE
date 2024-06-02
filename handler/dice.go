package handler

import (
	"fmt"
	"strconv"
	"time"

	"github.com/RaphaelHebert/DailyDices-BE/helper"
	"github.com/RaphaelHebert/DailyDices-BE/model"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gofiber/fiber/v2"
)

func Dices(ctx *fiber.Ctx) error {
	// Get uid from token
	userToken := ctx.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	uid := claims["uid"].(string)

	// Retrieve user
	user, err := helper.GetUser("_id", uid)
	if err != nil {
		fmt.Println("no user")
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	_id, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return ctx.SendStatus(400)
	}
	

	// Check user passed a number of dice
	n, err := strconv.Atoi(ctx.Query("n"))
	if err != nil {
		n = 3
	}

	dices := helper.RollDices(n)

	scoreId := primitive.NewObjectID()
	s := model.Score{Score: dices, Date: time.Now().Unix(), ID: scoreId}
	sl := append(user.Scores, s)

	// Update User with new score
	// Find the user and update its data
	query := bson.D{{Key: "_id", Value: _id}}
	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "scores", Value: sl},
			},
		},
	}		
	err = collection.FindOneAndUpdate(ctx.Context(), query, update).Err()
	
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return ctx.SendStatus(404)
		}
		return ctx.SendStatus(500)
	}

	// return the updated user
	return ctx.Status(fiber.StatusAccepted).JSON(s)
}