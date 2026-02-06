package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	COMPUTER = 0
	HUMAN    = 1
)

type Game struct {
	position *Position
	reader   *bufio.Scanner
}

func initGame() *Game {
	return &Game{
		position: initPosition(),
		reader:   bufio.NewScanner(os.Stdin),
	}
}

func (g *Game) askForPlayer(in io.Reader) int {
	g.reader = bufio.NewScanner(in)
	reader := g.readLine()
	if reader == "1" {
		return COMPUTER
	}
	if reader == "2" {
		return HUMAN
	}
	if reader == "q" {
		return -1
	}
	fmt.Println("Invalid input. Please enter 1, 2, or q")
	// fmt.Printf("\n%s\n", reader)

	return -1
}

func (g *Game) readLine() string {
	g.reader.Scan()
	return g.reader.Text()
}
