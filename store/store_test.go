package store

import (
	"bytes"
	"io/ioutil"
	"testing"

	. "github.com/acsellers/assert"
)

func TestStore(t *testing.T) {
	Within(t, func(test *Test) {
		s, e := NewStore(bytes.NewBufferString(`asdf
    asd
    free`))
		test.NoError(e)
		test.AreEqual(len(s.Words[4]), 2)
		test.AreEqual(len(s.Words[3]), 1)
		test.AreEqual(s.Lengths, []int{4, 3})
		test.AreEqual(s.WithLength(4, "farsdede"), []string{"asdf", "free"})
		test.AreEqual(s.WithLength(4, "asdff"), []string{"asdf"})
		test.AreEqual(s.WithLength(4, "isdff"), []string{})
		test.AreEqual(s.PrefilledAtLength(4, "farsdede", map[int]rune{1: 'r'}), []string{"free"})
		test.AreEqual(s.Longest("asdfree"), []string{"asdf", "free", "asd"})
	})
}

func BenchmarkSearch(b *testing.B) {
	f, e := ioutil.ReadFile("/usr/share/dict/words")
	if e != nil {
		panic(e)
	}
	s, e := NewStore(bytes.NewBuffer(f))
	if e != nil {
		panic(e)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.PrefilledAtLength(5, "eomhllgysnvj", map[int]rune{1: 'o'})
	}

}
