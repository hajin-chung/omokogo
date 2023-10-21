package main

type Stone struct {
	x       int
	y       int
	userNum int
}

const (
	GameStatusPlaying int = 0
	GameStatusDone        = 1
)

type Game struct {
	id      string
	userId1 string
	userId2 string
	status  int
	stones  []Stone
}

var _games []Game
