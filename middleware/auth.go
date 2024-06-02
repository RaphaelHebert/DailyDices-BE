package middleware

import (
	"fmt"
	"log"

	"github.com/gofor-little/env"
	"github.com/golang-jwt/jwt"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

// Protected protect routes
func Protected() fiber.Handler {
	err := env.Load(".env"); 
	if err != nil {
		log.Fatal("could not load env")
	}
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(env.Get("SECRET_KEY", ""))},
		ErrorHandler: jwtError,
	})
}
// Middleware iAdmin is ot be user after the Protected middleware as Protected will first check for correct token
func Admin() fiber.Handler {
	err := env.Load(".env"); 
	if err != nil {
		panic(err)
	}
    return func(ctx *fiber.Ctx) error {
		// Extract token
        authHeader := ctx.Get("Authorization")
		tokenString := authHeader[len("Bearer "):]

        // Parse the token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            // Ensure the token method conforms to "SigningMethodHMAC"
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            // Return the secret key
            return []byte(env.Get("SECRET_KEY", "")), nil
        })

        if err != nil {
            return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Invalid token",
            })
        }

        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            // Check if isAdmin is true
			fmt.Println(claims)
            if isAdmin, ok := claims["isadmin"].(bool); ok && isAdmin {
                return ctx.Next()
            }
            return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error": "Admin privileges required",
            })
        } else {
			fmt.Println("hello")
		}

        return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Invalid token claims",
        })
    }
}



func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}

