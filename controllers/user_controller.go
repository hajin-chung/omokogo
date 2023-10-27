package controllers

import (
	"omokogo/globals"

	"github.com/gofiber/fiber/v2"
)

func MeController(c *fiber.Ctx) error {
	sess, _ := globals.Store.Get(c)
	id, ok := sess.Get("id").(string)
	if ok != true {
		return c.Status(400).SendString("login please")
	}
	return c.SendString(id)
}

func TestController(c *fiber.Ctx) error {
	sess, _ := globals.Store.Get(c)
	_, ok := sess.Get("id").(string)
	if ok != true {
		return c.Status(400).SendString("login please")
	}
	return c.Render("pages/test", fiber.Map{
		"Title": "register",
	}, "layout")
}
