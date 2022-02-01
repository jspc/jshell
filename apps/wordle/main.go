package wordle

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/pterm/pterm"
)

var (
	maxAttempts   = 6
	wordleGobPath = filepath.Join(must(os.UserHomeDir()), ".config", "jshell", "wordle.gob")

	header = color.HiCyanString(` _  _  _  _____   ______ ______         _______     _____ _______ _     _
 |  |  | |     | |_____/ |     \ |      |______ ___   |   |______ |_____|
 |__|__| |_____| |    \_ |_____/ |_____ |______     __|__ ______| |     |

James' totally copyright/trademark/IP infringing wordle clone. Ssshhh! 🤫


`)

	runMenuTemplate = &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "➡️  {{ . | cyan }}",
		Inactive: "   {{ . | cyan }}",
		Selected: "➡️  {{ . | red | cyan }}",
		Details:  "Selected: {{ . }}",
	}
)

func must(i string, err error) string {
	if err != nil {
		panic(err)
	}

	return i
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

	prompt := promptui.Select{
		Label:     "Mode",
		Items:     []string{"See my stats", "Play today's game", "Back to main menu"},
		Templates: runMenuTemplate,
	}

	var i int
	for {
		fmt.Print("\033[H\033[2J")
		fmt.Println(header)

		i, _, err = prompt.Run()
		switch i {
		case 0:
			err = w.statsMode()
		case 1:
			err = w.gameMode()
		case 2:
			return nil
		}

		if err != nil {
			return
		}
	}
}

func (w *Wordle) statsMode() (err error) {
	var (
		played int
		won    int
		lost   int

		attempts = make([]int, 6)
	)

	played = len(w.Games)

	for _, g := range w.Games {
		if g.Complete {
			won += 1
			attempts[len(g.Rows)-1] += 1
		} else if len(g.Rows) == 6 {
			lost += 1
		}
	}

	fmt.Print("\033[H\033[2J")
	fmt.Println(header)

	color.HiRed("🚨🚨 Statz 🚨🚨")
	fmt.Printf("You've played %d matches, winning %d of them (with a win rate of %.2f %%)\n",
		played, won, float64(won/played*100))
	fmt.Printf("You abandoned %d games\n\n", played-won-lost)

	color.Cyan(fmt.Sprintf("Guess Distribution\n"))
	bars := make(pterm.Bars, 6)
	for i, count := range attempts {
		bars[i] = pterm.Bar{
			Label: fmt.Sprintf("%d", i+1),
			Value: count,
		}
	}

	pterm.DefaultBarChart.WithHorizontal().WithBars(bars).Render()

	fmt.Println("press enter to quit")
	fmt.Scanln()

	return
}

func (w *Wordle) gameMode() (err error) {
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

			if !isValidWord(input) {
				return errors.New("Not a valid word")
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
