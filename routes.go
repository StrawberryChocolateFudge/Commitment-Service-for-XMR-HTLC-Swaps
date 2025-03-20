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

	app.Get("/checkCommitment", getCheckCommitmentPage)

	//TODO: I need a page to request commitment info where the XMR address is visible and the expiry and the amount!
	app.Get("/getSecret", requestSecret)
	app.Post("/getSecret", processSecretRequest)

	//TODO: Make an Api for billing....
	app.Get("/apikeys", getApiKeysPage)
	app.Post("/apiKeys", newApikeys)

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
	return c.Render("views/index", fiber.Map{
		"ApiKey":                  "",
		"MoneroAddress":           "",
		"XmrAmount":               "0",
		"IsDollarsChecked":        "",
		"OneHourSelected":         "",
		"TwoHoursSelected":        "",
		"FourHoursSelected":       "selected",
		"EightHoursSelected":      "",
		"TwelveHoursSelected":     "",
		"TwentyFourHoursSelected": "",
		"FortyEightHoursSelected": "",
		"Confirmations":           "10",
		"ErrorOccured":            false,
		"ErrorTitle":              "",
		"ErrorMessage":            "",
		"IsPoseidonChecked":       false,
		"ViewKey":                 "",
	})
}

func GetCommitmentRequestError(reqbody *CommitmentRequest, errorMessage string) *fiber.Map {

	isPoseidonChecked := ""
	isDollars := ""

	if reqbody.IsPoseidon {
		isPoseidonChecked = "checked"
	}

	if reqbody.IsDollars {
		isDollars = "checked"
	}

	return &fiber.Map{
		"ApiKey":                  reqbody.ApiKey,
		"MoneroAddress":           reqbody.MoneroAddress,
		"XmrAmount":               reqbody.XmrAmount,
		"IsDollarsChecked":        isDollars,
		"OneHourSelected":         GetHoursSelected(reqbody.Expiry, 1),
		"TwoHoursSelected":        GetHoursSelected(reqbody.Expiry, 2),
		"FourHoursSelected":       GetHoursSelected(reqbody.Expiry, 4),
		"EightHoursSelected":      GetHoursSelected(reqbody.Expiry, 8),
		"TwelveHoursSelected":     GetHoursSelected(reqbody.Expiry, 12),
		"TwentyFourHoursSelected": GetHoursSelected(reqbody.Expiry, 24),
		"FortyEightHoursSelected": GetHoursSelected(reqbody.Expiry, 48),
		"Confirmations":           reqbody.Confirmations,
		"ErrorOccured":            true,
		"ErrorTitle":              "Error",
		"ErrorMessage":            errorMessage,
		"IsPoseidonChecked":       isPoseidonChecked,
		"ViewKey":                 reqbody.ViewKey,
	}
}

type CommitmentRequest struct {
	ApiKey        string  `form:"ApiKey"`
	MoneroAddress string  `form:"MoneroAddress"`
	XmrAmount     float64 `form:"XmrAmount"`
	IsDollars     bool    `form:"IsDollars"`
	Expiry        uint8   `form:"Expiry"`
	Confirmations uint64  `form:"Confirmations"`
	IsPoseidon    bool    `form:"IsPoseidon"`
	ViewKey       string  `form:"ViewKey"`
}

func processCommitmentRequest(c *fiber.Ctx) error {
	c.Accepts("application/x-www-form-urlencoded")

	reqbody := new(CommitmentRequest)

	if err := c.BodyParser(reqbody); err != nil {

		reqbody.XmrAmount = 0
		reqbody.Confirmations = 10
		reqbody.Expiry = 4
		reqbody.IsPoseidon = false
		reqbody.ApiKey = ""
		reqbody.MoneroAddress = ""
		reqbody.IsDollars = false
		reqbody.ViewKey = ""
		return c.Render("views/index", GetCommitmentRequestError(reqbody, "An error occured while parsing the request"))
	}

	var hashfunc string

	if reqbody.IsPoseidon {
		hashfunc = "poseidon"
	} else {
		hashfunc = "sha256"
	}

	fmt.Println(reqbody)
	//TODO: this is a test
	_, commitment := generateNewCommitment(hashfunc)

	addressValid, err := validateAddress(reqbody.MoneroAddress)

	if !addressValid {

		if err != nil {
			//TODO: log this proper later
			fmt.Printf("Address verification err %v", err)
		}

		return c.Render("views/index", GetCommitmentRequestError(reqbody, "Unable to verify Monero Address"))
	}

	if reqbody.XmrAmount <= 0 {
		return c.Render("views/index", GetCommitmentRequestError(reqbody, "Invalid XMR amount entered"))
	}

	return c.Render("views/ShowCommitment", fiber.Map{
		"Commitment": commitment,
	})
}

func GetHoursSelected(Expiry uint8, selectFor uint8) string {
	if Expiry == selectFor {
		return "selected"
	} else {
		return ""
	}
}

func getCheckCommitmentPage(c *fiber.Ctx) error {
	return c.Render("views/checkCommitment", fiber.Map{})
}

func requestSecret(c *fiber.Ctx) error {
	return c.Render("views/requestSecret", fiber.Map{})
}

func getApiKeysPage(c *fiber.Ctx) error {
	return c.Render("views/apikeys", fiber.Map{})
}

func newApikeys(c *fiber.Ctx) error {
	return c.Render("views/apikeydetails", fiber.Map{})
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
	fmt.Println(err.Error())
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
