package main

type RegisterReq struct {
	Name     string
	Password string
	Email    string
}

type LoginReq struct {
	Name     string
	Password string
}
type Status int

const (
	UserDisconnected Status = 0
	UserIdle                = 1
	UserQueue               = 2
	UserPlaying             = 3
)

type User struct {
	Id     string `db:"id"`
	Name   string `db:"name"`
	Email  string `db:"email"`
	Score  int    `db:"score"`
	Status Status `db:"status"`
	GameId string `db:"gameId"`
}

const (
	GamePlaying Status = 0
	GameDone           = 1
)

type Game struct {
	Id        string `db:"id"`
	WinnerIdx int    `db:"winnerIdx"`
	UserId1   string `db:"userId1"`
	UserId2   string `db:"userId2"`
	Status    Status `db:"status"`
}

type Stone struct {
	X int
	Y int
}
