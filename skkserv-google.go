package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/uyorum/go-skk-dictionary"
	"github.com/uyorum/go-skkserv-google/skkserv"
	"golang.org/x/text/encoding/japanese"
)

const AppName = "skkserv-google"
const AppVersion = "0.0.1"

var port_num *string
var verbose *bool
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

	port_num = flag.String("p", "1178", "Port number skkserv uses")
	verbose = flag.Bool("v", false, "Print request and respons to stdout")
	flag.Parse()

	dictionary_path_list = flag.Args()
}

type GoogleIMESKK struct {
	d *skkdictionary.SkkDictionary
}

func (s *GoogleIMESKK) ReturnTrans(b []byte) ([]byte, error) {
	var words, words_u []string
	var text_u string
	var err error

	text := skkserv.ParseRequest(b)

	// Whether used Google IME API
	api := false

	if skkdictionary.IsOkuriAri(text + " ") {
		str := s.d.Search(text + " ")
		if str == "" {
			return skkserv.BuildResponse([]string{}), nil
		}
		words = strings.Split(str[1:len(str)-1], "/")
	} else {
		text_u, err = decoder.String(text)
		if err != nil {
			return nil, err
		}
		words_u, err = TransliterateWithGoogle(text_u)
		if len(words_u) == 0 || err != nil {
			// Failed to communicate with server (may be offline)
			// use skk dictionary
			str := s.d.Search(text + " ")
			if str == "" {
				return skkserv.BuildResponse([]string{}), nil
			}
			words = strings.Split(str[1:len(str)-1], "/")
		} else {
			api = true
			for _, word_u := range words_u {
				word, err := encoder.String(word_u)
				if err != nil {
					word = ""
				}
				words = append(words, word)
			}
		}
	}

	if *verbose {
		go func() {
			if !api {
				for _, word := range words {
					word_u, err := decoder.String(word)
					if err != nil {
						word_u = ""
					}
					words_u = append(words_u, word_u)
				}
				text_u, err = decoder.String(text)
				if err != nil {
					return
				}
			}
			Log(text_u, words_u, api)
		}()
	}
	return skkserv.BuildResponse(words), nil
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

	var response_list string
	if len(response) == 0 {
		response_list = ""
	} else {
		response_list = "\"" + strings.Join(response, "\", \"") + "\""
	}

	log = log + request + "\": [" + response_list + "]}}"

	fmt.Println(log)
	return log
}

func main() {
	if len(dictionary_path_list) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	s := skkserv.NewSkkServ(AppName, AppVersion, &GoogleIMESKK{skkdictionary.NewSkkDictionary(dictionary_path_list[0])})
	s.Port = *port_num
	s.Run()
}
