package ttt

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

const (
	COMPUTER = 0
	HUMAN    = 1
)

type Game struct {
	position *Position
	in       *bufio.Scanner
	out      io.Writer
}

func initGame() *Game {
	return &Game{
		position: initPosition(),
		in:       bufio.NewScanner(os.Stdin),
		out:      os.Stdout,
	}
}

func (g *Game) askForPlayer(in io.Reader) int {
	g.in = bufio.NewScanner(in)
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

func (g *Game) askForMove() int {
	fmt.Fprint(g.out, "move: ")
	idx, _ := strconv.Atoi(g.readLine())
	if stringInSlice(idx, []int{0, 1, 2, 3, 4, 5, 6, 7, 8}) {
		if g.position.board[idx] == ' ' {
			return idx
		}
	}
	return -1
}

func (g *Game) Play() {

}

/********************************************************************/

func (g *Game) readLine() string {
	g.in.Scan()
	return g.in.Text()
}

func stringInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
