package main

import (
	"errors"
	"log"
	"omokogo/controllers"
	"omokogo/globals"
	"omokogo/hub"
	"omokogo/queries"
	"omokogo/utils"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

func main() {
	var err error

	err = utils.InitId()
	if err != nil {
		panic(err)
	}
	globals.InitStore()
	err = globals.LoadEnv()
	if err != nil {
		panic(err)
	}
	err = queries.InitDB()
	if err != nil {
		panic(err)
	}
	hub := hub.Init()

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
	app.Get("/", controllers.IndexController)
	app.Get("/login", controllers.LoginViewController)
	app.Post("/login", controllers.LoginController)
	app.Get("/register", controllers.RegisterViewController)
	app.Post("/register", controllers.RegisterController)
	app.Get("/me", controllers.MeController)
	app.Get("/test", controllers.TestController)
	app.Use("/ws", controllers.UpgradeWebsocket)
	app.Get(
		"/ws",
		controllers.WebsocketController,
		websocket.New(controllers.HandleWebsocket),
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
