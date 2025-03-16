package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber"
)

func router() {
	app := fiber.New()

	app.Get("/hello", hello)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))

}

// Handler
func hello(c *fiber.Ctx) error {
	return c.SendString("I made a ☕ for you!")
}
