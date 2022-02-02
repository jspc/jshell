package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/dave/jennifer/jen"
)

const (
	url = "http://www.mieliestronk.com/corncob_caps.txt"
)

var (
	// Store the index of found words against each distinct letter in a word.
	//
	// This will allow us to lookup words by the centre letter to minimise any
	// lookups we need to do in order to validate words
	letterIdxs = make(map[rune][]jen.Code)

	// wordlist contains any valid word
	wordlist = make([]jen.Code, 0)

	// pangrammable contains any valid word with 7 distinct letters
	pangrammable = make([]jen.Code, 0)

	// chars contains the distinct characters in our words. I assume this list
	// will contain *every* letter, but it's best to be safe
	chars    = make([]rune, 0)
	charLits = make([]jen.Code, 0)
)

func main() {
	// get wordlist
	r, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer r.Body.Close()

	buf := bufio.NewReader(r.Body)

	for {
		word, err := buf.ReadString(byte('\n'))
		if err != nil {
			break
		}

		word = strings.TrimSuffix(word, "\r\n")

		// We only care about words which are four letters or more
		if len(word) < 4 {
			continue
		}

		dl := distinctLetters(word)

		// Any word with more than 7 different letters (the number of
		// letters we generate) is unusable here
		if len(dl) > 7 {
			continue
		}

		wordlist = append(wordlist, jen.Lit(word))
		idx := len(wordlist) - 1

		if len(dl) == 7 {
			pangrammable = append(pangrammable, jen.Lit(idx))
		}

		for _, r := range dl {
			letterIdxs[r] = append(letterIdxs[r], jen.Lit(idx))

			if !contains(chars, r) {
				chars = append(chars, r)
				charLits = append(charLits, jen.LitRune(r))
			}
		}
	}

	if !errors.Is(err, io.EOF) {
		panic(err)
	}

	// generate file
	f := jen.NewFile("hex")
	f.HeaderComment("Code generated from generator/main.go DO NOT EDIT ")

	f.Comment("All of the valid words we have")
	f.Var().Id("validWords").Op("=").Index().String().Values(wordlist...)

	f.Comment("Indices of words which are valid pangrams (and, thus, valid letters to choose")
	f.Var().Id("pangrammableSeeds").Op("=").Index().Int().Values(pangrammable...)

	f.Comment("Mapping of words to 'middle characters'")
	f.Var().Id("mappings").Op("=").Map(jen.Rune()).Index().Int().Values(jen.DictFunc(func(d jen.Dict) {
		for letter, indexes := range letterIdxs {
			d[jen.LitRune(rune(letter))] = jen.Index().Int().Values(indexes...)
		}
	}))

	f.Comment("Distinct letters from our word list")
	f.Var().Id("chars").Op("=").Index().Rune().Values(charLits...)

	b := strings.Builder{}

	err = f.Render(&b)
	if err != nil {
		panic(err)
	}

	fmt.Println(b.String())
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
