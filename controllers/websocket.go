package controllers

import (
	"log"
	"omokogo/globals"

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
	sess, _ := globals.Store.Get(c)
	userId, ok := sess.Get("id").(string)
	if !ok {
		return c.Status(403).JSON(
			fiber.Map{"error": true, "message": "not authorized"})
	}

	c.Locals("userId", userId)
	return c.Next()
}

func HandleWebsocket(c *websocket.Conn) {
	var (
		mt      int
		msg     []byte
		message string
		err     error
	)

	for {
		mt, msg, err = c.ReadMessage()
		if err != nil {
			continue
		}
		message = string(msg[:])
		log.Printf("<<< %d: %s\n", mt, message)
		c.WriteMessage(websocket.TextMessage, msg)
	}
}
