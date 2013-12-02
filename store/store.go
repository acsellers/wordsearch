package store

import (
	"bufio"
	"io"
	"sort"
	"strings"
)

type Store struct {
	Words   map[int][]Word
	Count   int
	Lengths []int
}

type Word struct {
	Vowels, Consonants string
	Original           string
	Lengths            []int
}

func NewStore(source io.Reader) (Store, error) {
	s := Store{Words: make(map[int][]Word)}

	sc := bufio.NewScanner(source)
	for sc.Scan() {
		t := strings.TrimSpace(sc.Text())
		l := len([]rune(t))
		s.Words[l] = append(s.Words[l],
			Word{
				Vowels:     SortVowels(t),
				Consonants: SortConsonants(t),
				Original:   t,
			},
		)
		s.Count++
	}
	for l, _ := range s.Words {
		s.Lengths = append(s.Lengths, l)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(s.Lengths)))
	if err := sc.Err(); err != nil {
		return s, err
	}
	return s, nil
}

func SortConsonants(s string) string {
	ls := strings.Split(s, "")
	c := []string{}
	for _, l := range ls {
		switch l {
		case "a", "e", "i", "o", "u":
			// do nothing
		default:
			c = append(c, l)
		}
	}
	sort.Strings(c)
	return strings.Join(c, "")
}

func SortVowels(s string) string {
	ls := strings.Split(s, "")
	c := []string{}
	for _, l := range ls {
		switch l {
		case "a", "e", "i", "o", "u":
			c = append(c, l)
		}
	}

	sort.Strings(c)
	return strings.Join(c, "")
}

func (s Store) WithLength(l int, letters string) []string {
	matching := []string{}
	lv := []rune(SortVowels(letters))
	lc := []rune(SortConsonants(letters))
	var vi, ci int

WordLoop:
	for _, w := range s.Words[l] {
		vi, ci = 0, 0

		// check for vowel matching
		for _, wv := range w.Vowels {
			if vi == len(lv) {
				continue WordLoop
			}

			if lv[vi] == wv {
				vi++
			} else if wv < lv[vi] {
				continue WordLoop
			} else {

				for vi < len(lv) && lv[vi] < wv {
					vi++
				}

				if vi == len(lv) {
					continue WordLoop
				}
				if lv[vi] == wv {
					vi++
				} else {
					continue WordLoop
				}
			}
		}

		// check for consonant matching
		for _, wc := range w.Consonants {
			if ci == len(lc) {
				continue WordLoop
			}

			if lc[ci] == wc {
				ci++
			} else if wc < lc[ci] {
				continue WordLoop
			} else {

				for ci < len(lc) && lc[ci] < wc {
					ci++
				}

				if ci == len(lc) {
					continue WordLoop
				}
				if lc[ci] == wc {
					ci++
				} else {
					continue WordLoop
				}
			}
		}

		// passed all tests, append to matching
		matching = append(matching, w.Original)

	}

	return matching
}

// tries to return at least 10 matches for the letters, without reuse
// returning less than 10 is possible if there are few matches
// returning more than 10 is also possible
func (s Store) Longest(letters string) []string {
	matching := []string{}
	for _, length := range s.Lengths {
		if length <= len([]rune(letters)) {
			matching = append(matching, s.WithLength(length, letters)...)
		}
		if len(matching) > 10 {
			return matching
		}
	}
	return matching
}
func (s Store) PrefilledAtLength(l int, letters string, prefilled map[int]rune) []string {
	matching := []string{}

MatchLoop:
	for _, m := range s.WithLength(l, letters) {
		for i, l := range prefilled {
			if l != []rune(m)[i] {
				continue MatchLoop
			}
		}
		matching = append(matching, m)
	}
	return matching
}

func (s Store) PrefilledLongest(l int, letters string, prefilled map[int]rune) []string {
	matching := []string{}

	for _, length := range s.Lengths {
		if length <= len([]rune(letters)) {
		PrefillLoop:
			for _, w := range s.WithLength(length, letters) {
				for i, l := range prefilled {
					if l != []rune(w)[i] {
						continue PrefillLoop
					}
				}
				matching = append(matching, w)
			}
		}
		if len(matching) > 10 {
			return matching
		}
	}

	return matching
}
