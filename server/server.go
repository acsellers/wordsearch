package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/acsellers/wordsearch/store"
)

var (
	s                             store.Store
	indexPage                     string
	resultPrelude, resultPostlude string
	help                          = `This API uses HTML Form Variables and responds with
  JSON responses. There are two required form variables: letters
  and length. Length is either -1 for the longest possible word
  that could be made with the passed letters or a base 10 integer
  for all possible words to be made.
  `
)

func main() {
	f, e := ioutil.ReadFile("/usr/share/dict/words")
	if e != nil {
		panic(e)
	}
	s, e = store.NewStore(bytes.NewBuffer(f))
	if e != nil {
		panic(e)
	}

	ip, e := ioutil.ReadFile("assets/index.html")
	if e != nil {
		panic(e)
	}
	indexPage = string(ip)

	rp, e := ioutil.ReadFile("assets/result.html")
	if e != nil {
		panic(e)
	}
	rc := strings.Split(string(rp), "CONTENT")
	resultPrelude, resultPostlude = rc[0], rc[1]

	SetRoutes()

	log.Fatal(http.ListenAndServe(":8001", nil))
}

func NetError(w io.Writer, message string) {
	fmt.Fprintf(w, `{"error":"%s", "help":%s"}`, message, help)
}

func HtmlError(w io.Writer, message string) {
	fmt.Fprint(w, resultPrelude, message, resultPostlude)
}

func HtmlResult(w io.Writer, results []string) {
	fmt.Fprint(
		w,
		resultPrelude,
	)

	for len(results) > 0 {
		if len(results) >= 6 {
			fmt.Fprint(
				w,
				"<tr><td>",
				strings.Join(results[:6], "</td><td>"),
				"</td></tr>",
			)
			results = results[6:]
		} else {
			fmt.Fprint(
				w,
				"<tr><td>",
				strings.Join(results, "</td><td>"),
				"</td></tr>",
			)
			results = []string{}
		}
	}

	fmt.Fprint(
		w,
		resultPostlude,
	)
}
func JSONSearch(w http.ResponseWriter, r *http.Request) {
	if e := r.ParseForm(); e == nil {
		l := r.Form.Get("letters")
		if l == "" {
			NetError(w, "Letters were not given")
			return
		}
		wl := r.Form.Get("length")
		if wl == "" {
			NetError(w, "Length was not found")
			return
		}
		li, e := strconv.ParseInt(wl, 10, 32)
		if e != nil {
			NetError(w, "Length is not valid")
			return
		}
		var words []string
		if li == -1 {
			words = s.Longest(l)
		} else {
			words = s.WithLength(int(li), l)
		}
		io.WriteString(w, `{words:"`)
		io.WriteString(w, strings.Join(words, `","`))
		io.WriteString(w, `"}`)
	}
}
func Search(w http.ResponseWriter, r *http.Request) {
	if e := r.ParseForm(); e == nil {
		l := r.Form.Get("letters")
		if l == "" {
			HtmlError(w, "Letters were not given")
			return
		}
		wl := r.Form.Get("length")
		if wl == "" {
			HtmlError(w, "Length was not found")
			return
		}
		li, e := strconv.ParseInt(wl, 10, 32)
		if e != nil {
			HtmlError(w, "Length is not valid")
			return
		}
		var words []string
		if li == -1 {
			words = s.Longest(l)
		} else {
			words = s.WithLength(int(li), l)
		}
		HtmlResult(w, words)
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, indexPage)
}

func SetRoutes() {
	http.HandleFunc("/results", Search)
	http.HandleFunc("/results.json", JSONSearch)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	//http.Handle("/assets/", http.FileServer(http.Dir("assets")))
	http.HandleFunc("/", Index)
}
