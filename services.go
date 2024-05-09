package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/RaphaelHebert/DailyDices-BE/models"
	"github.com/gofor-little/env"
	"github.com/golang-jwt/jwt"
)

// Random return a number between 1 and 6 included
func random6() uint8 {
	return uint8(rand.Intn(6) + 1)
}

// rollDices takes the number i of dices to roll and returns a slice of i pseudo random numbers between 1 and 6 included.
func RollDices(i int) []models.Dice { 
 	res := make([]models.Dice, i)
	for c := range res {
		res[c] = models.Dice(random6())
	}
	return res
}

// ---------------------- TOKEN ----------------------------//
var secretKey = []byte(env.Get("SECRET_KEY", ""))

func CreateToken(username, email, uid string) (string, error) {
	var expirationTime = time.Minute * 60

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, 
        jwt.MapClaims{ 
			"username": username,
			"email": email, 
			"uid": uid,
			"exp": time.Now().Add(expirationTime).Unix(), 
        })

    tokenString, err := token.SignedString(secretKey)
    if err != nil {
    	return "", err
    }

 return tokenString, nil
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	   return secretKey, nil
	})
   
	if err != nil {
	   return err
	}
   
	if !token.Valid {
	   return fmt.Errorf("invalid token")
	}
   
	return nil
 }