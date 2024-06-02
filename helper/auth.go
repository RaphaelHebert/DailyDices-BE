package helper

import (
	"fmt"
	"log"
	"time"

	"github.com/RaphaelHebert/DailyDices-BE/model"
	"github.com/gofor-little/env"
	"github.com/golang-jwt/jwt/v5"
)
var secretKey []byte
var err error


func CreateToken(isAdmin bool, scores model.Scores, username, email, uid string) (string, error) {

	err = env.Load(".env"); 
	if err != nil {
		log.Fatal("no env found")
	}
	secretKey = []byte(env.Get("SECRET_KEY", ""))

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["isadmin"] = isAdmin
	claims["scores"] = scores
	claims["username"] = username
	claims["uid"] = uid
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString(secretKey)
    if err != nil {
    	return "", err
    }

 return t, nil
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