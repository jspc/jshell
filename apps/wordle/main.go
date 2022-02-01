package wordle

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

var (
	maxAttempts   = 6
	wordleGobPath = filepath.Join(must(os.UserHomeDir()), ".config", "jshell", "wordle.gob")

	header = color.HiCyanString(` _  _  _  _____   ______ ______         _______     _____ _______ _     _
 |  |  | |     | |_____/ |     \ |      |______ ___   |   |______ |_____|
 |__|__| |_____| |    \_ |_____/ |_____ |______     __|__ ______| |     |

James' totally copyright/trademark/IP infringing wordle clone. Ssshhh! 🤫


`)
)

func must(i string, err error) string {
	if err != nil {
		panic(err)
	}

	return i
}

// Character holds a specific character guess in a
// specific row.
//
// Value 'C' holds the specific character
//       'Position' holds whether the character is in the correct position
//       'Present' holds whether the character is present in the word, but wrong pos
type Character struct {
	C        byte
	Position bool
	Present  bool
}

func (c Character) String() string {
	if c.Position {
		return color.New(color.BgHiGreen, color.FgHiBlack).Sprint(string(c.C))
	}

	if c.Present {
		return color.New(color.BgYellow, color.FgHiWhite).Sprint(string(c.C))
	}

	return color.HiWhiteString(string(c.C))
}

func (c Character) Emoji() string {
	if c.Position {
		return "🟩"
	}

	if c.Present {
		return "🟨"
	}

	return "⬛"
}

// Row holds a slice of Character types representing a guess from a user
type Row struct {
	Characters []Character
}

func (r Row) String() string {
	cs := make([]string, len(r.Characters))
	for i := range cs {
		cs[i] = r.Characters[i].String()
	}

	return strings.Join(cs, "  ")
}

func (r Row) Emoji() string {
	sb := strings.Builder{}

	for _, c := range r.Characters {
		sb.WriteString(c.Emoji())
	}

	return sb.String()
}

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

// Wordle holds the state for a set of games
type Wordle struct {
	Games []*Game
}

func (w *Wordle) loadGames() (err error) {
	f, err := os.Open(wordleGobPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			w.Games = make([]*Game, 0)

			err = nil
		}

		return
	}

	d := gob.NewDecoder(f)

	return d.Decode(&w.Games)
}

func (w *Wordle) saveGames() (err error) {
	f, err := os.Create(wordleGobPath)
	if err != nil {
		return
	}

	e := gob.NewEncoder(f)

	return e.Encode(w.Games)
}

func (Wordle) Name() string {
	return "Wordle-ish"
}

func (Wordle) Description() string {
	return "James' totally copyright/trademark/IP infringing wordle clone. Ssshhh! 🤫"
}

func (Wordle) Cleanup() error {
	return nil
}

func (w *Wordle) Run() (err error) {
	err = w.loadGames()
	if err != nil {
		return
	}

	t := bod()
	var g *Game
	if len(w.Games) > 0 && w.Games[len(w.Games)-1].Date == t {
		g = w.Games[len(w.Games)-1]
	} else {
		g = NewGame()
		w.Games = append(w.Games, g)
	}

	for len(g.Rows) < 6 {
		fmt.Print("\033[H\033[2J")
		fmt.Println(g)

		if g.Complete {
			break
		}

		validate := func(input string) error {
			if len(input) != 5 {
				return errors.New("Guesses must be 5 letters long")
			}

			return nil
		}

		prompt := promptui.Prompt{
			Label:    "Guess",
			Validate: validate,
		}

		result, err := prompt.Run()
		if err != nil {
			fmt.Println(err)

			time.Sleep(time.Second * 2)

			continue
		}

		g.Guess(result)

		err = w.saveGames()
		if err != nil {
			panic(err)
		}
	}

	switch g.Complete {
	case true:
		fmt.Printf("Congrats! You got it in %d\n", len(g.Rows))

	case false:
		fmt.Println("Better luck next time")
	}

	fmt.Println("\n\nshare:\n----------\n")
	fmt.Println(g.Emoji())
	fmt.Println("\n----------\n")

	fmt.Println("press enter to quit")
	fmt.Scanln()

	return nil
}

func bod() time.Time {
	t := time.Now()
	year, month, day := t.Date()

	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func loadWord(seed int64) string {
	rand.Seed(seed)

	return words[rand.Intn(len(words))]
}

// hattip: https://stackoverflow.com/a/28890625
func formatDate(t time.Time) string {
	suffix := "th"
	switch t.Day() {
	case 1, 21, 31:
		suffix = "st"
	case 2, 22:
		suffix = "nd"
	case 3, 23:
		suffix = "rd"
	}
	return t.Format("Monday 2" + suffix + " January, 2006")
}
