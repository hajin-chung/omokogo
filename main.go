package main

import (
	"errors"
	"log"
	"omokogo/utils"
	"omokogo/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	utils.InitId()

	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})

	app.Use("/", logger.New())
	app.Static("/", "./public")
	app.Use("/ws/*", controllers.UpgradeWebsocket)

	err := app.Listen(":3000")
	if err != nil {
		log.Fatal("Failed to listen port")
	}
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	log.Println(err)
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	return c.Status(code).SendString("Internal Server Error")
}
