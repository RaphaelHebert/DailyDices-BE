package router

import (
	"github.com/RaphaelHebert/DailyDices-BE/handler"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {
	// app.Get("/", middleware.Protected(), hello)
	
	// User
	user := app.Group("/user")
	user.Get("/all", handler.GetAllUsers)
	user.Get("/", handler.GetUser)
	user.Put("/", handler.CreateUser)
	user.Delete("/", handler.DeleteUser)
	user.Put("/", handler.UpdateUser)

	// Auth
	app.Post("/login", handler.Login)
	app.Get("/roll-dices", handler.Dices)
}