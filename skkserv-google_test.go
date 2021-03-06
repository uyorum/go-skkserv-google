package main

import (
	"fmt"
	"testing"
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

type TestStringLog struct {
	request  string
	response []string
	api      bool
	log      string
}

var tests_for_googleime = []TestStringGoogle{
	{"かんすうがたげんご", []string{"関数型言語", "かんすうがたげんご", "カンスウガタゲンゴ", "ｶﾝｽｳｶﾞﾀｹﾞﾝｺﾞ"}},
}

var tests_for_request = []TestString{
	{"1じしょ \n", "1/辞書/地所/自署/字書/自書/\n"},
	{"1わたs ", "1/渡/\n"},
	{"1わあたs ", "4わあたs\n"},
}

var tests_for_log = []TestStringLog{
	{"じしょ", []string{"辞書", "地所", "自署", "字書", "自書"}, true, "{\"api\": {\"じしょ\": [\"辞書\", \"地所\", \"自署\", \"字書\", \"自書\"]}}"},
	{"わたs", []string{"渡"}, false, "{\"jisyo\": {\"わたs\": [\"渡\"]}}"},
	{"わあたs", []string{}, false, "{\"jisyo\": {\"わあたs\": []}}"},
}

func TestLog(t *testing.T) {
	for _, test := range tests_for_log {
		log := Log(test.request, test.response, test.api)
		if log != test.log {
			t.Errorf("Unexpected message: %s\nExpected: %s", log, test.log)
		}
	}
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

// func TestReturnTrans(t *testing.T) {
// 	port := ":55100"
// 	pwd, _ := os.Getwd()
// 	dictionary_path := pwd + "/" + DICTIONARY_FILENAME
// 	var server = skkserv.NewServer(port_num, &GoogleIMESKK{skkdictionary.NewSkkDictionary(dictionary_path)})
// 	go server.Run()
// 	conn, err := net.Dial("tcp", "localhost"+port_num)
// 	if err != nil {
// 		t.Errorf("Failed to connect to skkserv")
// 	}
// 	defer conn.Close()
// 	for _, test := range tests_for_request {
// 		fmt.Println(test.request)
// 		word, _ := encoder.String(test.request)
// 		fmt.Fprintf(conn, word)
// 		resp, _ := bufio.NewReader(conn).ReadString('\n')
// 		fmt.Println(decoder.String(resp))
// 		test.response, _ = encoder.String(test.response)
// 		if resp != test.response {
// 			t.Errorf("Unexpected response: %s", resp)
// 		}
// 	}
// }
