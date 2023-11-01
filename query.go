package main

import (
	"log"

	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func InitDB() error {
	dbc, err := sqlx.Connect("sqlite3", Env.DB_URL)
	if err != nil {
		return err
	}
	db = dbc
	return nil
}

func CheckUserName(name string) bool {
	user := User{}
	err := db.Get(
		&user, "SELECT id, name, email FROM user WHERE name = $1;",
		name,
	)
	return err == nil
}

func CheckUserEmail(email string) bool {
	user := User{}
	err := db.Get(
		&user, "SELECT id, name, email FROM user WHERE email = $1;",
		email,
	)
	return err == nil
}

func CreateUser(req RegisterReq) (User, error) {
	id := CreateId()
	hashedPassword := Hash(req.Password)
	_, err := db.Exec(
		"INSERT INTO user (id, name, password, email) VALUES($1, $2, $3, $4);",
		id, req.Name, hashedPassword, req.Email,
	)
	if err != nil {
		return User{}, err
	}
	user := User{
		Id:    id,
		Name:  req.Name,
		Email: req.Email,
	}
	return user, err
}

func LoginUser(req LoginReq) (User, error) {
	hashedPassword := Hash(req.Password)
	user := User{}
	err := db.Get(
		&user, "SELECT id, name, email FROM user WHERE name = $1 AND password = $2;",
		req.Name, hashedPassword,
	)
	if err != nil {
		return User{}, err
	}
	return user, err
}

func GetUser(userId string) (User, error) {
	user := User{}
	err := db.Get(
		&user, "SELECT id, name, email, score, status, IFNULL(gameId, '') gameId FROM user WHERE id = $1;",
		userId,
	)
	if err != nil {
		log.Printf("ERR no user %s found %s", userId, err.Error())
		return User{}, err
	}
	return user, err
}

func GetUserInQueue() ([]User, error) {
	queue := []User{}
	err := db.Select(&queue, "SELECT id, score FROM user WHERE status = $1", UserQueue)
	return queue, err
}

func SetUserStatus(userId string, status Status) error {
	_, err := db.Exec("UPDATE user SET status = $1 WHERE id = $2;", status, userId)
	if err != nil {
		log.Printf("ERR cannot change user %s status %s", userId, err.Error())
	}
	return err
}

func SetUserGameId(userId string, gameId string) error {
	_, err := db.Exec("UPDATE user SET gameId = $1 WHERE id = $2;", gameId, userId)
	if err != nil {
		log.Printf("ERR cannot set user %s gameId %s", userId, err.Error())
	}
	return err
}

func GetGame(gameId string) (Game, error) {
	game := Game{}
	err := db.Get(&game, "SELECT id, userId1, userId2 FROM game WHERE id = $1;", gameId)
	if err != nil {
		log.Printf("ERR cannot get game %s %s", gameId, err.Error())
		return Game{}, err
	}
	return game, err
}

func CreateGame(userId1 string, userId2 string) (Game, error) {
	gameId := CreateId()
	game := Game{
		Id:      gameId,
		UserId1: userId1,
		UserId2: userId2,
		Status:  GamePlaying,
	}
	_, err := db.Exec(
		"INSERT INTO game (id, userId1, userId2, status) VALUES ($1, $2, $3, $4)",
		game.Id, game.UserId1, game.UserId2, game.Status,
	)

	if err != nil {
		log.Printf("ERR cannot create game %s", err.Error())
	}

	return game, err
}

func GetStones(gameId string) ([]Stone, error) {
	stones := []Stone{}
	err := db.Select(&stones, "SELECT x, y FROM stone WHERE gameId = $1 ORDER BY placedAt ASC", gameId)

	if err != nil {
		log.Printf("ERR cannot get stones %s %s", gameId, err.Error())
	}
	return stones, err
}

func SetGameWinner(gameId string, winnerIdx int) error {
	_, err := db.Exec(
		"UPDATE game SET status = $1, winnerIdx = $2 WHERE id = $3",
		GameDone, winnerIdx, gameId,
	)

	if err != nil {
		log.Printf("ERR cannot set game winner %s %s", gameId, err.Error())
	}
	return err
}

func AppendStones(gameId string, stone Stone) error {
	_, err := db.Exec(
		"INSERT INTO stone (gameId, x, y) VALUES ($1, $2, $3);",
		gameId, stone.X, stone.Y,
	)
	return err
}
