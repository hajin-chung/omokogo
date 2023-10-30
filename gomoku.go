package main

import (
	"fmt"
)

func UserStateMessage(u *User) string {
	return fmt.Sprintf("STAT %s %d %s", u.Id, u.Status, u.GameId)
}

// TODO: implement
func CheckGameEnd(stones []Stone) (bool, int) {
	return false, 0
}

func GameStateMessage(g *Game, stones []Stone) string {
	var message = ""
	message += fmt.Sprintf("GST %s %s %s ", g.Id, g.UserId1, g.UserId2)
	for _, stone := range stones {
		message += fmt.Sprintf("(%d, %d) ", stone.X, stone.Y)
	}
	return message
}

// TODO: implement exponential matcher
func matchMaker(hub *Hub) {
	// ticker := time.NewTicker(10 * time.Second)
	// go func() {
	// 	for {
	// 		<-ticker.C
	//
	// 		var queue []string
	// 		for playerId, player := range g.players {
	// 			if player.status == PlayerQueue {
	// 				queue = append(queue, playerId)
	// 			}
	// 		}
	//
	// 		for len(queue) >= 2 {
	// 			playerId1 := queue[0]
	// 			playerId2 := queue[1]
	//
	// 			// create game
	// 			gameId := CreateId()
	// 			g.games[gameId] = &Game{
	// 				id:        gameId,
	// 				playerId1: playerId1,
	// 				playerId2: playerId2,
	// 				status:    GamePlaying,
	// 			}
	// 			g.players[playerId1].status = PlayerPlaying
	// 			g.players[playerId1].gameId = gameId
	// 			g.players[playerId2].status = PlayerPlaying
	// 			g.players[playerId2].gameId = gameId
	//
	// 			stateMessage := g.games[gameId].stateMessage()
	// 			hub.Write(playerId1, stateMessage)
	// 			hub.Write(playerId2, stateMessage)
	//
	// 			queue = queue[2:]
	// 		}
	// 	}
	// }()
}

// handle gomoku commands
func CommandHandler(hub *Hub, userId string, payload string) {
	user, err := GetUser(userId)
	if err != nil {
		hub.Write(userId, "ERR no user found")
		return
	}

	var ty string
	fmt.Sscanf(payload, "%s", &ty)
	switch ty {
	case "STAT":
		hub.Write(userId, UserStateMessage(&user))
	case "GST":
		var gameId string

		_, err := fmt.Sscanf(payload, "GST %s", &gameId)
		if err != nil {
			hub.Write(userId, "ERROR wrong command format")
			return
		}

		game, err := GetGame(gameId)
		if err != nil {
			hub.Write(userId, "ERR game not found")
			return
		}
		stones, err := GetStones(gameId)
		if err != nil {
			hub.Write(userId, "ERR cannot get stones")
			return
		}

		hub.Write(userId, GameStateMessage(&game, stones))
	case "ENQ":
		EnqueUser(hub, user)
	case "DEQ":
		DequeUser(hub, user)
	case "PLC":
		var x, y int
		_, err := fmt.Sscanf(payload, "PLC %d %d", &x, &y)
		if err != nil {
			hub.Write(userId, "ERR wrong command format")
			return
		}
		PlaceStone(hub, user, x, y)
	default:
		hub.Write(userId, "ERR unknown command")
	}
}

func EnqueUser(hub *Hub, user User) {
	if user.Status != UserIdle {
		hub.Write(user.Id, "ERR user is not idle")
		return
	}

	err := SetUserStatus(user.Id, UserQueue)
	if err != nil {
		hub.Write(user.Id, "ERR user status cannot update")
	} else {
		hub.Write(user.Id, "ENQ")
	}
}

func DequeUser(hub *Hub, user User) {
	if user.Status != UserQueue {
		hub.Write(user.Id, "ERR user is not in queue")
		return
	}

	SetUserStatus(user.Id, UserIdle)
	hub.Write(user.Id, "DEQ")
}

func PlaceStone(hub *Hub, user User, x int, y int) {
	if user.Status != UserPlaying {
		hub.Write(user.Id, "ERR player is not playing")
		return
	}

	game, err := GetGame(user.GameId)
	if err != nil || game.Status != GamePlaying {
		hub.Write(user.Id, "ERR game is not playing")
		return
	}

	var userIdx = 0
	if game.UserId1 == user.Id {
		userIdx = 0
	} else if game.UserId2 == user.Id {
		userIdx = 1
	} else {
		hub.Write(user.Id, "ERR user not in game")
		return
	}

	stones, err := GetStones(game.Id)
	if err != nil {
		hub.Write(user.Id, "ERR cannot get stones")
	}
	if len(stones)%2 == userIdx && CanPlace(stones, userIdx, x, y) {
		AppendStones(game.Id, Stone{x, y})
		hub.Write(game.UserId1, GameStateMessage(&game, stones))
		hub.Write(game.UserId2, GameStateMessage(&game, stones))

		// TODO: check if player won
		didEnd, _ := CheckGameEnd(stones)
		if didEnd {
		}
	} else {
		hub.Write(user.Id, "ERR cannot place")
	}
}

// TODO: implement
func CanPlace(stones []Stone, userIdx int, x int, y int) bool {
	return true
}
