package ttt_test

import (
	"bytes"
	"strings"
	"testing"

	ttt "github.com/jwc20/ssh-ttt"
)

var dummyPlayerStore = &ttt.StubPlayerStore{}
var dummyStdOut = &bytes.Buffer{}

type GameSpy struct {
	StartedWith  int
	FinishedWith string
	StartCalled  bool
	Moves        []int
	moveCount    int
	player       string
	over         bool
	winner       string
	board        string
}

func (g *GameSpy) Start(numberOfPlayers int) {
	g.StartCalled = true
	g.StartedWith = numberOfPlayers
	g.player = "X"
	g.board = "   |   |   \n-----------\n   |   |   \n-----------\n   |   |   \n"
}

func (g *GameSpy) Finish(winner string) {
	g.FinishedWith = winner
}

func (g *GameSpy) MakeMove(position int) error {
	g.Moves = append(g.Moves, position)
	g.moveCount++

	if g.moveCount >= 5 {
		g.over = true
		g.winner = "X"
	}

	if g.player == "X" {
		g.player = "O"
	} else {
		g.player = "X"
	}

	return nil
}

func (g *GameSpy) Board() string {
	return g.board
}

func (g *GameSpy) CurrentPlayer() string {
	return g.player
}

func (g *GameSpy) IsOver() bool {
	return g.over
}

func (g *GameSpy) Winner() string {
	return g.winner
}

func TestCLI(t *testing.T) {

	t.Run("start game and record moves", func(t *testing.T) {
		game := &GameSpy{}
		stdout := &bytes.Buffer{}

		in := userSends("1", "4", "2", "5", "3")
		cli := ttt.NewCLI(in, stdout, game)

		cli.PlayGame()

		assertGameStartedWith(t, game, 2)
		assertMovesEqual(t, game, []int{1, 4, 2, 5, 3})
	})

	t.Run("finish game with winner", func(t *testing.T) {
		game := &GameSpy{}
		stdout := &bytes.Buffer{}

		in := userSends("1", "4", "2", "5", "3")
		cli := ttt.NewCLI(in, stdout, game)

		cli.PlayGame()

		assertFinishCalledWith(t, game, "X")
	})

	t.Run("prints error for non-numeric input and continues", func(t *testing.T) {
		game := &GameSpy{}
		stdout := &bytes.Buffer{}

		in := userSends("abc", "1", "4", "2", "5", "3")
		cli := ttt.NewCLI(in, stdout, game)

		cli.PlayGame()

		assertOutputContains(t, stdout, ttt.BadMoveInputErrMsg)
		assertGameStartedWith(t, game, 2)
	})

	t.Run("prints error for out of range input", func(t *testing.T) {
		game := &GameSpy{}
		stdout := &bytes.Buffer{}

		in := userSends("0", "1", "4", "2", "5", "3")
		cli := ttt.NewCLI(in, stdout, game)

		cli.PlayGame()

		assertOutputContains(t, stdout, ttt.BadMoveInputErrMsg)
	})

	t.Run("prints prompt with current player", func(t *testing.T) {
		game := &GameSpy{}
		stdout := &bytes.Buffer{}

		in := userSends("1", "4", "2", "5", "3")
		cli := ttt.NewCLI(in, stdout, game)

		cli.PlayGame()

		assertOutputContains(t, stdout, "Player X, enter your move (1-9): ")
	})

	t.Run("prints win message", func(t *testing.T) {
		game := &GameSpy{}
		stdout := &bytes.Buffer{}

		in := userSends("1", "4", "2", "5", "3")
		cli := ttt.NewCLI(in, stdout, game)

		cli.PlayGame()

		assertOutputContains(t, stdout, "Player X wins!\n")
	})
}

func userSends(messages ...string) *strings.Reader {
	return strings.NewReader(strings.Join(messages, "\n") + "\n")
}

func assertGameStartedWith(t testing.TB, game *GameSpy, numberOfPlayers int) {
	t.Helper()
	if game.StartedWith != numberOfPlayers {
		t.Errorf("wanted Start called with %d but got %d", numberOfPlayers, game.StartedWith)
	}
}

func assertFinishCalledWith(t testing.TB, game *GameSpy, winner string) {
	t.Helper()
	if game.FinishedWith != winner {
		t.Errorf("wanted Finish called with %q but got %q", winner, game.FinishedWith)
	}
}

func assertGameNotStarted(t testing.TB, game *GameSpy) {
	t.Helper()
	if game.StartCalled {
		t.Error("game should not have started")
	}
}

func assertMovesEqual(t testing.TB, game *GameSpy, want []int) {
	t.Helper()
	if len(game.Moves) != len(want) {
		t.Fatalf("got %d moves %v, want %d moves %v", len(game.Moves), game.Moves, len(want), want)
	}
	for i, m := range game.Moves {
		if m != want[i] {
			t.Errorf("move %d: got %d, want %d", i, m, want[i])
		}
	}
}

func assertOutputContains(t testing.TB, stdout *bytes.Buffer, want string) {
	t.Helper()
	got := stdout.String()
	if !strings.Contains(got, want) {
		t.Errorf("output %q does not contain %q", got, want)
	}
}
