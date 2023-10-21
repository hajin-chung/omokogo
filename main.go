package main

import (
	"errors"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	InitId()
	InitStore()

	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})

	app.Use("/", logger.New())
	app.Static("/", "./public")

	app.Get("/user/new", CreateUserController)

	app.Use("/ws/*", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

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
