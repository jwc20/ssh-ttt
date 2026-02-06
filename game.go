package main

type Game struct {
	position *Position
}

func initGame() *Game {
	return &Game{
		position: initPosition(),
	}
}
