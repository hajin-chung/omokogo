package main

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func UpgradeWebsocket(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func WebsocketController(c *fiber.Ctx) error {
	sess, _ := Store.Get(c)
	userId, ok := sess.Get("id").(string)
	if !ok {
		return c.Status(403).JSON(
			fiber.Map{"error": true, "message": "not authorized"})
	}

	c.Locals("userId", userId)
	return c.Next()
}

func HandleWebsocket(c *websocket.Conn) {
	userId, ok := c.Locals("userId").(string)
	if !ok {
		return
	}
	hub := c.Locals("hub").(*Hub)
	hub.Add(userId, c)
}
