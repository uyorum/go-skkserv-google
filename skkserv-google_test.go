package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/akiym/go-skkserv"
	"github.com/uyorum/go-skk-dictionary"
)

const DICTIONARY_FILENAME = "SKK-JISYO.L"

type TestStringGoogle struct {
	request  string
	response []string
}

type TestString struct {
	request  string
	response string
}

var tests_for_googleime = []TestStringGoogle{
	{"かんすうがたげんご", []string{"関数型言語", "かんすうがたげんご", "カンスウガタゲンゴ", "ｶﾝｽｳｶﾞﾀｹﾞﾝｺﾞ"}},
}

var tests_for_request = []TestString{
	{"1じしょ \n", "1/辞書/地所/自署/字書/自書/\n"},
	{"1わたs ", "1/渡/\n"},
	{"1わあたs ", "4わあたs\n"},
}

func TestTransliterateWithGoogle(t *testing.T) {
	for _, test := range tests_for_googleime {
		fmt.Println(test.request)
		resp, err := TransliterateWithGoogle(test.request)
		if err != nil {
			t.Errorf("Error at query.")
		}
		fmt.Println(resp)
		for i, word := range resp {
			if word != test.response[i] {
				t.Errorf("Unexpected response: %s", word)
			}
		}
	}
}

func TestRequest(t *testing.T) {
	port_num := ":55100"
	pwd, _ := os.Getwd()
	dictionary_path := pwd + "/" + DICTIONARY_FILENAME
	var server = skkserv.NewServer(port_num, &GoogleIMESKK{skkdictionary.NewSkkDictionary(dictionary_path)})
	go server.Run()
	conn, err := net.Dial("tcp", "localhost"+port_num)
	if err != nil {
		t.Errorf("Failed to connect to skkserv")
	}
	defer conn.Close()
	for _, test := range tests_for_request {
		fmt.Println(test.request)
		word, _ := encoder.String(test.request)
		fmt.Fprintf(conn, word)
		resp, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println(decoder.String(resp))
		test.response, _ = encoder.String(test.response)
		if resp != test.response {
			t.Errorf("Unexpected response: %s", resp)
		}
	}
}
