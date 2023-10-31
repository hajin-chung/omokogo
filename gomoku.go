package main

import (
	"fmt"
	"log"
	"time"
)

func InitGomoku(hub *Hub) {
	go MatchMaker(hub)
}

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
					log.Printf("ERR cannot create game %s", err.Error())
					break
				}
				err = SetUserGameId(userId1, game.Id)
				if err != nil {
					log.Printf("ERR cannot set user gameId %s", err.Error())
					break
				}
				err = SetUserGameId(userId2, game.Id)
				if err != nil {
					log.Printf("ERR cannot set user gameId %s", err.Error())
					break
				}
				err = SetUserStatus(userId1, UserPlaying)
				if err != nil {
					log.Printf("ERR cannot set user status %s", err.Error())
					break
				}
				err = SetUserStatus(userId2, UserPlaying)
				if err != nil {
					log.Printf("ERR cannot set user status %s", err.Error())
					break
				}

				stateMessage := GameStateMessage(&game, []Stone{})
				hub.Write(userId1, stateMessage)
				hub.Write(userId2, stateMessage)

				sortedQueue = sortedQueue[2:]
			}
		}
	}()
}

func CommandHandler(hub *Hub, userId string, payload string) {
	user, err := GetUser(userId)
	if err != nil {
		log.Printf("ERR no user %s found %s", userId, err.Error())
		hub.Write(userId, "ERR no user found")
		return
	}

	// try to set user stauts to idle if commands come in
	if user.Status == UserDisconnected {
		err = SetUserStatus(user.Id, UserIdle)
		if err != nil {
			log.Printf("ERR cannot change user status %s", err.Error())
			hub.Write(userId, "ERR cannot change user status")
		}
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
			log.Printf("ERR cannot get game %s %s", gameId, err.Error())
			hub.Write(userId, "ERR cannot get game")
			return
		}
		stones, err := GetStones(gameId)
		if err != nil {
			log.Printf("ERR cannot get stones %s %s", gameId, err.Error())
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

	newStone := Stone {
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
func CanPlace(stones []Stone, userIdx int, newStone Stone) bool {
	return true
}
