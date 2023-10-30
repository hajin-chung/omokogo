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
	Id     string
	Name   string
	Email  string
	Score  int
	Status Status
	GameId string
}

const (
	GamePlaying Status = 0
	GameDone           = 1
)

type Game struct {
	Id string
	UserId1 string
	UserId2 string
	Status Status
}

type Stone struct {
	X int
	Y int
}

