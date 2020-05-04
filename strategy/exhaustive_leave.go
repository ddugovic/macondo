package strategy

import (
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alecthomas/mph"
	"github.com/rs/zerolog/log"

	"github.com/domino14/macondo/alphabet"
	"github.com/domino14/macondo/board"
	"github.com/domino14/macondo/move"
)

const (
	LeaveFilename = "leaves.idx.gz"
)

// ExhaustiveLeaveStrategy should apply an equity calculation for all leaves
// exhaustively.
type ExhaustiveLeaveStrategy struct {
	leaveValues *mph.CHD
}

func float32FromBytes(bytes []byte) float32 {
	bits := binary.BigEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}

func defaultForLexicon(lexiconName string) string {
	if strings.HasPrefix(lexiconName, "CSW") ||
		strings.HasPrefix(lexiconName, "TWL") ||
		strings.HasPrefix(lexiconName, "NWL") {

		return "default_english"
	}
	return ""
}

func (els *ExhaustiveLeaveStrategy) Init(lexiconName string, alph *alphabet.Alphabet,
	strategyDir, leavefile string) error {

	if leavefile == "" {
		leavefile = LeaveFilename
	}

	file, err := os.Open(filepath.Join(strategyDir, lexiconName, leavefile))
	if err != nil {
		defdir := defaultForLexicon(lexiconName)
		file, err = os.Open(filepath.Join(strategyDir, defdir, leavefile))
		if err != nil {
			return err
		}
		log.Info().Str("leavefile", leavefile).Str("dir", defdir).Msgf(
			"no lexicon-specific strategy")
	}
	defer file.Close()
	var gz *gzip.Reader
	if strings.HasSuffix(leavefile, ".gz") {
		gz, err = gzip.NewReader(file)
		defer gz.Close()
	}
	if gz != nil {
		log.Debug().Msg("reading from compressed file")
		els.leaveValues, err = mph.Read(gz)
	} else {
		els.leaveValues, err = mph.Read(file)
	}
	if err != nil {
		return err
	}
	log.Debug().Msgf("Size of MPH: %v", els.leaveValues.Len())
	return nil
}

func NewExhaustiveLeaveStrategy(lexiconName string,
	alph *alphabet.Alphabet, strategyDir, leavefile string) (*ExhaustiveLeaveStrategy, error) {

	strategy := &ExhaustiveLeaveStrategy{}

	err := strategy.Init(lexiconName, alph, strategyDir, leavefile)
	if err != nil {
		return nil, err
	}
	return strategy, nil
}

func (els ExhaustiveLeaveStrategy) Equity(play *move.Move, board *board.GameBoard,
	bag *alphabet.Bag, oppRack *alphabet.Rack) float64 {

	leave := play.Leave()
	score := play.Score()

	leaveAdjustment := 0.0
	otherAdjustments := 0.0

	// Use global placement and endgame adjustments; this is only when
	// not overriding this with an endgame player.
	if board.IsEmpty() {
		otherAdjustments += placementAdjustment(play)
	}
	if bag.TilesRemaining() == 0 {
		otherAdjustments += endgameAdjustment(play, oppRack, bag.LetterDistribution())
	} else {
		// the leave doesn't matter if the bag is empty
		leaveAdjustment = els.LeaveValue(leave)
	}

	// also need a pre-endgame adjustment that biases towards leaving
	// one in the bag, etc.
	return float64(score) + leaveAdjustment + otherAdjustments
}

func (els ExhaustiveLeaveStrategy) LeaveValue(leave alphabet.MachineWord) float64 {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic; leave was %v\n", leave.UserVisible(alphabet.EnglishAlphabet()))
			// Panic anyway; the recover was just to figure out which leave did it.
			panic("panicking anyway")
		}
	}()
	if len(leave) == 0 {
		return 0
	}
	if len(leave) > 1 {
		sort.Slice(leave, func(i, j int) bool {
			return leave[i] < leave[j]
		})
	}
	if len(leave) <= 6 {
		// log.Debug().Msgf("Need to look up leave for %v", leave.UserVisible(alphabet.EnglishAlphabet()))
		val := els.leaveValues.Get(leave.Bytes())
		// log.Debug().Msgf("Value was %v", val)
		return float64(float32FromBytes(val))
	}
	// Only will happen if we have a pass. Passes are very rare and
	// we should ignore this a bit since there will be a negative
	// adjustment already from the fact that we're scoring 0.
	return float64(0)
}
