package hex

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jspc/jshell/apps"
	"github.com/manifoldco/promptui"
)

var (
	hexGobPath = filepath.Join(must(os.UserHomeDir()), ".config", "jshell", "hex.gob")
	header     = `,,
||          ,
||/\\  _-_  \\ /
|| || || \\  \\
|| || ||/    /\\
\\ |/ \\,/  /  \;
  _/

`
)

func must(i string, err error) string {
	if err != nil {
		panic(err)
	}

	return i
}

type Hex struct {
	Games []*Game
}

func (Hex) Name() string        { return "Hex" }
func (Hex) Description() string { return "Find the hidden words in some characters 🪄" }
func (Hex) Cleanup() error      { return nil }

func (h *Hex) Run() (err error) {
	g, err := h.todaysGame()
	if err != nil {
		return
	}

	g.setWordlist()

	prompt := promptui.Prompt{
		Label:    "Guess (or ? for help)",
		Validate: g.validateGuesses,
	}

	for {
		fmt.Print("\033[H\033[2J")
		fmt.Println(g)

		result, err := prompt.Run()
		if err != nil {
			fmt.Println(err)

			time.Sleep(time.Second * 2)

			continue
		}

		switch result {
		case "?":
			h.help()

		case ":q", ":Q":
			goto byebye

		default:
			g.Guess(strings.ToUpper(result))

			err = h.saveGames()
			if err != nil {
				goto byebye
			}

		}
	}

	// super lazy label/ goto, rather than fixing my code
byebye:
	return
}

func (h *Hex) todaysGame() (g *Game, err error) {
	err = h.loadGames()
	if err != nil {
		return
	}

	t := apps.Bod()
	if len(h.Games) > 0 && h.Games[len(h.Games)-1].Date == t {
		return h.Games[len(h.Games)-1], nil
	}

	g = NewGame()
	h.Games = append(h.Games, g)

	return
}

func (h *Hex) loadGames() (err error) {
	//#nosec
	f, err := os.Open(hexGobPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			h.Games = make([]*Game, 0)

			err = nil
		}

		return
	}

	d := gob.NewDecoder(f)

	return d.Decode(&h.Games)
}

func (h *Hex) saveGames() (err error) {
	//#nosec
	f, err := os.Create(hexGobPath)
	if err != nil {
		return
	}

	e := gob.NewEncoder(f)

	return e.Encode(h.Games)
}

func (g *Hex) help() {
	fmt.Print("\033[H\033[2J")
	fmt.Println(header)
	fmt.Println()
	color.Magenta("Help\n====")
	fmt.Print(`Hex is a spelling game (gerrit? spell? hex? no?) where the aim is to find as many words
as possible, made of just the 7 letters given.

Further, each guess must contain the letter in the yellow square at least once.

Guesses must contain at least 4 letters, and can be as long as you want.

Each puzzle contains at least one 'pangram'- a word which contains every one of the given 7 letters
at least once.


`)

	color.Green("Keys\n====")
	fmt.Print(`Guesses can be made using the keyboard, pressing 'enter' to submit. Next to the guess prompt is either a red 'x' or a green tick.
This denotes whether the current guess is a valid guess (so: correct length, correct letters, real word, and not having been guessed before).

Guesses are case-insensitive.

To Quit the game, and return to the main menu, type ':q' (without quotes) and hit ener.

To see this help text again, type '?' (without quotes) and hit enter.


`)

	fmt.Println("Press enter to return to the game")

	//#nosec
	fmt.Scanln()
}

func (g *Game) validateGuesses(input string) error {
	input = strings.ToUpper(input)

	// control inputs
	if input == "?" || input == ":Q" {
		return nil
	}

	for _, c := range input {
		if !contains(g.Chars, c) {
			return fmt.Errorf("Letter '%c' is not in the letter set", c)
		}
	}

	if len(input) < 4 {
		return errors.New("Guesses must be at least 4 letters long")
	}

	if !contains([]rune(input), g.Chars[0]) {
		return errors.New("Word must contain the letter " + string(g.Chars[0]))
	}

	if g.guessed(input) {
		return errors.New("Word has already been guessed")
	}

	if !g.isValid(input) {
		return errors.New("Unknown word 🤷")
	}

	return nil
}
