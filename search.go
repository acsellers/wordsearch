package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"runtime"

	"github.com/acsellers/wordsearch/store"
)

func main() {
	f, e := ioutil.ReadFile("/usr/share/dict/words")
	if e != nil {
		panic(e)
	}
	s, e := store.NewStore(bytes.NewBuffer(f))
	if e != nil {
		panic(e)
	}
	s.PrefilledAtLength(5, "eomhllgysnvj", map[int]rune{1: 'o'})
}
