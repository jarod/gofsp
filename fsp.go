package main

import (
  "net"
  "bufio"
  "log"
	"os"
	"io"
	"io/ioutil"
)

const (
	DefaultPolicy = "<cross-domain-policy><allow-access-from domain=\"*\" to-ports=\"*\"/></cross-domain-policy>\x00"
)

type FspServer struct {
  listener *net.TCPListener 
	policyContent string
}

func NewFspServer() *FspServer {
	fsp := &FspServer{}
	fsp.policyContent = DefaultPolicy
  return fsp
}

func (fs *FspServer) Startup() {
  addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:843")
  if (err != nil) {
    log.Fatalln(err.String())
  }
  fs.listener, err = net.ListenTCP("tcp", addr)
  if err != nil {
    log.Fatalln(err.String())
  }
	log.Printf("Listening on %s", addr.String())

  for {
    conn, err := fs.listener.AcceptTCP()
		if (err != nil) {
			log.Fatalln(err.String())
		}
    go fs.read(conn)
  }
  fs.listener.Close()
}

func (fs *FspServer) LoadPolicy(src io.Reader) {
	data, err := ioutil.ReadAll(bufio.NewReader(src))
	if err != nil {
		log.Fatalln(err.String())
	}
	fs.policyContent = string(data) + "\x00"
}

func (fs *FspServer) read(conn *net.TCPConn) {
	defer conn.Close()

	r := bufio.NewReader(conn)
	for {
		_,  err := r.ReadString('\x00')
		if (err == os.EOF) {
			return
		} else if (err != nil) {
			log.Fatalln(err.String())
		}

		w := bufio.NewWriter(conn)
		w.WriteString(fs.policyContent)
		w.Flush()
		log.Printf("Sent policy file to %s", conn.RemoteAddr().String())
	}
}
