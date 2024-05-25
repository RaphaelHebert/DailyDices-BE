package router

import (
	"github.com/RaphaelHebert/DailyDices-BE/handler"
	"github.com/RaphaelHebert/DailyDices-BE/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {
	// app.Get("/", middleware.Protected(), hello)
	
	// User
	user := app.Group("/user")
	user.Get("/all", handler.GetAllUsers)
	user.Get("/:id", handler.GetUser)
	user.Post("/", handler.CreateUser)
	user.Delete("/:id", middleware.Protected(), handler.DeleteUser)
	user.Put("/:id", middleware.Protected(), handler.UpdateUser)

	// Auth
	app.Post("/login", handler.Login)
	app.Get("/token", middleware.Protected(), handler.Token)
	app.Get("/roll-dices", middleware.Protected(), handler.Dices)
	app.Get("/scores", middleware.Protected(),handler.Scores)
}