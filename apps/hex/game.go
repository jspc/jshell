package hex

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/fatih/color"
)

type Game struct {
	Date  time.Time
	Chars []rune
	Words []string

	wordlist map[string]int
}

func NewGame() (g *Game) {
	t := bod()

	g = &Game{
		Date:  t,
		Chars: selectChars(t),
		Words: make([]string, 0),
	}

	g.setWordlist()

	return
}

func (g *Game) setWordlist() {
	indices := mappings[g.Chars[0]]

	g.wordlist = make(map[string]int)
	for _, i := range indices {
		w := validWords[i]

		valid := true
		for _, c := range distinctLetters(w) {
			if !contains(g.Chars, c) {
				valid = false

				break
			}
		}

		if valid {
			g.wordlist[w] = 1
		}
	}
}

func (g Game) isValid(in string) bool {
	_, ok := g.wordlist[in]

	return ok
}

func (g *Game) isPangram(s string) bool {
	dl := distinctLetters(s)

	return len(dl) == len(g.Chars)
}

func (g Game) String() string {
	sb := strings.Builder{}

	sb.WriteString(header)
	sb.WriteString(color.MagentaString(formatDate(g.Date)))
	sb.WriteString("\n\n\n")

	sb.WriteString("letters: ")
	sb.WriteString(color.New(color.BgHiYellow, color.FgBlack).Sprint(string(g.Chars[0])))
	sb.WriteRune(' ')

	for _, c := range g.Chars[1:] {
		sb.WriteRune(c)
		sb.WriteRune(' ')
	}

	sb.WriteString("\nGuesses:\n\n")

	for idx, word := range g.Words {
		sb.WriteString(fmt.Sprintf("%d\t%s\n", idx+1, word))
	}

	return sb.String()

}

func (g *Game) Guess(s string) {
	sb := strings.Builder{}
	sb.WriteString(s)

	if g.isPangram(s) {
		sb.WriteString("\t😎 pangram 😎")
	}

	g.Words = append(g.Words, sb.String())
}

func bod() time.Time {
	t := time.Now()
	year, month, day := t.Date()

	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func selectChars(t time.Time) []rune {
	rand.Seed(t.UnixMicro())

	//#nosec
	vw := validWords[pangrammableSeeds[rand.Intn(len(pangrammableSeeds)-1)]]
	p := distinctLetters(vw)

	// shuffle chars and select first 7
	rand.Shuffle(len(p), func(i, j int) {
		p[i], p[j] = p[j], p[i]
	})

	return p
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

func distinctLetters(s string) []rune {
	runes := make([]rune, 0)

	for _, r := range s {
		if !contains(runes, r) {
			runes = append(runes, r)
		}
	}

	return runes
}

func contains(runes []rune, r rune) bool {
	for _, i := range runes {
		if r == i {
			return true
		}
	}

	return false
}
