package ttt_test

import (
	"strings"
	"testing"

	ttt "github.com/jwc20/ssh-ttt"
)

func TestCLI(t *testing.T) {

	t.Run("record chris win from user input", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &ttt.StubPlayerStore{}

		cli := ttt.NewCLI(playerStore, in)
		cli.PlayTTT()

		ttt.AssertPlayerWin(t, playerStore, "Chris")
	})

	t.Run("record cleo win from user input", func(t *testing.T) {
		in := strings.NewReader("Cleo wins\n")
		playerStore := &ttt.StubPlayerStore{}

		cli := ttt.NewCLI(playerStore, in)
		cli.PlayTTT()

		ttt.AssertPlayerWin(t, playerStore, "Cleo")
	})

}
