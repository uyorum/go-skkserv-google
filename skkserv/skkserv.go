package skkserv

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
)

type Handler interface {
	ReturnTrans([]byte) ([]byte, error)
}

type Connection interface {
	io.Reader
	Write([]byte) (int, error)
	Close() error
}

type SkkServ struct {
	Name    string
	Version string
	Addr    string
	Port    string
	Handler Handler
}

func BuildResponse(list []string) []byte {
	if len(list) == 0 {
		return []byte("4\n")
	}
	response := "1/" + strings.Join(list, "/") + "/\n"
	return []byte(response)
}

func ParseRequest(b []byte) string {
	return string(b[:len(b)-1])
}

func NewSkkServ(name string, version string, handler Handler) *SkkServ {
	s := SkkServ{
		Name:    name,
		Version: version,
		Handler: handler,
	}
	s.Addr = "0.0.0.0"
	s.Port = "1178"
	return &s
}

func (s SkkServ) Run() {
	ln, err := net.Listen("tcp", s.Addr+":"+s.Port)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s SkkServ) returnVersion() ([]byte, error) {
	return []byte(s.Name + " " + s.Version + " "), nil
}

func (s SkkServ) returnHost() ([]byte, error) {
	return []byte("127.0.0.1" + ": "), nil
}

func (s SkkServ) handleConnection(conn Connection) error {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		c, err := reader.ReadByte()
		if err != nil {
			return nil // end of buffer
		}
		switch c {
		case '0':
			return nil // disconnect
		case '1':
			buf, err := reader.ReadBytes(' ')
			if err != nil {
				return err
			}
			b, err := s.Handler.ReturnTrans(buf)
			if err != nil {
				conn.Write([]byte("0"))
				return err
			}
			conn.Write(b)
		case '2':
			b, err := s.returnVersion()
			if err != nil {
				return err
			}
			conn.Write(b)
		case '3':
			b, err := s.returnHost()
			if err != nil {
				return err
			}
			conn.Write(b)
		}
	}
}
