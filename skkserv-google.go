package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/akiym/go-skkserv"
	"github.com/uyorum/go-skk-dictionary"
	"golang.org/x/text/encoding/japanese"
)

var port_num *int
var dictionary_path_list []string

// UTF-8 to EUCJP
var encoder = japanese.EUCJP.NewEncoder()

// EUCJP to UTF-8
var decoder = japanese.EUCJP.NewDecoder()

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage of %s:
  %s [OPTIONS] /path/to/SKK-JISYO.L
Options
`, os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	port_num = flag.Int("p", 1178, "Port number skkserv uses")
	flag.Parse()

	dictionary_path_list = flag.Args()
}

type GoogleIMESKK struct {
	d *skkdictionary.SkkDictionary
}

func (s *GoogleIMESKK) Request(text string) ([]string, error) {
	var words []string
	var text_u string
	var err error

	if skkdictionary.IsOkuriAri(text + " ") {
		str := s.d.Search(text + " ")
		if str == "" {
			return nil, nil
		}
		words = strings.Split(str[1:len(str)-1], "/")
	} else {
		text_u, err = decoder.String(text)
		if err != nil {
			return nil, err
		}
		words, err = TransliterateWithGoogle(text_u)

		// Failed to communicate with server (may be offline)
		// use skk dictionary
		if len(words) == 0 || err != nil {
			str := s.d.Search(text + " ")
			if str == "" {
				return nil, nil
			}
			words = strings.Split(str[1:len(str)-1], "/")
			return words, nil
		}

		for i, word := range words {
			words[i], err = encoder.String(word)
			if err != nil {
				words[i] = ""
			}
		}
	}

	return words, nil
}

func TransliterateWithGoogle(text string) (words []string, err error) {
	v := url.Values{"langpair": {"ja-Hira|ja"}, "text": {text + ","}}
	resp, err := http.Get("http://www.google.com/transliterate?" + v.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	var w [][]interface{}
	if err := dec.Decode(&w); err != nil {
		return nil, err
	}
	for _, v := range w[0][1].([]interface{}) {
		words = append(words, v.(string))
	}
	return words, nil
}

func Log(request string, response []string, api bool) (log string) {
	if api {
		log = "{\"api\": {\""
	} else {
		log = "{\"jisyo\": {\""
	}

	response_list := strings.Join(response, "\", \"")
	log = log + request + "\": [\"" + response_list + "\"]}}"

	fmt.Println(log)
	return log
}

func main() {
	if len(dictionary_path_list) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	var server = skkserv.NewServer(":"+strconv.Itoa(*port_num), &GoogleIMESKK{skkdictionary.NewSkkDictionary(dictionary_path_list[0])})
	server.Run()
}
