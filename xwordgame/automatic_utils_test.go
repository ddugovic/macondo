package xwordgame

import (
	"os"
	"testing"

	"github.com/domino14/macondo/gaddag"
	"github.com/domino14/macondo/gaddagmaker"
)

var LexiconDir = os.Getenv("LEXICON_DIR")

func TestMain(m *testing.M) {
	if _, err := os.Stat("/tmp/gen_america.gaddag"); os.IsNotExist(err) {
		gaddagmaker.GenerateGaddag(LexiconDir+"America.txt", true, true)
		os.Rename("out.gaddag", "/tmp/gen_america.gaddag")
	}
	os.Exit(m.Run())
}
func TestCompVsCompStatic(t *testing.T) {
	gd := gaddag.LoadGaddag("/tmp/gen_america.gaddag")
	game := &XWordGame{}
	game.CompVsCompStatic(gd)
	if game.turn < 6 {
		t.Errorf("Expected game.turn < 6, got %v", game.turn)
	}
}