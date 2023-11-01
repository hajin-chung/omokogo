package main

import (
	"fmt"
	"log"
	"time"
)

var dy = []int{0, 1, 1, 1, 0, -1, -1, -1}
var dx = []int{1, 1, 0, -1, -1, -1, 0, 1}

func InitGomoku(hub *Hub) {
	go MatchMaker(hub)
}

func MatchMaker(hub *Hub) {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			<-ticker.C

			queue, err := GetUserInQueue()
			if err != nil {
				log.Printf("ERR cannot get users in queue %s\n", err.Error())
				continue
			}

			var sortedQueue []string
			for _, user := range queue {
				// TODO: consider scores
				sortedQueue = append(sortedQueue, user.Id)
			}

			for len(sortedQueue) >= 2 {
				userId1 := sortedQueue[0]
				userId2 := sortedQueue[1]
				log.Printf("Match made %s %s\n", userId1, userId2)

				// create game
				game, err := CreateGame(userId1, userId2)
				if err != nil {
					break
				}
				err = SetUserGameId(userId1, game.Id)
				if err != nil {
					break
				}
				err = SetUserGameId(userId2, game.Id)
				if err != nil {
					break
				}
				err = SetUserStatus(userId1, UserPlaying)
				if err != nil {
					break
				}
				err = SetUserStatus(userId2, UserPlaying)
				if err != nil {
					break
				}

				GameStateMessage(hub, userId1, game.Id)
				GameStateMessage(hub, userId2, game.Id)

				sortedQueue = sortedQueue[2:]
			}
		}
	}()
}

func CommandHandler(hub *Hub, userId string, payload string) {
	var ty string
	fmt.Sscanf(payload, "%s", &ty)
	switch ty {
	case "STAT":
		UserStateMessage(hub, userId)
	case "GST":
		var gameId string

		_, err := fmt.Sscanf(payload, "GST %s", &gameId)
		if err != nil {
			hub.Write(userId, "ERROR wrong command format")
			return
		}

		GameStateMessage(hub, userId, gameId)
	case "ENQ":
		hub.Write(userId, EnqueUser(userId))
	case "DEQ":
		hub.Write(userId, DequeUser(userId))
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

func UserStateMessage(hub *Hub, userId string) {
	user, err := GetUser(userId)
	if err != nil {
		hub.Write(userId, "ERR cannot get user")
		return
	}

	message := fmt.Sprintf("STAT %s %d %s", user.Id, user.Status, user.GameId)
	hub.Write(userId, message)
}

func GameStateMessage(hub *Hub, userId string, gameId string) {
	game, err := GetGame(gameId)
	if err != nil {
		hub.Write(userId, "ERR cannot get game")
		return
	}

	stones, err := GetStones(gameId)
	if err != nil {
		hub.Write(userId, "ERR cannot get stones")
		return
	}

	var message = ""
	message += fmt.Sprintf(
		"GST %s %d %d %s %s ",
		game.Id, game.Status, game.WinnerIdx, game.UserId1, game.UserId2,
	)
	for _, stone := range stones {
		message += fmt.Sprintf("(%d, %d) ", stone.X, stone.Y)
	}
	hub.Write(userId, message)
}

func EnqueUser(userId string) string {
	user, err := GetUser(userId)
	if err != nil {
		return "ERR cannot get user %s"
	}

	if user.Status != UserIdle {
		return "ERR user is not idle"
	}

	err = SetUserStatus(user.Id, UserQueue)
	if err != nil {
		return "ERR cannot set user status %s"
	}
	return "ENQ"
}

func DequeUser(userId string) string {
	user, err := GetUser(userId)
	if err != nil {
		return "ERR cannot get user"
	}

	if user.Status != UserQueue {
		return "ERR user is not in queue"
	}

	err = SetUserStatus(user.Id, UserIdle)
	if err != nil {
		return "ERR cannot set user status"
	}
	return "DEQ"
}

func PlaceStone(hub *Hub, userId string, x int, y int) {
	user, err := GetUser(userId)
	if err != nil {
		hub.Write(userId, "ERR cannot get user")
		return
	}

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

	newStone := Stone{
		X: x,
		Y: y,
	}

	stones, err := GetStones(game.Id)
	if err != nil {
		hub.Write(user.Id, "ERR cannot get stones")
	}
	if len(stones)%2 == userIdx && CanPlace(stones, userIdx, newStone) {
		AppendStones(game.Id, newStone)
		stones = append(stones, newStone)

		didEnd, winnerIdx := CheckGameEnd(stones)
		if didEnd {
			SetGameWinner(game.Id, winnerIdx)
		}

		GameStateMessage(hub, game.UserId1, game.Id)
		GameStateMessage(hub, game.UserId2, game.Id)
	} else {
		hub.Write(user.Id, "ERR cannot place")
	}
}

// TODO: implement
func CanPlace(stones []Stone, userIdx int, newStone Stone) bool {
	return true
}

func CheckGameEnd(stones []Stone) (bool, int) {
	board := [15][15]int{}
	var didGameEnd = false
	var winnerIdx = 0

	for i, stone := range stones {
		board[stone.Y][stone.X] = 2 - i%2
	}

	for i, stone := range stones {
		var playerIdx = 2 - i%2
		var cnt = 0
		for k := 0; k < 8; k++ {
			var yy = stone.Y + dy[k]
			var xx = stone.X + dx[k]
			if yy >= 0 && yy < 15 && xx >= 0 && xx < 15 && board[yy][xx] == playerIdx {
				cnt++
			} else {
				break
			}
		}

		// TODO: implement Renju Rule
		if cnt >= 5 {
			didGameEnd = true
			winnerIdx = playerIdx
			break
		}
	}
	return didGameEnd, winnerIdx
}
