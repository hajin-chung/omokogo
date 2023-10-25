package controllers

import "github.com/gofiber/fiber/v2"

func IndexController(c *fiber.Ctx) error {
	return c.Render("pages/index", fiber.Map{
		"Title": "hi",
	}, "layout")
}
