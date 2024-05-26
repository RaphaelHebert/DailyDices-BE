package main

import (
	"log"

	"github.com/RaphaelHebert/DailyDices-BE/config"
	"github.com/RaphaelHebert/DailyDices-BE/router"
	"github.com/gofor-little/env"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)



func main(){
	// deconnect from db
	err := env.Load(".env"); 
	if err != nil {
		panic(err)
	}
	defer config.Disconnect(config.Mc)

	app := fiber.New()

    app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173", // Adjust this to your Vite app's URL
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders:    "Content-Length, Content-Type",
		AllowCredentials: true,
	}))
	app.Use(logger.New())
	app.Use(requestid.New())

	router.SetupRoutes(app)
	// TODO: updated dev and prod environment management
	port := ":8080"
	if uri := env.Get("MONGODB_URI", ""); uri != "" {
		port = ":80"
	}
	log.Fatal(app.Listen(port))
}