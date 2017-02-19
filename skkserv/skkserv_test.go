package skkserv

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"reflect"
	"testing"
	"time"
)

func TestBuildResponse(t *testing.T) {
	cases := []struct {
		Input  []string
		Output string
	}{
		{[]string{"a"}, "1/a/\n"},
		{[]string{"a", "b"}, "1/a/b/\n"},
		{nil, "4\n"},
	}

	for _, tc := range cases {
		b := []byte(tc.Output)
		if got := BuildResponse(tc.Input); !reflect.DeepEqual(got, b) {
			t.Errorf("Expected: %s, got: %s", b, got)
		}
	}
}

func TestParseRequest(t *testing.T) {
	cases := []struct {
		Request []byte
		Parsed  string
	}{
		{[]byte("あ "), "あ"},
	}

	for _, tc := range cases {
		if got := ParseRequest(tc.Request); got != tc.Parsed {
			t.Errorf("Expected: %s, got: %s", tc.Parsed, got)
		}
	}
}

func TestRun(t *testing.T) {
	addr := ""
	port := "55555"
	s := SkkServ{
		Addr: addr,
		Port: port,
	}
	go s.Run()

	var conn Connection
	var err error

	for cnt := 1; cnt <= 100; cnt++ {
		conn, err = net.Dial("tcp", addr+":"+port)
		if err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	if err != nil {
		t.Errorf("Something wrong: %s", err)
	}
	defer conn.Close()
}

type TestConnection struct {
	io.Reader
	Written string
}

func (c *TestConnection) Write(b []byte) (int, error) {
	c.Written = string(b)
	return len(b), nil
}

func (c *TestConnection) Close() error {
	return nil
}

type TestHandler struct{}

func (h TestHandler) ReturnTrans(b []byte) ([]byte, error) {
	if text := string(b); text == "err " {
		return nil, fmt.Errorf("Error")
	}
	return b, nil
}

func TestHandleConnection(t *testing.T) {
	s := SkkServ{
		Name:    "test",
		Version: "version",
		Handler: &TestHandler{},
	}

	b1 := bytes.NewBufferString("1aaa ")
	b2 := bytes.NewBufferString("2")
	b3 := bytes.NewBufferString("3")
	b4 := bytes.NewBufferString("1err ")

	cases := []struct {
		Conn    TestConnection
		Written string
	}{
		{TestConnection{b1, ""}, "aaa "},
		{TestConnection{b2, ""}, "test version "},
		{TestConnection{b3, ""}, "127.0.0.1: "},
		{TestConnection{b4, ""}, "0"},
	}

	for _, tc := range cases {
		err := s.handleConnection(&tc.Conn)
		if err != nil && tc.Written != "0" {
			t.Errorf("Something wrong: %s", err)
		}
		if got := tc.Conn.Written; got != tc.Written {
			t.Errorf("Expected: %s, got: %s", tc.Written, got)
		}
	}
}
