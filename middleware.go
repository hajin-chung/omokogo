package main

import "github.com/gofiber/fiber/v2"

func ProtectedRoute(c *fiber.Ctx) error {
	sess, _ := Store.Get(c)
	userId, ok := sess.Get("id").(string)
	if !ok {
		return c.Redirect("/login")
	}

	c.Locals("userId", userId)
	return c.Next()
}
