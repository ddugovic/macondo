package movegen

import (
	"fmt"
	"strconv"

	"github.com/domino14/macondo/alphabet"
)

// MoveType is a type of move; a play, an exchange, pass, etc.
type MoveType uint8

const (
	MoveTypePlay MoveType = iota
	MoveTypeExchange
	MoveTypePass
	MoveTypePhonyTilesReturned

	MoveTypeEndgameTiles
	MoveTypeLostTileScore
)

// Move is a move. It can have a score, position, equity, etc. It doesn't
// have to be a scoring move.
type Move struct {
	action      MoveType
	score       int
	equity      int
	desc        string
	coords      string
	word        alphabet.MachineWord
	rowStart    uint8
	colStart    uint8
	vertical    bool
	bingo       bool
	tilesPlayed uint8
	alph        *alphabet.Alphabet
}

type ByScore []Move

func (a ByScore) Len() int           { return len(a) }
func (a ByScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByScore) Less(i, j int) bool { return a[i].score > a[j].score }

func (m Move) String() string {
	return fmt.Sprintf("<action: %v word: %v %v (%v) score: %v tp: %v vert: %v>",
		m.action, m.coords, m.word, m.word.UserVisible(m.alph), m.score,
		m.tilesPlayed, m.vertical)
}

func (m Move) uniqueKey() string {
	// Create a unique key for this play. In most cases, this will just be
	// the row/col position and the word string.
	// However, if only one tile has been played, we should only consider
	// horizontal plays
	if m.tilesPlayed != 1 {
		return fmt.Sprintf("%v%v", m.coords, m.word)
	}
	// Find the tile.
	var playedTile alphabet.MachineLetter
	var idx int
	var c alphabet.MachineLetter
	for idx, c = range m.word {
		if c != alphabet.PlayedThroughMarker {
			playedTile = c
			break
		}
	}
	var row, col uint8
	row = m.rowStart
	col = m.colStart
	// We want to get the coordinate of the tile that is on the board itself.
	if m.vertical {
		row, col = col, row
		row += uint8(idx)
	} else {
		col += uint8(idx)
	}
	key := fmt.Sprintf("%v-%v-%v", row, col, playedTile)
	return key
}

func toBoardGameCoords(row uint8, col uint8, vertical bool) string {
	colCoords := string(rune('A' + col))
	rowCoords := strconv.Itoa(int(row + 1))
	var coords string
	if vertical {
		coords = colCoords + rowCoords
	} else {
		coords = rowCoords + colCoords
	}
	return coords
}
