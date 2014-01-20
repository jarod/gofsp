package fsp

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"net"
	"time"
)

const (
	DefaultPolicy = "<cross-domain-policy><allow-access-from domain=\"*\" to-ports=\"*\"/></cross-domain-policy>\x00"
)

type Server struct {
	b []byte
}

func NewServer() *Server {
	s := &Server{}
	s.b = []byte(DefaultPolicy)
	return s
}

func (s *Server) ListenAndServe() {
	addr, err := net.ResolveTCPAddr("tcp", ":843")
	if err != nil {
		log.Fatalln(err)
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("fsp: Listening on %s", addr.String())

	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			log.Fatalln(err)
		}
		go s.handleConnection(conn)
	}
	l.Close()
}

func (s *Server) LoadPolicy(src io.Reader) {
	data, err := ioutil.ReadAll(bufio.NewReader(src))
	if err != nil {
		log.Fatalln(err)
	}
	s.b = append(data, []byte("\x00")...)
}

func (s *Server) handleConnection(conn *net.TCPConn) {
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(time.Second))
	r := bufio.NewReader(conn)
	_, err := r.ReadString('\x00')
	if err != nil {
		if err != io.EOF {
			log.Println(err)
		}
		return
	}
	conn.Write(s.b)
	log.Println("fsp: Sent policy file to", conn.RemoteAddr().String())
}
