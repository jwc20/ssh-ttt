package ttt_test

import (
	"testing"

	ttt "github.com/jwc20/ssh-ttt"
)

func TestGame_Start(t *testing.T) {
	t.Run("initializes with X as first player", func(t *testing.T) {
		game := ttt.NewTicTacToe(dummyPlayerStore)

		game.Start(2)

		if game.CurrentPlayer() != "X" {
			t.Errorf("got current player %q, want %q", game.CurrentPlayer(), "X")
		}

		if game.IsOver() {
			t.Error("game should not be over at start")
		}
	})
}

func TestGame_MakeMove(t *testing.T) {
	t.Run("alternates players after each move", func(t *testing.T) {
		game := ttt.NewTicTacToe(dummyPlayerStore)
		game.Start(2)

		assertCurrentPlayer(t, game, "X")

		game.MakeMove(1)
		assertCurrentPlayer(t, game, "O")

		game.MakeMove(2)
		assertCurrentPlayer(t, game, "X")
	})

	t.Run("returns error for occupied square", func(t *testing.T) {
		game := ttt.NewTicTacToe(dummyPlayerStore)
		game.Start(2)

		game.MakeMove(1)
		err := game.MakeMove(1)

		if err == nil {
			t.Error("expected error for occupied square")
		}
	})
}

func TestGame_WinDetection(t *testing.T) {
	t.Run("X wins with top row", func(t *testing.T) {
		game := ttt.NewTicTacToe(dummyPlayerStore)
		game.Start(2)

		game.MakeMove(1)
		game.MakeMove(4)
		game.MakeMove(2)
		game.MakeMove(5)
		game.MakeMove(3)

		assertGameOver(t, game)
		assertWinner(t, game, "X")
	})

	t.Run("O wins with left column", func(t *testing.T) {
		game := ttt.NewTicTacToe(dummyPlayerStore)
		game.Start(2)

		game.MakeMove(2)
		game.MakeMove(1)
		game.MakeMove(5)
		game.MakeMove(4)
		game.MakeMove(9)
		game.MakeMove(7)

		assertGameOver(t, game)
		assertWinner(t, game, "O")
	})

	t.Run("X wins with diagonal", func(t *testing.T) {
		game := ttt.NewTicTacToe(dummyPlayerStore)
		game.Start(2)

		game.MakeMove(1)
		game.MakeMove(2)
		game.MakeMove(5)
		game.MakeMove(3)
		game.MakeMove(9)

		assertGameOver(t, game)
		assertWinner(t, game, "X")
	})
}

func TestGame_Draw(t *testing.T) {
	t.Run("detects draw when board is full with no winner", func(t *testing.T) {
		game := ttt.NewTicTacToe(dummyPlayerStore)
		game.Start(2)

		game.MakeMove(1)
		game.MakeMove(2)
		game.MakeMove(3)
		game.MakeMove(5)
		game.MakeMove(4)
		game.MakeMove(6)
		game.MakeMove(8)
		game.MakeMove(7)
		game.MakeMove(9)

		assertGameOver(t, game)
		assertWinner(t, game, "")
	})
}

func TestGame_Finish(t *testing.T) {
	store := &ttt.StubPlayerStore{}
	game := ttt.NewTicTacToe(store)
	winner := "X"

	game.Finish(winner)
	ttt.AssertPlayerWin(t, store, winner)
}

func assertCurrentPlayer(t testing.TB, game *ttt.TicTacToe, want string) {
	t.Helper()
	if game.CurrentPlayer() != want {
		t.Errorf("got current player %q, want %q", game.CurrentPlayer(), want)
	}
}

func assertGameOver(t testing.TB, game *ttt.TicTacToe) {
	t.Helper()
	if !game.IsOver() {
		t.Error("expected game to be over")
	}
}

func assertWinner(t testing.TB, game *ttt.TicTacToe, want string) {
	t.Helper()
	if game.Winner() != want {
		t.Errorf("got winner %q, want %q", game.Winner(), want)
	}
}
