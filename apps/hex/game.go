package hex

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jspc/jshell/apps"
)

type Game struct {
	Date         time.Time
	Chars        []rune
	Words        []string
	Score        int
	TargetScore  int
	PangramCount int

	wordlist map[string]int
}

func NewGame() (g *Game) {
	t := apps.Bod()

	g = &Game{
		Date:  t,
		Chars: selectChars(t),
		Words: make([]string, 0),
	}

	return
}

func (g *Game) setWordlist() {
	indices := mappings[g.Chars[0]]
	targetScore := 0
	g.PangramCount = 0

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
			targetScore += g.score(w)

			if g.isPangram(w) {
				g.PangramCount += 1
			}
		}
	}

	g.TargetScore = targetScore
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
	pangramCount := 0
	words := make([]string, len(g.Words))

	for idx, w := range g.Words {
		words[idx] = fmt.Sprintf("%d\t%s", idx+1, w)

		if g.isPangram(w) {
			pangramCount += 1
			words[idx] += "\t😎 pangram 😎"
		}
	}

	if len(words) > 10 {
		words = words[len(words)-10:]
	}

	sb := strings.Builder{}

	sb.WriteString(header)
	sb.WriteString(color.MagentaString(apps.FormatDate(g.Date)))
	sb.WriteString("\n")
	sb.WriteString("Score: ")
	sb.WriteString(color.HiCyanString("%d ", g.Score))
	sb.WriteString("out of a maximum of: ")
	sb.WriteString(color.HiCyanString("%d\n", g.TargetScore))
	sb.WriteString("You have found ")
	sb.WriteString(color.HiCyanString("%d/%d", pangramCount, g.PangramCount))
	sb.WriteString(" pangram(s)!")

	sb.WriteString("\n\n")

	sb.WriteString("Letters: \n\t")
	sb.WriteString(color.New(color.BgHiYellow, color.FgBlack).Sprint(string(g.Chars[0])))
	sb.WriteRune(' ')

	for _, c := range g.Chars[1:] {
		sb.WriteRune(c)
		sb.WriteRune(' ')
	}

	sb.WriteString("\n\n\nGuesses:\n\n")

	for _, word := range words {
		sb.WriteString(word)
		sb.WriteString("\n")
	}

	return sb.String()

}

func (g *Game) Guess(s string) {
	g.Words = append(g.Words, s)
	g.Score += g.score(s)
}

func (g Game) guessed(s string) bool {
	for _, guess := range g.Words {
		if guess == s {
			return true
		}
	}

	return false
}

func (g Game) score(s string) int {
	score := 1
	if len(s) > 4 {
		score = len(s)
	}

	if g.isPangram(s) {
		score += 7

	}

	return score
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
