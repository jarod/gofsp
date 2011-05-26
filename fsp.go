package main

import (
  "net"
  "bufio"
  "log"
	"os"
)

const (
  ServePort = 843
)

type FspServer struct {
  listener *net.TCPListener 
}

func NewFspServer() *FspServer {
  return &FspServer{}
}

func (fs *FspServer) Startup() {
  addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:843")
  if (err != nil) {
    log.Fatalf("net.ResolveTCPAddr: %s", err.String())
  }
  fs.listener, err = net.ListenTCP("tcp", addr)
  if (err != nil) {
    log.Fatalf("net.ListenTCP: %s", err.String())
  }
	log.Printf("Listening on %s", addr.String())

  for {
    conn, err := fs.listener.AcceptTCP()
		if (err != nil) {
			log.Fatalf("AcceptTCP: %s", err.String())
		}
    go read(conn)
  }
  fs.listener.Close()
}

func read(conn *net.TCPConn) {
	defer conn.Close()

	r := bufio.NewReader(conn)
	for {
		_,  err := r.ReadString('\x00')
		if (err == os.EOF) {
			return
		} else if (err != nil) {
			log.Fatalf("ReadString: %s", err.String())
		}

		w := bufio.NewWriter(conn)
		w.WriteString(`<cross-domain-policy><allow-access-from domain="*" to-ports="*"/></cross-domain-policy>`)
		w.WriteString("\x00")
		w.Flush()
		log.Printf("Sent policy file to %s", conn.RemoteAddr().String())
	}
}
