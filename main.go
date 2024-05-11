package main

import (
	"fmt"
	"log"

	"github.com/RaphaelHebert/DailyDices-BE/db"
	"github.com/RaphaelHebert/DailyDices-BE/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func main(){
	app := fiber.New()

    app.Use(cors.New())
	app.Use(logger.New())
	app.Use(requestid.New())

	router.SetupRoutes(app)
	
	// to display dummy data
	fmt.Println("mockUUID:", db.MockUUID)

	err := app.Listen(":8080")
	if err != nil {
		log.Fatal(err)
	}
}