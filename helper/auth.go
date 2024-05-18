package helper

import (
	"fmt"
	"time"

	"github.com/gofor-little/env"
	"github.com/golang-jwt/jwt"
)

var secretKey = []byte(env.Get("SECRET_KEY", ""))

func CreateToken(username, email, uid string) (string, error) {
	var expirationTime = time.Minute * 600

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