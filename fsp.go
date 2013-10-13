package main

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

type FspServer struct {
	listener   *net.TCPListener
	policyData []byte
}

func NewFspServer() *FspServer {
	fsp := &FspServer{}
	fsp.policyData = []byte(DefaultPolicy)
	return fsp
}

func (fs *FspServer) Startup() {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:843")
	if err != nil {
		log.Fatalln(err)
	}
	fs.listener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Listening on %s", addr.String())

	for {
		conn, err := fs.listener.AcceptTCP()
		if err != nil {
			log.Fatalln(err)
		}
		go fs.handleConnection(conn)
	}
	fs.listener.Close()
}

func (fs *FspServer) LoadPolicy(src io.Reader) {
	data, err := ioutil.ReadAll(bufio.NewReader(src))
	if err != nil {
		log.Fatalln(err)
	}
	fs.policyData = append(data, []byte("\x00")...)
}

func (fs *FspServer) handleConnection(conn *net.TCPConn) {
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(time.Second * 2))
	r := bufio.NewReader(conn)
	_, err := r.ReadString('\x00')
	if err != nil {
		if err != io.EOF {
			log.Println(err)
		}
		return
	}
	conn.Write(fs.policyData)
	log.Println("Sent policy file to", conn.RemoteAddr().String())
}
