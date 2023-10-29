package main

import "fmt"

type Status int

const (
	PlayerIdle    Status = 0
	PlayerPlaying        = 1
	PlayerQueue          = 2
)

const (
	GamePlaying Status = 0
	GameDone           = 1
)

type Player struct {
	id     string
	status Status
	gameId string
}

type Stone struct {
	x int
	y int
}

type Game struct {
	id        string
	playerId1 string
	playerId2 string
	stones    []Stone
	status    Status
}

type Giwon struct {
	players map[string]Player
	games   map[string]Game
}

var giwon Giwon

// TODO: maybe refactor so that theres no global giwon var?
func InitGiwon() error {
	giwon = Giwon{
		players: make(map[string]Player),
		games:   make(map[string]Game),
	}
	return nil
}

// handle gomoku commands
func CommandHandler(hub *Hub, userId string, payload string) {
	var ty string
	fmt.Sscanf(payload, "%s", &ty)
	switch ty {
	case "ENQ":
		EnqueUser(hub, userId)
	case "DEQ":
		DequeUser(hub, userId)
	case "PLC":
		var x, y int
		_, err := fmt.Sscanf(payload, "PLC %d %d", &x, &y)
		if err != nil {
			hub.Write(userId, "ERR")
		}
		PlaceStone(hub, userId, x, y)
	default:
		hub.Write(userId, "ERR")
	}
}

func EnqueUser(hub *Hub, userId string) {
}

func DequeUser(hub *Hub, userId string) {
}

func PlaceStone(hub *Hub, userId string, x int, y int) {
}
