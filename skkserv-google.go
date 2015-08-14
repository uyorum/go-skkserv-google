package main

import (
	"github.com/akiym/go-skkserv"
	"code.google.com/p/go.text/encoding/japanese"
	"code.google.com/p/go.text/transform"
	"encoding/json"
	"net/http"
	"net/url"
	"io/ioutil"
	"strings"
	"flag"
	"strconv"
)

var portNum *int
var dics []string

func init() {
	portNum = flag.Int("p", 1178, "Port number skkserv uses")
	flag.Parse()
	dics = flag.Args()
}

type GoogleIMESKK struct{}

func (s *GoogleIMESKK) Request(text string) ([]string, error) {
	words, err := Transliterate(text)
	if err != nil {
		return nil, err
	}
	return words, nil
}

func utf8_to_eucjp(str string) (string, error) {
	ret, err := ioutil.ReadAll(transform.NewReader(strings.NewReader(str), japanese.EUCJP.NewEncoder()))
	if err != nil {
		return "", err
	}
	return string(ret), err
}

func eucjp_to_utf8(str string) (string, error) {
	ret, err := ioutil.ReadAll(transform.NewReader(strings.NewReader(str), japanese.EUCJP.NewDecoder()))
	if err != nil {
		return "", err
	}
	return string(ret), err
}

func Transliterate(text string) (words []string, err error) {
	text, err = eucjp_to_utf8(text)
	if err != nil {
		return nil, err
	}
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
		word := v.(string)
		result, err := utf8_to_eucjp(word)
		if err == nil {
			words = append(words, result)
		}
	}
	return words, nil
}

func main() {
	var server = skkserv.NewServer(":" + strconv.Itoa(*portNum), &GoogleIMESKK{})
	server.Run()
}
