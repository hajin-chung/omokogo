package main

import (
	"errors"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var err error

	err = InitId()
	if err != nil {
		panic(err)
	}
	InitStore()
	err = LoadEnv()
	if err != nil {
		panic(err)
	}
	err = InitDB()
	if err != nil {
		panic(err)
	}
	hub := InitHub()

	engine := html.New("./views", ".tmpl")
	app := fiber.New(fiber.Config{
		Views:        engine,
		ErrorHandler: ErrorHandler,
	})

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("hub", hub)
		return c.Next()
	})

	app.Use("/", logger.New())
	app.Static("/", "./public")
	app.Get("/", IndexController)
	app.Get("/me", ProtectedRoute, MeController)
	app.Use("/ws", UpgradeWebsocket)
	app.Get("/test", ProtectedRoute, TestController)
	app.Get("/login", LoginViewController)
	app.Get("/register", RegisterViewController)
	app.Post("/login", LoginController)
	app.Post("/register", RegisterController)
	app.Get(
		"/ws",
		WebsocketController,
		websocket.New(HandleWebsocket),
	)

	err = app.Listen(":3000")
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
