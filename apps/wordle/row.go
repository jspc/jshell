package wordle

import (
	"strings"
)

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
