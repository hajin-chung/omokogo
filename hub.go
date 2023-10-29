package main

import (
	"fmt"
	"log"

	"github.com/gofiber/contrib/websocket"
)

/*
=== HUB ===

handles multiple socket write and read
stores websocket connections by userId
handle socket close
*/

type Command struct {
	userId  string
	payload string
}

type Hub struct {
	clients map[string]map[*websocket.Conn]bool
	command chan Command
}

func InitHub() *Hub {
	hub := Hub{
		clients: make(map[string]map[*websocket.Conn]bool),
		command: make(chan Command),
	}
	go hub.Run()

	return &hub
}

func (h *Hub) Add(userId string, conn *websocket.Conn) error {
	if h.clients[userId] == nil {
		h.clients[userId] = make(map[*websocket.Conn]bool)
	}
	h.clients[userId][conn] = true

	conn.SetCloseHandler(func(code int, text string) error {
		delete(h.clients[userId], conn)
		return nil
	})

	h.Read(userId, conn)
	return nil
}

func (h *Hub) Read(userId string, conn *websocket.Conn) {
	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		if mt != websocket.TextMessage {
			continue
		}

		message := string(msg[:])
		c := Command{
			userId:  userId,
			payload: message,
		}

		log.Printf(">>> %s: %s\n", c.userId, c.payload)
		h.command <- c
	}
}

func (h *Hub) Write(userId string, payload string) {
	log.Printf("<<< %s: %s\n", userId, payload)
	for c := range h.clients[userId] {
		c.WriteMessage(websocket.TextMessage, []byte(payload))
	}
}

func (h *Hub) Run() {
	var ty string
	for {
		c := <-h.command
		fmt.Sscanf(c.payload, "%s", &ty)
		// handle hub related commands
		switch ty {
		case "ECH":
			h.Write(c.userId, c.payload)
		case "CONN":
			h.Write(c.userId, h.StatMessage(c.userId))
		default:
			CommandHandler(h, c.userId, c.payload)
		}
	}
}

func (h *Hub) StatMessage(userId string) string {
	return fmt.Sprintf("id: %s, connections: %d", userId, len(h.clients[userId]))
}
