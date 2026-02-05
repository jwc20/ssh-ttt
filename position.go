package main

import "strings"

type Position struct {
	board string
	turn  string
}

func initPosition() *Position {
	return &Position{
		board: strings.Repeat(" ", 9),
		turn:  "x",
	}
}
