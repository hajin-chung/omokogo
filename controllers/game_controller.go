package controllers

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/contrib/websocket"
)

// /ws/game/:id
func GameConnect(c *websocket.Conn) {
	gameId := c.Params("id")
	log.Println(gameId)
	// check game exists

	// check auth
	// userId := sess.Get("id")
	// if userId == nil { return }
	userId := "123213213"

	// sub to gameId

	errorChannel := make(chan error)
	go handleMessage(c, userId, gameId, errorChannel)
	err := <-errorChannel
	log.Println(err)
}

/*
Input Message Formats
1. place stone
	PLACE y x

Output Message Formats
1. error 
	ERROR message
*/

// handles messsage comming from websocket connections and
// sends individual messages that needs direct feedback (error messages)
func handleMessage(c *websocket.Conn, userId string, gameId string, errorChannel chan error) {
	var (
		mt      int
		msg     []byte
		message string
		err     error
	)

	for {
		mt, msg, err = c.ReadMessage()
		if err != nil {
			errorChannel <- err
			break
		}

		if mt != websocket.TextMessage {
			// TODO: it might be good to handle error?
			_ = c.WriteMessage(websocket.TextMessage, []byte("ERROR not a text message"))
			continue
		}

		message = string(msg[:])
		log.Printf("RECV %s : %s\n", userId, message)
	
		switch {
		// maybe store "PLACE" in some constant
		case strings.HasPrefix(message, "PLACE"):
			var x, y int
			_, err = fmt.Sscanf(message, "PLACE %d %d", &y, &x)
			if err != nil {
				_ = c.WriteMessage(websocket.TextMessage, []byte("ERROR wrong format"))
				continue
			}

			// err = giwon.Place(gameId, userId, y, x)
			// if err != nil {  }
		}
	}
}
