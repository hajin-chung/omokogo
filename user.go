package main

import "github.com/gofiber/contrib/websocket"

// TODO: mix some database stuff ^^
var _users []*User

const (
	UserStatusPlaying int = 0
	UserStatusIdle        = 1
	UserStatusQueue       = 2
)

type User struct {
	Id      string
	Name    string
	Score   float32
	Status  int
	GameId  string
	Sockets []*websocket.Conn
}

func NewUser(newUser *User) {
	_users = append(_users, newUser)
}
