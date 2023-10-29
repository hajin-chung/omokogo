package main

import (
	"fmt"
	"time"
)

type Status int

const (
	PlayerIdle    Status = 0
	PlayerQueue          = 1
	PlayerPlaying        = 2
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

func (p *Player) stateMessage() string {
	return fmt.Sprintf("STAT %s %d %s", p.id, p.status, p.gameId)
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

// TODO: implement
func (g *Game) canPlace(playerIdx int, x int, y int) bool {
	return true
}

func (g *Game) stateMessage() string {
	var message = ""
	message += fmt.Sprintf("GST %s %s %s ", g.id, g.playerId1, g.playerId2)
	for _, stone := range g.stones {
		message += fmt.Sprintf("(%d, %d) ", stone.x, stone.y)
	}
	return message
}

type Giwon struct {
	players map[string]*Player
	games   map[string]*Game
}

// TODO: implement exponential matcher
func (g *Giwon) matchMaker(hub *Hub) {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			<-ticker.C

			var queue []string
			for playerId, player := range g.players {
				if player.status == PlayerQueue {
					queue = append(queue, playerId)
				}
			}

			for len(queue) >= 2 {
				playerId1 := queue[0]
				playerId2 := queue[1]

				// create game
				gameId := CreateId()
				g.games[gameId] = &Game{
					id:        gameId,
					playerId1: playerId1,
					playerId2: playerId2,
					status:    GamePlaying,
				}
				g.players[playerId1].status = PlayerPlaying
				g.players[playerId1].gameId = gameId
				g.players[playerId2].status = PlayerPlaying
				g.players[playerId2].gameId = gameId

				stateMessage := g.games[gameId].stateMessage()
				hub.Write(playerId1, stateMessage)
				hub.Write(playerId2, stateMessage)

				queue = queue[2:]
			}
		}
	}()
}

var giwon Giwon

// TODO: maybe refactor so that theres no global giwon var?
func InitGiwon(hub *Hub) error {
	giwon = Giwon{
		players: make(map[string]*Player),
		games:   make(map[string]*Game),
	}
	giwon.matchMaker(hub)

	return nil
}

// handle gomoku commands
func CommandHandler(hub *Hub, userId string, payload string) {
	if giwon.players[userId] == nil {
		giwon.players[userId] = &Player{
			id:     userId,
			status: PlayerIdle,
			gameId: "",
		}
	}

	var ty string
	fmt.Sscanf(payload, "%s", &ty)
	switch ty {
	case "STAT":
		hub.Write(userId, giwon.players[userId].stateMessage())
	case "GST":
		var gameId string

		_, err := fmt.Sscanf(payload, "GST %s", &gameId)
		if err != nil {
			hub.Write(userId, "ERROR wrong command format")
			return
		}

		game := giwon.games[gameId]
		if game == nil {
			hub.Write(userId, "ERR game not found")
		} else {
			hub.Write(userId, game.stateMessage())
		}
	case "ENQ":
		EnqueUser(hub, userId)
	case "DEQ":
		DequeUser(hub, userId)
	case "PLC":
		var x, y int
		_, err := fmt.Sscanf(payload, "PLC %d %d", &x, &y)
		if err != nil {
			hub.Write(userId, "ERR wrong command format")
			return
		}
		PlaceStone(hub, userId, x, y)
	default:
		hub.Write(userId, "ERR unknown command")
	}
}

func EnqueUser(hub *Hub, userId string) {
	player := giwon.players[userId]
	if player.status != PlayerIdle {
		hub.Write(userId, "ERR player is not idle")
		return
	}

	player.status = PlayerQueue
	hub.Write(userId, "ENQ")
}

func DequeUser(hub *Hub, userId string) {
	player := giwon.players[userId]
	if player.status != PlayerQueue {
		hub.Write(userId, "ERR player is not in queue")
		return
	}

	player.status = PlayerIdle
	hub.Write(userId, "DEQ")
}

func PlaceStone(hub *Hub, userId string, x int, y int) {
	player := giwon.players[userId]
	if player.status != PlayerPlaying {
		hub.Write(userId, "ERR player is not playing")
		return
	}

	game := giwon.games[player.gameId]
	if game == nil || game.status != GamePlaying {
		hub.Write(userId, "ERR game is not playing")
		return
	}

	var playerIdx = 0
	if game.playerId1 == userId {
		playerIdx = 0
	} else if game.playerId2 == userId {
		playerIdx = 1
	} else {
		hub.Write(userId, "ERR player not in game")
		return
	}

	if len(game.stones)%2 == playerIdx && game.canPlace(playerIdx, x, y) {
		game.stones = append(game.stones, Stone{x, y})
		hub.Write(game.playerId1, game.stateMessage())
		hub.Write(game.playerId2, game.stateMessage())

		// TODO: check if player won
	} else {
		hub.Write(userId, "ERR cannot place")
	}
}
