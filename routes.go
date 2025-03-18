package main

import (
	"embed"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

//go:embed views/*
var viewsfs embed.FS

func router() {
	engine := html.NewFileSystem(http.FS(viewsfs), ".html")
	app := fiber.New(
		fiber.Config{
			Views:                   engine,
			EnableIPValidation:      true,
			EnableTrustedProxyCheck: true,
			ErrorHandler:            CustomErrorHandler,
		})

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	// Asking the browser to not cache any of the requests
	app.Use(cache.New(cache.Config{
		Expiration:   1 * time.Second, // Expire in 1 second if something is cached
		CacheControl: false,
	}))

	//These Apis return HTML
	app.Get("/", requestCommitment)
	app.Post("/", processCommitmentRequest)
	app.Get("/getSecret", requestSecret)
	app.Post("/getSecret", processSecretRequest)

	//TODO: Make an Api for billing....

	//These are the Json Apis
	app.Post("/v1/api/newCommitment", newCommitmentJson)
	app.Post("/v1/api/requestSecret", processRequestSecretJson)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))

}

func requestCommitment(c *fiber.Ctx) error {
	return c.Render("views/index", fiber.Map{})
}

type CommitmentRequest struct {
	MoneroAddress string    `form:"MoneroAddress"`
	XmrAmount     float64   `form:"XmrAmount"`
	ValidTill     time.Time `form:"ValidTill"`
	Confirmations uint64    `form:"Confirmations"`
}

func processCommitmentRequest(c *fiber.Ctx) error {
	c.Accepts("application/x-www-form-urlencoded")

	reqbody := new(CommitmentRequest)

	if err := c.BodyParser(reqbody); err != nil {
		//Error occured
		// TODO return the landing page, but with

	}

	fmt.Println(reqbody)
	//TODO: this is a test
	_, commitment := generateNewCommitment()
	//TODO: check the arguments and return an error if one of them is invalid
	//TODO: validate address
	//TODO; make sure the amount is not negative number
	//TODO: Make sure valid til is in the future...
	return c.Render("views/ShowCommitment", fiber.Map{
		"Commitment": commitment,
	})
}

func requestSecret(c *fiber.Ctx) error {
	return c.SendString("Request secret")
}

func processSecretRequest(c *fiber.Ctx) error {
	return c.SendString("processed secret request")
}

// Custom redirects for errors
func CustomErrorHandler(ctx *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusNotFound

	// Retrieve the custom status code if it's a *fiber.Error
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}
	if code == 404 {

		switch ctx.Path() {

		default:
			// Send normal error page with message
			return ctx.Status(code).Render("views/errors/404", &fiber.Map{
				"Message": err.Error(),
			})

		}

	}
	// else it's a 500 error
	return ctx.Status(code).Render("views/errors/500", &fiber.Map{
		"Message": err.Error(),
	})
}

func newCommitmentJson(c *fiber.Ctx) error {
	return c.SendString("new commitment json")
}

func processRequestSecretJson(c *fiber.Ctx) error {
	return c.SendString("asd")
}
