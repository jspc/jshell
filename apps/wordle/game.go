package wordle

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
)

// Game represents a day's game, including guesses and the day the
// game was generated for
type Game struct {
	Word     string
	Rows     []Row
	Date     time.Time
	Complete bool
}

func (g *Game) Guess(s string) (correct bool) {
	r := Row{
		Characters: make([]Character, len(g.Word)),
	}

	for i, c := range s {
		r.Characters[i].C = byte(c)

		if g.Word[i] == byte(c) {
			r.Characters[i].Position = true
		}

		if strings.Contains(g.Word, string(c)) {
			r.Characters[i].Present = true
		}
	}

	if g.Word == s {
		g.Complete = true
	}

	g.Rows = append(g.Rows, r)

	return g.Complete
}

func (g *Game) String() string {
	sb := strings.Builder{}

	sb.WriteString(header)
	sb.WriteString(color.MagentaString(formatDate(g.Date)))
	sb.WriteString("\n\n\n")

	for idx, row := range g.Rows {
		sb.WriteString(fmt.Sprintf("%d\t%s\n", idx+1, row.String()))
	}

	return sb.String()
}

func (g *Game) Emoji() string {
	sb := strings.Builder{}
	sb.WriteString("wordle-ish ")
	sb.WriteString(formatDate(g.Date))
	sb.WriteString("\n\n")

	for _, row := range g.Rows {
		sb.WriteString(row.Emoji())
		sb.WriteString("\n")
	}

	return sb.String()
}

func NewGame() *Game {
	t := bod()

	return &Game{
		Word: loadWord(t.Unix()),
		Rows: make([]Row, 0),
		Date: t,
	}
}
