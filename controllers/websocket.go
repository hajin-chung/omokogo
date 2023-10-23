package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/contrib/websocket"
)

func UpgradeWebsocket(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}
