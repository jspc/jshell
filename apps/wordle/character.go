package wordle

import (
	"github.com/fatih/color"
)

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
